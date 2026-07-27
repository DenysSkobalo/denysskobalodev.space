# denysskobalodev.space — Core Platform Infrastructure

An edge-native, monorepo platform designed for high performance, high availability, and unified session management across subdomains. Powered by Cloudflare Workers (V8 Isolates), Cloudflare D1, KV, and R2.

---

## Technology Stack

### Runtime & Edge Computing
- **Cloudflare Workers (V8 Isolates)** — Serverless execution layer with <5ms cold starts and zero container overhead.
- **Hono.js** — Ultra-lightweight (~14KB) web framework optimized for V8 Isolates with built-in `@hono/zod-openapi` validation.
- **Astro** — High-performance frontend framework utilizing server-side rendering (SSR) and static generation at the edge.

### Persistence & Storage
- **Cloudflare D1 (SQLite)** — Distributed relational database optimized for read-heavy workloads with primary write execution routing.
- **Cloudflare KV** — Global low-latency key-value storage for sliding-window rate limiting, session caching, and telemetry buffers.
- **Cloudflare R2** — S3-compatible object storage for static media assets, zero egress fees.

### Database ORM & Migrations
- **Drizzle ORM (`drizzle-orm/d1`)** — Edge-compatible, zero-overhead TypeScript ORM.
- **Drizzle Kit (`drizzle-kit`)** — Declarative SQL migration generation and management.

### Monorepo Architecture & Tooling
- **pnpm Workspaces** — Fast, disk-space efficient package manager with deterministic dependency trees.
- **Turborepo** — High-speed build system with intelligent caching for monorepo task pipelines.
- **TypeScript (Strict Mode)** — End-to-end type safety shared across apps and packages.

### CI/CD & Deployment
- **GitHub Actions** — Automated linting, typechecking, database migrations, and Wrangler deployment pipelines.
- **Wrangler CLI (`cloudflare/wrangler-action`)** — Cloudflare developer platform CLI tool for Workers and Pages deployment.

---

## Repository Structure

```text
denysskobalodev.space/
├── .github/
│   └── workflows/          # GitHub Actions deployment pipelines
├── apps/
│   ├── web/                # Astro-based Root Portfolio (denysskobalodev.space)
│   └── api-gateway/        # Hono.js Edge Gateway Worker (api.denysskobalodev.space)
├── packages/
│   ├── db/                 # Cloudflare D1 Schemas, Drizzle Models, & Migrations
│   ├── config/             # Shared TypeScript, ESLint, and Prettier configurations
│   └── shared/             # Shared utilities, crypto helpers, and global types
├── pnpm-workspace.yaml     # Workspace configuration
├── turbo.json              # Turborepo execution pipeline config
└── package.json            # Root configuration
```

---

## Architecture & Standards

### Cross-Subdomain Session Strategy
All subdomains (`*.denysskobalodev.space`) share a unified session model using `__Secure-session` HTTP-Only, SameSite=Lax cookies bound to `.denysskobalodev.space`.

### Read-after-Write Consistency
To handle D1 eventual consistency:
1. Mutations perform writes directly on the Primary D1 node and return the fresh resource payload in the response body.
2. The UI performs optimistic state hydrations without immediate refetching.
3. Subsequent reads pass an `X-D1-Sequence` tracking header to guarantee read consistency.

---

## Getting Started Locally

```bash
# Install dependencies
pnpm install

# Run database migrations locally
pnpm --filter @denysskobalo/db migrate:local

# Start development servers
pnpm turbo run dev
```
