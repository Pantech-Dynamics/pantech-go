# pantech-go

[![CI](https://github.com/Pantech-Dynamics/pantech-go/actions/workflows/ci.yml/badge.svg)](https://github.com/Pantech-Dynamics/pantech-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/Pantech-Dynamics/pantech-go.svg)](https://pkg.go.dev/github.com/Pantech-Dynamics/pantech-go)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

The official Go SDK for the [Pantech Dynamics public API](https://docs.pantechdynamics.com/api):
instances, VPC networks, firewalls, public IPs, load balancers, volumes, snapshots, managed
databases, managed Kubernetes, SSH keys, webhooks, billing, usage, and moving from another cloud.

- Every operation in the public API has a typed method, except agent sessions, which answer only
  AI apps connected through Pantech's MCP server.
- Writes carry an `Idempotency-Key` that is reused when a request is retried, so a retry never
  makes a change twice.
- Errors are typed (`*pantech.Error`) with the API's stable problem code.
- Lists are Go 1.23 iterators; asynchronous changes have waiting helpers.
- No dependencies outside the standard library.

## Install

```sh
go get github.com/Pantech-Dynamics/pantech-go
```

Go 1.23 or newer.

## Quickstart

```go
package main

import (
	"context"
	"fmt"
	"log"

	pantech "github.com/Pantech-Dynamics/pantech-go"
)

func main() {
	c, err := pantech.NewClient() // reads PANTECH_API_KEY
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	for vm, err := range c.Instances.All(ctx, nil) {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(vm.ID, vm.Name, vm.ObservedState)
	}
}
```

More in [`examples/`](examples): creating an instance, listing everything, following an
operation, a webhook receiver, and moving from AWS.

## Authentication

Create an API key in the console under **Organization › API keys**. A `read` key can call every
GET; a `write` key can also create, change and delete.

```go
c, err := pantech.NewClient(
	pantech.WithAPIKey(os.Getenv("PANTECH_API_KEY")), // the default source
	pantech.WithBaseURL("https://api.pantechdynamics.com/public/v1"), // or PANTECH_API_URL
)
```

| Option | What it does |
|---|---|
| `WithAPIKey(key)` | The `PAN_…` key. Default: `$PANTECH_API_KEY`. |
| `WithBaseURL(url)` | Base URL including `/public/v1`. Default: `$PANTECH_API_URL`, then production. |
| `WithHTTPClient(h)` | Your `*http.Client`. Its `Timeout` bounds each attempt (default 90 s). |
| `WithUserAgent(ua)` | Replaces the `pantech-go/<version>` User-Agent. |
| `WithMaxRetries(n)` | Retries after the first try. Default 3; 0 turns them off. |
| `WithHeader(k, v)` | A header on every request. |

The key is sent only in the `Authorization: Bearer` header. The SDK never logs it, and printing
a `Client` does not show it. `c.Me.Get(ctx)` describes the key; `c.Me.Revoke(ctx)` revokes it.

## Services

`c.Instances`, `c.Networks`, `c.Subnets`, `c.FirewallRules`, `c.PublicIPs`,
`c.PortForwardingRules`, `c.SecurityGroups`, `c.LoadBalancers`, `c.Volumes`, `c.Snapshots`,
`c.Databases`, `c.Kubernetes`, `c.SSHKeys`, `c.Webhooks`, `c.Migrate`, `c.Operations`, `c.Catalog` (regions,
plans, images, disk offerings, database engines, Kubernetes versions, price estimates),
`c.Billing` (credit, receipts, usage) and `c.Me`.

Request types follow the API's field names. Optional fields are pointers, and a nil field is not
sent; `pantech.Ptr` sets one inline:

```go
_, err := c.Volumes.Create(ctx, pantech.CreateVolumeRequest{
	Name:             "data",
	DiskOfferingSlug: "custom",
	SizeGB:           pantech.Ptr(50),
})
```

For anything the typed methods do not cover, `Do` sends a raw call with the same auth, retries
and error handling:

```go
var out map[string]any
err := c.Do(ctx, "GET", "/instances", url.Values{"q": {"web"}}, nil, &out)
```

## Errors

Every API error is a `*pantech.Error` with `Status`, `Code`, `Title`, `Detail`, `RequestID`,
`Errors` (the invalid fields of a `422`) and `RetryAfter`. Branch on `Code`, never on `Detail`:

```go
vm, err := c.Instances.Get(ctx, id)
switch {
case pantech.IsCode(err, "RESOURCE_NOT_FOUND"):
	// gone
case err != nil:
	if e, ok := pantech.AsError(err); ok {
		log.Printf("%s (quote request %s to support)", e.Detail, e.RequestID)
	}
	return err
}
```

`Detail` is plain language and safe to show to your users.

## Pagination

Lists return `{data, next_cursor}`. Paginated lists have two methods: `List` returns one
`*Page[T]`, and `All` iterates over every item, fetching pages as needed:

```go
page, err := c.Volumes.List(ctx, &pantech.ListParams{Limit: 100, Q: "db"})
// page.Data, page.NextCursor, page.HasMore()

for vol, err := range c.Volumes.All(ctx, &pantech.ListParams{ObservedState: "active"}) {
	if err != nil {
		return err
	}
	fmt.Println(vol.Name)
}
```

Lists the API returns in one page (regions, plans, images, SSH keys, webhooks, clusters…) return
a plain slice, and still follow `next_cursor` should the API ever add pages.

## Operations

Writes return `202 Accepted` with an operation. Wait for it:

```go
accepted, err := c.Instances.Stop(ctx, id)
if err != nil {
	return err
}
op, err := c.Operations.Wait(ctx, accepted.OperationID, &pantech.WaitOptions{Timeout: 10 * time.Minute})
var failed *pantech.OperationError
if errors.As(err, &failed) {
	fmt.Println(failed.Operation.Failure.Code, failed.Operation.Failure.Reason)
}
```

or follow it over Server-Sent Events with `c.Operations.Stream(ctx, id)`.

A new instance or database is an order, paid first: `c.Instances.WaitForOrder` and
`c.Databases.WaitForOrder` wait until it is `provisioned` and return an `*pantech.OrderError` when
it ends `payment_failed` or `failed`.

## Move from another cloud

`c.Migrate` connects an AWS account through a read-only role that only Pantech can assume (no
access keys are shared), scans what runs there, and plans the move with a price comparison and a
Terraform export. Nothing is created on Pantech until you apply the Terraform.

```go
conn, err := c.Migrate.CreateConnection(ctx, pantech.CreateCloudConnectionRequest{Name: "Production AWS"})
// Create the stack from *conn.Setup.QuickCreateURL (or conn.Setup.Template) in AWS, then:
_, err = c.Migrate.UpdateConnection(ctx, conn.ID, pantech.UpdateCloudConnectionRequest{RoleARN: &roleARN})
conn, err = c.Migrate.WaitForConnection(ctx, conn.ID, nil)

scan, err := c.Migrate.StartScan(ctx, conn.ID)
scan, err = c.Migrate.WaitForScan(ctx, scan.ID, nil)

plan, err := c.Migrate.CreateMigration(ctx, pantech.CreateMigrationRequest{ScanID: scan.ID, Name: "Move"})
fmt.Println(plan.Estimate.PantechMonthlyMinor, plan.Estimate.SourceMonthlyMinor)
tf, err := c.Migrate.ExportTerraform(ctx, plan.ID) // tf.Files: main.tf, variables.tf
```

A failed connection or scan returns a `*pantech.CloudFailureError` whose `Failure.Message` says
what to do. `c.Migrate.UpdateItem` excludes an item or changes its plan, size, range or role before
you export. The whole flow is in [`examples/migrate-from-aws`](examples/migrate-from-aws).

## Webhooks

Verify every delivery against the raw body before acting on it:

```go
body, _ := io.ReadAll(r.Body)
event, err := pantech.ParseWebhook(secret, r.Header.Get("Pantech-Signature"), body, 0)
if err != nil {
	http.Error(w, "invalid signature", http.StatusBadRequest)
	return
}
```

`pantech.VerifyWebhook(secret, header, body, tolerance)` only checks the signature. The header is
`t=<unix>,v1=<hex HMAC-SHA256(secret, t + "." + body)>`; the comparison is constant time, the
default tolerance is 5 minutes, and during the 24 hours after a secret rotation either `v1` value
verifies. `pantech.SignWebhook` signs a body the same way, for testing your receiver.

## Idempotency

Every POST, PUT, PATCH and DELETE gets a random `Idempotency-Key` (a UUIDv4), and the same key is
sent on every retry of that call, so the API returns the original result instead of acting twice.
To make a change safe to repeat across runs or processes, choose the key yourself:

```go
order, err := c.Instances.Create(ctx, req, pantech.WithIdempotencyKey("deploy-42-web-1"))
```

Every typed method, and `Do`, takes these per-call options: `WithIdempotencyKey(k)`,
`WithRequestHeader(k, v)` and `WithNoRetry()`.

## Retries

Network errors, `429`, `500`, `502`, `503` and `504` are retried up to 3 times with jittered
exponential backoff (0.5 s, 1 s, 2 s… up to 8 s), and a `Retry-After` is honoured (up to 60 s).
Routes the API does not replay are retried only on `429`, which the API refuses before acting:
creating or deleting an SSH key, opening an instance console, and testing, deleting, rotating the
secret of or redelivering a webhook. `WithNoRetry()` sends a call once. Cancelling the
context stops a call and its retries at once.

## Development

```sh
go vet ./... && go test -race ./...
scripts/sync-spec.sh   # refresh testdata/public-api.yaml from the published spec
```

The SDK follows the published [OpenAPI spec](https://docs.pantechdynamics.com/openapi/public-api.yaml).
See [CONTRIBUTING.md](CONTRIBUTING.md) for the workflow and how releases are made. Report security
issues as described in [SECURITY.md](SECURITY.md).

## License

MIT. See [LICENSE](LICENSE).
