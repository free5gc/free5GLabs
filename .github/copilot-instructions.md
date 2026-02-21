# free5GLabs - 5G Core Educational Repository

free5GLabs is a comprehensive educational repository containing hands-on labs for learning 5G Core (free5GC) development and deployment. It includes 8 labs covering network programming, concurrent programming, kernel networking, Docker deployment, service architecture, protocol analysis, Git workflows, and CI/CD practices.

**Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.**

## Working Effectively

### Environment Setup
- **REQUIRED**: Go 1.21+ (validated: Go 1.24.7 works perfectly)
- **REQUIRED**: Docker Engine + Docker Compose v2 (validated: Docker 28.0.4 + Compose v2.38.2)
- **REQUIRED**: Git with submodule support
- **REQUIRED for Lab3**: `bridge-utils` package for network management
- Install all requirements:
  ```bash
  # Install bridge-utils for lab3 network setup
  sudo apt update && sudo apt install -y bridge-utils
  
  # Verify installations
  go version  # Should show 1.21+
  docker --version && docker compose version
  brctl show  # Should show Docker bridge
  ```

### Repository Bootstrap
- **ALWAYS** initialize Git submodules first:
  ```bash
  git submodule update --init --recursive  # Takes ~10-15 seconds
  ```

### Build and Test Commands
#### Lab 0 (Network Programming with Go)
- **Build and test**: `cd lab0 && make test`
- **TIMING**: ~30 seconds including dependency downloads, ~5 seconds for subsequent runs
- **NEVER CANCEL**: Set timeout to 60+ seconds for first run
- **Validation**: Tests should show concurrent TCP connections being handled properly

#### Lab 4 (Service-Based Architecture)
- **Build**: `cd lab4/excersise && go mod tidy && go build`
- **TIMING**: ~17 seconds including dependency downloads
- **Dependencies**: Uses Gin web framework
- **NEVER CANCEL**: Set timeout to 45+ seconds for first build

#### Lab 6 (Git & GitHub + NF Example)
- **Build NF example**: `cd lab6/nf-example && make nf`
- **TIMING**: ~10 seconds including free5GC dependencies
- **Test binary**: `./bin/nf --help` (should show "SPYxFamily" themed help)
- **NEVER CANCEL**: Set timeout to 30+ seconds for first build

### Docker Operations (Lab 3)
- **Pull core images**: 
  ```bash
  cd lab3/exercise
  docker pull mongo:latest  # ~9 seconds
  docker pull free5gc/nrf:v4.2.0  # ~1-2 seconds
  ```
- **TIMING**: Individual pulls 1-15 seconds, full deployment setup 2-5 minutes
- **NEVER CANCEL**: Set timeout to 15+ minutes for complete lab3 deployment
- **Network setup**: Requires `bridge-utils` for `brctl` commands

## Validation Scenarios

### Lab 0 - TCP Server Validation
**CRITICAL**: After implementing TCP functions, ALWAYS validate with manual test:
- Run: `make test`
- Expected: Multiple concurrent TCP connections handled successfully
- Logs should show: "TCP is listening", "new client accepted", "Handle Request"
- **Must show concurrent connection handling (10+ connections simultaneously)**

### Lab 3 - Docker Deployment Validation
**CRITICAL**: After Docker setup changes, ALWAYS validate full deployment:
- Test image pulls: `docker pull free5gc/nrf:v4.2.0`
- Verify network bridges: `brctl show` (should show docker0 and potentially br-free5gc)
- **Manual validation**: Deploy core services and verify connectivity between NFs
- **Check logs**: Each NF should start successfully and register with NRF

### Lab 4 - HTTP API Validation
**CRITICAL**: After HTTP/API changes, ALWAYS validate server functionality:
- Build and run the API server
- Test basic HTTP endpoints with curl or similar
- **Validate CORS configuration** (uses gin-contrib/cors)

### Lab 6 - NF Development Validation
**CRITICAL**: After NF code changes, ALWAYS validate build and execution:
- Build: `make nf`
- Test: `./bin/nf --help` and `./bin/nf -c config/example.yaml`
- **Verify free5GC integration patterns** (consumer/processor/context usage)

## Common Commands and Expected Timing

