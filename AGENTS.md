# ANR Platform - Agent Instructions

## 1. Source of Truth

Before making architectural decisions or writing code, read:

`docs/ANR_PLATFORM_CONTEXT.md`

Treat that document as the source of truth for ANR Platform.

If a generic Agent Skill conflicts with a decision already locked in
`docs/ANR_PLATFORM_CONTEXT.md`, the ANR Platform context takes precedence.

Do not redesign locked architecture unless the user explicitly requests it.

---

## 2. Agent Skills

Workspace Agent Skills are located at:

`.agents/skills/`

Use only the skills relevant to the current task.

Do not blindly apply every installed skill.

### Go Backend

For Go backend architecture and implementation, prioritize:

- `go-clean-architecture`
- `golang-project-layout`
- `golang-code-style`
- `golang-naming`
- `golang-design-patterns`
- `golang-structs-interfaces`
- `golang-dependency-injection`
- `golang-dependency-management`
- `golang-context`
- `golang-error-handling`
- `golang-database`
- `golang-security`
- `golang-testing`
- `golang-lint`
- `golang-observability`
- `golang-performance`
- `golang-troubleshooting`

### Infrastructure

For infrastructure work, use relevant installed skills for:

- Docker
- Docker Compose
- PostgreSQL
- Redis
- Nginx
- security
- observability
- Prometheus / PromQL

### CI/CD and Git

For repository and delivery workflows, prioritize:

- `conventional-git`
- `golang-continuous-integration`

and other relevant installed CI/CD skills.

---

## 3. Locked Repository Structure

Do not reorganize the root repository structure.

The project root remains:

anr-platform/
├── .agents/
│   └── skills/
├── .github/
│   └── workflows/
├── admin/
├── backend/
├── docs/
├── infrastructure/
│   ├── cache/
│   ├── database/
│   ├── docker/
│   ├── gateway/
│   ├── management/
│   ├── monitoring/
│   ├── scripts/
│   ├── security/
│   └── storage/
├── landing/
├── tools/
├── .env.example
├── .gitignore
├── AGENTS.md
├── docker-compose.yml
├── Makefile
└── README.md

Do not introduce alternative root layouts such as `apps/`, `services/`,
or another infrastructure root unless explicitly requested.

---

## 4. Backend Architecture

The backend is written in Go.

Use Clean Architecture.

Core principles:

- business logic must not depend on HTTP, database, Redis, or framework code
- dependencies point inward
- domain logic remains infrastructure-independent
- use cases orchestrate application behavior
- interfaces/adapters translate between external systems and application logic
- infrastructure implements external concerns
- dependency wiring happens at the application boundary

Do not introduce microservices unless explicitly requested.

The initial backend is a modular monolith designed so components can be
separated later if scaling requires it.

---

## 5. Authentication

Authentication is first-party and implemented by the Go backend.

Do not introduce Firebase Auth, Supabase Auth, Auth0, Clerk, or another
external authentication platform unless explicitly requested.

Security-sensitive implementation must use the relevant security skills.

---

## 6. Data Layer

Primary database:

PostgreSQL

Cache / ephemeral data:

Redis

Keep persistence implementation outside the core domain.

Repository interfaces belong on the inward-facing side of the architecture.
Concrete PostgreSQL and Redis implementations belong in infrastructure/adapters.

---

## 7. Application Modules

ANR applications follow the naming convention:

`anr-NNN-slug`

Example:

`anr-001-wallpaper`

Common platform capabilities should be shared rather than duplicated across
application modules.

The Admin Dashboard must be capable of aggregating information across ANR apps.

---

## 8. Infrastructure

Production infrastructure initially runs on one VPS using Docker.

Core infrastructure includes:

- Docker
- Docker Compose
- PostgreSQL
- Redis
- Nginx
- Portainer

Monitoring may include:

- Prometheus
- Grafana
- Uptime Kuma

Keep the system capable of scaling later without prematurely introducing
Kubernetes or unnecessary distributed architecture.

---

## 9. Working Rules

Before implementing a task:

1. Read the relevant section of `docs/ANR_PLATFORM_CONTEXT.md`.
2. Identify the relevant skills under `.agents/skills/`.
3. Apply those skills while respecting locked ANR decisions.
4. Inspect existing code before creating new abstractions.
5. Preserve established naming and package conventions.
6. Prefer simple production-ready solutions over speculative abstractions.

When asked to design something first:

- do not start implementation
- present the proposed structure
- explain responsibilities
- wait for approval before creating files

When asked to implement:

- make the smallest coherent change
- keep boundaries explicit
- avoid unrelated refactoring
- add or update tests where appropriate

---

## 10. Conflict Priority

When instructions conflict, use this project-level priority:

1. Explicit instruction from the user
2. `docs/ANR_PLATFORM_CONTEXT.md`
3. `AGENTS.md`
4. Task-relevant `.agents/skills`
5. Generic conventions or agent defaults

Never allow a generic skill to silently override a locked ANR architectural decision.

---

## 11. Current Development Order

Follow the current ANR Platform roadmap:

1. Go Backend
2. Docker Compose
3. PostgreSQL
4. Redis
5. Nginx Gateway
6. Landing
7. Admin Dashboard
8. CI/CD
9. Monitoring

Do not jump to later milestones unless explicitly requested.

---

## 12. Agent Behavior

Do not redesign architecture unnecessarily.

Do not introduce technologies simply because an installed skill mentions them.

Do not assume optional Go libraries are approved merely because their skills
are installed.

For example, Wire, Fx, Dig, gRPC, GraphQL, Cobra, Viper, or other libraries
require an actual project need or explicit architectural decision before adoption.

Use installed skills as implementation guidance, not as mandatory dependencies.

When uncertain about an architectural decision, inspect
`docs/ANR_PLATFORM_CONTEXT.md` before making assumptions.
