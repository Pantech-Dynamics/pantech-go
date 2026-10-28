package pantech

import (
	"bufio"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"testing"
)

// operationMethods maps every operationId in the public spec to the SDK method
// that calls it. TestSpecCoverage fails when the spec has an operation that is
// missing here, so a new endpoint cannot ship without a method. Method
// expressions make a renamed or deleted method a compile error.
var operationMethods = map[string]any{
	"getMe":                             (*MeService).Get,
	"revokeMe":                          (*MeService).Revoke,
	"getCredit":                         (*BillingService).GetCredit,
	"createCreditPurchase":              (*BillingService).CreateCreditPurchase,
	"getCreditPurchase":                 (*BillingService).GetCreditPurchase,
	"listReceipts":                      (*BillingService).ListReceipts,
	"getReceipt":                        (*BillingService).GetReceipt,
	"getReceiptPdf":                     (*BillingService).GetReceiptPDF,
	"getUsage":                          (*BillingService).GetUsage,
	"getUsageByResource":                (*BillingService).GetUsageByResource,
	"listRegions":                       (*CatalogService).ListRegions,
	"listPlans":                         (*CatalogService).ListPlans,
	"listImages":                        (*CatalogService).ListImages,
	"listDiskOfferings":                 (*CatalogService).ListDiskOfferings,
	"quoteDiskOffering":                 (*CatalogService).QuoteDiskOffering,
	"listDatabaseEngines":               (*CatalogService).ListDatabaseEngines,
	"listKubernetesVersions":            (*CatalogService).ListKubernetesVersions,
	"estimatePricing":                   (*CatalogService).EstimatePrice,
	"listInstances":                     (*InstancesService).List,
	"createInstance":                    (*InstancesService).Create,
	"getInstance":                       (*InstancesService).Get,
	"deleteInstance":                    (*InstancesService).Delete,
	"startInstance":                     (*InstancesService).Start,
	"stopInstance":                      (*InstancesService).Stop,
	"rebootInstance":                    (*InstancesService).Reboot,
	"rebuildInstance":                   (*InstancesService).Rebuild,
	"renameInstance":                    (*InstancesService).Rename,
	"resizeInstance":                    (*InstancesService).Resize,
	"changeInstanceSSHKey":              (*InstancesService).ChangeSSHKey,
	"restoreInstance":                   (*InstancesService).Restore,
	"changeInstanceSecurityGroup":       (*InstancesService).ChangeSecurityGroup,
	"attachInstancePrivateNetwork":      (*InstancesService).AttachPrivateNetwork,
	"detachInstancePrivateNetwork":      (*InstancesService).DetachPrivateNetwork,
	"openInstanceConsole":               (*InstancesService).OpenConsole,
	"issueVolumeKeyCredential":          (*InstancesService).IssueVolumeKeyCredential,
	"getInstanceMetrics":                (*InstancesService).GetMetrics,
	"getInstanceUsage":                  (*InstancesService).GetUsage,
	"listInstanceVolumes":               (*InstancesService).ListVolumes,
	"listInstanceSnapshots":             (*InstancesService).ListSnapshots,
	"createInstanceSnapshot":            (*InstancesService).CreateSnapshot,
	"getInstanceSnapshotSchedule":       (*InstancesService).GetSnapshotSchedule,
	"putInstanceSnapshotSchedule":       (*InstancesService).PutSnapshotSchedule,
	"openSSHAccess":                     (*InstancesService).OpenSSHAccess,
	"getSSHAccess":                      (*InstancesService).GetSSHAccess,
	"revokeSSHAccess":                   (*InstancesService).RevokeSSHAccess,
	"listInstanceOrders":                (*InstancesService).ListOrders,
	"getInstanceOrder":                  (*InstancesService).GetOrder,
	"listNetworks":                      (*NetworksService).List,
	"createNetwork":                     (*NetworksService).Create,
	"quoteNetwork":                      (*NetworksService).Quote,
	"getNetwork":                        (*NetworksService).Get,
	"deleteNetwork":                     (*NetworksService).Delete,
	"listSubnets":                       (*SubnetsService).List,
	"createSubnet":                      (*SubnetsService).Create,
	"getSubnet":                         (*SubnetsService).Get,
	"deleteSubnet":                      (*SubnetsService).Delete,
	"listFirewallRules":                 (*FirewallRulesService).List,
	"createFirewallRule":                (*FirewallRulesService).Create,
	"getFirewallRule":                   (*FirewallRulesService).Get,
	"deleteFirewallRule":                (*FirewallRulesService).Delete,
	"listPublicIPs":                     (*PublicIPsService).List,
	"createPublicIP":                    (*PublicIPsService).Create,
	"quotePublicIP":                     (*PublicIPsService).Quote,
	"getPublicIP":                       (*PublicIPsService).Get,
	"deletePublicIP":                    (*PublicIPsService).Delete,
	"attachPublicIP":                    (*PublicIPsService).Attach,
	"detachPublicIP":                    (*PublicIPsService).Detach,
	"listPortForwardingRules":           (*PortForwardingRulesService).List,
	"createPortForwardingRule":          (*PortForwardingRulesService).Create,
	"getPortForwardingRule":             (*PortForwardingRulesService).Get,
	"deletePortForwardingRule":          (*PortForwardingRulesService).Delete,
	"listSecurityGroups":                (*SecurityGroupsService).List,
	"createSecurityGroup":               (*SecurityGroupsService).Create,
	"getSecurityGroup":                  (*SecurityGroupsService).Get,
	"deleteSecurityGroup":               (*SecurityGroupsService).Delete,
	"replaceSecurityGroupRules":         (*SecurityGroupsService).ReplaceRules,
	"listLoadBalancers":                 (*LoadBalancersService).List,
	"createLoadBalancer":                (*LoadBalancersService).Create,
	"getLoadBalancer":                   (*LoadBalancersService).Get,
	"updateLoadBalancer":                (*LoadBalancersService).Update,
	"deleteLoadBalancer":                (*LoadBalancersService).Delete,
	"listVolumes":                       (*VolumesService).List,
	"createVolume":                      (*VolumesService).Create,
	"getVolume":                         (*VolumesService).Get,
	"deleteVolume":                      (*VolumesService).Delete,
	"attachVolume":                      (*VolumesService).Attach,
	"detachVolume":                      (*VolumesService).Detach,
	"resizeVolume":                      (*VolumesService).Resize,
	"getVolumeSnapshotSchedule":         (*VolumesService).GetSnapshotSchedule,
	"putVolumeSnapshotSchedule":         (*VolumesService).PutSnapshotSchedule,
	"listVolumeSnapshots":               (*VolumesService).ListSnapshots,
	"createVolumeSnapshot":              (*VolumesService).CreateSnapshot,
	"getVolumeUsage":                    (*VolumesService).GetUsage,
	"listSnapshots":                     (*SnapshotsService).List,
	"quoteSnapshot":                     (*SnapshotsService).Quote,
	"getSnapshot":                       (*SnapshotsService).Get,
	"deleteSnapshot":                    (*SnapshotsService).Delete,
	"restoreSnapshot":                   (*SnapshotsService).Restore,
	"getSnapshotUsage":                  (*SnapshotsService).GetUsage,
	"listDatabases":                     (*DatabasesService).List,
	"createDatabase":                    (*DatabasesService).Create,
	"getDatabase":                       (*DatabasesService).Get,
	"deleteDatabase":                    (*DatabasesService).Delete,
	"replaceDatabaseAccessRules":        (*DatabasesService).ReplaceAccessRules,
	"setDatabaseSecurityGroups":         (*DatabasesService).SetSecurityGroups,
	"changeDatabasePassword":            (*DatabasesService).ChangePassword,
	"resetDatabasePassword":             (*DatabasesService).ResetPassword,
	"startDatabase":                     (*DatabasesService).Start,
	"stopDatabase":                      (*DatabasesService).Stop,
	"resizeDatabaseStorage":             (*DatabasesService).ResizeStorage,
	"listDatabaseSnapshots":             (*DatabasesService).ListSnapshots,
	"createDatabaseSnapshot":            (*DatabasesService).CreateSnapshot,
	"getDatabaseSnapshot":               (*DatabasesService).GetSnapshot,
	"deleteDatabaseSnapshot":            (*DatabasesService).DeleteSnapshot,
	"listDatabaseOrders":                (*DatabasesService).ListOrders,
	"getDatabaseOrder":                  (*DatabasesService).GetOrder,
	"listKubernetesClusters":            (*KubernetesService).ListClusters,
	"createKubernetesCluster":           (*KubernetesService).CreateCluster,
	"getKubernetesCluster":              (*KubernetesService).GetCluster,
	"scaleKubernetesCluster":            (*KubernetesService).UpdateCluster,
	"deleteKubernetesCluster":           (*KubernetesService).DeleteCluster,
	"upgradeKubernetesCluster":          (*KubernetesService).UpgradeCluster,
	"resizeKubernetesClusterPlan":       (*KubernetesService).ResizePlan,
	"upgradeKubernetesClusterNetwork":   (*KubernetesService).UpgradeNetwork,
	"startKubernetesCluster":            (*KubernetesService).StartCluster,
	"stopKubernetesCluster":             (*KubernetesService).StopCluster,
	"getKubernetesKubeconfig":           (*KubernetesService).GetKubeconfig,
	"listKubernetesClusterNodes":        (*KubernetesService).ListNodes,
	"getKubernetesWorkloads":            (*KubernetesService).GetWorkloads,
	"applyKubernetesManifest":           (*KubernetesService).ApplyManifest,
	"getKubernetesLiveManifests":        (*KubernetesService).GetLiveManifests,
	"exposeKubernetesService":           (*KubernetesService).ExposeService,
	"syncKubernetesRouting":             (*KubernetesService).SyncRouting,
	"listKubernetesNodePools":           (*KubernetesService).ListNodePools,
	"createKubernetesNodePool":          (*KubernetesService).CreateNodePool,
	"scaleKubernetesNodePool":           (*KubernetesService).ScaleNodePool,
	"deleteKubernetesNodePool":          (*KubernetesService).DeleteNodePool,
	"listKubernetesAddons":              (*KubernetesService).ListAddons,
	"toggleKubernetesAddon":             (*KubernetesService).ToggleAddon,
	"configureKubernetesEndpointAccess": (*KubernetesService).ConfigureEndpointAccess,
	"listSSHKeys":                       (*SSHKeysService).List,
	"createSSHKey":                      (*SSHKeysService).Create,
	"getSSHKey":                         (*SSHKeysService).Get,
	"deleteSSHKey":                      (*SSHKeysService).Delete,
	"listWebhookEventTypes":             (*WebhooksService).ListEventTypes,
	"listWebhooks":                      (*WebhooksService).List,
	"createWebhook":                     (*WebhooksService).Create,
	"getWebhook":                        (*WebhooksService).Get,
	"updateWebhook":                     (*WebhooksService).Update,
	"deleteWebhook":                     (*WebhooksService).Delete,
	"testWebhook":                       (*WebhooksService).Test,
	"rotateWebhookSecret":               (*WebhooksService).RotateSecret,
	"listWebhookDeliveries":             (*WebhooksService).ListDeliveries,
	"getWebhookDelivery":                (*WebhooksService).GetDelivery,
	"redeliverWebhook":                  (*WebhooksService).Redeliver,
	"listCloudConnections":              (*MigrateService).ListConnections,
	"createCloudConnection":             (*MigrateService).CreateConnection,
	"getCloudConnection":                (*MigrateService).GetConnection,
	"updateCloudConnection":             (*MigrateService).UpdateConnection,
	"deleteCloudConnection":             (*MigrateService).DeleteConnection,
	"listCloudScans":                    (*MigrateService).ListScans,
	"startCloudScan":                    (*MigrateService).StartScan,
	"getCloudScan":                      (*MigrateService).GetScan,
	"listMigrations":                    (*MigrateService).ListMigrations,
	"createMigration":                   (*MigrateService).CreateMigration,
	"getMigration":                      (*MigrateService).GetMigration,
	"renameMigration":                   (*MigrateService).RenameMigration,
	"deleteMigration":                   (*MigrateService).DeleteMigration,
	"updateMigrationItem":               (*MigrateService).UpdateItem,
	"exportMigrationTerraform":          (*MigrateService).ExportTerraform,
	"getOperation":                      (*OperationsService).Get,
	"streamOperation":                   (*OperationsService).Stream,
}

