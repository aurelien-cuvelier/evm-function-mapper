# EVM Function Signature Finder

A CLI tool for analyzing Ethereum Virtual Machine (EVM) bytecode. This tool extracts function signatures from the bytecode and attempts to decode them using the [4bytes.directory](https://www.4byte.directory/) API.

> Note: This project is experimental and primarily developed for personal use, and is still incomplete and under active development. While it is open-source, no guarantees or warranties are provided. Contributions are welcome, but use it at your own risk.

---

## Features

- Extract function selectors from EVM bytecode.
- Decode function selectors into human-readable signatures using the 4bytes API.

---

## Limitations

- **Binary Search Dispatchers:** The tool currently does **not support** function dispatchers implemented with binary search. Only linear dispatchers are supported.

---

## Usage

To run the CLI tool, use the following commands:

```bash
go run cmd/main.go [flags]

```

### Available Flags:

- `--rpc`: Specify the RPC endpoint for fetching on-chain bytecode.
- `--address`: Specify the smart contract address (used with `-rpc`).
- `--bytecode`: Provide the compiled contract bytecode directly

## Example Commands

### Analyze On-Chain Bytecode

```bash
go run cmd/main.go --rpc "https://mainnet.infura.io/v3/YOUR-PROJECT-ID" --address "0xContractAddress"
```

### Analyze Local Bytecode

```bash
go run cmd/main.go --bytecode "0x6003600501..."
```
