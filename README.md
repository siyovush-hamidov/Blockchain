# 🔒 PoW-Blockchain

A simple Proof-of-Work (PoW) blockchain built from scratch in Go.

---

## Overview

This project implements a basic PoW blockchain where nodes mine blocks and clients interact with the network. It’s a practical way to dive into blockchain fundamentals using Go.

---

## Source

This project is inspired by and based on the concepts from [number571's blockchain tutorial](https://github.com/number571/blockchain/blob/master/_example/blockchain.pdf).

---

## Getting Started
### Build the Project
1. Clone the repository:
```bash
git clone https://github.com/siyovush-hamidov/PoW-Blockchain
cd PoW-Blockchain
```
2. Build using the Makefile:
```bash
make build
```
This compiles `node` and `client` executables into the project directory.

---

## Usage

### Start the First Node

```bash
./node -serve::8080 -newuser:node1.key -newchain:chain1.db -loadaddr:addrlist.json
```
This command runs a node on port `8080`, creates a new user key (`node1.key`), and initializes a blockchain database (`chain1.db`).
### Start the Second Node. Open a new terminal
```bash
./node -serve::9090 -newuser:node2.key -newchain:chain2.db -loadaddr:addrlist.json
```
This command runs another node on port `9090` with its own key and database.
### Run the Client. Also requires another terminal
```bash
./client -loaduser:node1.key -loadaddr:addrlist.json
```
This connects to the nodes using `node1.key` and interacts with the blockchain.
### Client Commands
After running the client, use these commands:
- `/user balance` — Check your balance across nodes.
- `/chain tx <receiver> <value>` — Send a transaction (e.g., /chain tx aaa 3).
- `/chain print` — Display the blockchain.
- `/chain balance` — Check addresses' balance.
- `/chain size` — Get the number of blocks.
- `/chain block <number>` — View a specific block.
- `/exit` — Exit the client.