### Go Operations
| Command | Location | Time | Notes |
|---------|----------|------|-------|
| `make test` | lab0/ | 30s first, 5s subsequent | NEVER CANCEL: Use 60s timeout |
| `go build` | lab4/excersise/ | 17s first | NEVER CANCEL: Use 45s timeout |
| `make nf` | lab6/nf-example/ | 10s first | NEVER CANCEL: Use 30s timeout |
| `gofmt -l .` | Any Go dir | <1s | All code is pre-formatted |

### Docker Operations
| Command | Location | Time | Notes |
|---------|----------|------|-------|
| `docker pull mongo` | Any | 9s | NEVER CANCEL: Use 120s timeout |
| `docker pull free5gc/*` | Any | 1-2s | Images are optimized |
| Full deployment | lab3/exercise/ | 2-5 min | NEVER CANCEL: Use 15+ min timeout |

### Git Operations
| Command | Location | Time | Notes |
|---------|----------|------|-------|
| `git submodule update --init --recursive` | Root | 15s | Required for lab6 |

## Validation Requirements

### Pre-Change Validation
- **ALWAYS** run these commands before making changes:
  ```bash
  # Test current Go code builds
  cd lab0 && make test  # Should pass or show expected unimplemented errors
  cd ../lab4/excersise && go mod tidy && go build
  cd ../lab6/nf-example && make nf
  
  # Verify Docker works
  docker --version && docker compose version
  brctl show
  ```

### Post-Change Validation  
- **ALWAYS** run these commands after making changes:
  ```bash
  # Test builds still work
  gofmt -l lab0/ lab4/excersise/ lab6/nf-example/  # Should show no output
  
  # Run any affected lab tests
  cd lab0 && make test  # If you changed lab0
  cd lab4/excersise && go build  # If you changed lab4
  cd lab6/nf-example && make nf  # If you changed lab6
  ```

### Critical Validation Scenarios
1. **TCP Concurrent Connections** (Lab 0): Must handle 10+ simultaneous connections
2. **Docker Network Setup** (Lab 3): All bridges and interfaces properly configured
3. **HTTP API Functionality** (Lab 4): All endpoints responding correctly
4. **NF Integration** (Lab 6): Follows free5GC patterns and builds successfully

## Repository Structure Reference

```
free5GLabs/
├── lab0/                    # Network Programming (Go + tests)
│   ├── Makefile            # make test command
│   ├── go.mod              # Go 1.21.9 module
│   ├── tcp.go              # Exercise implementation
│   └── tcp_test.go         # Test validation
├── lab1/                    # Concurrent Programming (documentation)
├── lab2/                    # Linux Kernel Networking (documentation)  
├── lab3/                    # free5GC Docker Deployment
│   └── exercise/
│       ├── deploy_exercise.yaml  # Full 5G core stack
│       └── config/         # NF configurations
├── lab4/                    # HTTP Service Architecture
│   └── excersise/
│       ├── go.mod          # Gin web framework
│       └── todo.go         # API implementation
├── lab5/                    # Protocol Analysis (documentation)
├── lab6/                    # Git & GitHub + NF Development
│   └── nf-example/         # Submodule: free5GC NF example
│       ├── Makefile        # make nf command  
│       └── cmd/main.go     # NF implementation
└── lab7/                    # CI/CD (documentation)
```

## Troubleshooting

### Go Build Issues
- **Problem**: `go mod` outdated
- **Solution**: Run `go mod tidy` before building
- **Problem**: Go version too old
- **Solution**: Install Go 1.21+ from go.dev

### Docker Issues
- **Problem**: Permission denied
- **Solution**: Add user to docker group or use sudo
- **Problem**: Network conflicts in lab3
- **Solution**: `docker network prune` and `brctl show` to verify bridges

### Git Submodule Issues
- **Problem**: lab6/nf-example empty
- **Solution**: `git submodule update --init --recursive`

## Key Learning Objectives by Lab

- **Lab 0**: TCP server implementation with goroutines and channels
- **Lab 1**: Concurrency patterns, race conditions, synchronization
- **Lab 2**: Linux kernel networking, interfaces, routing
- **Lab 3**: Docker containerization, free5GC deployment, network setup
- **Lab 4**: RESTful APIs, HTTP protocols, service architecture  
- **Lab 5**: Packet capture, protocol analysis, tcpdump
- **Lab 6**: Git workflows, branching, free5GC NF development patterns
- **Lab 7**: GitHub Actions, CI/CD pipelines, automated testing

Always prioritize hands-on validation over theoretical understanding - these labs are designed for practical 5G development skills.