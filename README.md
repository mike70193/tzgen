# tzgen

Go binding to Tezos smart contracts, using code generation.

## Installation

```bash
go install github.com/jeanschmitt/tzgen@latest
```

## Usage

### From a deployed contract

```bash
tzgen --name Hello --pkg contracts --address KT1K3ZqbYq1bCwpSPNX9xBgQd8CaYxRVXd4P -o ./contracts/Hello.go
```

The endpoint is `https://ghostnet.smartpy.io` by default, but can be overridden with `--endpoint`.

### From a micheline file

```bash
tzgen --name Hello --pkg contracts --src ./Hello.json -o ./contracts/Hello.go
```

## Renaming structs

Some structs don't have annotations in the contract's script.
In this case, an auto-generated name is given.

It is possible to give a configuration map to tzgen, to map these auto-generated names to the one you want.

To do so, pass a yaml to tzgen with the `-f` flag.

Example of a fixup file:

```yaml
FA2NFTRecord3:
  name: OperatorForAll
  fields:
    Field0: Addr
    Field1: Owner

FA2NFTRecord5:
  name: BalanceOfRequest

FA2NFTRequest:
  equals: FA2NFTRecord5
```
