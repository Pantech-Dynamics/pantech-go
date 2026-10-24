package pantech

import (
	"context"
	"testing"
)

func TestAllFollowsTheCursor(t *testing.T) {
	rec := newRecorder(t,
		ok(`{"data":[{"id":"vm_1"},{"id":"vm_2"}],"next_cursor":"c2"}`),
		ok(`{"data":[{"id":"vm_3"}],"next_cursor":"c3"}`),
		ok(`{"data":[{"id":"vm_4"}],"next_cursor":null}`),
	)
	c := newTestClientAt(t, rec.srv.URL)
	var ids []string
	for vm, err := range c.Instances.All(context.Background(), &ListParams{Limit: 2, Q: "web"}) {
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, vm.ID)
	}
	if len(ids) != 4 || ids[3] != "vm_4" {
		t.Fatalf("ids = %v", ids)
	}
	wantQueries := []string{"limit=2&q=web", "cursor=c2&limit=2&q=web", "cursor=c3&limit=2&q=web"}
	for i, want := range wantQueries {
		if got := rec.request(i).URL.RawQuery; got != want {
			t.Errorf("page %d query = %q, want %q", i, got, want)
		}
	}
}

func TestAllStopsWhenTheCaller(t *testing.T) {
	rec := newRecorder(t, ok(`{"data":[{"id":"a"},{"id":"b"}],"next_cursor":"c2"}`))
	c := newTestClientAt(t, rec.srv.URL)
	for range c.Snapshots.All(context.Background(), nil) {
		break
	}
	if rec.count() != 1 {
		t.Fatalf("%d requests after an early break", rec.count())
	}
}

func TestAllStopsOnARepeatedCursor(t *testing.T) {
	rec := newRecorder(t, ok(`{"data":[{"id":"a"}],"next_cursor":"same"}`))
	c := newTestClientAt(t, rec.srv.URL)
	n := 0
	for _, err := range c.Volumes.All(context.Background(), nil) {
		if err != nil {
			t.Fatal(err)
		}
		n++
	}
	if n != 2 || rec.count() != 2 {
		t.Fatalf("%d items from %d requests", n, rec.count())
	}
}

func TestAllYieldsTheError(t *testing.T) {
	rec := newRecorder(t,
		ok(`{"data":[{"id":"a"}],"next_cursor":"c2"}`),
		status(400, `{"code":"INVALID_CURSOR","detail":"Start again without a cursor."}`),
	)
	c := newTestClientAt(t, rec.srv.URL)
	var items int
	var last error
	for _, err := range c.Networks.All(context.Background(), nil) {
		if err != nil {
			last = err
			continue
		}
		items++
	}
	if items != 1 || !IsCode(last, "INVALID_CURSOR") {
		t.Fatalf("items %d, err %v", items, last)
	}
}

func TestNestedAllPassesThePathParameter(t *testing.T) {
	rec := newRecorder(t, ok(`{"data":[{"id":"snap_1"}],"next_cursor":null}`))
	c := newTestClientAt(t, rec.srv.URL)
	for _, err := range c.Databases.AllSnapshots(context.Background(), "db_1", &ListParams{ObservedState: "active"}) {
		if err != nil {
			t.Fatal(err)
		}
	}
	r := rec.request(0)
	if r.URL.Path != "/public/v1/databases/db_1/snapshots" || r.URL.RawQuery != "observed_state=active" {
		t.Fatalf("sent %s?%s", r.URL.Path, r.URL.RawQuery)
	}
}

func TestSinglePageListsFollowTheCursorAnyway(t *testing.T) {
	rec := newRecorder(t,
		ok(`{"data":[{"slug":"individual"}],"next_cursor":"more"}`),
		ok(`{"data":[{"slug":"starter"}],"next_cursor":null}`),
	)
	c := newTestClientAt(t, rec.srv.URL)
	plans, err := c.Catalog.ListPlans(context.Background(), &ListPlansParams{Placement: "vpc"})
	if err != nil {
		t.Fatal(err)
	}
	if len(plans) != 2 || plans[1].Slug != "starter" {
		t.Fatalf("plans = %+v", plans)
	}
	if q := rec.request(1).URL.RawQuery; q != "cursor=more&placement=vpc" {
		t.Fatalf("second query = %q", q)
	}
}

func TestEmptyListIsNotNil(t *testing.T) {
	rec := newRecorder(t, ok(`{"data":[],"next_cursor":null}`))
	c := newTestClientAt(t, rec.srv.URL)
	keys, err := c.SSHKeys.List(context.Background())
	if err != nil || keys == nil || len(keys) != 0 {
		t.Fatalf("keys %v err %v", keys, err)
	}
}
