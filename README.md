# Web3 Wallet Domain Engine

> The executable domain and evidence layer behind the public wallet reliability portfolio.

[Interactive Reliability Lab](https://wallet-reliability-lab.vercel.app) · [Engineering Portfolio](https://xiuqiu-site.vercel.app) · [Catalog Inspector](https://qianqiu0404.github.io/web3-wallet-engineer-lab/)

This repository models wallet-backend state machines, idempotency, fund invariants, risk decisions, nonce allocation, and deterministic failure recovery without real keys or funds. It owns executable Go evidence and the versioned Scenario Catalog; the polished interaction experience lives in the separate Wallet Reliability Lab repository.

## Repository role

```text
xiuqiu-site                     portfolio and evidence index
        │
        ├── wallet-reliability-lab      interactive explanation layer
        │             │
        │             └── pins Scenario Catalog v1
        │
        └── web3-wallet-engineer-lab    this repository
                      ├── Go domain API
                      ├── fund invariants
                      ├── deterministic fault models
                      └── versioned scenario contract
```

The GitHub Pages site is a technical Catalog Inspector, not the primary product demo.

## Quick start

```bash
go test ./...
go run ./cmd/api

cd web
npm ci
npm test
npm run dev
```

The API listens on `:8080`. Example requests are in [examples/requests.http](examples/requests.http).

## Scenario Catalog v1

`scenarios/catalog.json` is the canonical machine-readable recovery catalog. `scenarios/catalog.schema.json` defines the public v1 contract. The Go model and technical web inspector both consume the same catalog.

The catalog contains one normal withdrawal baseline and six failure models:

1. duplicate withdrawal idempotency and payload conflicts;
2. broadcast timeout with an unknown chain result;
3. canonical chain success followed by local persistence failure;
4. deposit reorg and reversal entries;
5. nonce gaps and replacement transactions;
6. fail-closed behavior when risk or signer dependencies are unavailable.

Every scenario defines its injected fault, fund invariant, first stop-loss action, recovery basis, current boundary, deterministic timeline, and named Go test.

## What is verified

- deposit address ownership before crediting;
- idempotent deposit replay and conflict rejection;
- withdrawal review, risk rejection, nonce allocation, simulated broadcast, and confirmation;
- one canonical fund result across duplicate, unknown, compensation, reorg, replacement, and dependency-outage paths;
- audit logs, health checks, and Prometheus-style text metrics;
- Catalog v1 semantic and schema-contract checks.

## API outline

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/healthz` | Health check |
| GET | `/metrics` | Text metrics |
| POST | `/api/users` | Create user |
| POST | `/api/users/{id}/addresses` | Allocate deposit address |
| POST | `/api/deposits/simulate` | Simulate idempotent deposit credit |
| POST | `/api/withdrawals` | Request withdrawal |
| POST | `/api/admin/withdrawals/{id}/approve` | Approve and simulate broadcast |
| POST | `/api/chain/tx/{id}/confirm` | Simulate canonical confirmation |

## Verification

```bash
go vet ./...
go test -race ./...
npm --prefix web ci
npm --prefix web test
npm --prefix web run build
```

## Boundaries

- no real blockchain RPC, private key, signer, or asset;
- no PostgreSQL transaction, distributed lock, queue, or multi-instance coordination;
- no production ledger, custody operation, or claim of production readiness;
- the technical Pages inspector does not deploy the Go API;
- private wallet-service repositories are not dependencies of this public engine.

## License

MIT © xiuqiu
