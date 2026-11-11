package pantech

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

const repositoryJSON = `{"id":"crr_1","namespace":"acme","name":"api/server","full_name":"registry.pantechdynamics.com/acme/api/server",
"description":"","visibility":"private","immutable_tags":false,"retention":{"keep_last":50,"untagged_after_days":null,"keep_pattern":"v*"},
"tag_count":31,"storage_bytes":1048576000,"last_pushed_at":"2026-11-10T10:00:00Z","status":"active","suspended_reason":null,
"created_at":"2026-11-01T09:00:00Z","updated_at":"2026-11-10T10:00:00Z","version":4}`

const tokenJSON = `{"id":"prt_1","name":"github-actions","permission":"push","repositories":["api/server"],"expires_at":null,
"last_used_at":null,"created_at":"2026-11-10T09:00:00Z","revoked_at":null,"secret":"prt_1_s3cretvalue"}`

func TestRegistryNamespaceNotSetIsACode(t *testing.T) {
	rec := newRecorder(t, status(404, `{"title":"Not found","code":"REGISTRY_NAMESPACE_NOT_SET","detail":"Claim a namespace first."}`))
	c := newTestClientAt(t, rec.srv.URL)
	_, err := c.Registry.GetNamespace(context.Background())
	if !IsCode(err, CodeRegistryNamespaceNotSet) || rec.request(0).URL.Path != "/public/v1/registry" {
		t.Fatalf("err %v", err)
	}
}

func TestRegistrySetNamespace(t *testing.T) {
	rec := newRecorder(t, ok(`{"name":"acme","registry_host":"registry.pantechdynamics.com","create_on_push":false,"repository_count":0,
"storage_bytes":0,"storage_limit_bytes":53687091200,"repository_limit":100,"created_at":"2026-11-10T09:00:00Z","updated_at":"2026-11-10T09:00:00Z","version":1}`))
	c := newTestClientAt(t, rec.srv.URL)
	ns, err := c.Registry.SetNamespace(context.Background(), SetRegistryNamespaceRequest{Name: "acme", CreateOnPush: Ptr(false)})
	if err != nil {
		t.Fatal(err)
	}
	r := rec.request(0)
	if r.Method != "PUT" || r.URL.Path != "/public/v1/registry" || r.Header.Get("Idempotency-Key") == "" {
		t.Fatalf("sent %s %s", r.Method, r.URL.Path)
	}
	if rec.bodies[0] != `{"name":"acme","create_on_push":false}` {
		t.Fatalf("body %s", rec.bodies[0])
	}
	if ns.RegistryHost != DefaultRegistryHost || ns.StorageLimitBytes != 53687091200 {
		t.Fatalf("namespace %+v", ns)
	}
}

func TestRegistryCreateRepository(t *testing.T) {
	rec := newRecorder(t, status(201, repositoryJSON))
	c := newTestClientAt(t, rec.srv.URL)
	repo, err := c.Registry.CreateRepository(context.Background(), CreateRegistryRepositoryRequest{
		Name: "api/server", Retention: &RegistryRetention{KeepLast: Ptr(50), KeepPattern: Ptr("v*")},
	})
	if err != nil {
		t.Fatal(err)
	}
	if repo.ID != "crr_1" || repo.Retention == nil || *repo.Retention.KeepLast != 50 || repo.Retention.UntaggedAfterDays != nil {
		t.Fatalf("repo %+v", repo)
	}
	// Retention is sent whole: a rule left nil is an explicit null.
	want := `{"name":"api/server","retention":{"keep_last":50,"untagged_after_days":null,"keep_pattern":"v*"}}`
	if rec.bodies[0] != want {
		t.Fatalf("body %s", rec.bodies[0])
	}
}

func TestRegistryUpdateRepositorySendsOnlyWhatIsSet(t *testing.T) {
	rec := newRecorder(t, ok(repositoryJSON))
	c := newTestClientAt(t, rec.srv.URL)
	if _, err := c.Registry.UpdateRepository(context.Background(), "crr_1", UpdateRegistryRepositoryRequest{Visibility: Ptr(RegistryPublic)}); err != nil {
		t.Fatal(err)
	}
	if r := rec.request(0); r.Method != "PATCH" || r.URL.Path != "/public/v1/registry/repositories/crr_1" || rec.bodies[0] != `{"visibility":"public"}` {
		t.Fatalf("sent %s %s %s", r.Method, r.URL.Path, rec.bodies[0])
	}
}