// aheadOfSpec are operations the SDK already supports that the spec in
// testdata does not list yet. Remove an entry once the spec has it.
var aheadOfSpec = map[string]bool{}

// notInSDK are spec operations the SDK deliberately has no method for yet,
// each with the reason. Agent sessions answer only AI apps connected through
// Pantech's MCP server; a standard API key gets 403
// AGENT_SESSIONS_NOT_FOR_API_KEYS. Remove an entry when its method is added.
var notInSDK = map[string]string{
	"getInstanceAgent":       "agent sessions: MCP-connected AI apps only",
	"startAgentSession":      "agent sessions: MCP-connected AI apps only",
	"listAgentSessions":      "agent sessions: MCP-connected AI apps only",
	"getAgentSession":        "agent sessions: MCP-connected AI apps only",
	"messageAgentSession":    "agent sessions: MCP-connected AI apps only",
	"callAgentTool":          "agent sessions: MCP-connected AI apps only",
	"getAgentToolCall":       "agent sessions: MCP-connected AI apps only",
	"listAgentSessionEvents": "agent sessions: MCP-connected AI apps only",
	"decideAgentApproval":    "agent sessions: MCP-connected AI apps only",
	"terminateAgentSession":  "agent sessions: MCP-connected AI apps only",
	"getAgentSessionReport":  "agent sessions: MCP-connected AI apps only",
}

