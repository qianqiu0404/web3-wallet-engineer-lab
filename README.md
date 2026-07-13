# Web3 Wallet Engineer Lab

A small, runnable Go wallet-backend lab for explaining and testing custody-wallet domain boundaries without real keys or funds.

It models users, deposit addresses, deposits, withdrawals, review, risk rules, nonce allocation, hot/cold wallets, collection planning, simulated chain transactions, audit logs, and Prometheus-style metrics. The implementation intentionally uses an in-memory store and simulated chain adapter so the state transitions remain easy to inspect.

## Quick start

```bash
go test ./...
go run ./cmd/api
```

The API listens on `:8080`:

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/metrics
```

Example requests are in [examples/requests.http](examples/requests.http).

## What is verified

- Deposit address ownership is checked before crediting.
- Replaying the same `(chain, tx_hash)` and payload returns the original deposit.
- Reusing the same transaction key with different amount or metadata is rejected as an idempotency conflict.
- Withdrawal risk rejection, review, nonce allocation, simulated broadcast, and confirmation are covered by tests.
- Audit logs, health checks, and text metrics are exposed by the HTTP service.

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

The failure playbook and roadmap distinguish runnable evidence from production design. See [docs/testing-and-roadmap.md](docs/testing-and-roadmap.md) and [docs/interview-assets.md](docs/interview-assets.md).

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