func TestRegistryAllTagsFollowsTheCursor(t *testing.T) {
	tag := `{"name":"v1","digest":"sha256:aa","media_type":"application/vnd.oci.image.index.v1+json","size_bytes":10,"platforms":["linux/amd64"],"pushed_at":"2026-11-10T10:00:00Z","last_pulled_at":null,"immutable":false}`
	rec := newRecorder(t,
		ok(`{"data":[`+tag+`],"next_cursor":"c2"}`),
		ok(`{"data":[`+strings.Replace(tag, `"v1"`, `"v0"`, 1)+`],"next_cursor":null}`),
	)
	c := newTestClientAt(t, rec.srv.URL)
	var names []string
	for tg, err := range c.Registry.AllTags(context.Background(), "crr_1", &ListRegistryParams{Q: "v"}) {
		if err != nil {
			t.Fatal(err)
		}
		names = append(names, tg.Name)
	}
	if len(names) != 2 || rec.request(1).URL.Query().Get("cursor") != "c2" || rec.request(1).URL.Query().Get("q") != "v" {
		t.Fatalf("names %v query %s", names, rec.request(1).URL.RawQuery)
	}
	if rec.request(0).URL.Path != "/public/v1/registry/repositories/crr_1/tags" {
		t.Fatalf("path %s", rec.request(0).URL.Path)
	}
}

func TestRegistryDeleteTagEscapesTheTag(t *testing.T) {
	rec := newRecorder(t, status(202, ``))
	c := newTestClientAt(t, rec.srv.URL)
	if err := c.Registry.DeleteTag(context.Background(), "crr_1", "v1/../x"); err != nil {
		t.Fatal(err)
	}
	r := rec.request(0)
	if r.Method != "DELETE" || r.URL.EscapedPath() != "/public/v1/registry/repositories/crr_1/tags/v1%2F..%2Fx" {
		t.Fatalf("sent %s %s", r.Method, r.URL.EscapedPath())
	}
}

func TestRegistryTagImmutableIsACode(t *testing.T) {
	rec := newRecorder(t, status(409, `{"title":"Tag immutable","code":"TAG_IMMUTABLE","detail":"This repository's tags cannot be deleted."}`))
	c := newTestClientAt(t, rec.srv.URL)
	if err := c.Registry.DeleteTag(context.Background(), "crr_1", "v1"); !IsCode(err, CodeTagImmutable) {
		t.Fatalf("err %v", err)
	}
}

func TestRegistryCreateTokenSecretDoesNotPrint(t *testing.T) {
	rec := newRecorder(t, status(201, tokenJSON))
	c := newTestClientAt(t, rec.srv.URL)
	tok, err := c.Registry.CreateToken(context.Background(), CreateRegistryTokenRequest{Name: "github-actions", Permission: RegistryPush, Repositories: []string{"api/server"}})
	if err != nil {
		t.Fatal(err)
	}
	if tok.Secret != "prt_1_s3cretvalue" {
		t.Fatalf("secret %q", tok.Secret)
	}
	if rec.bodies[0] != `{"name":"github-actions","permission":"push","repositories":["api/server"]}` {
		t.Fatalf("body %s", rec.bodies[0])
	}
	holder := struct{ Token *RegistryToken }{tok}
	for _, s := range []string{fmt.Sprint(tok), fmt.Sprintf("%v", *tok), fmt.Sprintf("%+v", *tok), fmt.Sprintf("%#v", tok), fmt.Sprintf("%+v", holder)} {
		if strings.Contains(s, "s3cretvalue") {
			t.Fatalf("printed the secret: %s", s)
		}
	}
	if !strings.Contains(fmt.Sprint(*tok), "prt_1") {
		t.Fatalf("String %q", fmt.Sprint(*tok))
	}
	// Marshalled on purpose (to save it), the secret is kept.
	data, _ := json.Marshal(tok)
	if !strings.Contains(string(data), "s3cretvalue") {
		t.Fatalf("json %s", data)
	}
}

func TestRegistryListAndDeleteTokens(t *testing.T) {
	rec := newRecorder(t, ok(`{"data":[`+strings.Replace(tokenJSON, `,"secret":"prt_1_s3cretvalue"`, "", 1)+`],"next_cursor":null}`), status(204, ""))
	c := newTestClientAt(t, rec.srv.URL)
	toks, err := c.Registry.ListTokens(context.Background())
	if err != nil || len(toks) != 1 || toks[0].Secret != "" || rec.request(0).URL.Path != "/public/v1/registry/tokens" {
		t.Fatalf("tokens %v err %v", toks, err)
	}
	if err := c.Registry.DeleteToken(context.Background(), "prt_1"); err != nil {
		t.Fatal(err)
	}
	if r := rec.request(1); r.Method != "DELETE" || r.URL.Path != "/public/v1/registry/tokens/prt_1" {
		t.Fatalf("sent %s %s", r.Method, r.URL.Path)
	}
}
