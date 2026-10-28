# Contributing to pantech-go

Thanks for helping improve the Go SDK for the Pantech Dynamics public API. Bug reports, fixes and
new methods are all welcome. Please read the [Code of Conduct](CODE_OF_CONDUCT.md) first, and
report security issues privately as described in [SECURITY.md](SECURITY.md), never in a public
issue.

## Development setup

You need Go 1.23 or newer. The module has no dependencies outside the standard library.

```sh
git clone https://github.com/Pantech-Dynamics/pantech-go
cd pantech-go
go build ./...
```

## Checks

Run these before opening a pull request; CI runs the same ones on the oldest supported Go and the
newest.

```sh
gofmt -l .               # must print nothing
go vet ./...
go test -race ./...
go build ./...           # includes examples/
```

## Layout

```
client.go          NewClient, options, Do, RequestOptions, headers, body cap
retry.go           retry policy, backoff, Retry-After, the routes that do not replay
errors.go          *Error (problem+json), FieldError, IsCode, AsError
pagination.go      Page[T], the iter.Seq2 paginator
wait.go            Operations.Wait, Instances/Databases.WaitForOrder
stream.go          Operations.Stream (Server-Sent Events)
webhook.go         VerifyWebhook, ParseWebhook, SignWebhook
types.go           request and response structs, one per schema in the spec
params.go          query parameter structs
<service>.go       one file per service: instances.go, networks.go, ...
coverage_test.go   operationId -> method mapping and the spec tests
testdata/          the vendored copy of the public API spec
examples/          runnable programs; they must build
```

## Keeping in step with the API spec

The SDK follows the published OpenAPI spec at
<https://docs.pantechdynamics.com/openapi/public-api.yaml>. A copy lives in
`testdata/public-api.yaml`, and a weekly workflow checks it against the published one.

1. Refresh the copy:

   ```sh
   scripts/sync-spec.sh                  # the published spec
   scripts/sync-spec.sh path/or/url.yaml # or another copy
   ```

2. Run `go test -run 'TestSpec' ./...`. It names each operation that has no method.
3. For each one:
   - add or update the structs in `types.go`. Field names are the spec's in Go case
     (`instance_id` becomes `InstanceID`) and JSON tags are the spec's names. In a request an
     optional field is a pointer with `omitempty` (set it with `Ptr`); in a response a nullable
     field is a pointer. `date-time` is `time.Time`; integers in minor units or bytes are `int64`.
   - add the method to its service file in the same shape as its neighbours: path parameters as
     strings, query parameters as a `*...Params` struct, the body as a value, and
     `opts ...RequestOption` last. A paginated list gets `List` returning `*Page[T]` and an
     `All...` iterator; a `202` returns `*Accepted`.
   - add the row to `operationMethods` in `coverage_test.go`.
   - an operation the SDK deliberately leaves out goes in `notInSDK` in `coverage_test.go`, with
     the reason.
4. A route that does not replay on the same Idempotency-Key goes in `noReplay` in `retry.go`, with
   a test in `retry_test.go`.

## Guidelines

- Every request goes through `Client.Do`, which owns auth, the Idempotency-Key, retries and error
  mapping (`Stream` is the one exception).
- `Do`'s signature, the option names and `*Error`'s fields are public API: do not change them in a
  minor release.
- Keep the module dependency-free and on `go 1.23`.
- Never log or print the API key, and never put it in a URL or an error.
- Add tests for every change in behaviour.
- A change users will notice gets a line under `## Unreleased` in `CHANGELOG.md`, in the same pull
  request.

## Releasing

Maintainers release by renaming `## Unreleased` in `CHANGELOG.md` to `## vX.Y.Z — YYYY-MM-DD` (Go
tags need the `v`) and bumping `Version` in `client.go` to match; a test keeps the two in step. On
push to `main`, the auto-release workflow tags the newest version heading if it is new and
publishes a GitHub release whose notes are that section. Go modules need no build artifacts: the
tag is the release.
