Project Blueprint: ViraMail (Next-Gen AI-Native Email Platform)
1. Executive Summary
ViraMail is a high-performance, bilingual (FA/EN), and AI-centric email ecosystem built for the 2026 standards. It leverages a microservices architecture to provide a seamless bridge between legacy protocols (SMTP/IMAP) and the future of communication (JMAP/MCP/Agentic AI).

2. Technical Philosophy & Constraints
Language: Go (Golang) 1.2x+. Focus on Concurrency, Memory Safety, and Low Latency.

Architecture Pattern: Microservices with Clean Architecture (Domain-Driven Design).

Database-per-Service: Strict isolation. No cross-service database access.

Inter-Service Communication: * Synchronous: gRPC (Protobuf) for high-speed internal calls.

Asynchronous: NATS or Kafka for event-driven workflows (e.g., Mail Ingested -> AI Summary Trigger).

Multilingual Support: Native i18n for Persian (RTL) and English (LTR). Every service must carry LangContext.

3. Microservices Definition (Independent Binaries)
Every service resides in /cmd and must be independently deployable via Docker.

smtp-ingress-service: Handles ESMTP/LMTP. Responsible for TLS termination, SPF/DKIM/DMARC/ARC validation, and initial mail queuing.

jmap-api-service: The modern gateway. Implements RFC 8620. Handles stateless JSON over HTTP and WebSockets for real-time sync.

legacy-bridge-service: Implements IMAP4rev2 and POP3 as a translation layer over the internal storage.

ai-processor-service: The "Intelligence Brain". Implements MCP (Model Context Protocol). Handles LLM-based spam detection, bilingual summarization, and AI Agent hooks.

storage-service: Aggregator for Metadata (PostgreSQL) and Blob data (S3-compatible). Handles Zstd compression and deduplication.

admin-service: RESTful API for domain management, user provisioning, and i18n configuration.

4. Protocol & Feature Matrix
A. Transport & Security (The "Shield")
Protocols: SMTP, ESMTP, LMTP, MTA-STS, DANE.

Trust: Full implementation of SPF, DKIM, DMARC, ARC, and BIMI.

Encryption: TLS 1.3 mandated for all inter-server and client-server traffic.

B. Access & Sync (The "Core")
Primary: JMAP (JSON-based, mobile-optimized).

Legacy: IMAP4rev2, POP3.

Real-time: WebSockets for instant push across FA/EN interfaces.

C. Intelligence & Integrations (The "Future")
MCP Integration: Native server-side support for AI Agents (Claude, GPT, Gemini) to access mail context securely.

AI Anti-Spam: Semantic analysis using LLMs instead of static Regex.

Full-Text Search: Bilingual (Persian/English) indexing via Tantivy/Meilisearch.

i18n: Global RTL/LTR support with Persian Vazirmatn font integration in UI.

5. Directory Structure (Standard Layout)
Plaintext
/cmd                    # Microservices entry points (main.go for each)
/internal               # Private business logic
  /smtp, /jmap, /ai     # Service-specific logic
  /domain               # Shared entities (Internal)
/pkg                    # Shared libraries (gRPC-proto, Utils, i18n-helpers)
/web                    # Next.js 15 Bilingual Dashboard (RTL/LTR)
/api/proto              # Protobuf definitions for gRPC
/deploy                 # Docker-compose & K8s manifests
6. Implementation Strategy for AI
Phase 1: Setup the Shared Repository (pkg/), gRPC definitions, and the docker-compose environment.

Phase 2: Develop the storage-service and smtp-ingress-service to enable mail receiving.

Phase 3: Build the jmap-api-service for data retrieval.

Phase 4: Integrate ai-processor-service with MCP and i18n layers.

Phase 5: Complete the Frontend with Next.js and Tailwind RTL.