# Lab 7: CI/CD with GitHub Actions

## Introduction

In Lab 7, you will learn how to use GitHub Actions to automate the CI/CD process.

## Goals of this lab

- Understand What is CI/CD
- Understand how to use GitHub Actions

## What is CI/CD?

CI/CD is a part of DevOps practices that automates the integration and delivery of code changes. CI/CD pipeline automates the process of integrating code changes, building the code, testing the code, and deploying the code to production.

### Continuous Integration (CI)

Continuous Integration (CI) is a software development practice where developers integrate code changes into a shared repository frequently. Each integration is verified by an automated build and automated tests. This practice helps to detect and fix integration errors quickly.

A stable CI pipeline ensures that the code is always in a deployable state.
For the CI pipeline to be successful, it should be:
- reproducible: The testing environment should be consistent across all stages of the pipeline, it means that the same test should pass on the developer's machine and the CI server.
- reliable: The CI pipeline should be reliable, it should not fail randomly.
- readable: The CI pipeline should be easy to read and understand.

![alt text](image.png)

In free5GC, we leverage GitHub Actions to perform the CI process. The CI process includes:
- Linting: Check the code style and formatting.
- Unit Testing: Run unit tests to ensure the code works as expected.
- Build: Build the code to ensure that the code can be compiled successfully.
- Integration Testing: Run integration tests to ensure that the code works well with other components.

Any Pull Request (PR) to the free5GC repository will trigger the CI process. The PR will be merged only if the CI process is successful (and should be approved by the project owner).

Besides, In our main repository (free5gc/free5gc), We implement the integration test to guarantee the functionality of the code.
When the new code is merged into the main branch, The integration test will be triggered automatically and do the following tests:
- ./test_ci.sh TestNasReroute
- ./test_ci.sh TestRegistration
- ./test_ci.sh TestGUTIRegistration
- ./test_ci.sh TestServiceRequest
- ./test_ci.sh TestXnHandover
- ./test_ci.sh TestDeregistration
- ./test_ci.sh TestPDUSessionReleaseRequest
- ./test_ci.sh TestPaging
- ./test_ci.sh TestN2Handover
- ./test_ci.sh TestReSynchronization
- ./test_ci.sh TestDuplicateRegistration
- ./test_ci.sh TestEAPAKAPrimeAuthentication
- ./test_ci.sh TestMultiAmfRegistration

> The test procedure is described in the workflow file [`.github/workflows/test.yml`](https://github.com/free5gc/free5gc/blob/main/.github/workflows/test.yml).

### Run your CI pipeline locally

We recommend to use [act](https://github.com/nektos/act) to test your CI pipeline before pushing your code to the repository.

It can avoid the unnecessary CI pipeline execution and save the CI resources. Also, it can help you to debug the CI pipeline in the early stage.

### Continuous Deployment (CD)

Continuous Deployment (CD) is a software development practice where code changes are automatically deployed to production. The CD pipeline automates the process of deploying code changes to production.

The delivery process of the 5G Core Network is more complex than the traditional web application.
It is because of we need to setup the needed infratructure (machine, network, etc) to run the 5G Core Network, You can imagine that:
> If you want to deploy the 5G Core Network to the production Kubernetes cluster, You need to setup the Kubernetes cluster, the network, the storage, etc.
> It is hard to maintain the configuration between the development environment and the production environment.