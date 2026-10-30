# Changelog

What changed in each version of pantech-go. Every pull request that changes what the SDK does
adds a line under **Unreleased**. Releasing renames that heading to the version and its date
(`## v0.2.0 — 2026-11-02`); a push to main whose newest heading is an untagged version tags it and
publishes a GitHub release with that section as its notes.

## Unreleased

### Added

- `KubernetesNodePool.FailureCode` and `FailureReason`: why a pool's workers last failed to build
  (for example `PROVISIONING_CAPACITY_UNAVAILABLE`), while the pool stays `provisioning` or
  `updating` and is retried. Nil while nothing has failed, and on API versions that do not return
  them yet.

## v0.2.0 — 2026-10-28

### Added

- `Migrate`, for moving from another cloud: `CreateConnection`, `ListConnections`,
  `GetConnection`, `UpdateConnection` (name, role ARN, retry) and `DeleteConnection` for AWS
  connections through a read-only role; `StartScan`, `ListScans` and `GetScan` (with
  `GetCloudScanParams{Include: "inventory"}`); `CreateMigration`, `ListMigrations`,
  `GetMigration`, `RenameMigration`, `DeleteMigration`, `UpdateItem` and `ExportTerraform` for
  migration plans, their price comparison and their Terraform.
- `Migrate.WaitForConnection` and `Migrate.WaitForScan`, which return a `*CloudFailureError` when
  the connection or scan fails.
- `KubernetesCluster.NodeHealth`, `NodeHealthReason`, `ReadyNodes` and `NodesCheckedAt`, and
  `CreateLoadBalancerRequest.Protocol`.
- `examples/migrate-from-aws`.

### Changed

- The vendored spec is the current public API. Agent sessions, which answer only AI apps connected
  through Pantech's MCP server, have no methods yet.

## v0.1.0 — 2026-10-24

### Added

- First release: a typed method for every operation in the public API (instances and their
  orders, VPC networks, subnets, firewall rules, public IPs, port forwarding, security groups,
  load balancers, volumes, snapshots, databases and their orders, Kubernetes clusters, node pools
  and add-ons, SSH keys, webhooks, operations, catalogue, billing, usage, and `Me`, including
  `Me.Revoke` for `DELETE /me`).
- `Client.Do` for raw calls with the same auth, retries and errors, and per-call options
  `WithIdempotencyKey`, `WithRequestHeader` and `WithNoRetry`.
- An `Idempotency-Key` (UUIDv4) on every write, reused across retries; retries with jittered
  backoff on network errors, 429 (honouring `Retry-After`), 500, 502, 503 and 504, except on the
  routes the API does not replay.
- `*pantech.Error` with the problem's status, code, title, detail, request id, field errors and
  `Retry-After`, and `pantech.IsCode`.
- `iter.Seq2` iterators over paginated lists, `Operations.Wait`, `Operations.Stream`,
  `Instances.WaitForOrder` and `Databases.WaitForOrder`.
- `VerifyWebhook`, `ParseWebhook` and `SignWebhook` for `Pantech-Signature`.
