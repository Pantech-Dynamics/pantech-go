# Changelog

What changed in each version of pantech-go. Every pull request that changes what the SDK does
adds a line under **Unreleased**. Releasing renames that heading to the version and its date
(`## v0.2.0 — 2026-11-02`); a push to main whose newest heading is an untagged version tags it and
publishes a GitHub release with that section as its notes.

## v0.4.0 — 2026-11-10

### Changed

- The vendored spec is the current public API, with DNS, domains, the container registry, IAM and
  certificates; every new operation has a method.

### Added

- `DNS`, for Pantech DNS: `ListZones`/`AllZones`, `CreateZone`, `GetZone`, `UpdateZone`,
  `DeleteZone` and `VerifyZone` for hosted zones; `ListRecords`/`AllRecords`, `CreateRecord`,
  `GetRecord`, `UpdateRecord`, `DeleteRecord` and `BatchRecords` (up to 1000 changes, all or
  nothing) for record sets; `ImportZone` and `ExportZone` for BIND zone files.
  `DNS.WaitForZone` and `DNS.WaitForRecord` return a `*DNSFailureError` when the zone or record
  fails.
- `Domains`: `ListTLDs`, `Search`, `CheckTransfer`, `List`/`All`, `Get`, `Register`, `Transfer`,
  `Renew`, `Update` (auto-renew, transfer lock, privacy), `SetNameServers`, `GetContact`,
  `UpdateContact`, `RequestAuthCode`, `GetAuthCode`, `ListOrders` and `GetOrder`.
  `Domains.WaitForOrder` returns a `*DomainOrderError` (with `Refunded`) when an order fails, and
  `Domains.WaitForAuthCode` waits for a requested auth code. `TLDOffer.RegisterPrice`,
  `TransferPrice` and `RenewPrice` compute an order's amount.
- `Domains.RetryOrderPayment` (`POST /domain-orders/{order_id}/retry-payment`): an order waiting
  in `awaiting_payment` tries its payment again at once and waits another hour; it never charges
  twice. `409 CodeDomainOrderNotAwaitingPayment` when the order is not waiting for a payment.
  `DomainOrder.PaymentURL` is now set only when the bank asks the customer to confirm a saved-card
  payment, and `DomainPaymentRequiredError` reads well without one.
- `WithIfMatch(version)`, which sends `If-Match` so a change applies only to the version you read
  (`412 VERSION_CONFLICT` otherwise).
- Constants for the problem codes of DNS and domains, such as `CodeDNSRecordExists`,
  `CodeVersionConflict` and `CodeDomainNotAvailable`.
- `DomainContact`, `TransferDomainRequest` and `DomainAuthCode` do not print their contents with
  `fmt`, so the registrant contact and auth codes stay out of logs.
- `examples/dns-records`.
- `Certificates`, for the certificate manager: `List`/`All`, `Request`, `Import`, `Get`,
  `Update`, `Verify`, `Reimport`, `Export`, `Revoke` and `Delete`, and `WaitForCertificate`,
  which returns a `*CertificateFailureError` when the certificate fails, expires or is revoked.
  `ExportedCertificate.FullChainPEM`. The private key in `ExportedCertificate`,
  `ImportCertificateRequest` and `ReimportCertificateRequest` never prints with `fmt`. Constants
  for the certificate problem codes, such as `CodeCertificateNotExportable`.
- `DNSZone.DelegationDeadlineAt`, and `DNS.QuoteZones` (`DNSZoneQuote`): five zones are free,
  each further one is billed by the hour.
- Card payment for domain orders: `DomainOrder.PaymentURL` and `PaymentExpiresAt`, the
  `awaiting_payment` status (`DomainOrderAwaitingPayment`). `Domains.WaitForOrder` returns a
  `*DomainPaymentRequiredError` (matching `pantech.ErrPaymentRequired`) carrying the payment link
  instead of waiting on the user, and `Domains.WaitForPayment` waits through the payment until
  the order succeeds or fails.
- `TLDOffer.PrivacyAlwaysOn` for TLDs whose registry always hides contacts (.ng), and
  `CodePrivacyAlwaysOn`.
- HTTPS on appliance load balancers: `CertificateID`, `TLSPolicy` and `RedirectHTTP` on
  `CreateLoadBalancerRequest` and `UpdateLoadBalancerRequest`; `CertificateID`,
  `CertificateFingerprintSHA256`, `TLSPolicy` and `RedirectHTTP` on `LoadBalancer`; field codes such
  as `CodeLoadBalancerTLSNotSupported` and `CodeCertificateNotUsable`.
- `Registry`, for the private container registry: `GetNamespace` and `SetNamespace`;
  `ListRepositories`/`AllRepositories`, `CreateRepository`, `GetRepository`,
  `UpdateRepository` (description, visibility, immutable tags, retention) and
  `DeleteRepository`; `ListTags`/`AllTags`, `GetTag` and `DeleteTag`; `ListTokens`,
  `CreateToken` and `DeleteToken` for registry tokens. A created token's `Secret` is not printed
  by `fmt`. Constants for the registry's problem codes, such as `CodeRegistryNamespaceNotSet` and
  `CodeTagImmutable`, and `DefaultRegistryHost`.
- `IAM`, for teams and access policies: `ListActions`; `GetSettings` and `UpdateSettings` (the
  member baseline); `ListTeams`/`AllTeams`, `CreateTeam`, `GetTeam`, `UpdateTeam`, `DeleteTeam`,
  `ListTeamMembers`/`AllTeamMembers`, `AddTeamMember` and `RemoveTeamMember`;
  `ListPolicies`/`AllPolicies` (`Managed` filter), `CreatePolicy`, `GetPolicy`, `UpdatePolicy`
  (with `WithIfMatch`) and `DeletePolicy`; `ListAttachments`, `Attach` and `Detach`;
  `ListMembers`/`AllMembers` (`Q` matches email or name), `GetEffectiveAccess` and `Simulate`. Typed `PolicyStatement` and `PolicyConditions`, built with
  `pantech.Allow` and `pantech.Deny` and `WithSid`, `WithSourceIP` and `WithMFA`. Constants for
  the IAM problem codes, such as `CodeAccessDenied` and `CodeInvalidPolicy`.

## v0.3.0 — 2026-11-06

### Added

- `KubernetesVersionList.FixedNodeSizes` (`KubernetesFixedNodeSize`): the zones that build
  Kubernetes nodes only at fixed sizes, one entry per size. There, a cluster, a node plan resize
  and a node pool take only a plan of one of those sizes.

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
