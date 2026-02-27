ViraMail - Next-Gen Email Platform (Phase 1 MVP)

This repository contains an initial skeleton for ViraMail per `INSTRUCTIONS.md`.

Phases:
- Phase 1: Shared pkg, gRPC proto stubs, smtp-ingress-service, storage stub, docker-compose, Ubuntu installer script.
- Phase 2: storage-service full impl, persistence, S3, Postgres, DKIM/SPF verification.
- Phase 3: jmap-api-service and legacy-bridge-service.
- Phase 4: ai-processor-service + MCP integration.
- Phase 5: Next.js bilingual frontend and production infra.

How to run (dev):
1. Install Go 1.20+
2. Run: go run ./cmd/smtp-ingress-service
