package pantech

import (
	"context"
	"net/url"
)

// KubernetesService manages Kubernetes clusters, node pools and add-ons.
type KubernetesService struct{ c *Client }

// ListClusters lists Kubernetes clusters: GET /kubernetes-clusters.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *KubernetesService) ListClusters(ctx context.Context, opts ...RequestOption) ([]KubernetesCluster, error) {
	return listAll[KubernetesCluster](ctx, s.c, "/kubernetes-clusters", nil, opts...)
}

// CreateCluster creates a Kubernetes cluster: POST /kubernetes-clusters.
// Scope: write.
func (s *KubernetesService) CreateCluster(ctx context.Context, req CreateKubernetesClusterRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/kubernetes-clusters", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCluster gets a Kubernetes cluster: GET /kubernetes-clusters/{cluster_id}.
// Scope: read.
func (s *KubernetesService) GetCluster(ctx context.Context, clusterID string, opts ...RequestOption) (*KubernetesCluster, error) {
	var out KubernetesCluster
	if err := s.c.Do(ctx, "GET", "/kubernetes-clusters/"+url.PathEscape(clusterID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateCluster changes a Kubernetes cluster: PATCH /kubernetes-clusters/{cluster_id}.
// Scope: write.
func (s *KubernetesService) UpdateCluster(ctx context.Context, clusterID string, req UpdateKubernetesClusterRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "PATCH", "/kubernetes-clusters/"+url.PathEscape(clusterID), nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteCluster deletes a Kubernetes cluster: DELETE /kubernetes-clusters/{cluster_id}.
// Scope: write.
func (s *KubernetesService) DeleteCluster(ctx context.Context, clusterID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "DELETE", "/kubernetes-clusters/"+url.PathEscape(clusterID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpgradeCluster upgrades a Kubernetes cluster: POST /kubernetes-clusters/{cluster_id}/upgrade.
// Scope: write.
func (s *KubernetesService) UpgradeCluster(ctx context.Context, clusterID string, req UpgradeKubernetesClusterRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/kubernetes-clusters/"+url.PathEscape(clusterID)+"/upgrade", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ResizePlan resizes worker node compute plan: POST /kubernetes-clusters/{cluster_id}/resize-plan.
// Scope: write.
func (s *KubernetesService) ResizePlan(ctx context.Context, clusterID string, req ResizeKubernetesClusterPlanRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/kubernetes-clusters/"+url.PathEscape(clusterID)+"/resize-plan", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpgradeNetwork upgrades Kubernetes cluster networking: POST /kubernetes-clusters/{cluster_id}/upgrade-network.
// Scope: write.
func (s *KubernetesService) UpgradeNetwork(ctx context.Context, clusterID string, req UpgradeKubernetesNetworkRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/kubernetes-clusters/"+url.PathEscape(clusterID)+"/upgrade-network", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// StartCluster starts a Kubernetes cluster: POST /kubernetes-clusters/{cluster_id}/start.
// Scope: write.
func (s *KubernetesService) StartCluster(ctx context.Context, clusterID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/kubernetes-clusters/"+url.PathEscape(clusterID)+"/start", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// StopCluster stops a Kubernetes cluster: POST /kubernetes-clusters/{cluster_id}/stop.
// Scope: write.
func (s *KubernetesService) StopCluster(ctx context.Context, clusterID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/kubernetes-clusters/"+url.PathEscape(clusterID)+"/stop", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetKubeconfig downloads a kubeconfig: POST /kubernetes-clusters/{cluster_id}/kubeconfig.
// Scope: write.
func (s *KubernetesService) GetKubeconfig(ctx context.Context, clusterID string, opts ...RequestOption) (*Kubeconfig, error) {
	var out Kubeconfig
	if err := s.c.Do(ctx, "POST", "/kubernetes-clusters/"+url.PathEscape(clusterID)+"/kubeconfig", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListNodes lists a cluster's nodes: GET /kubernetes-clusters/{cluster_id}/nodes.
// Scope: read.
func (s *KubernetesService) ListNodes(ctx context.Context, clusterID string, opts ...RequestOption) (*KubernetesClusterNodes, error) {
	var out KubernetesClusterNodes
	if err := s.c.Do(ctx, "GET", "/kubernetes-clusters/"+url.PathEscape(clusterID)+"/nodes", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetWorkloads gets a cluster's services and pods: GET /kubernetes-clusters/{cluster_id}/workloads.
// Scope: read.
func (s *KubernetesService) GetWorkloads(ctx context.Context, clusterID string, opts ...RequestOption) (*KubernetesWorkloads, error) {
	var out KubernetesWorkloads
	if err := s.c.Do(ctx, "GET", "/kubernetes-clusters/"+url.PathEscape(clusterID)+"/workloads", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ApplyManifest applies a manifest to a cluster: POST /kubernetes-clusters/{cluster_id}/apply-manifest.
// Scope: write.
func (s *KubernetesService) ApplyManifest(ctx context.Context, clusterID string, req ApplyKubernetesManifestRequest, opts ...RequestOption) (*KubernetesManifestApplied, error) {
	var out KubernetesManifestApplied
	if err := s.c.Do(ctx, "POST", "/kubernetes-clusters/"+url.PathEscape(clusterID)+"/apply-manifest", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetLiveManifests exports a cluster's workloads as YAML: GET /kubernetes-clusters/{cluster_id}/live-manifests.
// Scope: read.
func (s *KubernetesService) GetLiveManifests(ctx context.Context, clusterID string, opts ...RequestOption) (*KubernetesLiveManifests, error) {
	var out KubernetesLiveManifests
	if err := s.c.Do(ctx, "GET", "/kubernetes-clusters/"+url.PathEscape(clusterID)+"/live-manifests", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ExposeService exposes a Kubernetes service: POST /kubernetes-clusters/{cluster_id}/expose-service.
// Scope: write.
func (s *KubernetesService) ExposeService(ctx context.Context, clusterID string, req ExposeKubernetesServiceRequest, opts ...RequestOption) (*KubernetesServiceExposed, error) {
	var out KubernetesServiceExposed
	if err := s.c.Do(ctx, "POST", "/kubernetes-clusters/"+url.PathEscape(clusterID)+"/expose-service", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// SyncRouting syncs a cluster's public routing: POST /kubernetes-clusters/{cluster_id}/sync-routing.
// Scope: write.
func (s *KubernetesService) SyncRouting(ctx context.Context, clusterID string, opts ...RequestOption) (*KubernetesRoutingSynced, error) {
	var out KubernetesRoutingSynced
	if err := s.c.Do(ctx, "POST", "/kubernetes-clusters/"+url.PathEscape(clusterID)+"/sync-routing", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListNodePools lists node pools: GET /kubernetes-clusters/{cluster_id}/node-pools.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *KubernetesService) ListNodePools(ctx context.Context, clusterID string, opts ...RequestOption) ([]KubernetesNodePool, error) {
	return listAll[KubernetesNodePool](ctx, s.c, "/kubernetes-clusters/"+url.PathEscape(clusterID)+"/node-pools", nil, opts...)
}

// CreateNodePool creates a node pool: POST /kubernetes-clusters/{cluster_id}/node-pools.
// Scope: write.
func (s *KubernetesService) CreateNodePool(ctx context.Context, clusterID string, req CreateKubernetesNodePoolRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/kubernetes-clusters/"+url.PathEscape(clusterID)+"/node-pools", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ScaleNodePool scales a node pool: PATCH /kubernetes-clusters/{cluster_id}/node-pools/{pool_id}.
// Scope: write.
func (s *KubernetesService) ScaleNodePool(ctx context.Context, clusterID string, poolID string, req ScaleKubernetesNodePoolRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "PATCH", "/kubernetes-clusters/"+url.PathEscape(clusterID)+"/node-pools/"+url.PathEscape(poolID), nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteNodePool deletes a node pool: DELETE /kubernetes-clusters/{cluster_id}/node-pools/{pool_id}.
// Scope: write.
func (s *KubernetesService) DeleteNodePool(ctx context.Context, clusterID string, poolID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "DELETE", "/kubernetes-clusters/"+url.PathEscape(clusterID)+"/node-pools/"+url.PathEscape(poolID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListAddons lists cluster add-ons: GET /kubernetes-clusters/{cluster_id}/addons.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *KubernetesService) ListAddons(ctx context.Context, clusterID string, opts ...RequestOption) ([]KubernetesClusterAddon, error) {
	return listAll[KubernetesClusterAddon](ctx, s.c, "/kubernetes-clusters/"+url.PathEscape(clusterID)+"/addons", nil, opts...)
}

// ToggleAddon toggles a cluster add-on: POST /kubernetes-clusters/{cluster_id}/addons/{addon_name}/toggle.
// Scope: write.
func (s *KubernetesService) ToggleAddon(ctx context.Context, clusterID string, addonName string, req ToggleKubernetesAddonRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/kubernetes-clusters/"+url.PathEscape(clusterID)+"/addons/"+url.PathEscape(addonName)+"/toggle", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ConfigureEndpointAccess configures API server endpoint access: POST /kubernetes-clusters/{cluster_id}/endpoint-access.
// Scope: write.
func (s *KubernetesService) ConfigureEndpointAccess(ctx context.Context, clusterID string, req ConfigureKubernetesEndpointAccessRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/kubernetes-clusters/"+url.PathEscape(clusterID)+"/endpoint-access", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