// specOperation is one operation read from the spec.
type specOperation struct {
	id, method, path string
}

var (
	specPathLine   = regexp.MustCompile(`^  "?(/[^":]*)"?:\s*$`)
	specMethodLine = regexp.MustCompile(`^    (get|put|post|patch|delete):\s*$`)
	specOpIDLine   = regexp.MustCompile(`^\s+operationId:\s*"?([A-Za-z0-9_]+)"?\s*$`)
)

// readSpecOperations is a deliberately tiny line-based reader for the spec's
// paths, methods and operationIds, so the SDK needs no YAML dependency. It
// relies on the spec's two-space indentation (paths at 2, methods at 4).
func readSpecOperations(t *testing.T, file string) []specOperation {
	t.Helper()
	f, err := os.Open(file)
	if err != nil {
		t.Fatalf("opening the spec: %v", err)
	}
	defer f.Close()
	var ops []specOperation
	var path, method string
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	for sc.Scan() {
		line := sc.Text()
		if m := specPathLine.FindStringSubmatch(line); m != nil {
			path, method = strings.TrimPrefix(m[1], "/public/v1"), ""
			continue
		}
		if strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "   ") && strings.TrimSpace(line) != "" {
			path = "" // another key at path level, such as a components section
		}
		if m := specMethodLine.FindStringSubmatch(line); m != nil && path != "" {
			method = strings.ToUpper(m[1])
			continue
		}
		if m := specOpIDLine.FindStringSubmatch(line); m != nil && method != "" {
			ops = append(ops, specOperation{id: m[1], method: method, path: path})
			method = ""
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("reading the spec: %v", err)
	}
	if len(ops) == 0 {
		t.Fatalf("found no operations in %s", file)
	}
	return ops
}

