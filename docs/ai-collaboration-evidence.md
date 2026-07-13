# AI collaboration evidence

This repository was prepared as a public, reviewable engineering artifact through a human-controlled AI collaboration loop:

1. The human selected one bounded task: make duplicate deposit ingestion idempotent.
2. The source lab was copied out of a private Obsidian workspace and sanitized; local paths and private-note metadata were removed.
3. The implementation added a `(chain, tx_hash)` index, stable replay behavior, conflict rejection, and unit tests.
4. Automated checks run `gofmt`, `go test ./...`, path/secret scans, and GitHub Actions.
5. The human remains responsible for repository visibility, technical claims, review, and acceptance.

AI assistance does not turn this lab into production experience. The README and tests explicitly describe the in-memory and simulated-chain boundaries.
