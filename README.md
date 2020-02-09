# raft-sqlite : Sqlite Raft backend

![Master](https://github.com/shettyh/raft-sqlite/workflows/Master/badge.svg)

This repository provides the raftsqlite package. The package exports the SQLStore which is an implementation of both a LogStore and StableStore.

It is meant to be used as a backend for the `raft` [package
here](https://github.com/hashicorp/raft).

This implementation uses Sqlite. This implementation is inspired by [raft-boltdb](https://github.com/hashicorp/raft-boltdb) and [Github Orchestrator] (https://github.com/github/orchestrator)