func TestSpecCoverage(t *testing.T) {
	ops := readSpecOperations(t, "testdata/public-api.yaml")
	inSpec := map[string]bool{}
	for _, op := range ops {
		inSpec[op.id] = true
		if _, skip := notInSDK[op.id]; skip {
			if _, ok := operationMethods[op.id]; ok {
				t.Errorf("%s has a method and is listed in notInSDK: remove it from notInSDK", op.id)
			}
			continue
		}
		if _, ok := operationMethods[op.id]; !ok {
			t.Errorf("spec operation %s (%s %s) has no SDK method: add one and list it in operationMethods", op.id, op.method, op.path)
		}
	}
	for id := range notInSDK {
		if !inSpec[id] {
			t.Errorf("notInSDK lists %s, which is not in the spec: remove it", id)
		}
	}
	for id := range operationMethods {
		if !inSpec[id] && !aheadOfSpec[id] {
			t.Errorf("operationMethods lists %s, which is not in the spec: remove it or add it to aheadOfSpec", id)
		}
	}
	t.Logf("%d operations in the spec, %d mapped, %d not in the SDK", len(ops), len(operationMethods), len(notInSDK))
}

// TestSpecRoutes calls every mapped method against a recording server and
// checks it sends the spec's method and path.
func TestSpecRoutes(t *testing.T) {
	var mu sync.Mutex
	var got []*http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		got = append(got, r)
		mu.Unlock()
		if strings.HasSuffix(r.URL.Path, "/events") {
			w.Header().Set("Content-Type", "text/event-stream")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()
	c := newTestClientAt(t, srv.URL)

	routes := map[string]specOperation{}
	for _, op := range readSpecOperations(t, "testdata/public-api.yaml") {
		routes[op.id] = op
	}

	for id, fn := range operationMethods {
		op, ok := routes[id]
		if !ok {
			continue
		}
		mu.Lock()
		got = nil
		mu.Unlock()
		callWithPlaceholders(t, c, fn)
		mu.Lock()
		reqs := got
		mu.Unlock()
		if len(reqs) != 1 {
			t.Errorf("%s: made %d requests, want 1", id, len(reqs))
			continue
		}
		r := reqs[0]
		// QuoteMeta escapes the braces of {param}, hence the escaped pattern.
		want := regexp.MustCompile("^/public/v1" + regexp.MustCompile(`\\\{[^}]+\\\}`).ReplaceAllString(regexp.QuoteMeta(op.path), `[^/]+`) + "$")
		if r.Method != op.method || !want.MatchString(r.URL.Path) {
			t.Errorf("%s: sent %s %s, want %s %s", id, r.Method, r.URL.Path, op.method, op.path)
		}
		if isWrite(r.Method) && r.Header.Get("Idempotency-Key") == "" {
			t.Errorf("%s: write sent without an Idempotency-Key", id)
		}
	}
}

// callWithPlaceholders calls a method expression with its service from c and
// placeholder arguments, and drains an iterator if it returns one.
func callWithPlaceholders(t *testing.T, c *Client, fn any) {
	t.Helper()
	v := reflect.ValueOf(fn)
	ft := v.Type()
	recv := serviceFor(t, c, ft.In(0))
	args := []reflect.Value{recv, reflect.ValueOf(context.Background())}
	n := ft.NumIn()
	if ft.IsVariadic() {
		n-- // ...RequestOption: Call fills it in
	}
	for i := 2; i < n; i++ {
		switch in := ft.In(i); in.Kind() {
		case reflect.String:
			args = append(args, reflect.ValueOf("id_"+string(rune('a'+i))))
		case reflect.Int:
			args = append(args, reflect.ValueOf(1))
		default:
			args = append(args, reflect.Zero(in))
		}
	}
	out := v.Call(args)
	if len(out) == 1 && out[0].Kind() == reflect.Func { // iter.Seq2
		yield := reflect.MakeFunc(out[0].Type().In(0), func([]reflect.Value) []reflect.Value {
			return []reflect.Value{reflect.ValueOf(false)}
		})
		out[0].Call([]reflect.Value{yield})
	}
}

func serviceFor(t *testing.T, c *Client, typ reflect.Type) reflect.Value {
	t.Helper()
	cv := reflect.ValueOf(c).Elem()
	for i := 0; i < cv.NumField(); i++ {
		if f := cv.Field(i); f.Type() == typ && f.CanInterface() {
			return f
		}
	}
	t.Fatalf("Client has no field of type %s", typ)
	return reflect.Value{}
}
