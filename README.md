# Self-stabilizing randomized Byzantine-tolerant Binary Consensus

## About
Binary consensus is a problem in which a set of processors must agree on a single binary value. In asynchronous systems, where a subset of the processors may be malicious, this challenge gets more challenging. We study malicious and more serious problems in this work: transient faults. These are temporary violations of the system's operating assumptions that might cause the system's state to change unexpectedly, making recovery impossible without human intervention. We implement an existing protocol for randomized Byzantine-tolerant binary consensus algorithm that is loosely-self-stabilizing using the Go programming language and the ZeroMQ communication framework.

## Contribution
We present the first, to our best knowledge, implementation and experimental validation and evaluation of a self-stabilizing randomized Byzantine-tolerant algorithm, namely of the algorithm by [Georgiou et al.](https://arxiv.org/pdf/2103.14649.pdf) We use the Go programming language together with the ZeroMQ message-passing library. Also, we perform the experimental validation to make sure of the correctness of our implementation using unit tests. We then proceed to compare this algorithm with the original non-stabilizing binary consensus algorithm by [Mostefaoui et al.](https://dl.acm.org/doi/pdf/10.1145/2611462.2611468) implemented as part of a [degree project](https://github.com/v-petrou/BFTWithoutSignatures). Moreover, we evaluated the performance overhead which is caused to the presence of Byzantine and transient faults.

## Requirments
A Linux build environment (tested on Ubuntu 20.04 LTS and CentOS 8) with the following components:
- [`Go Programming language`](https://go.dev/)

- [`ZeroMQ library (libzmq)`](https://zeromq.org/download/)

- [`ZeroMQ plugin for Go (pebbe/zmq4)`](https://zeromq.org/languages/go/#pebbe-zmq4)

## Installation
Clone the repository:
```bash
git clone git@github.com:constandinos/self-stabilizing-binary-consensus.git
cd self-stabilizing-binary-consensus
```

Give execute permission to bash scripts:
```bash
chmod +x scripts/*
```

## Usage
**Execution:**
```bash
./scripts/run.sh <N> <M> <CLIENTS> <REMOTE> <BYZANTINE_SCENARIO> <SELF_STABILIZING> <CORRUPTION> <DEBUG> <OPTIMIZATION>
```
Arguments explanation:
- `N` Number of processors (tested with N=4 to N=16, do not set N>16)
- `M` Predefined system parameter that ensures that round r bounded by M (tested with M=6)
- `CLIENTS` Number of clients (tested with CLIENTS=1)
- `REMOTE` REMOTE=0: Execution on localhost | REMOTE=1: Execution on cluster
- `BYZANTINE_SCENARIO` BYZANTINE_SCENARIO=0: Normal | BYZANTINE_SCENARIO=1: Idle attack | BYZANTINE_SCENARIO=2: Inverse attack | BYZANTINE_SCENARIO=3: Half&Half attack | BYZANTINE_SCENARIO=4: Random attack
- `SELF_STABILIZING` SELF_STABILIZING=0: Execution of non-self-stabilizing algorithm | SELF_STABILIZING=1: Execution of self-stabilizing algorithm
- `CORRUPTION` CORRUPTION=0: No corruptions | CORRUPTION=1: Corrupt initial state | CORRUPTION=2: Random corruptions
- `DEBUG` DEBUG=0: No logs | DEBUG=1: Write logs for monitoring of the system's operation
- `OPTIMIZATION` OPTIMIZATION=0: No algorithm optimization | OPTIMIZATION=1: Algorithm optimization (sends messages every 2nd iteration)

**Example:**
```bash
./scripts/run.sh 4 6 1 0 0 1 1 1 1
```
Arguments explanation:
- `N=4`
- `M=6`
- `CLIENTS=1`
- `REMOTE=0`
- `BYZANTINE_SCENARIO=0`
- `SELF_STABILIZING=1`
- `CORRUPTION=1`
- `DEBUG=1`
- `OPTIMIZATION=1`

**Termination:**
```bash
./scripts/kill.sh
```

**Notes:**
- For better experimental results please use `clustercg0` machine as localhost
- You can find the output on `logs/out` and error on `logs/error`

## References
- [Loosely-self-stabilizing Byzantine-tolerant Binary Consensus for Signature-free Message-passing Systems (studied algorithm)](https://arxiv.org/pdf/2103.14649.pdf)
- [BFTWithoutSignatures (baseline project)](https://github.com/v-petrou/BFTWithoutSignatures)
