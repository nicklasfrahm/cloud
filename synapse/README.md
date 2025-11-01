# Synapse

Synapse is a developer platform frontend for managing releases and environment overlays stored in Git. It's a minimal Next.js + React + Tailwind Typescript app and includes a small, pluggable client interface so the UI can talk to either GitHub or local files for testing.

Key goals included in this repo:

- Dashboard showing clusters and services.
- Per-service pages with overlay YAMLs and promote/diff workflows.
- A `SynapseClient` interface with interchangeable transports (`GitHubClient` and `FilesystemClient`).

## Deployment config structure

Deploy and release data follows this layout in the monorepo:

```
deploy/clusters/$cluster/$tenant/$release.yaml
deploy/services/$release/
	├── 00-base.yaml
	├── 10-env-dev.yaml
	├── 10-env-stg.yaml
	├── 10-env-prd.yaml
	└── 30-tenant-mytenant.yaml
```

Each `release.yaml` references an OCI chart and a tag, for example:

```yaml
repository: oci://ghcr.io/nicklasfrahm/charts
release: sample
chart: sample-chart
tag: 0.3.0
```

## Frontend requirements

- Single-page app (Next.js + React + TailwindCSS).
- `/` (Dashboard): shows cluster and service counts, search bar for services, click a service to view details.
- `/services/:release`: shows overlays (00-base.yaml, 10-env-\*.yaml, etc.), with "Show diff" and "Promote" actions between environments. "Promote" shows a diff and opens a PR labeled `synapse`.

## SynapseClient interface

The UI uses an interface so transports can be swapped during development or testing.

```
interface SynapseClient {
	listClusters(): Promise<string[]>;
	listServices(): Promise<Service[]>;
	getService(release: string): Promise<Service | null>;
	getDiff(release: string, src: string, dst: string): Promise<string>;
	promote(
		release: string,
		src: string,
		dst: string,
		title: string
	): Promise<{ ok: boolean; prUrl?: string; error?: string }>;
}
```

Transports implemented here:

- `GitHubClient` (uses GitHub API)
- `FilesystemClient` (reads from local `deploy/` data)

## Local API endpoints (for dev/testing)

Mock API endpoints are provided for the frontend to call during development:

- `GET /api/synapse/clusters`
- `GET /api/synapse/services`
- `GET /api/synapse/services/[release]`
- `POST /api/synapse/diff` (body: { release, src, dst })
- `POST /api/synapse/promote` (body: { release, src, dst, title })

## Getting started (dev)

Install and run the Next.js app in `synapse/`:

```bash
cd synapse
npm install
npm run dev
```

Open http://localhost:3000 and explore the dashboard.

## Notes & next steps

- The repo includes examples and a working `FilesystemClient` so you can develop without hitting external APIs.
- Small follow-ups that add value:
  - Add unit tests for the `SynapseClient` implementations.
  - Wire up GitHub PR creation with a smaller, documented service account.

---

If you'd like, I can also add a short CONTRIBUTING section and a tiny smoke test that hits the mock API endpoints.
