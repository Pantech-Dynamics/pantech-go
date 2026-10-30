package pantech

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// A pool whose workers cannot be built says why; one with nothing failing
// reads back with no failure.
func TestListNodePoolsReadsTheBuildFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/v1/kubernetes-clusters/k8s_1/node-pools" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"data":[
			{"id":"k8spool_1","cluster_id":"k8s_1","name":"db","status":"provisioning","observed_workers":0,"workers":2,
			 "failure_code":"PROVISIONING_CAPACITY_UNAVAILABLE","failure_reason":"There is no room for this size in the region right now."},
			{"id":"k8spool_2","cluster_id":"k8s_1","name":"web","status":"ready","observed_workers":1,"workers":1,
			 "failure_code":null,"failure_reason":null}],"next_cursor":null}`)
	}))
	defer srv.Close()

	pools, err := newTestClientAt(t, srv.URL).Kubernetes.ListNodePools(context.Background(), "k8s_1")
	if err != nil {
		t.Fatal(err)
	}
	if len(pools) != 2 {
		t.Fatalf("pools = %d, want 2", len(pools))
	}
	if p := pools[0]; p.FailureCode == nil || *p.FailureCode != "PROVISIONING_CAPACITY_UNAVAILABLE" || p.FailureReason == nil || *p.FailureReason == "" {
		t.Fatalf("failed pool = code %v reason %v", p.FailureCode, p.FailureReason)
	}
	if p := pools[1]; p.FailureCode != nil || p.FailureReason != nil {
		t.Fatalf("ready pool = code %v reason %v", *p.FailureCode, *p.FailureReason)
	}
}
