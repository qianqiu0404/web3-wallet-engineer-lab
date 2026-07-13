# Wallet Reliability Lab

A public Web3 wallet-backend lab for testing state machines, fund invariants,
and failure recovery without real keys or funds. The repository keeps the
original Go domain API and adds a deterministic Vue experiment console backed
by the same versioned scenario catalog used by Go validation.

## Quick start

```bash
go test ./...
go run ./cmd/api

cd web
npm ci
npm test
npm run dev
```

The API listens on `:8080`:

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/metrics
```

Example requests are in [examples/requests.http](examples/requests.http).

## Reliability experiments

The normal withdrawal flow is scenario `00`. Six failure experiments cover:

1. duplicate withdrawal idempotency and conflicting payloads;
2. broadcast timeout with an unknown chain result;
3. canonical chain success followed by a local DB or outbox failure;
4. deposit reorg and reversal entries;
5. nonce gaps and replacement transactions;
6. fail-closed behavior when risk or signer dependencies are unavailable.

Each experiment includes the injected fault, fund invariant, first stop-loss
action, recovery basis, current boundary, deterministic timeline, and a named
Go test. Scenario copy lives in `scenarios/catalog.json`; both the web console
and Go tests validate that catalog.

## What is verified

- Deposit address ownership is checked before crediting.
- Replaying the same `(chain, tx_hash)` and payload returns the original deposit.
- Reusing the same transaction key with different amount or metadata is rejected as an idempotency conflict.
- Withdrawal risk rejection, review, nonce allocation, simulated broadcast, and confirmation are covered by tests.
- Audit logs, health checks, and text metrics are exposed by the HTTP service.
- The six failure experiments have executable Go invariants and web catalog tests.

```text
Client / Admin API
  -> HTTP handler
  -> Wallet service + risk rules
  -> In-memory store
  -> Simulated chain transaction
```

## Important limits

This is an engineering learning lab, not a production wallet:

- no real blockchain RPC, private key, signer, or asset;
- no PostgreSQL transaction, distributed lock, message queue, or multi-instance coordination;
- no production ledger, reorg compensation, or custody operations;
- chain confirmations and broadcasts are simulated.

The failure playbook and roadmap distinguish runnable evidence from production
design. See [docs/testing-and-roadmap.md](docs/testing-and-roadmap.md).

## Verification

```bash
go test ./...
npm --prefix web ci
npm --prefix web test
npm --prefix web run build
```

The static web console is deployed through GitHub Pages. The optional local Go
API is not deployed by Pages and never receives a key or RPC credential.

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

## License

MIT
