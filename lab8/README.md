# Lab 8: OSS Vulnerabilities

## Introduction

In Lab 8, you will learn how an open source project like free5GC prevents vulnerabilities as early as possible in the Software Development Life Cycle (SDLC).

GitHub already provides [Dependabot](https://docs.github.com/en/code-security/dependabot) out of the box, which makes dependency vulnerability management much easier for OSS maintainers. However, a tool alone is not enough: without a well-defined workflow, it is still hard to handle a critical security risk at the first moment it appears. What we really need is a pipeline that tells us **what we ship**, **what is broken in what we ship**, and **who fixes it, when**.

## Goals of this lab

- Understand what "shift-left security" means in the SDLC
- Understand what an SBOM is, and why free5GC needs it at different layers
- Understand how static analysis, dynamic analysis, and fuzzing complement each other
- Understand how LLM agents can be used to accelerate vulnerability discovery and remediation
- Understand how a community security reporting (disclosure) pipeline works

## Shift Left: Why Earlier Is Cheaper

The later a vulnerability is found, the more expensive it is to fix. A bug caught by a linter on a developer's laptop costs a few minutes; the same bug found in a deployed 5G Core costs a CVE, a patch release, an upgrade for every operator, and a lot of reputation.

So the principle is simple: **push every security check as close to the developer as possible**.

```
design -> code -> PR (CI) -> release -> deploy -> production
  ^        ^        ^           ^          ^          ^
threat   linter   SAST /      SBOM /    image /   monitoring /
model    / IDE    unit test   signing   config     CVE watch
                  / fuzz      / scan     scan
```

In Lab 7 you have already built a CI pipeline (linting, unit test, build, integration test). Lab 8 is about adding the **security stages** into that same pipeline.

## Part 1: Know What You Ship (SBOM)

An **SBOM** (Software Bill of Materials) is a machine-readable inventory of every component inside your software: direct dependencies, transitive dependencies, versions, licenses, and hashes. Common formats are [SPDX](https://spdx.dev/) and [CycloneDX](https://cyclonedx.org/).

> Why do we care? When a new CVE is published (think about Log4Shell), the very first question is: *"Are we affected?"*
> Without an SBOM, answering that question takes days of manual grepping. With an SBOM, it is one query.

The most basic thing to do is to **generate an SBOM in CI/CD for every build**, and publish it together with the release artifacts.

### free5GC has more than one layer

free5GC is not a single binary; it is a set of repositories and deliverables, and each layer has a different risk profile and a different owner. They must be managed **separately**:

| Layer | Example | What the SBOM covers | Typical risk |
| --- | --- | --- | --- |
| Library / pkg | `free5gc/openapi`, `free5gc/nas`, `free5gc/ngap`, `free5gc/util` | Go module dependencies | A vulnerable parser is inherited by every NF that imports it |
| Network Function | `free5gc/amf`, `free5gc/smf`, `free5gc/upf` ... | Go modules plus the libraries above | NF-specific logic, SBI exposure |
| Main repo | `free5gc/free5gc` | All NFs as submodules, build toolchain | The version combination that is actually released |
| Container | `free5gc/free5gc-compose` | Base image, OS packages, Compose configuration | Outdated base image, root containers, exposed ports |
| Orchestration | `free5gc/free5gc-helm` | Chart dependencies, Kubernetes manifests | Over-privileged pods, host network, secrets in values |

A CVE in a base image is *not* the same problem as a CVE in `golang.org/x/net`: different severity, different fix, different release train. That is why the SBOM should be produced **per layer** and aggregated afterwards.

### Suggested CI/CD flow

1. Generate the SBOM on every build (for example with `syft`, `cyclonedx-gomod`, or `docker sbom`).
2. Attach the SBOM as a build artifact / release asset, so downstream users can consume it.
3. Scan the SBOM against a vulnerability database (for example with `grype`, `osv-scanner`, `govulncheck`, or `trivy`).
4. Fail the pipeline (or open an issue automatically) when a `HIGH` / `CRITICAL` finding appears.
5. Let Dependabot handle the routine version bumps, and let the SBOM scan be the safety net for what Dependabot cannot see (vendored code, base images, chart dependencies).

> [!NOTE]
> An SBOM only tells you whether the **material** you used is safe. It says nothing about the bugs you wrote yourself. That is what Part 2 is for.

## Part 2: Know What You Wrote (SAST, DAST, Fuzzing)

To kill coding mistakes in the cradle, we introduce three complementary techniques into CI/CD.

### Static Analysis (SAST)

Static analysis inspects source code without running it. It is fast, runs on every PR, and catches whole classes of bugs: nil dereference, unchecked errors, injection, hardcoded credentials, unsafe type assertions.

Useful for a Go project like free5GC:
- `go vet`, `staticcheck`, `golangci-lint` (you already run linting in Lab 7 -- just enable the security linters such as `gosec`)
- [CodeQL](https://codeql.github.com/), which GitHub can run natively on every PR
- Secret scanning with push protection, so supported secrets are blocked before they reach the Git history

The cost of SAST is **false positives**. Triage them, suppress them explicitly with a reason, and keep the signal-to-noise ratio high -- a pipeline that always fails is a pipeline that everyone ignores.

### Dynamic Analysis (DAST)

Dynamic analysis observes the system while it is running. In a 5G Core this is especially valuable, because most of the interesting bugs only appear when real signaling flows through the NFs.

Examples:
- Run the integration test (Lab 7) with the Go race detector (`-race`) and with sanitizers enabled
- Probe the SBI (Lab 4) with an API security scanner: malformed JSON, missing OAuth2 token, wrong content type, oversized payload
- Replay malformed NAS / NGAP / PFCP messages against a running deployment (Lab 5 teaches you how to capture the good ones first)

### Fuzzing

Fuzzing feeds automatically generated, mutated inputs into a function and looks for crashes, hangs, or assertion failures. It is extremely effective on **protocol decoders**, and a 5G Core is full of them: NAS, NGAP, PFCP, GTP-U, ASN.1/PER, and JSON in the SBI.

Go has a built-in fuzzer since Go 1.18, so a fuzz target can live right next to the unit tests you wrote in Lab 0.

Real case: CVE-2022-43677 in free5GC was found exactly this way. Read
[Fuzz Testing in Go: Discovering Vulnerabilities and Analyzing a Real Case (CVE-2022-43677)](https://free5gc.org/blog/20230809/main/)
before you continue -- it is the best reference for this part of the lab.

Practical tips:
- Every crash found by the fuzzer should become a **regression unit test** (the corpus file goes into `testdata/`)
- Run a short fuzz session on every PR, and a long one nightly -- fuzzing is a marathon, not a sprint
- Fuzz the *decoding* path first; that is where attacker-controlled bytes arrive

## Part 3: The LLM Era -- Security Agents

In the LLM era there are many new variations we can play with. From the easiest to the hardest:

### Level 1: Security guard skill plus automated code review

Put a "security guard" skill / instruction file in each repository (for example under `.github/`), describing the project's threat model, the protocol layers, the sensitive paths (SBI handlers, NAS decoders, configuration parsing), and the coding rules that must never be broken.

Then let a coding agent perform an **automated code review** on every PR. Unlike a linter, an agent can reason about *design* problems: a missing authorization check, a token validated in the wrong place, an error path that leaks internal state, a log line that prints a SUPI.

This is the cheapest thing to build and it gives value on day one.

### Level 2: CVE pipeline with specialized agents

Build an automated pipeline that continuously collects 5GC-related CVEs -- not only free5GC's own advisories, but also other implementations (Open5GS, OAI, commercial cores) and the shared dependencies. Public sources include the NVD feed, the GitHub Advisory Database, and OSV.

Then chain a few specialized agents together:

- **Collector agent**: fetch and normalize new CVEs, keeping the ones relevant to 5GC / Go / the SBOM inventory from Part 1
- **Analyzer agent**: map the CVE onto our codebase -- is the vulnerable function reachable? which NF, which layer, which release?
- **Hotfix agent**: draft a patch plus a regression test, and open a PR
- **Verifier agent**: run the CI (lint, unit test, integration test) and summarize the remaining risk for the maintainer

> [!CAUTION]
> The agent proposes, a **human maintainer decides**. Never let an agent merge a security patch by itself, and never let it publish an advisory automatically. LLMs hallucinate, and a wrong security fix is worse than no fix.

### Level 3: Red team / Blue team agents

The hardest (and the most fun) one: build a **red team agent** that aggregates common exploitation techniques against a 5G Core -- SBI access without a valid OAuth2 token, NAS replay, malformed NGAP, PFCP session hijacking, signaling storms -- and executes them against a lab deployment. In parallel, a **blue team agent** watches the logs and metrics, detects the attack, and proposes the mitigation.

This turns security testing into a continuous, self-improving exercise instead of a one-off penetration test.

> [!WARNING]
> Only run offensive tooling against systems you own or are explicitly authorized to test, and stay within the agreed scope. Unauthorized testing may be illegal.

## Part 4: Use the Power of the Community

The last piece is process, not technology. An open source project needs a complete and *visible* security reporting pipeline:

1. **A published policy**: a `SECURITY.md` telling people what is in scope, which versions are supported, and where to report -- privately, through GitHub Security Advisories or a security contact address, never as a public issue.
2. **An acknowledgement SLA**: reply to the reporter quickly, even if the full analysis takes longer.
3. **Triage**: reproduce the issue, assess the impact with CVSS, and decide the affected layers using the SBOM from Part 1.
4. **Coordinated fix**: develop the patch privately, and prepare the release for every affected branch at once.
5. **Disclosure**: request a CVE, publish the advisory, credit the reporter, and announce it on the forum / blog / release notes.
6. **Follow-up**: add a regression test (and a fuzz corpus entry) so that the same bug can never come back.

Researchers are a free, world-class security team -- but only if reporting to you is easy and being credited is guaranteed.

## Exercise

1. Pick one free5GC repository (a `pkg` library, an NF, or `free5gc-compose`) and generate an SBOM for it. How many transitive dependencies are there? How many of them do you recognize?
2. Scan that SBOM with an OSS vulnerability scanner. List the findings, and for each one answer: *is the vulnerable code path actually reachable from free5GC?*
3. Enable a security linter (for example `gosec`) or CodeQL on your fork, and open a PR that fixes one real finding.
4. Write a Go fuzz target for one decoding function (NAS, NGAP, or PFCP). Run it for at least 10 minutes and report what you found -- including "nothing", which is also a valid result.
5. Draft a `SECURITY.md` for your own repository. Who receives the report? How long until the reporter gets an answer?

### Question

- Why is it not enough to manage a single SBOM for the whole free5GC project? Give a concrete example where the severity differs between two layers.
- Dependabot opens a PR for a dependency bump and the CI is green. Is it safe to merge? What else would you check?
- SAST, DAST, and fuzzing all find bugs. Give one class of vulnerability that **only** one of them can realistically find.
- What are the risks of letting an LLM agent write security patches, and how would you mitigate them?

## Reference

- [Fuzz Testing in Go: Discovering Vulnerabilities and Analyzing a Real Case (CVE-2022-43677)](https://free5gc.org/blog/20230809/main/)
- [Web security: CSRF vulnerability in webconsole](https://free5gc.org/blog/20230823/20230823/)
- [Authentication Mechanism in NRF: What Is OAuth?](https://free5gc.org/blog/20230802/20230802/)
- [Go Fuzzing](https://go.dev/doc/security/fuzz/)
- [GitHub Dependabot](https://docs.github.com/en/code-security/dependabot)
- [GitHub Security Advisories](https://docs.github.com/en/code-security/security-advisories)
- [CycloneDX](https://cyclonedx.org/) / [SPDX](https://spdx.dev/)
- [OSV: Open Source Vulnerability database](https://osv.dev/)
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
