package pantech

import "time"

// Accepted is the API's accepted object.
// Returned with `202 Accepted` by every asynchronous write.
type Accepted struct {
	// The operation carrying out the change.
	OperationID string `json:"operation_id"`
	// The resource being created or changed, e.g. the new snapshot's id for a snapshot create.
	ResourceID string `json:"resource_id"`
	// The operation's status when the request was accepted, usually `submitting`.
	Status string `json:"status"`
}

// AttachVolumeRequest is the API's attach volume request object.
type AttachVolumeRequest struct {
	// A running instance in the volume's zone, e.g. `vm_…`.
	InstanceID string `json:"instance_id"`
}

// ChangeSecurityGroupRequest is the API's change security group request object.
type ChangeSecurityGroupRequest struct {
	// One of your security groups in the instance's zone, e.g. `sg_…`.
	SecurityGroupID string `json:"security_group_id"`
}

// CreateFirewallRuleRequest is the API's create firewall rule request object.
type CreateFirewallRuleRequest struct {
	// Evaluation order, 100 to 9999, unique on the subnet; lower numbers are evaluated first and the first match decides.
	Number int `json:"number"`
	// `ingress` (traffic into the subnet) or `egress` (traffic out of it). One of: ingress, egress.
	Direction *string `json:"direction,omitempty"`
	// `tcp`, `udp`, `icmp` or `all`. One of: tcp, udp, icmp, all.
	Protocol string `json:"protocol"`
	// First port, 1 to 65535.
	PortStart *int `json:"port_start,omitempty"`
	// Last port, 1 to 65535, at least `port_start`.
	PortEnd *int `json:"port_end,omitempty"`
	// ICMP type, for `icmp` only, e.g. `8` for echo request.
	ICMPType *int `json:"icmp_type,omitempty"`
	// ICMP code, for `icmp` only.
	ICMPCode *int `json:"icmp_code,omitempty"`
	// One IPv4 CIDR: the source for `ingress`, the destination for `egress`.
	CIDR string `json:"cidr"`
	// `allow` or `deny`. One of: allow, deny.
	Action *string `json:"action,omitempty"`
}

// CreateInstanceRequest is the API's create instance request object.
type CreateInstanceRequest struct {
	// 1 to 255 characters of lowercase letters a-z, digits 0-9 and hyphens (-), starting and ending with a letter or digit, unique among your instances.
	Name string `json:"name"`
	// The size, from `GET /public/v1/plans`.
	PlanSlug string `json:"plan_slug"`
	// The operating system, from `GET /public/v1/images`.
	ImageSlug string `json:"image_slug"`
	// Optional.
	Encrypted *bool `json:"encrypted,omitempty"`
	// An SSH key from `GET /public/v1/ssh-keys`, installed for the default user.
	SSHKeyID string `json:"ssh_key_id"`
	// A region from `GET /public/v1/regions`, e.g. `af-abj`.
	Region *string `json:"region,omitempty"`
	// Place the instance in this VPC subnet, with a private address only.
	SubnetID *string `json:"subnet_id,omitempty"`
	// The subnet's VPC.
	NetworkID *string `json:"network_id,omitempty"`
	// Standard instances only: one of your security groups.
	SecurityGroupID *string `json:"security_group_id,omitempty"`
	// Your own labels, as string keys and values.
	Tags *map[string]string `json:"tags,omitempty"`
	// Turn on Dynamic IP once the instance runs.
	DynamicIP *CreateInstanceDynamicIP `json:"address_rotation,omitempty"`
}

// CreateInstanceDynamicIP configures Dynamic IP on a new instance.
type CreateInstanceDynamicIP struct {
	// How often the rotating public address changes. One of: every_12_hours, daily, weekly, monthly.
	Schedule string `json:"schedule"`
	// `keep` (default) lets open connections finish on the old address; `move` moves them at once. One of: keep, move.
	Connections *string `json:"connections,omitempty"`
}

// PutDynamicIPRequest enables or updates Dynamic IP on an instance.
type PutDynamicIPRequest struct {
	// How often the rotating public address changes. One of: every_12_hours, daily, weekly, monthly.
	Schedule string `json:"schedule"`
	// `keep` (default) lets open connections finish on the old address; `move` moves them at once. One of: keep, move.
	Connections *string `json:"connections,omitempty"`
	// Allow rotating the address on demand. Requires on-demand rotation enabled for your organization.
	OnDemand *bool `json:"on_demand,omitempty"`
}

// DynamicIP is an instance's Dynamic IP configuration and status.
type DynamicIP struct {
	InstanceID string `json:"instance_id"`
	// How often the rotating public address changes. One of: every_12_hours, daily, weekly, monthly.
	Schedule string `json:"schedule"`
	// How open connections behave on rotation. One of: keep, move.
	Connections string `json:"connections"`
	// Whether on-demand rotation is allowed.
	OnDemand bool `json:"on_demand"`
	// Current Dynamic IP state. One of: enabling, active, degraded, disabling, off, failed.
	State string `json:"state"`
	// The rotating public address in use now; nil before the first one is ready or while degraded onto the instance's static address.
	CurrentAddress *string `json:"current_address"`
	// When the next scheduled rotation runs; nil when not scheduled.
	NextRotationAt *time.Time `json:"next_rotation_at"`
	// When the address last rotated; nil before the first rotation.
	LastRotatedAt *time.Time `json:"last_rotated_at"`
	// Percentage added to the instance's hourly price while Dynamic IP is on.
	SurchargePercent float64 `json:"surcharge_percent"`
	// On-demand rotations used in the last 24 hours.
	OnDemandUsed int `json:"on_demand_used"`
	// Maximum on-demand rotations allowed per 24 hours.
	OnDemandLimit int `json:"on_demand_limit"`
	// Why Dynamic IP is `failed` or `degraded`; nil otherwise.
	FailureReason *string `json:"failure_reason"`
	// Recent Dynamic IP changes, newest first.
	Events []DynamicIPEvent `json:"events"`
}

// DynamicIPEvent is one recorded change in an instance's Dynamic IP history.
type DynamicIPEvent struct {
	// What triggered the change. One of: enabled, scheduled, on_demand, failover, handover, disabled.
	Kind        string    `json:"kind"`
	FromAddress *string   `json:"from_address"`
	ToAddress   *string   `json:"to_address"`
	Degraded    bool      `json:"degraded"`
	OccurredAt  time.Time `json:"occurred_at"`
}

// CreateNetworkRequest is the API's create network request object.
type CreateNetworkRequest struct {
	// 1 to 255 characters, unique among your VPC networks.
	Name string `json:"name"`
	// The VPC's private IPv4 range, e.g. `10.0.0.0/16`.
	CIDR string `json:"cidr"`
	// A region from `GET /public/v1/regions` that offers VPC networks.
	Region *string `json:"region,omitempty"`
}

// CreateLoadBalancerRequest is the API's create load balancer request object.
type CreateLoadBalancerRequest struct {
	// 1 to 63 characters.
	Name string `json:"name"`
	// The zone to place the load balancer in.
	ZoneID *string `json:"zone_id,omitempty"`
	// `vpc_native` (a VPC zone: needs `public_ip_id` and `subnet_id`) or `appliance` (a standard zone: needs `zone_id`; a load balancer VM with its own public address, billed as that VM). One of: vpc_native, appliance.
	Kind *string `json:"kind,omitempty"`
	// An active public IP with `purpose` `load_balancer` (required for `vpc_native`).
	PublicIPID *string `json:"public_ip_id,omitempty"`
	// A subnet of the public IP's private network.
	SubnetID *string `json:"subnet_id,omitempty"`
	// `roundrobin` takes turns, `leastconn` picks the target with the fewest connections, `source` keeps a client on one target. One of: roundrobin, leastconn, source.
	Algorithm *string `json:"algorithm,omitempty"`
	// The TCP port clients connect to, 1 to 65535.
	PublicPort int `json:"public_port"`
	// The port on each target.
	PrivatePort *int `json:"private_port,omitempty"`
	// The listener protocol, `tcp` when left out. `tcp` passes traffic
	// through unchanged; `https` ends TLS on the load balancer with
	// CertificateID and speaks plain HTTP to the instances (`appliance`
	// only: `vpc_native` answers LOAD_BALANCER_TLS_NOT_SUPPORTED). One of: tcp, https.
	Protocol *string `json:"protocol,omitempty"`
	// `https` only, and required there: an issued certificate (`cert_…`).
	// While used, it cannot be deleted or revoked; renewals reach the load
	// balancer by themselves.
	CertificateID *string `json:"certificate_id,omitempty"`
	// `https` only; `modern` when left out. `modern`: TLS 1.2 and 1.3 with
	// forward-secret AEAD ciphers. `tls13`: TLS 1.3 only. One of: modern, tls13.
	TLSPolicy *string `json:"tls_policy,omitempty"`
	// `https` only: also answer port 80, redirecting to HTTPS with a 301.
	RedirectHTTP *bool `json:"redirect_http,omitempty"`
	// Source addresses allowed to connect, as IPv4 CIDRs (at most 20).
	CIDRList *[]string `json:"cidr_list,omitempty"`
	// Instances to spread traffic across (at most 50).
	InstanceIDs *[]string `json:"instance_ids,omitempty"`
}

// UpdateLoadBalancerRequest is the API's update load balancer request object.
type UpdateLoadBalancerRequest struct {
	// A new name, 1 to 63 characters.
	Name *string `json:"name,omitempty"`
	// One of: roundrobin, leastconn, source.
	Algorithm *string `json:"algorithm,omitempty"`
	// The whole new set of targets (at most 50).
	InstanceIDs *[]string `json:"instance_ids,omitempty"`
	// `https` only: serve this certificate instead.
	CertificateID *string `json:"certificate_id,omitempty"`
	// `https` only. One of: modern, tls13.
	TLSPolicy *string `json:"tls_policy,omitempty"`
	// `https` only: turn the port 80 redirect on or off.
	RedirectHTTP *bool `json:"redirect_http,omitempty"`
}

// LoadBalancerMember is the API's load balancer member object.
type LoadBalancerMember struct {
	InstanceID string `json:"instance_id"`
	// The instance's name; null when it is gone.
	InstanceName *string `json:"instance_name"`
	// `present` while it should receive traffic; `deleted` once you removed it and that has not finished. One of: present, deleted.
	DesiredState string `json:"desired_state"`
	// `active` once it receives traffic. One of: pending, active, deleting, deleted, failed, unknown.
	ObservedState string `json:"observed_state"`
}

// LoadBalancer is the API's load balancer object.
type LoadBalancer struct {
	// e.g. `lb_...`.
	ID     string `json:"id"`
	Name   string `json:"name"`
	ZoneID string `json:"zone_id"`
	// One of: vpc_native, appliance.
	Kind       string `json:"kind"`
	PublicIPID string `json:"public_ip_id"`
	// The address clients connect to; null until it is known.
	PublicIPAddress *string `json:"public_ip_address"`
	NetworkID       *string `json:"network_id"`
	SubnetID        string  `json:"subnet_id"`
	// One of: tcp, http, https.
	Protocol string `json:"protocol"`
	// One of: roundrobin, leastconn, source.
	Algorithm   string `json:"algorithm"`
	PublicPort  int    `json:"public_port"`
	PrivatePort int    `json:"private_port"`
	// The certificate an `https` listener serves; nil for `tcp`.
	CertificateID *string `json:"certificate_id"`
	// The SHA-256 of the leaf the load balancer serves now, to compare with
	// the certificate's FingerprintSHA256 after a renewal; nil until it
	// first serves one.
	CertificateFingerprintSHA256 *string `json:"certificate_fingerprint_sha256"`
	// `modern` or `tls13` for `https`; nil for `tcp`.
	TLSPolicy *string `json:"tls_policy"`
	// Whether port 80 redirects to HTTPS.
	RedirectHTTP bool `json:"redirect_http"`
	// Source addresses allowed to connect; empty is anywhere.
	CIDRList []string             `json:"cidr_list"`
	Members  []LoadBalancerMember `json:"members"`
	// What you asked for: `present`, or `deleted` once you deleted it. One of: present, deleted.
	DesiredState string `json:"desired_state"`
	// `active` once it serves traffic. One of: pending, provisioning, active, deleting, deleted, failed, unknown.
	ObservedState string `json:"observed_state"`
	// False while a change you made has not reached the load balancer yet.
	InSync    bool       `json:"in_sync"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// KubernetesVersion is the API's kubernetes version object.
type KubernetesVersion struct {
	// e.g. `k8sv_...`.
	ID      string `json:"id"`
	ZoneID  string `json:"zone_id"`
	Version string `json:"version"`
	// `withdrawn` versions take no new clusters. One of: available, withdrawn.
	Status string `json:"status"`
	// The fewest vCPUs a node plan needs for this version.
	MinCPU int `json:"min_cpu"`
	// The least memory, in MB, a node plan needs for this version.
	MinMemoryMB int `json:"min_memory_mb"`
}

// KubernetesVersionList is the API's kubernetes version list object.
type KubernetesVersionList struct {
	Data []KubernetesVersion `json:"data"`
	// Zones where a cluster can have 3 control nodes.
	HAZoneIDs []string `json:"ha_zone_ids"`
	// Zones that build Kubernetes nodes only at fixed sizes, one entry per
	// size, so a zone may appear more than once. There, a cluster, a node
	// plan resize and a node pool take only a plan of exactly one entry's
	// VCPU and MemoryMB; anything else is 422 KUBERNETES_PLAN_NOT_IN_ZONE.
	// A zone not listed takes any plan that meets the version's minimum.
	FixedNodeSizes []KubernetesFixedNodeSize `json:"fixed_node_sizes"`
}

// KubernetesFixedNodeSize is one node size a zone builds, from
// KubernetesVersionList.FixedNodeSizes.
type KubernetesFixedNodeSize struct {
	ZoneID string `json:"zone_id"`
	// The vCPUs of one node size the zone builds; 0 when the zone takes no
	// plan now.
	VCPU int `json:"vcpu"`
	// The memory, in MB, of that node size; 0 when the zone takes no plan now.
	MemoryMB int `json:"memory_mb"`
}

// CreateKubernetesClusterRequest is the API's create kubernetes cluster request object.
type CreateKubernetesClusterRequest struct {
	// 1 to 40 lowercase letters, digits and hyphens, starting with a letter.
	Name   string `json:"name"`
	ZoneID string `json:"zone_id"`
	// From the Kubernetes versions list, in `zone_id`.
	KubernetesVersionID string `json:"kubernetes_version_id"`
	// Required in a VPC zone: the subnet the nodes go in.
	SubnetID *string `json:"subnet_id,omitempty"`
	// The plan slug every node uses, from the plans list.
	NodePlan string `json:"node_plan"`
	// 1, or 3 for a highly available control plane where the zone offers it. One of: 1, 3.
	ControlNodes *int `json:"control_nodes,omitempty"`
	// Worker nodes, 1 to 10.
	Workers     *int                   `json:"workers,omitempty"`
	Autoscaling *KubernetesAutoscaling `json:"autoscaling,omitempty"`
	// VPC zone only: the IPv4 CIDRs allowed to reach the API server (at most 20).
	APIAllowedCIDRs *[]string `json:"api_allowed_cidrs,omitempty"`
	// Cluster architecture. One of: appliance.
	Kind *string `json:"kind,omitempty"`
	// CNI backend plugin (`flannel`, `wireguard`, or `calico`).
	CNIPlugin *string `json:"cni_plugin,omitempty"`
	// IPv4 CIDR range allocated for cluster pods.
	PodCIDR *string `json:"pod_cidr,omitempty"`
	// IPv4 CIDR range allocated for cluster services.
	ServiceCIDR *string `json:"service_cidr,omitempty"`
	// Who can reach the API server. One of: private, public, dual.
	EndpointAccess *string `json:"endpoint_access,omitempty"`
}

// KubernetesAutoscaling is the API's kubernetes autoscaling object.
// The platform's autoscaler: about once a minute it adds a worker when a pod that a worker could run is waiting for room, and removes the newest worker once it has run nothing but DaemonSet pods for 10 minutes, one step at a time and always within min_workers and max_workers.
type KubernetesAutoscaling struct {
	// The cluster autoscaler adds and removes workers within the range as pods need them.
	Enabled bool `json:"enabled"`
	// 1 to 10 when enabled; 0 when not.
	MinWorkers int `json:"min_workers"`
	// min_workers to 10 when enabled; 0 when not.
	MaxWorkers int `json:"max_workers"`
}

// UpdateKubernetesClusterRequest is the API's update kubernetes cluster request object.
// Send any of the fields; those left out stay as they are.
type UpdateKubernetesClusterRequest struct {
	// A fixed worker count, 1 to 10.
	Workers     *int                   `json:"workers,omitempty"`
	Autoscaling *KubernetesAutoscaling `json:"autoscaling,omitempty"`
	// VPC zone only: the whole new list of IPv4 CIDRs allowed to reach the API server.
	APIAllowedCIDRs *[]string `json:"api_allowed_cidrs,omitempty"`
}

// UpgradeKubernetesClusterRequest is the API's upgrade kubernetes cluster request object.
type UpgradeKubernetesClusterRequest struct {
	// One of the cluster's `available_upgrades`.
	KubernetesVersionID string `json:"kubernetes_version_id"`
}

// ResizeKubernetesClusterPlanRequest is the API's resize kubernetes cluster plan request object.
type ResizeKubernetesClusterPlanRequest struct {
	// The new node plan slug, from the plans list.
	NodePlan string `json:"node_plan"`
}

// UpgradeKubernetesNetworkRequest is the API's upgrade kubernetes network request object.
type UpgradeKubernetesNetworkRequest struct {
	// The new CNI plugin (`flannel`, `wireguard`, or `calico`).
	CNIPlugin *string `json:"cni_plugin,omitempty"`
	// The new IPv4 CIDR range allocated for cluster pods.
	PodCIDR *string `json:"pod_cidr,omitempty"`
	// The new IPv4 CIDR range allocated for cluster services.
	ServiceCIDR *string `json:"service_cidr,omitempty"`
	// The new allowed IPv4 CIDRs reaching the API server.
	APIAllowedCIDRs *[]string `json:"api_allowed_cidrs,omitempty"`
}

// KubernetesNode is the API's kubernetes node object.
type KubernetesNode struct {
	VCPU     int `json:"vcpu"`
	MemoryMB int `json:"memory_mb"`
	DiskGB   int `json:"disk_gb"`
}

// KubernetesCluster is the API's kubernetes cluster object.
type KubernetesCluster struct {
	// e.g. `k8s_...`.
	ID                  string `json:"id"`
	Name                string `json:"name"`
	ZoneID              string `json:"zone_id"`
	KubernetesVersionID string `json:"kubernetes_version_id"`
	KubernetesVersion   string `json:"kubernetes_version"`
	// The private network (VPC) the nodes are in; null in a standard zone.
	NetworkID    *string        `json:"network_id"`
	SubnetID     *string        `json:"subnet_id"`
	NodePlanID   string         `json:"node_plan_id"`
	Node         KubernetesNode `json:"node"`
	ControlNodes int            `json:"control_nodes"`
	Workers      int            `json:"workers"`
	// control_nodes plus workers: the VMs billed.
	Nodes int `json:"nodes"`
	// One of: running, stopped, deleted.
	DesiredState string `json:"desired_state"`
	// `updating` while a scale or upgrade runs. One of: pending, provisioning, running, updating, stopping, stopped, retired, deleted, failed, unknown.
	ObservedState string `json:"observed_state"`
	// False while a scale or upgrade has not finished.
	InSync        bool    `json:"in_sync"`
	FailureCode   *string `json:"failure_code"`
	FailureReason *string `json:"failure_reason"`
	// Versions this cluster can be upgraded to now.
	AvailableUpgrades []KubernetesVersion   `json:"available_upgrades"`
	Autoscaling       KubernetesAutoscaling `json:"autoscaling"`
	// Sources allowed to reach the API server; empty is any address.
	APIAllowedCIDRs []string `json:"api_allowed_cidrs"`
	// The API server's URL, once the cluster has one: the public address when the endpoint is public or dual and one is routed, else the private one.
	Endpoint *string `json:"endpoint"`
	// The public address the API server is reachable on; null until one is routed, and for a private endpoint.
	PublicIP *string `json:"public_ip"`
	// The public port the API server is reachable on; null until one is routed, and for a private endpoint.
	PublicPort *int `json:"public_port"`
	// What Kubernetes last reported of the nodes while the cluster runs; null while it is not running. One of: healthy, degraded, unknown.
	NodeHealth *string `json:"node_health"`
	// Why NodeHealth is `degraded` when the cause is the cluster itself rather than a node; null otherwise.
	NodeHealthReason *string `json:"node_health_reason"`
	// Nodes Kubernetes reports Ready, as of NodesCheckedAt; null when NodeHealth is unknown or null.
	ReadyNodes *int `json:"ready_nodes"`
	// When the nodes were last read; null if never.
	NodesCheckedAt *time.Time `json:"nodes_checked_at"`
	// Storage your workloads' persistent volume claims use, billed with the nodes' disks.
	VolumeStorageGB int `json:"volume_storage_gb"`
	// Cluster architecture kind, e.g. `appliance`.
	Kind string `json:"kind"`
	// CNI backend plugin (`flannel`, `wireguard`, or `calico`).
	CNIPlugin string `json:"cni_plugin"`
	// IPv4 CIDR range allocated for cluster pods.
	PodCIDR string `json:"pod_cidr"`
	// IPv4 CIDR range allocated for cluster services.
	ServiceCIDR string `json:"service_cidr"`
	// API server endpoint network access mode (`public`, `private`, or `dual`). One of: public, private, dual.
	EndpointAccess string     `json:"endpoint_access"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

// Kubeconfig is the API's kubeconfig object.
type Kubeconfig struct {
	// A kubeconfig file (YAML).
	Kubeconfig string `json:"kubeconfig"`
}

// KubernetesNodePoolTaint is the API's kubernetes node pool taint object.
type KubernetesNodePoolTaint struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	// One of: NoSchedule, PreferNoSchedule, NoExecute.
	Effect string `json:"effect"`
}

// KubernetesNodePool is the API's kubernetes node pool object.
type KubernetesNodePool struct {
	// `provisioning` (workers being built), `ready` (every worker has joined), `updating` (a scale in progress), `deleting`, or `unavailable` (a pool made while node pools were a preview, with no node size: delete it and create it again). A pool whose workers cannot be built stays `provisioning` or `updating` and is retried; FailureCode and FailureReason say why. One of: provisioning, ready, updating, deleting, unavailable.
	Status string `json:"status"`
	// The pool's workers that have joined the cluster.
	ObservedWorkers int `json:"observed_workers"`
	// Why the pool's workers last failed to build, e.g. `PROVISIONING_CAPACITY_UNAVAILABLE` (the region has no room for them now) or `PROVISIONING_LIMIT_EXCEEDED`. Nil while nothing has failed; cleared once the pool is `ready` or is scaled.
	FailureCode *string `json:"failure_code"`
	// A sentence saying what failed and what to do. Nil when FailureCode is.
	FailureReason *string `json:"failure_reason"`
	// e.g. `k8spool_...`.
	ID                 string                    `json:"id"`
	ClusterID          string                    `json:"cluster_id"`
	Name               string                    `json:"name"`
	NodePlanID         string                    `json:"node_plan_id"`
	Workers            int                       `json:"workers"`
	AutoscalingEnabled bool                      `json:"autoscaling_enabled"`
	MinWorkers         int                       `json:"min_workers"`
	MaxWorkers         int                       `json:"max_workers"`
	Labels             map[string]string         `json:"labels"`
	Taints             []KubernetesNodePoolTaint `json:"taints"`
	CreatedAt          *time.Time                `json:"created_at"`
	UpdatedAt          *time.Time                `json:"updated_at"`
}

// CreateKubernetesNodePoolRequest is the API's create kubernetes node pool request object.
type CreateKubernetesNodePoolRequest struct {
	// Unique name for the node pool.
	Name string `json:"name"`
	// Compute plan ID or slug for the worker VMs.
	PlanID string `json:"plan_id"`
	// Worker nodes, 1 to 10.
	Workers            *int  `json:"workers,omitempty"`
	AutoscalingEnabled *bool `json:"autoscaling_enabled,omitempty"`
	// With autoscaling_enabled: fewest workers, 1 to 10.
	MinWorkers *int `json:"min_workers,omitempty"`
	// With autoscaling_enabled: most workers, min_workers to 10.
	MaxWorkers *int                       `json:"max_workers,omitempty"`
	Labels     *map[string]string         `json:"labels,omitempty"`
	Taints     *[]KubernetesNodePoolTaint `json:"taints,omitempty"`
}

// ScaleKubernetesNodePoolRequest is the API's scale kubernetes node pool request object.
type ScaleKubernetesNodePoolRequest struct {
	// Worker nodes, 1 to 10.
	Workers            *int  `json:"workers,omitempty"`
	AutoscalingEnabled *bool `json:"autoscaling_enabled,omitempty"`
	MinWorkers         *int  `json:"min_workers,omitempty"`
	MaxWorkers         *int  `json:"max_workers,omitempty"`
}

// KubernetesClusterAddon is the API's kubernetes cluster addon object.
type KubernetesClusterAddon struct {
	// e.g. `k8saddon_...`.
	ID        string `json:"id"`
	ClusterID string `json:"cluster_id"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Status    string `json:"status"`
	Enabled   bool   `json:"enabled"`
	// JSON-encoded add-on configuration.
	Configuration string     `json:"configuration"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

// ToggleKubernetesAddonRequest is the API's toggle kubernetes addon request object.
type ToggleKubernetesAddonRequest struct {
	// Whether the add-on should be enabled or disabled.
	Enabled bool `json:"enabled"`
}

// ConfigureKubernetesEndpointAccessRequest is the API's configure kubernetes endpoint access request object.
type ConfigureKubernetesEndpointAccessRequest struct {
	// Endpoint access mode. One of: public, private, dual.
	EndpointAccess string `json:"endpoint_access"`
	// IPv4 CIDRs allowed to connect to the cluster API.
	APIAllowedCIDRs *[]string `json:"api_allowed_cidrs,omitempty"`
}

// CreatePortForwardingRuleRequest is the API's create port forwarding rule request object.
type CreatePortForwardingRuleRequest struct {
	// An instance in the public IP's VPC.
	InstanceID string `json:"instance_id"`
	// `tcp` or `udp`. One of: tcp, udp.
	Protocol string `json:"protocol"`
	// First public port, 1 to 65535.
	PublicPortStart int `json:"public_port_start"`
	// Last public port.
	PublicPortEnd *int `json:"public_port_end,omitempty"`
	// First port on the instance.
	PrivatePortStart *int `json:"private_port_start,omitempty"`
	// Last port on the instance; the private range must be as wide as the public one.
	PrivatePortEnd *int `json:"private_port_end,omitempty"`
}

// CreatePublicIPRequest is the API's create public iprequest object.
type CreatePublicIPRequest struct {
	// The VPC network to allocate the address in.
	NetworkID string `json:"network_id"`
	// `static_nat` maps every port to one instance (`instance_id`). One of: static_nat, port_forwarding, load_balancer.
	Purpose *string `json:"purpose,omitempty"`
	// For `static_nat` only: an instance in `network_id` without a public IP.
	InstanceID *string `json:"instance_id,omitempty"`
}

// CreateSSHKeyRequest is the API's create sshkey request object.
type CreateSSHKeyRequest struct {
	// 1 to 255 characters, unique among your SSH keys.
	Name string `json:"name"`
	// One OpenSSH public key line, e.g. `ssh-ed25519 AAAA… deploy@ci`.
	PublicKey *string `json:"public_key,omitempty"`
}

// CreateSecurityGroupRequest is the API's create security group request object.
type CreateSecurityGroupRequest struct {
	// 1 to 255 characters, unique among your live security groups.
	Name string `json:"name"`
	// The group's rules; at least one.
	Rules *[]SecurityGroupRule `json:"rules,omitempty"`
}

// ResizeDatabaseStorageRequest is the API's resize database storage request object.
type ResizeDatabaseStorageRequest struct {
	// The new data disk size in GB: larger than the current `data_volume_size_gb`, at most the zone's `max_gb` (2000 today) and a multiple of its `step_gb` (10 today); see `storage` in `GET /public/v1/database-engines`.
	StorageGB int `json:"storage_gb"`
}

// CreateSnapshotRequest is the API's create snapshot request object.
type CreateSnapshotRequest struct {
	// 1 to 255 characters, unique among this resource's snapshots, e.g. `pre-upgrade`.
	Name string `json:"name"`
}

// CreateSubnetRequest is the API's create subnet request object.
type CreateSubnetRequest struct {
	// 1 to 255 characters.
	Name string `json:"name"`
	// An IPv4 range inside the VPC's CIDR that does not overlap its other subnets, e.g. `10.0.1.0/24`.
	CIDR string `json:"cidr"`
}

// CreateVolumeRequest is the API's create volume request object.
type CreateVolumeRequest struct {
	// 1 to 255 characters.
	Name string `json:"name"`
	// The volume type, from `GET /public/v1/disk-offerings`.
	DiskOfferingSlug string `json:"disk_offering_slug"`
	// Size in GB.
	SizeGB *int `json:"size_gb,omitempty"`
	// A region from `GET /public/v1/regions`.
	Region *string `json:"region,omitempty"`
	// Attach the new volume at once to this running instance in the same zone.
	InstanceID *string `json:"instance_id,omitempty"`
	// With `instance_id`: an absolute path to mount the volume at, e.g. `/data`.
	MountPoint *string `json:"mount_point,omitempty"`
	// Optional.
	Encrypted *bool `json:"encrypted,omitempty"`
}

// ChangeDatabasePasswordRequest is the API's change database password request object.
type ChangeDatabasePasswordRequest struct {
	// The new admin password: 16 to 128 printable ASCII characters, with no spaces, single quotes (`'`), double quotes (`"`) or backslashes (`\`).
	Password string `json:"password"`
}

// CreateDatabaseRequest is the API's create database request object.
type CreateDatabaseRequest struct {
	// 1 to 255 characters, unique among your databases.
	Name string `json:"name"`
	// `postgresql`, `mysql` or `mariadb`; see `GET /public/v1/database-engines`. One of: postgresql, mysql, mariadb.
	Engine string `json:"engine"`
	// A version line offered for `engine` in `zone_id`, from `GET /public/v1/database-engines`, e.g. `18` for PostgreSQL or `8.4` for MySQL.
	Version string `json:"version"`
	// The size, from `GET /public/v1/plans`: its vCPU and memory, and its `disk_gb` as the smallest data disk.
	PlanSlug string `json:"plan_slug"`
	// Where to create it: a zone from the version's `zones`, e.g. `af-abj-1` (standard) or `af-abj-2` (VPC).
	ZoneID string `json:"zone_id"`
	// Required in a VPC zone (`af-abj-2`): one of your subnets there, which the database joins.
	SubnetID *string `json:"subnet_id,omitempty"`
	// The initial allow-list: IPv4 CIDRs with a prefix of `/8` or longer, never `0.0.0.0/0`, at most 50, no range twice.
	AccessRules *[]string `json:"access_rules,omitempty"`
	// The admin login's name: 3 to 32 characters, lowercase letters, digits and underscores, starting with a letter.
	AdminUsername *string `json:"admin_username,omitempty"`
	// Optional admin password: 16 to 128 printable ASCII characters, with no spaces, single quotes (`'`), double quotes (`"`) or backslashes (`\`).
	Password *string `json:"password,omitempty"`
	// Data disk size in GB.
	StorageGB *int `json:"storage_gb,omitempty"`
}

// Database is the API's database object.
type Database struct {
	// The admin login to connect as: the one chosen at create, or `dbadmin`.
	AdminUsername string `json:"admin_username"`
	// Your security groups attached to the database, in the order attached (at most 5).
	SecurityGroupIDs []string `json:"security_group_ids"`
	// The allow-list as the database applies it: its own `access_rules` plus the ranges its security groups add, each range once under its first source.
	EffectiveAccessRules []DatabaseEffectiveAccessRule `json:"effective_access_rules"`
	// Rules of the attached security groups that add nothing to the database, and why.
	IgnoredSecurityGroupRules []DatabaseIgnoredSecurityGroupRule `json:"ignored_security_group_rules"`
	// The database's id, `db_` followed by 26 characters.
	ID string `json:"id"`
	// Its name, unique among your databases.
	Name string `json:"name"`
	// `postgresql`, `mysql` or `mariadb`. One of: postgresql, mysql, mariadb.
	Engine string `json:"engine"`
	// The engine's version line, e.g. `18` or `8.4`.
	Version string `json:"version"`
	// The TCP port to connect to: 5432 for PostgreSQL, 3306 for MySQL and MariaDB.
	Port int `json:"port"`
	// The plan setting its vCPU and memory.
	PlanID string `json:"plan_id"`
	// The data disk's size in GB, as in place and billed.
	DataVolumeSizeGB int `json:"data_volume_size_gb"`
	// The size a storage resize in flight is growing the data disk to; null when none.
	PendingDataVolumeSizeGB *int `json:"pending_data_volume_size_gb"`
	// The zone it runs in, e.g. `af-abj-2`.
	ZoneID string `json:"zone_id"`
	// Its VPC subnet in a VPC zone; null in a standard zone.
	SubnetID *string `json:"subnet_id"`
	// The name to connect to, e.g. `db-7k2q9x4m1a.af-abj-2.db.pantechdynamics.com`.
	Hostname *string `json:"hostname"`
	// The database's private address, never public.
	PrivateIP *string `json:"private_ip"`
	// What you asked for: `running`, `stopped` or `deleted`. One of: running, stopped, deleted.
	DesiredState string `json:"desired_state"`
	// What the platform last observed. One of: pending, provisioning, running, stopping, stopped, deleting, deleted, failed, unknown.
	ObservedState string `json:"observed_state"`
	// Increases by one with every change you request.
	Generation int `json:"generation"`
	// The `generation` last applied: equal to `generation` once the last change (access rules, security groups, password) is in effect.
	ObservedGeneration int `json:"observed_generation"`
	// The database's own allow-list, as set at create or with *Replace access rules*.
	AccessRules []DatabaseAccessRule `json:"access_rules"`
	// A stable code for why the last change failed; null when nothing failed.
	FailureCode *string `json:"failure_code"`
	// When it was created (RFC 3339, UTC).
	CreatedAt *time.Time `json:"created_at"`
	// When it last changed (RFC 3339, UTC).
	UpdatedAt *time.Time `json:"updated_at"`
}

// DatabaseAccessRule is the API's database access rule object.
type DatabaseAccessRule struct {
	// An allowed IPv4 range, in canonical form.
	CIDR string `json:"cidr"`
	// Always `tcp`. One of: tcp.
	Protocol string `json:"protocol"`
	// The engine's port: 5432 or 3306.
	Port int `json:"port"`
}

// DatabaseEffectiveAccessRule is the API's database effective access rule object.
type DatabaseEffectiveAccessRule struct {
	// An allowed IPv4 range, in canonical form.
	CIDR string `json:"cidr"`
	// Always `tcp`. One of: tcp.
	Protocol string `json:"protocol"`
	// The engine's port: 5432 or 3306.
	Port int `json:"port"`
	// Where the range comes from: `access_rule`, or `security_group:` followed by the group's id (e.g. `security_group:sg_…`).
	Source string `json:"source"`
}

// DatabaseIgnoredSecurityGroupRule is the API's database ignored security group rule object.
type DatabaseIgnoredSecurityGroupRule struct {
	// The attached group the rule belongs to.
	SecurityGroupID string `json:"security_group_id"`
	// The rule's direction. One of: ingress, egress.
	Direction string `json:"direction"`
	// The rule's protocol. One of: tcp, udp, icmp, all.
	Protocol string `json:"protocol"`
	// `80`, `3000-4000`, or null for all ports.
	PortRange *string `json:"port_range"`
	// The rule's CIDR.
	CIDR string `json:"cidr"`
	// Why it adds nothing: `egress_rule`, `protocol_not_tcp` (udp or icmp), `port_not_covered` (its ports exclude the engine's), `cidr_too_wide` (wider than /8, including `0.0.0.0/0`) or `cidr_not_ipv4`. One of: egress_rule, protocol_not_tcp, port_not_covered, cidr_too_wide, cidr_not_ipv4.
	Reason string `json:"reason"`
}

// SetDatabaseSecurityGroupsRequest is the API's set database security groups request object.
type SetDatabaseSecurityGroupsRequest struct {
	// The complete set to attach: at most 5 distinct ids of your own security groups in the database's zone, e.g. `sg_…`.
	SecurityGroupIDs []string `json:"security_group_ids"`
}

// DatabaseAccessRuleRequest is the API's database access rule request object.
type DatabaseAccessRuleRequest struct {
	// One IPv4 CIDR to allow, with a prefix of `/8` to `/32`; never `0.0.0.0/0`.
	CIDR string `json:"cidr"`
}

// DatabaseEngineVersion is the API's database engine version object.
type DatabaseEngineVersion struct {
	Version string   `json:"version"`
	EolDate *string  `json:"eol_date"`
	Zones   []string `json:"zones"`
}

// DatabaseEngineStorage is the API's database engine storage object.
type DatabaseEngineStorage struct {
	// The zone, e.g. `af-abj-2`.
	ZoneID string `json:"zone_id"`
	// The zone's minimum data disk in GB.
	MinGB int `json:"min_gb"`
	// True when `min_gb` is 0, so the chosen plan's `disk_gb` is the minimum.
	MinIsPlanDisk bool `json:"min_is_plan_disk"`
	// The largest data disk in GB.
	MaxGB int `json:"max_gb"`
	// Above the base size (the minimum), `storage_gb` must be a multiple of this.
	StepGB int `json:"step_gb"`
	// Price of one GB of data disk for a 730-hour month, in minor units of `currency`; null when it could not be priced.
	PricePerGBMonthMinor *int64 `json:"price_per_gb_month_minor"`
	// `NGN` or `USD`; null when not priced.
	Currency *string `json:"currency"`
}

// DatabaseEngine is the API's database engine object.
type DatabaseEngine struct {
	// The value `engine` takes in *Create a database*: `postgresql`, `mysql` or `mariadb`.
	Engine string `json:"engine"`
	// The engine's name for display.
	DisplayName string `json:"display_name"`
	// The port its databases listen on: 5432 or 3306.
	Port int `json:"port"`
	// Version lines offered for new databases, long-term-support lines only.
	Versions []DatabaseEngineVersion `json:"versions"`
	// The data disk sizes a new database of this engine may have, per zone: use them for `storage_gb` in *Create a database* and *Grow a database's storage*.
	Storage []DatabaseEngineStorage `json:"storage"`
}

// DatabaseOrder is the API's database order object.
// How a new database is paid for and provisioned.
type DatabaseOrder struct {
	// The admin login the database is created with.
	AdminUsername string `json:"admin_username"`
	// The order's id, `ord_…`.
	ID string `json:"id"`
	// The same value as `id`.
	OrderID string `json:"order_id"`
	// The card payment attempt; null when the order was paid from credit or not charged yet.
	AttemptID *string `json:"attempt_id"`
	// The id the database has, or will have once provisioned.
	DatabaseID string `json:"database_id"`
	// `awaiting_payment`, then `paid`, `provisioning` and `provisioned` (final: the database exists). One of: awaiting_payment, paid, provisioning, provisioned, payment_failed, failed.
	Status string `json:"status"`
	// The database's create operation, set once provisioning starts; null before.
	OperationID *string `json:"operation_id"`
	// Why the order failed, when it did: e.g. `insufficient_credit`, `card_declined`, `payment_expired` (not paid within an hour), `database_name_taken`, `provisioning_failed`.
	FailureCode *string `json:"failure_code"`
	// The initial payment in minor units of `currency` (kobo or cents).
	AmountMinor int64 `json:"amount_minor"`
	// `NGN` or `USD`: your organization's billing currency. One of: NGN, USD.
	Currency string `json:"currency"`
	// When it was created (RFC 3339, UTC).
	CreatedAt *time.Time `json:"created_at"`
	// When it last changed (RFC 3339, UTC).
	UpdatedAt *time.Time `json:"updated_at"`
}

// DatabaseOrderAccepted is the API's database order accepted object.
// Returned with `202 Accepted` by `POST /public/v1/databases`, with `Cache-Control: no-store`: the order that pays for and then provisions the database.
type DatabaseOrderAccepted struct {
	// The admin login the database is created with: the one you sent, or `dbadmin`.
	AdminUsername string `json:"admin_username"`
	// Follow it at `GET /public/v1/database-orders/{order_id}` (also the `Location` header).
	OrderID string `json:"order_id"`
	// The id the database will have.
	DatabaseID string `json:"database_id"`
	// Same as `database_id`.
	ResourceID string `json:"resource_id"`
	// The order's status now, usually `awaiting_payment`. One of: awaiting_payment, paid, provisioning, provisioned, payment_failed, failed.
	Status string `json:"status"`
	// Set once provisioning starts.
	OperationID *string `json:"operation_id"`
	// Null for a new order; on a retry, the order's `failure_code` if it has failed since.
	FailureCode *string `json:"failure_code"`
	// The initial payment, in minor units (kobo or cents): the price of an instance of the same plan with no public address.
	AmountMinor int64 `json:"amount_minor"`
	// `NGN` or `USD`: your organization's billing currency. One of: NGN, USD.
	Currency string `json:"currency"`
	// True only when the password was generated for you and this request created the order.
	PasswordReturned bool `json:"password_returned"`
	// The generated admin password (32 letters and digits), in this response only.
	Password *string `json:"password"`
}

// DatabasePasswordAccepted is the API's database password accepted object.
// Returned with `202 Accepted` by `reset-password`, with `Cache-Control: no-store`.
type DatabasePasswordAccepted struct {
	// The operation carrying out the change.
	OperationID string `json:"operation_id"`
	// The database.
	ResourceID string `json:"resource_id"`
	// The operation's status when the request was accepted, usually `submitting`.
	Status string `json:"status"`
	// False on a retry with the same Idempotency-Key; `password` is then null.
	PasswordReturned bool `json:"password_returned"`
	// The generated admin password (32 letters and digits), in this response only.
	Password *string `json:"password"`
}

// DiskOffering is the API's disk offering object.
type DiskOffering struct {
	// Use this in `disk_offering_slug`.
	Slug string `json:"slug"`
	// Its name for display.
	Name string `json:"name"`
	// Fixed size; null for a custom-size offering.
	SizeGB *int `json:"size_gb"`
	// When true, send `size_gb` on create; the price is then per GB.
	CustomSize bool `json:"custom_size"`
	// The storage class, e.g. `shared` or `local`; null when not stated.
	StorageType *string `json:"storage_type"`
	// `NGN` or `USD`: your organization's currency.
	Currency string `json:"currency"`
	// Price per hour in minor units of `currency`: for the whole volume, or per GB for a custom-size offering.
	HourlyPriceMinor int64 `json:"hourly_price_minor"`
}

// Failure is the API's failure object.
// Why it failed, when it did; otherwise null.
type Failure struct {
	// A stable failure code to branch on, e.g. `INSTANCE_BUSY`.
	Code string `json:"code"`
	// A readable sentence saying what went wrong and what to do; safe to show to your users.
	Reason string `json:"reason"`
}

// FirewallRule is the API's firewall rule object.
type FirewallRule struct {
	// e.g. `aclr_...`.
	ID string `json:"id"`
	// The subnet the rule applies to, as a whole.
	SubnetID string `json:"subnet_id"`
	// Evaluation order, 100-9999; lowest first.
	Number int `json:"number"`
	// `ingress` (into the subnet) or `egress` (out of it). One of: ingress, egress.
	Direction string `json:"direction"`
	// `tcp`, `udp`, `icmp` or `all`. One of: tcp, udp, icmp, all.
	Protocol string `json:"protocol"`
	// First port, for `tcp`/`udp`; null otherwise.
	PortStart *int `json:"port_start"`
	// Last port, for `tcp`/`udp`; null otherwise.
	PortEnd *int `json:"port_end"`
	// ICMP type, for `icmp`; null for every type.
	ICMPType *int `json:"icmp_type"`
	// ICMP code, for `icmp`; null for every code.
	ICMPCode *int `json:"icmp_code"`
	// Source (ingress) or destination (egress), e.g. `0.0.0.0/0`.
	CIDR string `json:"cidr"`
	// `allow` or `deny`. One of: allow, deny.
	Action string `json:"action"`
	// What you asked for: `present`, or `deleted` once you deleted it. One of: present, deleted.
	DesiredState string `json:"desired_state"`
	// What the platform last observed: `pending` and `provisioning` while it is being created, `active` once ready, `deleting`/`deleted` on the way out, `failed` if it could not be made, `unknown` when it could not be checked. One of: pending, provisioning, active, deleting, deleted, failed, unknown.
	ObservedState string `json:"observed_state"`
	// When it was created (RFC 3339, UTC).
	CreatedAt *time.Time `json:"created_at"`
	// When it last changed (RFC 3339, UTC).
	UpdatedAt *time.Time `json:"updated_at"`
}

// FirewallRuleFields is the API's firewall rule fields object.
type FirewallRuleFields struct {
	// Evaluation order, 100-9999; lowest first.
	Number int `json:"number"`
	// `ingress` (into the subnet) or `egress` (out of it). One of: ingress, egress.
	Direction string `json:"direction"`
	// `tcp`, `udp`, `icmp` or `all`. One of: tcp, udp, icmp, all.
	Protocol string `json:"protocol"`
	// First port, for `tcp`/`udp`; null otherwise.
	PortStart *int `json:"port_start"`
	// Last port, for `tcp`/`udp`; null otherwise.
	PortEnd *int `json:"port_end"`
	// ICMP type, for `icmp`; null for every type.
	ICMPType *int `json:"icmp_type"`
	// ICMP code, for `icmp`; null for every code.
	ICMPCode *int `json:"icmp_code"`
	// Source (ingress) or destination (egress), e.g. `0.0.0.0/0`.
	CIDR string `json:"cidr"`
	// `allow` or `deny`. One of: allow, deny.
	Action string `json:"action"`
}

// FirewallRuleList is the API's firewall rule list object.
type FirewallRuleList struct {
	Data []FirewallRule `json:"data"`
	// Always null.
	NextCursor *string `json:"next_cursor"`
	// The platform's own rules around yours.
	SystemRules []FirewallRuleFields `json:"system_rules"`
}

// GeneratedSSHKey is the API's generated sshkey object.
type GeneratedSSHKey struct {
	SSHKey
	// The private key, in this response only.
	PrivateKey string `json:"private_key"`
}

// Image is the API's image object.
type Image struct {
	// The image's id.
	ID string `json:"id"`
	// Use this in `image_slug`.
	Slug string `json:"slug"`
	// Its name for display, e.g. `Ubuntu`.
	Name string `json:"name"`
	// Its version, e.g. `24.04`.
	Version string `json:"version"`
	// Zones this image can be deployed in.
	Zones []string `json:"zones"`
	// True for an encrypted image: the instance's root disk is encrypted on our storage, with its own key, and the instance needs our key servers to start.
	Encrypted bool `json:"encrypted"`
	// On an encrypted image, the slug of the plain image it is the encrypted version of (e.g. `ubuntu-24-04` for `ubuntu-24-04-encrypted`); null otherwise.
	EncryptedVariantOf *string `json:"encrypted_variant_of"`
}

// InstanceSpec is the API's instance spec object.
// The capacity this instance was created with.
type InstanceSpec struct {
	VCPU     int `json:"vcpu"`
	MemoryMB int `json:"memory_mb"`
	DiskGB   int `json:"disk_gb"`
}

// Instance is the API's instance object.
type Instance struct {
	// e.g. `vm_06gf...`.
	ID string `json:"id"`
	// Its name, unique among your instances.
	Name string `json:"name"`
	// The plan's id.
	PlanID string `json:"plan_id"`
	// The plan's slug, as `plan_slug` takes it (e.g. `starter`); null if the plan has since been retired.
	PlanSlug *string `json:"plan_slug"`
	// The image's id.
	ImageID string `json:"image_id"`
	// The image's slug, as `image_slug` takes it (e.g. `ubuntu-24-04`); null if the image has since been retired.
	ImageSlug *string `json:"image_slug"`
	// True when the instance runs an encrypted image: its root disk is encrypted on our storage, with its own key.
	Encrypted bool `json:"encrypted"`
	// e.g. `af-abj`.
	Region *string `json:"region"`
	// e.g. `af-abj-1`.
	Zone *string `json:"zone"`
	// The VPC, for an instance in one.
	NetworkID *string `json:"network_id"`
	// Its VPC subnet; null for a standard instance.
	SubnetID *string `json:"subnet_id"`
	// The security group, for a standard instance.
	SecurityGroupID *string `json:"security_group_id"`
	// The instance's own public address (standard instances).
	PublicIPv4 *string `json:"public_ipv4"`
	// Its private address: in its VPC subnet, or on the zone's guest network.
	PrivateIPv4 *string `json:"private_ipv4"`
	// Its interface on the zone's private database network: `none`, `attaching`, `attached` or `detaching`. One of: none, attaching, attached, detaching.
	PrivateNetworkState string `json:"private_network_state"`
	// Platform VM deletion is refused while true.
	DeletionProtection bool `json:"deletion_protection"`
	// That interface's address once attached (kept while it is being removed); null otherwise.
	PrivateNetworkIP *string `json:"private_network_ip"`
	// What you asked for: `running`, `stopped` or `deleted`. One of: running, stopped, deleted.
	DesiredState string `json:"desired_state"`
	// What the platform last observed. One of: pending, provisioning, running, stopping, stopped, deleting, deleted, failed, unknown.
	ObservedState string `json:"observed_state"`
	// The capacity this instance was created with.
	Spec *InstanceSpec `json:"spec"`
	// Your own labels, as sent at create.
	Tags map[string]string `json:"tags"`
	// Why the last change failed, when it did; otherwise null.
	Failure *Failure `json:"failure"`
	// When it was created (RFC 3339, UTC).
	CreatedAt *time.Time `json:"created_at"`
	// When it last changed (RFC 3339, UTC).
	UpdatedAt *time.Time `json:"updated_at"`
}

// InstanceOrder is the API's instance order object.
type InstanceOrder struct {
	// e.g. `ord_...`.
	ID string `json:"id"`
	// The id the instance has, or will have once provisioned.
	InstanceID string `json:"instance_id"`
	// `awaiting_payment`, then `paid`, `provisioning` and `provisioned` (final: the instance exists). One of: awaiting_payment, paid, provisioning, provisioned, payment_failed, failed.
	Status string `json:"status"`
	// The instance's create operation, set once provisioning starts.
	OperationID *string `json:"operation_id"`
	// Why the order failed, e.g. `insufficient_credit`, `card_declined`, `payment_expired`; null otherwise.
	FailureCode *string `json:"failure_code"`
	// The initial payment in minor units of `currency`.
	AmountMinor int64 `json:"amount_minor"`
	// `NGN` or `USD`. One of: NGN, USD.
	Currency string `json:"currency"`
	// When it was created (RFC 3339, UTC).
	CreatedAt *time.Time `json:"created_at"`
	// When it last changed (RFC 3339, UTC).
	UpdatedAt *time.Time `json:"updated_at"`
}

// MeAPIKey is the API's me apikey object.
type MeAPIKey struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// `standard` for a key made under Organization › API keys; `ai_app` for the key an AI app was given when a member connected it. One of: standard, ai_app.
	Kind string `json:"kind"`
	// The AI app's name as it introduced itself (not verified by Pantech), e.g. `Research assistant`.
	ClientName *string `json:"client_name"`
	// `delete` appears only on an AI app's key, together with `write`. One of: read, write, delete.
	Scopes    []string   `json:"scopes"`
	ExpiresAt *time.Time `json:"expires_at"`
}

// Me is the API's me object.
type Me struct {
	// The organization the key belongs to.
	OrganizationID string `json:"organization_id"`
	// The organization's project; every resource belongs to it.
	ProjectID string   `json:"project_id"`
	APIKey    MeAPIKey `json:"api_key"`
}

// MetricsPoint is the API's metrics point object.
type MetricsPoint struct {
	BucketStart                 time.Time `json:"bucket_start"`
	BucketEnd                   time.Time `json:"bucket_end"`
	CPUPercent                  *float64  `json:"cpu_percent"`
	MemoryPercent               *float64  `json:"memory_percent"`
	NetworkInBytes              *int64    `json:"network_in_bytes"`
	NetworkOutBytes             *int64    `json:"network_out_bytes"`
	DiskReadPeakBytesPerSecond  *int64    `json:"disk_read_peak_bytes_per_second"`
	DiskWritePeakBytesPerSecond *int64    `json:"disk_write_peak_bytes_per_second"`
}

// MetricsSeries is the API's metrics series object.
type MetricsSeries struct {
	// The window the points cover: `1h`, `24h`, `7d` or `30d`. One of: 1h, 24h, 7d, 30d.
	Interval string `json:"interval"`
	// False when no samples exist yet, for example just after the instance was created.
	Available bool `json:"available"`
	// When the newest sample was taken; null when none.
	SampledAt *string `json:"sampled_at"`
	// One point per bucket, oldest first.
	Points []MetricsPoint `json:"points"`
}

// MonthlyQuote is the API's monthly quote object.
type MonthlyQuote struct {
	// `NGN` or `USD`.
	Currency string `json:"currency"`
	// Price for a full month in minor units of `currency`.
	MonthlyMinor int64 `json:"monthly_minor"`
	// Hours in the month the price assumes: 730.
	NominalMonthHours int `json:"nominal_month_hours"`
}

// Network is the API's network object.
// A VPC: a private network of your own.
type Network struct {
	// e.g. `net_...`.
	ID string `json:"id"`
	// Its name, unique among your VPC networks.
	Name string `json:"name"`
	// e.g. `10.0.0.0/16`.
	CIDR string `json:"cidr"`
	// e.g. `af-abj`.
	Region *string `json:"region"`
	// e.g. `af-abj-2`.
	Zone *string `json:"zone"`
	// Instances in this VPC; null when it could not be counted.
	InstanceCount *int `json:"instance_count"`
	// What you asked for: `present`, or `deleted` once you deleted it. One of: present, deleted.
	DesiredState string `json:"desired_state"`
	// What the platform last observed: `pending` and `provisioning` while it is being created, `active` once ready, `deleting`/`deleted` on the way out, `failed` if it could not be made, `unknown` when it could not be checked. One of: pending, provisioning, active, deleting, deleted, failed, unknown.
	ObservedState string `json:"observed_state"`
	// When it was created (RFC 3339, UTC).
	CreatedAt *time.Time `json:"created_at"`
	// When it last changed (RFC 3339, UTC).
	UpdatedAt *time.Time `json:"updated_at"`
}

// Operation is the API's operation object.
type Operation struct {
	// e.g. `op_...`.
	ID string `json:"id"`
	// e.g. `instance`, `volume`, `network`, `security_group`, `database`. One of: instance, network, subnet, security_group, snapshot, volume, ssh_access_grant, relay_vm, network_acl_rule, public_ip, port_forwarding_rule, database.
	ResourceType string `json:"resource_type"`
	// The resource the operation acts on, e.g. a `vm_…`, `db_…` or `snap_…` id.
	ResourceID string `json:"resource_id"`
	// What it does, e.g. `create_instance`, `stop_instance`.
	Kind string `json:"kind"`
	// `submitting` and `submitted` while it runs; `succeeded` and `failed` are final. One of: submitting, submitted, succeeded, failed.
	Status string `json:"status"`
	// Why it failed, when `status` is `failed`; otherwise null.
	Failure *Failure `json:"failure"`
	// When it was created (RFC 3339, UTC).
	CreatedAt *time.Time `json:"created_at"`
	// When it last changed (RFC 3339, UTC).
	UpdatedAt *time.Time `json:"updated_at"`
}

// OptionalQuote is the API's optional quote object.
type OptionalQuote struct {
	// False when no price applies; `unpriced_reason` then says why.
	Priced bool `json:"priced"`
	// Why no price applies; null when priced.
	UnpricedReason *string `json:"unpriced_reason"`
	// `NGN` or `USD`; null when not priced.
	Currency *string `json:"currency"`
	// Price for a full month in minor units; null when not priced.
	MonthlyMinor *int64 `json:"monthly_minor"`
	// Hours in the month the price assumes: 730.
	NominalMonthHours int `json:"nominal_month_hours"`
}

// OrderAccepted is the API's order accepted object.
// An instance is bought through an order: the initial payment is taken first (from credit, else the default card), then the instance is provisioned.
type OrderAccepted struct {
	// Follow it at `GET /public/v1/instance-orders/{order_id}` (also the `Location` header).
	OrderID string `json:"order_id"`
	// The id the instance will have.
	InstanceID string `json:"instance_id"`
	// The order's status now, usually `awaiting_payment`. One of: awaiting_payment, paid, provisioning, provisioned, payment_failed, failed.
	Status string `json:"status"`
	// Set once provisioning starts.
	OperationID *string `json:"operation_id"`
	// Null for a new order; on a retry, the order's `failure_code` if it has failed since.
	FailureCode *string `json:"failure_code"`
	// The initial payment, in minor units (kobo or cents).
	AmountMinor int64 `json:"amount_minor"`
	// `NGN` or `USD`. One of: NGN, USD.
	Currency string `json:"currency"`
}

// PlanPrice is the API's plan price object.
// In the organization's currency; null when it cannot be priced.
type PlanPrice struct {
	Currency string `json:"currency"`
	// A full month running.
	MonthlyEstimateMinor int64 `json:"monthly_estimate_minor"`
	// A full month stopped (storage only).
	StorageFloorMinor int64 `json:"storage_floor_minor"`
	// Taken when the instance is ordered.
	InitialPaymentMinor int64 `json:"initial_payment_minor"`
}

// Plan is the API's plan object.
type Plan struct {
	// The plan's id.
	ID string `json:"id"`
	// Use this in `plan_slug`.
	Slug string `json:"slug"`
	// Its name for display.
	Name string `json:"name"`
	// Virtual CPUs.
	VCPU int `json:"vcpu"`
	// Memory in MB.
	MemoryMB int `json:"memory_mb"`
	// Root disk in GB; for a database, the smallest data disk.
	DiskGB int `json:"disk_gb"`
	// In the organization's currency; null when it cannot be priced.
	Price *PlanPrice `json:"price"`
	// Why `price` is null; null when priced.
	UnpricedReason *string `json:"unpriced_reason"`
}

// PortForwardingRule is the API's port forwarding rule object.
type PortForwardingRule struct {
	// e.g. `pfr_...`.
	ID         string `json:"id"`
	PublicIPID string `json:"public_ip_id"`
	InstanceID string `json:"instance_id"`
	// The target instance's name; null when it is gone.
	InstanceName *string `json:"instance_name"`
	// `tcp` or `udp`. One of: tcp, udp.
	Protocol         string `json:"protocol"`
	PublicPortStart  int    `json:"public_port_start"`
	PublicPortEnd    int    `json:"public_port_end"`
	PrivatePortStart int    `json:"private_port_start"`
	PrivatePortEnd   int    `json:"private_port_end"`
	// What you asked for: `present`, or `deleted` once you deleted it. One of: present, deleted.
	DesiredState string `json:"desired_state"`
	// What the platform last observed: `pending` and `provisioning` while it is being created, `active` once ready, `deleting`/`deleted` on the way out, `failed` if it could not be made, `unknown` when it could not be checked. One of: pending, provisioning, active, deleting, deleted, failed, unknown.
	ObservedState string `json:"observed_state"`
	// When it was created (RFC 3339, UTC).
	CreatedAt *time.Time `json:"created_at"`
	// When it last changed (RFC 3339, UTC).
	UpdatedAt *time.Time `json:"updated_at"`
}

// PricingEstimateMachine is the API's pricing estimate machine object.
type PricingEstimateMachine struct {
	PlanSlug                  string  `json:"plan_slug"`
	Quantity                  int     `json:"quantity"`
	RunningHours              int     `json:"running_hours"`
	Priced                    bool    `json:"priced"`
	UnpricedReason            *string `json:"unpriced_reason"`
	UnitMonthlyEstimateMinor  *int64  `json:"unit_monthly_estimate_minor"`
	TotalMonthlyEstimateMinor *int64  `json:"total_monthly_estimate_minor"`
}

// PricingEstimateExtra is the API's pricing estimate extra object.
type PricingEstimateExtra struct {
	Requested      int     `json:"requested"`
	Priced         bool    `json:"priced"`
	UnpricedReason *string `json:"unpriced_reason"`
	AmountMinor    *int64  `json:"amount_minor"`
}

// PricingEstimate is the API's pricing estimate object.
type PricingEstimate struct {
	// `NGN` or `USD`, as requested.
	Currency string `json:"currency"`
	// Hours in the month prices assume: 730.
	NominalMonthHours int                      `json:"nominal_month_hours"`
	Machines          []PricingEstimateMachine `json:"machines"`
	ExtraStorageGB    *PricingEstimateExtra    `json:"extra_storage_gb"`
	Networks          *PricingEstimateExtra    `json:"networks"`
	SnapshotGB        *PricingEstimateExtra    `json:"snapshot_gb"`
	// False when any line could not be priced.
	PricingComplete bool `json:"pricing_complete"`
	// The sum of every priced line, in minor units.
	PricedSubtotalMinor int64 `json:"priced_subtotal_minor"`
}

// PricingEstimateMachineRequest is the API's pricing estimate machine request object.
type PricingEstimateMachineRequest struct {
	PlanSlug     string `json:"plan_slug"`
	Quantity     int    `json:"quantity"`
	RunningHours int    `json:"running_hours"`
	// One of: standard, vpc.
	Placement *string `json:"placement,omitempty"`
}

// PricingEstimateRequest is the API's pricing estimate request object.
type PricingEstimateRequest struct {
	// The currency to price in: `NGN` or `USD`. One of: NGN, USD.
	Currency string `json:"currency"`
	// Instances to price: up to 50 lines, each a plan, a quantity and running hours per month.
	Machines *[]PricingEstimateMachineRequest `json:"machines,omitempty"`
	// Extra volume storage in GB.
	ExtraStorageGB *int `json:"extra_storage_gb,omitempty"`
	// Number of VPC networks.
	NetworkCount *int `json:"network_count,omitempty"`
	// Snapshot storage in GB.
	SnapshotGB *int `json:"snapshot_gb,omitempty"`
}

// AttachPublicIPRequest is the API's attach public iprequest object.
type AttachPublicIPRequest struct {
	// An instance in the address's private network.
	InstanceID string `json:"instance_id"`
}

// PublicIP is the API's public ip object.
type PublicIP struct {
	// e.g. `pip_...`.
	ID        string `json:"id"`
	NetworkID string `json:"network_id"`
	// The VPC's name.
	NetworkName *string `json:"network_name"`
	// `static_nat` maps every port to one instance; `port_forwarding` carries per-port rules. One of: static_nat, port_forwarding, load_balancer.
	Purpose string `json:"purpose"`
	// The public IPv4 address; null until allocated.
	Address *string `json:"address"`
	// static_nat: the instance it maps to; null while it is detached (still yours, and billed once).
	InstanceID *string `json:"instance_id"`
	// The mapped instance's name; null when none.
	InstanceName *string `json:"instance_name"`
	Region       *string `json:"region"`
	Zone         *string `json:"zone"`
	// What you asked for: `present`, or `deleted` once you deleted it. One of: present, deleted.
	DesiredState string `json:"desired_state"`
	// What the platform last observed: `pending` and `provisioning` while it is being created, `active` once ready, `deleting`/`deleted` on the way out, `failed` if it could not be made, `unknown` when it could not be checked. One of: pending, provisioning, active, deleting, deleted, failed, unknown.
	ObservedState string `json:"observed_state"`
	// False while an attach or detach has not reached the address yet.
	InSync bool `json:"in_sync"`
	// When it was created (RFC 3339, UTC).
	CreatedAt *time.Time `json:"created_at"`
	// When it last changed (RFC 3339, UTC).
	UpdatedAt *time.Time `json:"updated_at"`
}

// PutSnapshotScheduleRequest is the API's put snapshot schedule request object.
type PutSnapshotScheduleRequest struct {
	// How often: `daily`, `weekly` or `monthly`. One of: daily, weekly, monthly.
	Frequency string `json:"frequency"`
	// How many successful automatic snapshots to keep, 1 to 168; older ones are deleted.
	RetentionCount *int `json:"retention_count,omitempty"`
	// False pauses the schedule without forgetting it.
	Enabled *bool `json:"enabled,omitempty"`
}

// RebuildRequest is the API's rebuild request object.
type RebuildRequest struct {
	// The image to install, from `GET /public/v1/images`.
	ImageSlug *string `json:"image_slug,omitempty"`
}

// RegionPlacement is the API's region placement object.
type RegionPlacement struct {
	// `standard` instances or `vpc` networks. One of: standard, vpc.
	Kind              string  `json:"kind"`
	Zone              string  `json:"zone"`
	Available         bool    `json:"available"`
	UnavailableReason *string `json:"unavailable_reason"`
	// The range of the zone's private database network, e.g. `10.250.0.0/20`, when instances there can add a private network interface (`POST /public/v1/instances/{instance_id}/private-network`); null otherwise, and in every VPC zone.
	PrivateNetworkCIDR *string `json:"private_network_cidr"`
}

// Region is the API's region object.
type Region struct {
	// e.g. `af-abj`.
	Code string `json:"code"`
	// Its name for display.
	Name       string            `json:"name"`
	Placements []RegionPlacement `json:"placements"`
}

// RenameRequest is the API's rename request object.
type RenameRequest struct {
	// The new name: 1 to 255 characters of lowercase letters a-z, digits 0-9 and hyphens (-), starting and ending with a letter or digit, unique among your instances.
	Name string `json:"name"`
}

// ReplaceDatabaseAccessRulesRequest is the API's replace database access rules request object.
type ReplaceDatabaseAccessRulesRequest struct {
	// The complete new allow-list, at most 50 rules with no range twice.
	Rules []DatabaseAccessRuleRequest `json:"rules"`
}

// ReplaceSecurityGroupRulesRequest is the API's replace security group rules request object.
type ReplaceSecurityGroupRulesRequest struct {
	// The complete new rule list; at least one rule.
	Rules []SecurityGroupRule `json:"rules"`
}

// ChangeSSHKeyRequest is the API's change sshkey request object.
type ChangeSSHKeyRequest struct {
	// The SSH key the instance should accept from now on, from `GET /public/v1/ssh-keys`.
	SSHKeyID string `json:"ssh_key_id"`
}

// ResizeRequest is the API's resize request object.
type ResizeRequest struct {
	// A bigger plan than the instance's current one, from `GET /public/v1/plans`.
	PlanSlug string `json:"plan_slug"`
}

// ResizeVolumeRequest is the API's resize volume request object.
type ResizeVolumeRequest struct {
	// The volume type after the resize, from `GET /public/v1/disk-offerings`.
	DiskOfferingSlug string `json:"disk_offering_slug"`
	// The new size in GB, for a custom-size offering: larger than the current size.
	SizeGB *int `json:"size_gb,omitempty"`
}

// RestoreRequest is the API's restore request object.
type RestoreRequest struct {
	// A snapshot of this instance's root disk.
	SnapshotID string `json:"snapshot_id"`
	// The instance's current name, exactly, to confirm: everything written since the snapshot is lost.
	ConfirmName string `json:"confirm_name"`
}

// RestoreSnapshotRequest is the API's restore snapshot request object.
type RestoreSnapshotRequest struct {
	// Name of the new volume, 1 to 255 characters.
	Name string `json:"name"`
	// The new volume's type, from `GET /public/v1/disk-offerings`.
	DiskOfferingSlug string `json:"disk_offering_slug"`
	// For a custom-size offering: at least the snapshot's size in GB.
	SizeGB *int `json:"size_gb,omitempty"`
}

// SSHAccessGrant is the API's sshaccess grant object.
type SSHAccessGrant struct {
	// e.g. `sshg_...`.
	ID         string `json:"id"`
	InstanceID string `json:"instance_id"`
	// `pending` while the port is being opened, `active` once you can connect, then `revoking`, `revoked`, `expired` or `failed`. One of: pending, active, revoking, revoked, expired, failed.
	Status string `json:"status"`
	// Connect here once `active`.
	Host *string `json:"host"`
	// The port to connect to once `active`, usually 22.
	Port *int `json:"port"`
	// When the port closes by itself, 15 minutes after it opened.
	ExpiresAt *time.Time `json:"expires_at"`
	// When it was revoked; null otherwise.
	RevokedAt *time.Time `json:"revoked_at"`
	// When it was created (RFC 3339, UTC).
	CreatedAt *time.Time `json:"created_at"`
}

// SSHKey is the API's sshkey object.
type SSHKey struct {
	// e.g. `sshk_...`.
	ID string `json:"id"`
	// Its name, unique among your SSH keys.
	Name string `json:"name"`
	// The key's SHA256 fingerprint.
	Fingerprint string `json:"fingerprint"`
	// The OpenSSH public key.
	PublicKey string `json:"public_key"`
	// When it was created (RFC 3339, UTC).
	CreatedAt *time.Time `json:"created_at"`
}

// SecurityGroup is the API's security group object.
// A reusable firewall for standard instances.
type SecurityGroup struct {
	// e.g. `sg_...`.
	ID string `json:"id"`
	// Its name, unique among your live security groups.
	Name string `json:"name"`
	// Its rules.
	Rules []SecurityGroupRule `json:"rules"`
	// What you asked for: `present`, or `deleted` once you deleted it. One of: present, deleted.
	DesiredState string `json:"desired_state"`
	// What the platform last observed: `pending` and `provisioning` while it is being created, `active` once ready, `deleting`/`deleted` on the way out, `failed` if it could not be made, `unknown` when it could not be checked. One of: pending, provisioning, active, deleting, deleted, failed, unknown.
	ObservedState string `json:"observed_state"`
	// When it was created (RFC 3339, UTC).
	CreatedAt *time.Time `json:"created_at"`
	// When it last changed (RFC 3339, UTC).
	UpdatedAt *time.Time `json:"updated_at"`
}

// SecurityGroupRule is the API's security group rule object.
type SecurityGroupRule struct {
	// `ingress` (into the instance) or `egress` (out of it). One of: ingress, egress.
	Direction string `json:"direction"`
	// `tcp`, `udp`, `icmp` or `all`. One of: tcp, udp, icmp, all.
	Protocol string `json:"protocol"`
	// For `tcp`/`udp`: one port (`22`) or an inclusive range (`8000-8100`).
	PortRange string `json:"port_range"`
	// One IPv4 CIDR: the source for `ingress`, the destination for `egress`, e.g. `203.0.113.4/32`.
	CIDR string `json:"cidr"`
}

// Snapshot is the API's snapshot object.
type Snapshot struct {
	// e.g. `snap_...`.
	ID string `json:"id"`
	// Its name, unique among its source's snapshots.
	Name string `json:"name"`
	// Set for a snapshot of an instance's root disk.
	InstanceID *string `json:"instance_id"`
	// The source instance's name; null when it is gone or not an instance snapshot.
	SourceInstanceName *string `json:"source_instance_name"`
	// Set for a snapshot of a volume.
	VolumeID *string `json:"volume_id"`
	// The source volume's name; null when it is gone or not a volume snapshot.
	VolumeName *string `json:"volume_name"`
	// Set for a snapshot of a managed database's data disk.
	DatabaseID *string `json:"database_id"`
	// The source's region, e.g. `af-abj`.
	Region *string `json:"region"`
	// `automatic` when taken by a snapshot schedule. One of: manual, automatic.
	Trigger string `json:"trigger"`
	// Storage used; 0 until complete.
	SizeBytes int64 `json:"size_bytes"`
	// What you asked for: `present`, or `deleted` once you deleted it. One of: present, deleted.
	DesiredState string `json:"desired_state"`
	// What the platform last observed: `pending` and `provisioning` while it is being created, `active` once ready, `deleting`/`deleted` on the way out, `failed` if it could not be made, `unknown` when it could not be checked. One of: pending, provisioning, active, deleting, deleted, failed, unknown.
	ObservedState string `json:"observed_state"`
	// When the snapshot finished; null until then.
	CompletedAt *time.Time `json:"completed_at"`
	// When it was created (RFC 3339, UTC).
	CreatedAt *time.Time `json:"created_at"`
	// When it last changed (RFC 3339, UTC).
	UpdatedAt *time.Time `json:"updated_at"`
}

// SnapshotSchedule is the API's snapshot schedule object.
type SnapshotSchedule struct {
	// The instance whose root disk is snapshotted; null for a volume's schedule.
	InstanceID *string `json:"instance_id"`
	// The volume snapshotted; null for an instance's schedule.
	VolumeID *string `json:"volume_id"`
	// `daily`, `weekly` or `monthly`. One of: daily, weekly, monthly.
	Frequency string `json:"frequency"`
	// Successful copies kept; older ones are deleted.
	RetentionCount int `json:"retention_count"`
	// False while paused.
	Enabled bool `json:"enabled"`
	// When the next automatic snapshot is due; null while paused.
	NextRunAt *time.Time `json:"next_run_at"`
}

// Subnet is the API's subnet object.
type Subnet struct {
	// e.g. `snet_...`.
	ID        string `json:"id"`
	NetworkID string `json:"network_id"`
	// Its name.
	Name string `json:"name"`
	// Inside the VPC's CIDR.
	CIDR string `json:"cidr"`
	// e.g. `af-abj`.
	Region *string `json:"region"`
	// e.g. `af-abj-2`.
	Zone *string `json:"zone"`
	// Instances in this subnet; null when they could not be counted.
	InstanceCount *int `json:"instance_count"`
	// What you asked for: `present`, or `deleted` once you deleted it. One of: present, deleted.
	DesiredState string `json:"desired_state"`
	// What the platform last observed: `pending` and `provisioning` while it is being created, `active` once ready, `deleting`/`deleted` on the way out, `failed` if it could not be made, `unknown` when it could not be checked. One of: pending, provisioning, active, deleting, deleted, failed, unknown.
	ObservedState string `json:"observed_state"`
	// When it was created (RFC 3339, UTC).
	CreatedAt *time.Time `json:"created_at"`
	// When it last changed (RFC 3339, UTC).
	UpdatedAt *time.Time `json:"updated_at"`
}

// ResourceCost is the API's resource cost object.
type ResourceCost struct {
	ResourceID          string `json:"resource_id"`
	ResourceName        string `json:"resource_name"`
	ResourceKind        string `json:"resource_kind"`
	Type                string `json:"type"`
	DaysActive          int    `json:"days_active"`
	WindowDays          int    `json:"window_days"`
	RateUnit            string `json:"rate_unit"`
	RateUnitPriceMicros int64  `json:"rate_unit_price_micros"`
	CostMinor           int64  `json:"cost_minor"`
}

// UsageByResource is the API's usage by resource object.
type UsageByResource struct {
	Currency    string         `json:"currency"`
	PeriodStart string         `json:"period_start"`
	PeriodEnd   string         `json:"period_end"`
	Data        []ResourceCost `json:"data"`
}

// CreateCreditPurchaseRequest is the API's create credit purchase request object.
type CreateCreditPurchaseRequest struct {
	// How much credit to buy, in minor units of your organization's currency: kobo for NGN, cents for USD.
	AmountMinor int64 `json:"amount_minor"`
}

// ReceiptFact is the API's receipt fact object.
type ReceiptFact struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// Receipt is the API's receipt object.
// A receipt for money that moved on your account.
type Receipt struct {
	// e.g. `rcp_...`.
	ID string `json:"id"`
	// `PTD-RCT-<year>-<number>`, sequential within the year.
	Number string `json:"number"`
	// `payment` a card payment, `invoice_payment` an invoice paid by card or from credit, `credit_purchase` a top-up, `wallet_debit` credit taken for a new resource, `refund` money returned to your credit. One of: payment, invoice_payment, credit_purchase, wallet_debit, refund.
	Kind       string `json:"kind"`
	IssuerName string `json:"issuer_name"`
	// One of: NGN, USD.
	Currency    string `json:"currency"`
	AmountMinor int64  `json:"amount_minor"`
	// When the receipt was issued.
	IssuedAt *time.Time `json:"issued_at"`
	// When the money moved.
	OccurredAt *time.Time `json:"occurred_at"`
	// What it was for, e.g. `Account credit top-up`.
	Description string `json:"description"`
	// Your organization, as printed on the receipt.
	BilledTo string `json:"billed_to"`
	// The details printed on the receipt, e.g. the payment method (`Visa •••• 4242`) and the date.
	Facts []ReceiptFact `json:"facts"`
	// The code printed on the receipt, which Pantech uses to confirm a copy is genuine.
	Verification string `json:"verification"`
	// The account transaction the receipt is for, e.g. `ltx_...`.
	TransactionID string `json:"transaction_id"`
	// The size of the PDF.
	PDFBytes int64 `json:"pdf_bytes"`
}

// CreditPurchaseReceipt is the API's credit purchase receipt object.
// The receipt, once the purchase has settled and the receipt is issued (usually within a minute); null until then.
type CreditPurchaseReceipt struct {
	// e.g. `rcp_...`.
	ID string `json:"id"`
	// e.g. `PTD-RCT-2026-000123`.
	Number string `json:"number"`
}

// CreditPurchase is the API's credit purchase object.
// One credit top-up.
type CreditPurchase struct {
	// e.g. `cpr_...`.
	ID string `json:"id"`
	// `pending` until paid; `settled` once the credit is in your balance. One of: pending, settled, failed, expired, refunded, charged_back.
	Status string `json:"status"`
	// The credit bought, in minor units of `currency`.
	AmountMinor int64 `json:"amount_minor"`
	// One of: NGN, USD.
	Currency string `json:"currency"`
	// `checkout` for every purchase made through the API. One of: checkout, saved_card, bank_transfer.
	Method string `json:"method"`
	// An identifier for the payment processor that takes the payment. Treat it
	// as an opaque label: values can be added.
	Provider          string     `json:"provider"`
	CheckoutURL       *string    `json:"checkout_url"`
	CheckoutExpiresAt *time.Time `json:"checkout_expires_at"`
	// Why the payment failed, when `status` is `failed`.
	FailureCode *string `json:"failure_code"`
	// When the credit was added; null until then.
	SettledAt *time.Time `json:"settled_at"`
	// When it was created (RFC 3339, UTC).
	CreatedAt *time.Time `json:"created_at"`
	// The receipt, once the purchase has settled and the receipt is issued (usually within a minute); null until then.
	Receipt *CreditPurchaseReceipt `json:"receipt"`
}

// CreditOutlook is the API's credit outlook object.
// Credit against the estimated next invoice.
type CreditOutlook struct {
	// One of: NGN, USD.
	Currency string `json:"currency"`
	// Credit you can spend now.
	AvailableMinor int64 `json:"available_minor"`
	// Credit purchases that haven't settled yet.
	PendingMinor int64 `json:"pending_minor"`
	// Initial payments for new instances that haven't been used yet.
	PrepaidMinor int64 `json:"prepaid_minor"`
	// An estimate of your next invoice, including VAT.
	EstimatedNextInvoiceMinor int64 `json:"estimated_next_invoice_minor"`
	// When the current billing period closes and the invoice is issued.
	NextInvoiceAt time.Time `json:"next_invoice_at"`
	// `covered` when credit and prepaid instance payments cover the estimate, `low` when they don't, `empty` when no credit is left. One of: covered, low, empty.
	Standing string `json:"standing"`
	// Whether a default payment method will be charged for anything credit doesn't cover.
	HasPaymentMethod bool `json:"has_payment_method"`
}

// UsageTotals is the API's usage totals object.
type UsageTotals struct {
	RunningHours           float64 `json:"running_hours"`
	StorageGBHours         float64 `json:"storage_gb_hours"`
	AllocatedMemoryMBHours float64 `json:"allocated_memory_mb_hours"`
	NetworkBytesSent       float64 `json:"network_bytes_sent"`
	NetworkBytesReceived   float64 `json:"network_bytes_received"`
}

// UsageEntry is the API's usage entry object.
type UsageEntry struct {
	ResourceID  *string    `json:"resource_id"`
	UsageType   int        `json:"usage_type"`
	Usage       string     `json:"usage"`
	PeriodStart *time.Time `json:"period_start"`
	PeriodEnd   *time.Time `json:"period_end"`
}

// UsageSummary is the API's usage summary object.
type UsageSummary struct {
	// Start of the measured window; null when nothing was measured.
	WindowStart *time.Time `json:"window_start"`
	// End of the measured window; null when nothing was measured.
	WindowEnd *time.Time   `json:"window_end"`
	Totals    UsageTotals  `json:"totals"`
	Entries   []UsageEntry `json:"entries"`
}

// VolumeMonthlyCost is the API's volume monthly cost object.
// What it costs per month at its size; null when it cannot be priced.
type VolumeMonthlyCost struct {
	Currency    string `json:"currency"`
	AmountMinor int64  `json:"amount_minor"`
}

// VolumeEncryption is the API's volume encryption object.
// What the platform can truthfully say about this volume's encryption.
type VolumeEncryption struct {
	// One of: none, requested, encrypted.
	Status string `json:"status"`
}

// Volume is the API's volume object.
type Volume struct {
	// e.g. `vol_...`.
	ID string `json:"id"`
	// Its name.
	Name string `json:"name"`
	// Its size in GB.
	SizeGB int `json:"size_gb"`
	// See `GET /public/v1/disk-offerings`.
	DiskOfferingSlug string `json:"disk_offering_slug"`
	// The storage class, e.g. `shared`; null when not stated.
	StorageType *string `json:"storage_type"`
	Region      *string `json:"region"`
	Zone        *string `json:"zone"`
	// The instance it is attached to; null when detached.
	AttachedInstanceID *string `json:"attached_instance_id"`
	// That instance's name; null when detached.
	AttachedInstanceName *string `json:"attached_instance_name"`
	// Where an attach in progress is moving it.
	DesiredInstanceID *string `json:"desired_instance_id"`
	// Where it is mounted, when mounted at create; null otherwise.
	MountPoint *string `json:"mount_point"`
	// Set when restored from a snapshot.
	SourceSnapshotID *string `json:"source_snapshot_id"`
	// What it costs per month at its size; null when it cannot be priced.
	MonthlyCost *VolumeMonthlyCost `json:"monthly_cost"`
	// What you asked for: `present`, or `deleted` once you deleted it. One of: present, deleted.
	DesiredState string `json:"desired_state"`
	// What the platform last observed: `pending` and `provisioning` while it is being created, `active` once ready, `deleting`/`deleted` on the way out, `failed` if it could not be made, `unknown` when it could not be checked. One of: pending, provisioning, active, deleting, deleted, failed, unknown.
	ObservedState string `json:"observed_state"`
	// What the platform can truthfully say about this volume's encryption.
	Encryption *VolumeEncryption `json:"encryption"`
	// When it was created (RFC 3339, UTC).
	CreatedAt *time.Time `json:"created_at"`
	// When it last changed (RFC 3339, UTC).
	UpdatedAt *time.Time `json:"updated_at"`
}

// Webhook is the API's webhook object.
// A URL Pantech sends your organization's events to.
type Webhook struct {
	// e.g. `whk_...`.
	ID string `json:"id"`
	// Where events are sent: an `https` URL on the public internet.
	URL string `json:"url"`
	// Your note about it; may be empty.
	Description string `json:"description"`
	// What it subscribes to: event types, `<resource>.*`, or `*`.
	Events []string `json:"events"`
	// Whether events are sent to it.
	Enabled bool `json:"enabled"`
	// Why it is disabled: `disabled_by_user`, or `failing` after three days of failed deliveries. One of: disabled_by_user, failing.
	DisabledReason *string `json:"disabled_reason"`
	// When its deliveries started failing with no success since (RFC 3339, UTC); null while they succeed.
	FailingSince *time.Time `json:"failing_since"`
	// When it was created (RFC 3339, UTC).
	CreatedAt time.Time `json:"created_at"`
	// When it last changed (RFC 3339, UTC).
	UpdatedAt time.Time `json:"updated_at"`
}

// WebhookWithSecret is the API's webhook with secret object.
type WebhookWithSecret struct {
	Webhook
	// The signing secret (`whsec_…`).
	Secret *string `json:"secret"`
}

// CreateWebhookRequest is the API's create webhook request object.
type CreateWebhookRequest struct {
	// An `https` URL on the public internet, at most 2,048 characters, with no user name or password.
	URL string `json:"url"`
	// A note for yourself, at most 255 characters.
	Description *string `json:"description,omitempty"`
	// 1 to 100 entries: event types from `GET /webhook-event-types`, `<resource>.*`, or `*`.
	Events []string `json:"events"`
}

// UpdateWebhookRequest is the API's update webhook request object.
type UpdateWebhookRequest struct {
	// A new `https` URL.
	URL *string `json:"url,omitempty"`
	// A new note, at most 255 characters.
	Description *string `json:"description,omitempty"`
	// The new subscriptions, replacing the old ones.
	Events *[]string `json:"events,omitempty"`
	// `false` stops deliveries; `true` resumes them and clears the failure history.
	Enabled *bool `json:"enabled,omitempty"`
}

// WebhookDelivery is the API's webhook delivery object.
// One event sent, or waiting to be sent, to one webhook.
type WebhookDelivery struct {
	// e.g. `whd_...`.
	ID string `json:"id"`
	// The webhook it is for.
	WebhookID string `json:"webhook_id"`
	// The event's `id`, the same in every delivery of it: deduplicate on it.
	EventID string `json:"event_id"`
	// e.g. `instance.created`.
	EventType string `json:"event_type"`
	// `pending` until it succeeds, fails eight times, or is cancelled because the webhook was disabled or deleted. One of: pending, succeeded, failed, canceled.
	Status string `json:"status"`
	// How many times it has been sent.
	Attempts int `json:"attempts"`
	// When it is next sent, while `pending`.
	NextAttemptAt *time.Time `json:"next_attempt_at"`
	// When it was last sent.
	LastAttemptAt *time.Time `json:"last_attempt_at"`
	// The HTTP status your endpoint last answered; null if it did not answer.
	LastStatusCode *int `json:"last_status_code"`
	// Why the last attempt failed: `http_<status>`, `timeout`, `connection_failed`, `dns_failed`, `tls_failed`, `blocked_address` (the URL resolved to an address that is not on the public internet), `invalid_url`, `endpoint_disabled` or `endpoint_deleted`.
	LastError *string `json:"last_error"`
	// The first 1,024 characters of your endpoint's last response body.
	LastResponse *string `json:"last_response"`
	// How long the last attempt took, in milliseconds.
	LastDurationMS *int `json:"last_duration_ms"`
	// When it succeeded.
	DeliveredAt *time.Time `json:"delivered_at"`
	// When the event was queued for this webhook.
	CreatedAt time.Time `json:"created_at"`
}

// WebhookDeliveryDetail is the API's webhook delivery detail object.
type WebhookDeliveryDetail struct {
	WebhookDelivery
	Payload WebhookEvent `json:"payload"`
}

// WebhookEvent is the API's webhook event object.
// The body of every delivery.
type WebhookEvent struct {
	// The event's id (`evt_…`).
	ID string `json:"id"`
	// e.g. `instance.created`, `invoice.paid`, `webhook.test`.
	Type string `json:"type"`
	// The shape of this body, as a date: `2026-10-22`.
	APIVersion string `json:"api_version"`
	// When it happened (RFC 3339, UTC).
	CreatedAt time.Time `json:"created_at"`
	// Your organization.
	OrganizationID string `json:"organization_id"`
	// The project (`acc_…`) the resource is in.
	ProjectID string `json:"project_id"`
	// What happened; depends on `type`.
	Data map[string]any `json:"data"`
}

// WebhookEventType is the API's webhook event type object.
type WebhookEventType struct {
	// e.g. `instance.created`.
	Type string `json:"type"`
	// What it is about, for `<resource>.*` subscriptions.
	Resource string `json:"resource"`
	// What it means.
	Description string `json:"description"`
}

// DiskOfferingQuote is the API's disk offering quote object.
type DiskOfferingQuote struct {
	DiskOfferingSlug  string `json:"disk_offering_slug"`
	SizeGB            int    `json:"size_gb"`
	Currency          string `json:"currency"`
	HourlyPriceMinor  int64  `json:"hourly_price_minor"`
	MonthlyPriceMinor int64  `json:"monthly_price_minor"`
}

// InstanceConsole is the API's instance console object.
type InstanceConsole struct {
	URL string `json:"url"`
}

// VolumeKeyCredential is the API's volume key credential object.
type VolumeKeyCredential struct {
	InstanceID string `json:"instance_id"`
	// The secret, `pvk_` followed by 43 characters.
	Credential string `json:"credential"`
}

// KubernetesClusterNode is the API's kubernetes cluster node object.
type KubernetesClusterNode struct {
	// The node's name, which is its id inside Kubernetes.
	ID   string `json:"id"`
	Name string `json:"name"`
	// One of: control, worker.
	Role string `json:"role"`
	// One of: Ready, NotReady.
	Status string `json:"status"`
	// The node's internal address.
	IPAddress      string     `json:"ip_address"`
	VCPU           int        `json:"vcpu"`
	MemoryMB       int        `json:"memory_mb"`
	DiskGB         int        `json:"disk_gb"`
	KubeletVersion string     `json:"kubelet_version"`
	CreatedAt      *time.Time `json:"created_at"`
}

// KubernetesClusterNodes is the API's kubernetes cluster nodes object.
type KubernetesClusterNodes struct {
	Nodes []KubernetesClusterNode `json:"nodes"`
}

// KubernetesServicePort is the API's kubernetes service port object.
type KubernetesServicePort struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Port     int    `json:"port"`
	// Left out when the target is a named container port.
	TargetPort int `json:"target_port"`
	NodePort   int `json:"node_port"`
}

// KubernetesWorkloadService is the API's kubernetes workload service object.
type KubernetesWorkloadService struct {
	Name       string                  `json:"name"`
	Namespace  string                  `json:"namespace"`
	Type       string                  `json:"type"`
	ClusterIP  string                  `json:"cluster_ip"`
	ExternalIP *string                 `json:"external_ip"`
	Ports      []KubernetesServicePort `json:"ports"`
	// The public port routed to this service; null when it is not exposed.
	PublicPort  *int       `json:"public_port"`
	ExternalURL *string    `json:"external_url"`
	CreatedAt   *time.Time `json:"created_at"`
}

// KubernetesPod is the API's kubernetes pod object.
type KubernetesPod struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Status    string `json:"status"`
	// Ready containers over all containers.
	Ready    string `json:"ready"`
	Restarts int    `json:"restarts"`
	NodeName string `json:"node_name"`
	PodIP    string `json:"pod_ip"`
	// As kubectl's AGE column.
	Age       string     `json:"age"`
	CreatedAt *time.Time `json:"created_at"`
}

// KubernetesWorkloads is the API's kubernetes workloads object.
type KubernetesWorkloads struct {
	ClusterID string `json:"cluster_id"`
	// The public address exposed services and the API server are reached on; null until routing has been synced.
	PublicIP *string `json:"public_ip"`
	// The cluster's API server port on `public_ip`; null for a private endpoint or before routing is synced.
	APIPort  *int                        `json:"api_port"`
	Services []KubernetesWorkloadService `json:"services"`
	Pods     []KubernetesPod             `json:"pods"`
}

// KubernetesAppliedResource is the API's kubernetes applied resource object.
type KubernetesAppliedResource struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
	// Empty for a cluster-wide object.
	Namespace string `json:"namespace"`
	// One of: created, configured, unchanged.
	Status string `json:"status"`
}

// KubernetesManifestApplied is the API's kubernetes manifest applied object.
type KubernetesManifestApplied struct {
	Applied   bool   `json:"applied"`
	ClusterID string `json:"cluster_id"`
	// One line per object, as kubectl prints it.
	Output    string                      `json:"output"`
	Resources []KubernetesAppliedResource `json:"resources"`
}

// ApplyKubernetesManifestRequest is the API's apply kubernetes manifest request object.
type ApplyKubernetesManifestRequest struct {
	Manifest string `json:"manifest"`
	// Where namespaced objects without one go; default `default`.
	Namespace *string `json:"namespace,omitempty"`
}

// KubernetesLiveManifests is the API's kubernetes live manifests object.
type KubernetesLiveManifests struct {
	ClusterID string `json:"cluster_id"`
	// Empty when you have deployed nothing.
	Manifest      string    `json:"manifest"`
	ResourceCount int       `json:"resource_count"`
	RetrievedAt   time.Time `json:"retrieved_at"`
}

// KubernetesServiceExposed is the API's kubernetes service exposed object.
type KubernetesServiceExposed struct {
	Success     bool   `json:"success"`
	ClusterID   string `json:"cluster_id"`
	ServiceName string `json:"service_name"`
	Namespace   string `json:"namespace"`
	PublicIP    string `json:"public_ip"`
	PublicPort  int    `json:"public_port"`
	NodePort    int    `json:"node_port"`
	Protocol    string `json:"protocol"`
	ExternalURL string `json:"external_url"`
	// One of: active.
	Status string `json:"status"`
}

// ExposeKubernetesServiceRequest is the API's expose kubernetes service request object.
type ExposeKubernetesServiceRequest struct {
	ServiceName string `json:"service_name"`
	// Default `default`.
	Namespace  *string `json:"namespace,omitempty"`
	PublicPort int     `json:"public_port"`
	// Leave out to keep the Service's own NodePort, or let the cluster choose one.
	NodePort *int `json:"node_port,omitempty"`
	// One of: TCP, UDP.
	Protocol *string `json:"protocol,omitempty"`
}

// KubernetesSyncedPort is the API's kubernetes synced port object.
type KubernetesSyncedPort struct {
	// One of: kubernetes-api, service.
	Purpose     string `json:"purpose"`
	PublicPort  int    `json:"public_port"`
	PrivatePort int    `json:"private_port"`
	Protocol    string `json:"protocol"`
	// `control node`, or the service as namespace/name.
	Target string `json:"target"`
	// One of: active.
	Status string `json:"status"`
}

// KubernetesRoutingSynced is the API's kubernetes routing synced object.
type KubernetesRoutingSynced struct {
	Success   bool   `json:"success"`
	ClusterID string `json:"cluster_id"`
	PublicIP  string `json:"public_ip"`
	// Null for a private endpoint.
	APIPort     *int                   `json:"api_port"`
	SyncedPorts []KubernetesSyncedPort `json:"synced_ports"`
	Message     string                 `json:"message"`
}

// CloudFailure is the API's cloud failure object.
// Why the last check or scan failed; null when it didn't.
type CloudFailure struct {
	// One of: assume_role_denied, permission_missing, provider_unavailable, scan_too_large, internal,
	// agent_offline, token_rejected, invalid_upload.
	Code string `json:"code"`
	// What to do about it.
	Message string `json:"message"`
}

// CloudConnectionSetup is the API's cloud connection setup object.
// What you run in your own cloud to connect it. Every setup grants read-only
// access, listed one permission at a time in Permissions; none creates a key
// Pantech holds.
type CloudConnectionSetup struct {
	// `federated`: Pantech reads your cloud with 5-minute tokens it signs for this connection only
	// (AWS, Google Cloud, Azure, Huawei Cloud). `agent`: you run the Pantech Migrate agent, which
	// keeps your read-only token on your machine. One of: federated, agent.
	Access string `json:"access"`
	// Pantech's identity provider URL, which your cloud trusts; nil for agent providers.
	IssuerURL *string `json:"issuer_url"`
	// The audience every token carries; nil for agent providers.
	Audience *string `json:"audience"`
	// The subject your trust policy requires: this connection's id.
	Subject string `json:"subject"`
	// Opens AWS's create-stack page with this template and connection filled in; null when not offered.
	QuickCreateURL *string `json:"quick_create_url"`
	// The main setup file: AWS's CloudFormation template (JSON), Google Cloud's gcloud script,
	// Azure's Azure CLI script, Huawei Cloud's steps, or, for agent providers, how to make a
	// read-only token.
	Template string `json:"template"`
	// One of: cloudformation, shell, markdown, agent.
	TemplateFormat string `json:"template_format"`
	// Further setup files, such as a Terraform version of the setup or JSON documents to paste.
	Files []CloudSetupFile `json:"files"`
	// Every read-only permission the setup grants.
	Permissions []string `json:"permissions"`
	// How to run the Pantech Migrate agent; nil for federated providers.
	Agent *CloudAgentSetup `json:"agent"`
}

// CloudSetupFile is one further file of a connection's setup.
type CloudSetupFile struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// CloudAgentSetup is the API's cloud agent setup object.
// How to run the Pantech Migrate agent; null for federated providers.
type CloudAgentSetup struct {
	// The agent's container image.
	Image string `json:"image"`
	// Where the agent's binaries are published.
	DownloadURL *string `json:"download_url"`
	APIURL      string  `json:"api_url"`
	// The environment variables the agent reads your read-only credential from, on your machine.
	// It is never sent to Pantech.
	TokenVariables []string `json:"token_variables"`
	TokenHowTo     string   `json:"token_how_to"`
	// Contains `<ENROLLMENT_CODE>` where the one-time code goes.
	EnrollCommand    string `json:"enroll_command"`
	RunCommand       string `json:"run_command"`
	ContainerCommand string `json:"container_command"`
}

// CloudAgentStatus is the API's cloud agent status object.
// The connection's enrolled agent; null when none is.
type CloudAgentStatus struct {
	ID string `json:"id"`
	// Whether it polled in the last 2 minutes.
	Online              bool      `json:"online"`
	AgentVersion        string    `json:"agent_version"`
	Platform            string    `json:"platform"`
	EnrolledAt          time.Time `json:"enrolled_at"`
	LastSeenAt          time.Time `json:"last_seen_at"`
	CredentialExpiresAt time.Time `json:"credential_expires_at"`
}

// CloudGCPFederation is the API's cloud GCP federation object.
// What the Google Cloud setup printed.
type CloudGCPFederation struct {
	ProjectNumber       string `json:"project_number"`
	PoolID              string `json:"pool_id"`
	ProviderID          string `json:"provider_id"`
	ServiceAccountEmail string `json:"service_account_email"`
	// The projects to scan, 1 to 20.
	ProjectIDs []string `json:"project_ids"`
}

// CloudAzureFederation is the API's cloud Azure federation object.
// What the Azure setup printed.
type CloudAzureFederation struct {
	TenantID string `json:"tenant_id"`
	// The app registration's application (client) ID.
	ClientID string `json:"client_id"`
	// The subscriptions to scan, 1 to 20.
	SubscriptionIDs []string `json:"subscription_ids"`
}

// CloudHuaweiFederation is the API's cloud Huawei federation object.
// Your Huawei Cloud identity provider and the projects to scan.
type CloudHuaweiFederation struct {
	IdentityProviderID string `json:"identity_provider_id"`
	// 1 to 20 region projects.
	Projects []CloudHuaweiProject `json:"projects"`
}

// CloudHuaweiProject is one Huawei Cloud region project to scan.
type CloudHuaweiProject struct {
	Region string `json:"region"`
	ID     string `json:"id"`
}

// CloudConnection is the API's cloud connection object.
// A connection to an account in another cloud: read-only, through tokens
// Pantech signs (federated providers) or through the Pantech Migrate agent.
type CloudConnection struct {
	ID string `json:"id"`
	// One of: aws, gcp, azure, huawei, digitalocean, hetzner, linode, vultr, contabo, openstack.
	Provider string `json:"provider"`
	Name     string `json:"name"`
	// `awaiting_role` (federated providers: waiting for what the setup printed), `awaiting_agent`
	// (agent providers: no agent enrolled), `verifying`, `connected` or `failed`.
	// One of: awaiting_role, awaiting_agent, verifying, connected, failed.
	Status  string        `json:"status"`
	Failure *CloudFailure `json:"failure"`
	RoleARN *string       `json:"role_arn"`
	// Set for AWS connections only; see ExternalAccountID.
	AWSAccountID *string `json:"aws_account_id"`
	// The source account: an AWS account, a Google Cloud pool's project number, an Azure tenant, or
	// what the agent reported.
	ExternalAccountID *string `json:"external_account_id"`
	// What the Google Cloud setup printed.
	GCP *CloudGCPFederation `json:"gcp"`
	// What the Azure setup printed.
	Azure *CloudAzureFederation `json:"azure"`
	// Your Huawei Cloud identity provider and the projects to scan.
	Huawei *CloudHuaweiFederation `json:"huawei"`
	// The connection's enrolled agent; nil when none is.
	Agent      *CloudAgentStatus `json:"agent"`
	VerifiedAt *time.Time        `json:"verified_at"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	// What you run in your own cloud to connect it.
	Setup CloudConnectionSetup `json:"setup"`
}

// CreateCloudConnectionRequest is the API's create cloud connection request object.
type CreateCloudConnectionRequest struct {
	// `aws`, the default, or another provider ListProviders lists.
	// One of: aws, gcp, azure, huawei, digitalocean, hetzner, linode, vultr, contabo, openstack.
	Provider *string `json:"provider,omitempty"`
	// 1 to 100 characters, for example `Production AWS`.
	Name string `json:"name"`
}

// UpdateCloudConnectionRequest is the API's update cloud connection request object.
type UpdateCloudConnectionRequest struct {
	Name *string `json:"name,omitempty"`
	// AWS: the stack's RoleArn output.
	RoleARN *string `json:"role_arn,omitempty"`
	// Google Cloud: what the setup printed.
	GCP *CloudGCPFederation `json:"gcp,omitempty"`
	// Azure: what the setup printed.
	Azure *CloudAzureFederation `json:"azure,omitempty"`
	// Huawei Cloud: the identity provider and the projects to scan.
	Huawei *CloudHuaweiFederation `json:"huawei,omitempty"`
	// Check the connection again.
	Retry *bool `json:"retry,omitempty"`
}

// CloudScanCounts is the API's cloud scan counts object.
type CloudScanCounts struct {
	Instances          int `json:"instances"`
	Volumes            int `json:"volumes"`
	Networks           int `json:"networks"`
	Subnets            int `json:"subnets"`
	SecurityGroups     int `json:"security_groups"`
	PublicIPs          int `json:"public_ips"`
	LoadBalancers      int `json:"load_balancers"`
	Databases          int `json:"databases"`
	Buckets            int `json:"buckets"`
	Users              int `json:"users"`
	SSHKeys            int `json:"ssh_keys"`
	KubernetesClusters int `json:"kubernetes_clusters"`
}

// CloudCostMonth is the API's cloud cost month object.
type CloudCostMonth struct {
	// YYYY-MM.
	Month       string `json:"month"`
	AmountMinor int64  `json:"amount_minor"`
}

// CloudCost is the API's cloud cost object.
// The source account's whole bill, every service.
type CloudCost struct {
	Currency string `json:"currency"`
	// Average of the months, in minor units of currency.
	AverageMonthlyMinor int64            `json:"average_monthly_minor"`
	Months              []CloudCostMonth `json:"months"`
}

// CloudInstance is a server a scan found.
type CloudInstance struct {
	ID               string     `json:"id"`
	Name             *string    `json:"name"`
	Region           string     `json:"region"`
	Zone             *string    `json:"zone"`
	InstanceType     string     `json:"instance_type"`
	State            string     `json:"state"`
	VCPU             *int       `json:"vcpu"`
	MemoryMiB        *int       `json:"memory_mib"`
	Architecture     string     `json:"architecture"`
	Platform         string     `json:"platform"`
	VolumeIDs        []string   `json:"volume_ids"`
	NetworkID        *string    `json:"network_id"`
	SubnetID         *string    `json:"subnet_id"`
	SecurityGroupIDs []string   `json:"security_group_ids"`
	PublicIP         *string    `json:"public_ip"`
	PrivateIP        *string    `json:"private_ip"`
	KeyName          *string    `json:"key_name"`
	LaunchedAt       *time.Time `json:"launched_at"`
	OSName           *string    `json:"os_name"`
}

// CloudVolume is a disk a scan found.
type CloudVolume struct {
	ID         string  `json:"id"`
	Region     string  `json:"region"`
	Zone       *string `json:"zone"`
	Type       string  `json:"type"`
	SizeGB     int     `json:"size_gb"`
	AttachedTo *string `json:"attached_to"`
	Device     *string `json:"device"`
	Encrypted  bool    `json:"encrypted"`
}

// CloudNetwork is a network a scan found.
type CloudNetwork struct {
	ID        string  `json:"id"`
	Region    string  `json:"region"`
	Name      *string `json:"name"`
	CIDR      string  `json:"cidr"`
	IsDefault bool    `json:"is_default"`
}

// CloudSubnet is a subnet a scan found.
type CloudSubnet struct {
	ID           string  `json:"id"`
	Region       string  `json:"region"`
	NetworkID    string  `json:"network_id"`
	Name         *string `json:"name"`
	CIDR         string  `json:"cidr"`
	Zone         *string `json:"zone"`
	MapsPublicIP bool    `json:"maps_public_ip"`
}

// CloudSecurityGroupRule is one rule of a firewall group a scan found.
type CloudSecurityGroupRule struct {
	// One of: ingress, egress.
	Direction     string  `json:"direction"`
	Protocol      string  `json:"protocol"`
	FromPort      *int    `json:"from_port"`
	ToPort        *int    `json:"to_port"`
	CIDR          *string `json:"cidr"`
	IPv6CIDR      *string `json:"ipv6_cidr"`
	SourceGroupID *string `json:"source_group_id"`
	PrefixListID  *string `json:"prefix_list_id"`
}

// CloudSecurityGroup is a firewall group a scan found.
type CloudSecurityGroup struct {
	ID          string                   `json:"id"`
	Region      string                   `json:"region"`
	NetworkID   *string                  `json:"network_id"`
	Name        string                   `json:"name"`
	Description *string                  `json:"description"`
	Rules       []CloudSecurityGroupRule `json:"rules"`
}

// CloudPublicIP is a public address a scan found.
type CloudPublicIP struct {
	ID         string  `json:"id"`
	Region     string  `json:"region"`
	Address    string  `json:"address"`
	InstanceID *string `json:"instance_id"`
}

// CloudLoadBalancer is a load balancer a scan found.
type CloudLoadBalancer struct {
	ID        string              `json:"id"`
	Region    string              `json:"region"`
	Name      string              `json:"name"`
	Type      string              `json:"type"`
	Scheme    string              `json:"scheme"`
	NetworkID *string             `json:"network_id"`
	Listeners []MigrationListener `json:"listeners"`
}

// CloudDatabase is a database a scan found.
type CloudDatabase struct {
	ID                 string `json:"id"`
	Region             string `json:"region"`
	Engine             string `json:"engine"`
	EngineVersion      string `json:"engine_version"`
	InstanceClass      string `json:"instance_class"`
	VCPU               *int   `json:"vcpu"`
	MemoryMiB          *int   `json:"memory_mib"`
	StorageGB          int    `json:"storage_gb"`
	MultiAZ            bool   `json:"multi_az"`
	PubliclyAccessible bool   `json:"publicly_accessible"`
}

// CloudBucket is a storage bucket a scan found, by name.
type CloudBucket struct {
	Name      string     `json:"name"`
	Region    *string    `json:"region"`
	CreatedAt *time.Time `json:"created_at"`
}

// CloudUser is a person a scan found.
type CloudUser struct {
	ID        *string    `json:"id"`
	Name      string     `json:"name"`
	Email     *string    `json:"email"`
	IsAdmin   bool       `json:"is_admin"`
	CreatedAt *time.Time `json:"created_at"`
	// Where the person came from: `iam_user`, `identity_center`, `gcp_iam`, `entra`, `team`.
	Source *string `json:"source"`
	// One of: admin, write, read.
	AccessLevel *string `json:"access_level"`
	// The source cloud's own role names.
	Roles []string `json:"roles"`
}

// CloudSSHKey is an SSH key a scan found.
type CloudSSHKey struct {
	Name        string  `json:"name"`
	Region      string  `json:"region"`
	Fingerprint *string `json:"fingerprint"`
	PublicKey   *string `json:"public_key"`
}

// CloudInventory is the API's cloud inventory object.
// Everything a scan found. Names and tags are the account owner's own text.
type CloudInventory struct {
	Instances          []CloudInstance          `json:"instances"`
	Volumes            []CloudVolume            `json:"volumes"`
	Networks           []CloudNetwork           `json:"networks"`
	Subnets            []CloudSubnet            `json:"subnets"`
	SecurityGroups     []CloudSecurityGroup     `json:"security_groups"`
	PublicIPs          []CloudPublicIP          `json:"public_ips"`
	LoadBalancers      []CloudLoadBalancer      `json:"load_balancers"`
	Databases          []CloudDatabase          `json:"databases"`
	Buckets            []CloudBucket            `json:"buckets"`
	Users              []CloudUser              `json:"users"`
	SSHKeys            []CloudSSHKey            `json:"ssh_keys"`
	KubernetesClusters []CloudKubernetesCluster `json:"kubernetes_clusters"`
}

// CloudKubernetesCluster is a Kubernetes cluster a scan found. The cluster's
// certificate authority is never kept.
type CloudKubernetesCluster struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Region string `json:"region"`
	// Who runs the control plane: eks, gke, aks, doks, lke or self_managed.
	Provider       string          `json:"provider"`
	Version        *string         `json:"version"`
	EndpointPublic bool            `json:"endpoint_public"`
	NodePools      []CloudNodePool `json:"node_pools"`
	Addons         []string        `json:"addons"`
	CreatedAt      *time.Time      `json:"created_at"`
}

// CloudNodePool is a node pool of a Kubernetes cluster a scan found.
type CloudNodePool struct {
	Name         string  `json:"name"`
	InstanceType *string `json:"instance_type"`
	VCPU         *int    `json:"vcpu"`
	MemoryMiB    *int    `json:"memory_mib"`
	Architecture string  `json:"architecture"`
	GPU          bool    `json:"gpu"`
	Windows      bool    `json:"windows"`
	Desired      int     `json:"desired"`
	Autoscaling  bool    `json:"autoscaling"`
	Min          int     `json:"min"`
	Max          int     `json:"max"`
}

// CloudScan is the API's cloud scan object.
// A read of everything a connection's role can see.
type CloudScan struct {
	ID           string `json:"id"`
	ConnectionID string `json:"connection_id"`
	// `worker` (Pantech read your cloud) or `agent` (your agent did). One of: worker, agent.
	Runner string `json:"runner"`
	// One of: queued, running, succeeded, failed.
	Status            string        `json:"status"`
	Failure           *CloudFailure `json:"failure"`
	CreatedAt         time.Time     `json:"created_at"`
	StartedAt         *time.Time    `json:"started_at"`
	FinishedAt        *time.Time    `json:"finished_at"`
	Provider          *string       `json:"provider"`
	AWSAccountID      *string       `json:"aws_account_id"`
	ExternalAccountID *string       `json:"external_account_id"`
	// Regions read.
	Regions []string         `json:"regions"`
	Counts  *CloudScanCounts `json:"counts"`
	// The source account's whole bill, every service.
	Cost *CloudCost `json:"cost"`
	// Parts not read, e.g. `cost_unavailable` or `region_denied:eu-south-1`.
	Warnings []string `json:"warnings"`
	// Kinds cut at 2,000 items.
	Truncated []string `json:"truncated"`
	// Everything the scan found; only with GetCloudScanParams{Include: "inventory"}.
	Inventory *CloudInventory `json:"inventory"`
}

// MigrationListener is the API's migration listener object.
type MigrationListener struct {
	Protocol string `json:"protocol"`
	Port     int    `json:"port"`
}

// MigrationRule is the API's migration rule object.
type MigrationRule struct {
	// One of: ingress, egress.
	Direction string `json:"direction"`
	// One of: tcp, udp, icmp, all.
	Protocol  string  `json:"protocol"`
	PortRange *string `json:"port_range"`
	CIDR      string  `json:"cidr"`
}

// MigrationFlag is the API's migration flag object.
type MigrationFlag struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	// A blocking flag makes the item unsupported.
	Blocking bool `json:"blocking"`
}

// MigrationItemSource is the API's migration item source object.
type MigrationItemSource struct {
	ID        string  `json:"id"`
	Name      *string `json:"name"`
	Region    *string `json:"region"`
	Summary   string  `json:"summary"`
	VCPU      *int    `json:"vcpu"`
	MemoryMiB *int    `json:"memory_mib"`
	DiskGB    *int    `json:"disk_gb"`
}

// MigrationItemTarget is the API's migration item target object.
// What the item becomes on Pantech; fields that don't apply to its kind are null or empty.
type MigrationItemTarget struct {
	PlanSlug      *string         `json:"plan_slug"`
	DiskGB        *int            `json:"disk_gb"`
	CIDR          *string         `json:"cidr"`
	NetworkKey    *string         `json:"network_key"`
	SubnetKey     *string         `json:"subnet_key"`
	InstanceKey   *string         `json:"instance_key"`
	Rules         []MigrationRule `json:"rules"`
	Engine        *string         `json:"engine"`
	EngineVersion *string         `json:"engine_version"`
	StorageGB     *int            `json:"storage_gb"`
	// One of: admin, member.
	Role      *string             `json:"role"`
	Email     *string             `json:"email"`
	PublicKey *string             `json:"public_key"`
	Listeners []MigrationListener `json:"listeners"`
	// The image a server is planned on when its operating system decides it (Windows Server).
	ImageSlug *string `json:"image_slug"`
	// The Pantech invite sent to this person from this plan.
	InviteID  *string    `json:"invite_id"`
	InvitedAt *time.Time `json:"invited_at"`
	// The Pantech bucket a bucket item becomes; nil for every other kind.
	Bucket *MigrationBucketTarget `json:"bucket"`
	// What a kubernetes_cluster item becomes on Pantech; nil for every other kind.
	Kubernetes *MigrationKubernetesTarget `json:"kubernetes"`
}

// MigrationBucketTarget is the API's migration bucket target object.
// The Pantech bucket a bucket item becomes.
type MigrationBucketTarget struct {
	// A valid Pantech bucket name made from the source's; nil when none could be made.
	Name       *string `json:"name"`
	Region     string  `json:"region"`
	Versioning bool    `json:"versioning"`
	// The source bucket was public. The Pantech bucket stays private unless the copy asks.
	PublicRead bool `json:"public_read"`
	// The source's own daily metric, read by the scan; nil when unknown.
	SizeBytes *int64 `json:"size_bytes"`
	// The source's own daily metric, read by the scan; nil when unknown.
	ObjectCount *int64 `json:"object_count"`
	// The day (YYYY-MM-DD) the size and object count describe: show them as "as of" this day.
	SizeAsOf *string `json:"size_as_of"`
}

// MigrationKubernetesTarget is the API's migration Kubernetes target object.
// What a kubernetes_cluster item becomes on Pantech.
type MigrationKubernetesTarget struct {
	VersionID *string             `json:"version_id"`
	Version   *string             `json:"version"`
	NodePools []MigrationNodePool `json:"node_pools"`
	Addons    []string            `json:"addons"`
	// The customer consented to a lower Kubernetes version than the source runs.
	AllowDowngrade bool `json:"allow_downgrade"`
}

// MigrationNodePool is a node pool a kubernetes_cluster item plans on Pantech.
type MigrationNodePool struct {
	Name        string  `json:"name"`
	PlanSlug    *string `json:"plan_slug"`
	Count       int     `json:"count"`
	Autoscaling bool    `json:"autoscaling"`
	Min         int     `json:"min"`
	Max         int     `json:"max"`
	SourceName  string  `json:"source_name"`
}

// MigrationItem is the API's migration item object.
type MigrationItem struct {
	// Stable across edits: kind, region and the source id.
	Key string `json:"key"`
	// One of: network, subnet, security_group, instance, volume, public_ip, load_balancer, database, ssh_key, user, bucket,
	// kubernetes_cluster.
	Kind   string              `json:"kind"`
	Source MigrationItemSource `json:"source"`
	// One of: planned, excluded, unsupported.
	Status string          `json:"status"`
	Flags  []MigrationFlag `json:"flags"`
	// What the item becomes on Pantech.
	Target MigrationItemTarget `json:"target"`
}

// MigrationSummary is the API's migration summary object.
type MigrationSummary struct {
	Planned     int `json:"planned"`
	Excluded    int `json:"excluded"`
	Unsupported int `json:"unsupported"`
	// Planned items per kind.
	ByKind map[string]int `json:"by_kind"`
}

// MigrationEstimate is the API's migration estimate object.
// Pantech's monthly price for the plan next to the source's monthly bill, computed when read.
type MigrationEstimate struct {
	Currency string `json:"currency"`
	// False when part of the plan could not be priced; then PantechMonthlyMinor is nil.
	PantechPriced       bool   `json:"pantech_priced"`
	PantechMonthlyMinor *int64 `json:"pantech_monthly_minor"`
	// What the Pantech price includes: servers, volumes, networks.
	Covers []string `json:"covers"`
	// Planned items the price leaves out, e.g. databases.
	NotCovered         []string `json:"not_covered"`
	SourceKnown        bool     `json:"source_known"`
	SourceMonthlyMinor *int64   `json:"source_monthly_minor"`
	SourceCurrency     *string  `json:"source_currency"`
	SourceMonths       int      `json:"source_months"`
	// `whole_account`: the source figure is the whole cloud bill. One of: whole_account.
	SourceScope string `json:"source_scope"`
}

// Migration is the API's migration object.
// A plan to move what a scan found to Pantech.
type Migration struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Region       string `json:"region"`
	ConnectionID string `json:"connection_id"`
	ScanID       string `json:"scan_id"`
	// `planning`. Copying arrives in a later release. One of: planning.
	Stage   string           `json:"stage"`
	Version int              `json:"version"`
	Summary MigrationSummary `json:"summary"`
	// Nil in a list.
	Estimate *MigrationEstimate `json:"estimate"`
	// Nil in a list.
	Items     []MigrationItem `json:"items"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// CreateMigrationRequest is the API's create migration request object.
type CreateMigrationRequest struct {
	ScanID string `json:"scan_id"`
	Name   string `json:"name"`
	// A Pantech region code; the first region when left out.
	Region *string `json:"region,omitempty"`
}

// RenameMigrationRequest is the API's rename migration request object.
type RenameMigrationRequest struct {
	Name string `json:"name"`
}

// UpdateMigrationItemRequest is the API's update migration item request object.
type UpdateMigrationItemRequest struct {
	// The item's key.
	Key       string  `json:"key"`
	Excluded  *bool   `json:"excluded,omitempty"`
	PlanSlug  *string `json:"plan_slug,omitempty"`
	StorageGB *int    `json:"storage_gb,omitempty"`
	CIDR      *string `json:"cidr,omitempty"`
	// One of: member, admin.
	Role  *string `json:"role,omitempty"`
	Email *string `json:"email,omitempty"`
	// On a kubernetes_cluster item: consent to a lower Kubernetes version than the source runs.
	AllowVersionDowngrade *bool `json:"allow_version_downgrade,omitempty"`
}

// MigrationTerraformFile is one file of a Terraform export.
type MigrationTerraformFile struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// MigrationTerraform is the API's migration terraform object.
type MigrationTerraform struct {
	Files []MigrationTerraformFile `json:"files"`
}

// CloudProvider is the API's cloud provider object: a cloud Pantech Migrate
// can connect.
type CloudProvider struct {
	// One of: aws, gcp, azure, huawei, digitalocean, hetzner, linode, vultr, contabo, openstack.
	Slug string `json:"slug"`
	Name string `json:"name"`
	// One of: federated, agent.
	Access string `json:"access"`
	// Whether a scan reads the source's spend. When false the comparison says the spend is not available.
	CostReadable bool `json:"cost_readable"`
	// Whether a scan lists the people with access. When false, add them to the plan by hand.
	PeopleReadable bool `json:"people_readable"`
	// Whether this deployment can connect it now.
	Available bool `json:"available"`
}

// CloudAgentEnrollment is the API's cloud agent enrollment object.
type CloudAgentEnrollment struct {
	// The one-time enrollment code, shown only now. It works once, for an hour.
	Code          string    `json:"code"`
	ExpiresAt     time.Time `json:"expires_at"`
	EnrollCommand string    `json:"enroll_command"`
}

// SendMigrationInvitesRequest is the API's send migration invites request object.
type SendMigrationInvitesRequest struct {
	// The `user` items to invite, by key: 1 to 50.
	Keys []string `json:"keys"`
}

// MigrationInviteResult is what happened to one person of a plan.
type MigrationInviteResult struct {
	Key   string  `json:"key"`
	Email *string `json:"email"`
	Role  *string `json:"role"`
	// `invited`; `already_invited` (from this plan or by hand); `already_member`; `role_not_allowed`
	// (an admin can't invite an admin: change the person's role to member or ask an owner);
	// `skipped` (see Reason); `failed` (try again).
	// One of: invited, already_invited, already_member, role_not_allowed, skipped, failed.
	Outcome string `json:"outcome"`
	// Why a person was skipped: `email_missing`, `excluded`, `not_a_person`, `unknown_item`.
	Reason   *string `json:"reason"`
	InviteID *string `json:"invite_id"`
}

// MigrationInvites is the API's migration invites object: what happened to
// each person, and the plan.
type MigrationInvites struct {
	Results   []MigrationInviteResult `json:"results"`
	Migration Migration               `json:"migration"`
}

// BucketImportFailure is the API's bucket import failure object.
// Why it failed; null when it didn't. Message says what to do.
type BucketImportFailure struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// CopyAccessRequest is the API's copy access request object.
type CopyAccessRequest struct {
	// The source buckets Pantech may copy, 1 to 50, each once. Required the first time.
	Buckets []string `json:"buckets,omitempty"`
	// The copy stack's `CopyRoleArn` output. Setting it queues a verification.
	RoleARN *string `json:"role_arn,omitempty"`
	// Verify the role again.
	Retry *bool `json:"retry,omitempty"`
}

// CopyAccessBucket is a bucket chosen for copying.
type CopyAccessBucket struct {
	Name   string  `json:"name"`
	Region *string `json:"region"`
	// One of: enabled, suspended, off.
	Versioning *string `json:"versioning"`
	// Verified: an import may copy this bucket.
	Ready   bool                 `json:"ready"`
	Failure *BucketImportFailure `json:"failure"`
}

// CopyAccessSetup is what the copy role trusts, and its template.
type CopyAccessSetup struct {
	IssuerURL string `json:"issuer_url"`
	Audience  string `json:"audience"`
	// The copy role's token subject: the connection id followed by `:copy`.
	Subject string `json:"subject"`
	// The CloudFormation template (JSON text) for a role that may read only the chosen buckets;
	// nil until buckets are chosen.
	Template *string `json:"template"`
}

// CopyAccess is the API's copy access object: a connection's opt-in access
// for copying the chosen buckets, separate from its read-only role.
type CopyAccess struct {
	ConnectionID string `json:"connection_id"`
	// One of: not_requested, awaiting_role, verifying, verified, failed.
	Status     string               `json:"status"`
	Failure    *BucketImportFailure `json:"failure"`
	RoleARN    *string              `json:"role_arn"`
	Buckets    []CopyAccessBucket   `json:"buckets"`
	VerifiedAt *time.Time           `json:"verified_at"`
	Setup      CopyAccessSetup      `json:"setup"`
}

// ResyncBucketImportRequest is the API's resync bucket import request object.
type ResyncBucketImportRequest struct {
	// Off unless set. After this pass has listed the whole source, delete the objects this import
	// copied that the source no longer has, and only those still holding the bytes the import
	// wrote. Objects the import did not write are never deleted. Needs the delete scope.
	MirrorDeletes *bool `json:"mirror_deletes,omitempty"`
}

// StartBucketImportRequest is the API's start bucket import request object.
// Name exactly one of DestinationBucketName and DestinationBucketID.
type StartBucketImportRequest struct {
	ConnectionID string `json:"connection_id"`
	// The source bucket's name. It must have verified copy access.
	SourceBucket string `json:"source_bucket"`
	// A new Pantech bucket to make and copy into.
	DestinationBucketName *string `json:"destination_bucket_name,omitempty"`
	// An existing Pantech bucket of the organization to copy into.
	DestinationBucketID *string `json:"destination_bucket_id,omitempty"`
	// Versioning on a new bucket.
	Versioning *bool `json:"versioning,omitempty"`
	// Make a new bucket publicly readable. Off unless asked, even when the source was public.
	PublicRead *bool `json:"public_read,omitempty"`
	// The plan the move belongs to; its size for the bucket is checked against the storage left.
	MigrationID *string `json:"migration_id,omitempty"`
}

// BucketImportProgress is a bucket import's progress in its current pass.
type BucketImportProgress struct {
	ObjectsCopied int64 `json:"objects_copied"`
	// Already in the Pantech bucket unchanged.
	ObjectsSkipped int64 `json:"objects_skipped"`
	ObjectsFailed  int64 `json:"objects_failed"`
	BytesCopied    int64 `json:"bytes_copied"`
	BytesSkipped   int64 `json:"bytes_skipped"`
}

// BucketImportObjectFailure is an object that failed a bucket import's pass.
type BucketImportObjectFailure struct {
	Key string `json:"key"`
	// One of: object_archived, object_access_denied, invalid_key, object_too_large, source_changed,
	// verification_failed, transfer_failed, destination_rejected.
	Reason  string    `json:"reason"`
	Message string    `json:"message"`
	At      time.Time `json:"at"`
}

// BucketImportEvent is one entry of a bucket import's audit trail.
type BucketImportEvent struct {
	// One of: started, cancel_requested, resync_requested, succeeded, failed, cancelled,
	// mirror_deleted, mirror_skipped.
	Action string `json:"action"`
	// Nil for Pantech's own steps.
	ActorUserID *string   `json:"actor_user_id"`
	Detail      *string   `json:"detail"`
	OccurredAt  time.Time `json:"occurred_at"`
}

// BucketImportSource is the bucket a bucket import copies from.
type BucketImportSource struct {
	Provider string  `json:"provider"`
	Bucket   string  `json:"bucket"`
	Region   *string `json:"region"`
}

// BucketImportDestination is the Pantech bucket a bucket import copies into.
type BucketImportDestination struct {
	BucketID   string `json:"bucket_id"`
	BucketName string `json:"bucket_name"`
	PublicRead bool   `json:"public_read"`
}

// BucketImport is the API's bucket import object: one source bucket copied
// into a Pantech bucket.
type BucketImport struct {
	ID           string                  `json:"id"`
	ConnectionID string                  `json:"connection_id"`
	MigrationID  *string                 `json:"migration_id"`
	Source       BucketImportSource      `json:"source"`
	Destination  BucketImportDestination `json:"destination"`
	// One of: queued, running, succeeded, failed, cancelled.
	Status          string               `json:"status"`
	Failure         *BucketImportFailure `json:"failure"`
	CancelRequested bool                 `json:"cancel_requested"`
	// 1 for the first copy; each resync adds one.
	Pass     int                  `json:"pass"`
	Progress BucketImportProgress `json:"progress"`
	// Bytes copied over every pass.
	TotalBytesCopied int64 `json:"total_bytes_copied"`
	// The current pass was asked to delete what this import copied and the source no longer has.
	MirrorDeletes bool `json:"mirror_deletes"`
	// Objects deleted by mirror passes, over every pass.
	ObjectsDeleted int64 `json:"objects_deleted"`
	// When a queued import waiting on an unreachable side is tried again.
	NextAttemptAt *time.Time `json:"next_attempt_at"`
	CreatedAt     time.Time  `json:"created_at"`
	StartedAt     *time.Time `json:"started_at"`
	FinishedAt    *time.Time `json:"finished_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	// The first 100 objects that failed this pass; nil in a list.
	Failures []BucketImportObjectFailure `json:"failures"`
	// The audit trail; nil in a list.
	Events []BucketImportEvent `json:"events"`
}

// ServerCopyAWSSource is how a server copy reads AWS EBS snapshots.
type ServerCopyAWSSource struct {
	// The copy stack's CopyRoleArn output, never the discover role.
	CopyRoleARN string `json:"copy_role_arn"`
	// The copy role may snapshot the server's volumes itself.
	CreateSnapshots *bool `json:"create_snapshots,omitempty"`
	// The first pass's snapshot for each volume id, tagged pantech-migrate-connection.
	Snapshots map[string]string `json:"snapshots,omitempty"`
}

// StartServerCopyRequest is the API's start server copy request object.
type StartServerCopyRequest struct {
	// The planned server's item key.
	ItemKey string `json:"item_key"`
	// `agent` (the default) or `aws_snapshot`. One of: agent, aws_snapshot.
	Source *string `json:"source,omitempty"`
	// The new instance's name; the source's when empty.
	Name *string `json:"name,omitempty"`
	// Overrides the item's plan.
	PlanSlug        *string `json:"plan_slug,omitempty"`
	SubnetID        *string `json:"subnet_id,omitempty"`
	SecurityGroupID *string `json:"security_group_id,omitempty"`
	// Optional. The copied disk keeps its own users and keys.
	SSHKeyID *string `json:"ssh_key_id,omitempty"`
	// For `aws_snapshot`.
	AWS *ServerCopyAWSSource `json:"aws,omitempty"`
}

// ServerCopyActRequest is the API's server copy act request object, sent to
// sync, cut over or cancel a server copy.
type ServerCopyActRequest struct {
	// The next pass's snapshot for each volume id, for an AWS copy that makes none itself.
	Snapshots map[string]string `json:"snapshots,omitempty"`
}

// ServerCopyTemplate is the API's server copy template object: the AWS copy
// role for one planned server.
type ServerCopyTemplate struct {
	// CloudFormation JSON.
	Template    string            `json:"template"`
	SnapshotTag map[string]string `json:"snapshot_tag"`
}

// ServerCopyFailure is why a server copy failed.
type ServerCopyFailure struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ServerCopyInstance is the Pantech instance a server copy builds.
type ServerCopyInstance struct {
	Name            string  `json:"name"`
	PlanSlug        string  `json:"plan_slug"`
	Region          string  `json:"region"`
	SubnetID        *string `json:"subnet_id"`
	SecurityGroupID *string `json:"security_group_id"`
	SSHKeyID        *string `json:"ssh_key_id"`
	// Set once the instance is being created.
	ID      *string `json:"id"`
	OrderID *string `json:"order_id"`
}

// ServerCopyDisk is one disk a server copy copies.
type ServerCopyDisk struct {
	Index int `json:"index"`
	// One of: root, data.
	Role       string  `json:"role"`
	Source     string  `json:"source"`
	SizeBytes  int64   `json:"size_bytes"`
	SnapshotID *string `json:"snapshot_id"`
	VolumeID   *string `json:"volume_id"`
}

// ServerCopyFlag is a pre-flight finding of a server copy.
type ServerCopyFlag struct {
	Code     string `json:"code"`
	Blocking bool   `json:"blocking"`
	Message  string `json:"message"`
}

// ServerCopyPass is the pass a server copy is running.
type ServerCopyPass struct {
	Number int `json:"number"`
	// One of: full, delta, final.
	Kind      string     `json:"kind"`
	StartedAt *time.Time `json:"started_at"`
}

// ServerCopyProgress is a server copy's progress in its pass.
type ServerCopyProgress struct {
	ChunksTotal   int64  `json:"chunks_total"`
	ChunksDone    int64  `json:"chunks_done"`
	BytesTotal    int64  `json:"bytes_total"`
	BytesUploaded int64  `json:"bytes_uploaded"`
	WireBytes     int64  `json:"wire_bytes"`
	ETASeconds    *int64 `json:"eta_seconds"`
}

// ServerCopy is the API's server copy object: one planned server's disks
// copied into a new Pantech instance.
type ServerCopy struct {
	// `sco_…`.
	ID           string `json:"id"`
	MigrationID  string `json:"migration_id"`
	ItemKey      string `json:"item_key"`
	ConnectionID string `json:"connection_id"`
	// One of: agent, aws_snapshot.
	Source string `json:"source"`
	// One of: waiting, copying, synced, importing, creating_instance, completed, failed, cancelled.
	Status  string             `json:"status"`
	Failure *ServerCopyFailure `json:"failure"`
	// One of: linux, windows.
	OSFamily string `json:"os_family"`
	// One of: bios, uefi.
	BootMode *string            `json:"boot_mode"`
	Instance ServerCopyInstance `json:"instance"`
	Disks    []ServerCopyDisk   `json:"disks"`
	Flags    []ServerCopyFlag   `json:"flags"`
	Pass     *ServerCopyPass    `json:"pass"`
	// One of: full, delta, final.
	PendingPass *string            `json:"pending_pass"`
	Progress    ServerCopyProgress `json:"progress"`
	SyncedAt    *time.Time         `json:"synced_at"`
	ExpiresAt   time.Time          `json:"expires_at"`
	CompletedAt *time.Time         `json:"completed_at"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// StartDatabaseImportRequest is the API's start database import request object.
type StartDatabaseImportRequest struct {
	// The plan item's key: one of the migration's `database` items.
	ItemKey string `json:"item_key"`
	// One of your Pantech databases with the same engine as the source. Create it first: it bills
	// as usual; the move itself is free.
	DestinationDatabaseID string `json:"destination_database_id"`
	// The database to create and fill on the destination; it must be empty. Defaults to the
	// source's name, made a valid name: 1 to 63 lowercase letters, digits and underscores, starting
	// with a letter.
	TargetDatabase *string `json:"target_database,omitempty"`
	// `dump` copies once (stop writes to the source first). `replicate` keeps following the
	// source's changes until you cut over (PostgreSQL sources only). One of: dump, replicate.
	Mode *string `json:"mode,omitempty"`
	// The source's name in your Migrate agent's configuration, 1 to 64 of letters, digits, `.`, `_`
	// and `-`. Its address and login stay on the agent and are never sent to Pantech.
	SourceRef string `json:"source_ref"`
}

// DatabaseImportFailure is the API's database import failure object.
// Why the import stopped; null while it has not.
type DatabaseImportFailure struct {
	// `statement_refused`, `chunk_invalid`, `access_revoked`, `destination_unavailable`, or one the
	// agent reported: `source_unreachable`, `source_auth_failed`, `source_permission_denied`,
	// `source_not_configured`, `dump_tool_missing`, `source_changed`, `internal`.
	Code    string `json:"code"`
	Message string `json:"message"`
}

// DatabaseImportIssue is one finding of a database import's pre-flight.
type DatabaseImportIssue struct {
	// For example `source_version_newer`, `extension_unavailable`, `size_exceeds_storage`,
	// `size_not_measured`, `table_without_primary_key`, `replication_check_failed`,
	// `mode_not_available`, `target_not_empty`, `destination_not_ready`,
	// `destination_not_supported`.
	Code     string `json:"code"`
	Blocking bool   `json:"blocking"`
	Message  string `json:"message"`
	// A copy-paste fix for the source's provider (an RDS parameter group setting, a Cloud SQL flag,
	// an Azure server parameter, SQL).
	Fix *string `json:"fix"`
	// What the issue is about (an extension, tables, a check).
	Subject *string `json:"subject"`
}

// DatabaseImportSource is what the Migrate agent measured on the source.
// Numbers it could not measure are nil, never estimated.
type DatabaseImportSource struct {
	Engine  string  `json:"engine"`
	Version *string `json:"version"`
	// `aws_rds`, `gcp_cloudsql`, `azure`, `digitalocean` or `self_hosted`, as the agent detected it.
	Provider   *string  `json:"provider"`
	Database   *string  `json:"database"`
	SizeBytes  *int64   `json:"size_bytes"`
	Extensions []string `json:"extensions"`
	TableCount int      `json:"table_count"`
}

// DatabaseImportReadiness is the API's database import readiness object.
// The latest pre-flight; null until the agent ran one.
type DatabaseImportReadiness struct {
	CheckedAt time.Time `json:"checked_at"`
	// A blocking issue stops the import until it is fixed and readiness is checked again.
	Blocked bool                  `json:"blocked"`
	Issues  []DatabaseImportIssue `json:"issues"`
	Source  DatabaseImportSource  `json:"source"`
}

// DatabaseImportDestination is the API's database import destination object.
// What your Pantech database reported over the import's own connection; null
// in a list. Numbers are nil until measured.
type DatabaseImportDestination struct {
	// The import login's state on the destination.
	// One of: not_prepared, requested, active, revoking, revoked, failed.
	AccessState        string  `json:"access_state"`
	Measured           bool    `json:"measured"`
	Engine             *string `json:"engine"`
	ServerVersion      *string `json:"server_version"`
	DiskAvailableBytes *int64  `json:"disk_available_bytes"`
	DiskTotalBytes     *int64  `json:"disk_total_bytes"`
	// Tables already in the target database; an import needs an empty one.
	TargetTables      *int     `json:"target_tables"`
	TrustedExtensions []string `json:"trusted_extensions"`
}

// DatabaseImportVerification is the API's database import verification object.
// The table's comparison; null until it was verified.
type DatabaseImportVerification struct {
	SourceRows      int64 `json:"source_rows"`
	DestinationRows int64 `json:"destination_rows"`
	// Key ranges whose rows were checksummed on both sides.
	SamplesChecked    int  `json:"samples_checked"`
	MismatchedSamples int  `json:"mismatched_samples"`
	Passed            bool `json:"passed"`
}

// DatabaseImportTable is one table's progress in a database import.
type DatabaseImportTable struct {
	Schema        *string `json:"schema"`
	Name          string  `json:"name"`
	HasPrimaryKey bool    `json:"has_primary_key"`
	// One of: pending, copying, copied, verified, verify_failed.
	Status       string                      `json:"status"`
	PagesCopied  int64                       `json:"pages_copied"`
	RowsCopied   int64                       `json:"rows_copied"`
	BytesCopied  int64                       `json:"bytes_copied"`
	Verification *DatabaseImportVerification `json:"verification"`
	UpdatedAt    time.Time                   `json:"updated_at"`
}

// DatabaseImport is the API's database import object: one of a plan's
// databases moved into a Pantech database through the Migrate agent.
type DatabaseImport struct {
	ID                    string `json:"id"`
	MigrationID           string `json:"migration_id"`
	ItemKey               string `json:"item_key"`
	DestinationDatabaseID string `json:"destination_database_id"`
	TargetDatabase        string `json:"target_database"`
	// One of: dump, replicate.
	Mode      string `json:"mode"`
	SourceRef string `json:"source_ref"`
	// The source's engine family. One of: postgres, mysql, mariadb.
	Engine string `json:"engine"`
	// One of: awaiting_agent, preflight, ready, blocked, copying, replicating, cutting_over,
	// verifying, verifying_failed, completed, cancelling, cancelled, failed.
	Status           string                 `json:"status"`
	Failure          *DatabaseImportFailure `json:"failure"`
	CutoverRequested bool                   `json:"cutover_requested"`
	CancelRequested  bool                   `json:"cancel_requested"`
	// A new pre-flight was asked for and the agent has not run it yet.
	ReadinessRequested bool `json:"readiness_requested"`
	// The import login and its loopback-only access were removed from the destination.
	AccessRevoked bool `json:"access_revoked"`
	// A Migrate agent holds this import now.
	AgentConnected bool `json:"agent_connected"`
	// Replication lag in bytes of source WAL not yet applied; nil in a dump or before replication
	// starts.
	LagBytes *int64 `json:"lag_bytes"`
	// The last applied source position (replicate only).
	Position    *string                    `json:"position"`
	Readiness   *DatabaseImportReadiness   `json:"readiness"`
	Destination *DatabaseImportDestination `json:"destination"`
	// Per-table progress; nil in a list.
	Tables []DatabaseImportTable `json:"tables"`
	// Always false: moving data is free. The destination database bills as usual.
	Charged    bool       `json:"charged"`
	Version    int        `json:"version"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	FinishedAt *time.Time `json:"finished_at"`
}

// KubernetesMoveFailure is the API's Kubernetes move failure object.
// Why the move, or its last step, failed; null when nothing failed. Message
// is safe to show.
type KubernetesMoveFailure struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// KubernetesMoveFlag is a finding of a Kubernetes move.
type KubernetesMoveFlag struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// KubernetesMoveRef names one Kubernetes object.
type KubernetesMoveRef struct {
	// API group; empty for the core API.
	Group     string `json:"group"`
	Kind      string `json:"kind"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

// KubernetesMoveOptions are what a Kubernetes move was started with.
type KubernetesMoveOptions struct {
	IncludeSecrets       bool     `json:"include_secrets"`
	CopyVolumes          bool     `json:"copy_volumes"`
	PublishLoadBalancers bool     `json:"publish_load_balancers"`
	Namespaces           []string `json:"namespaces"`
}

// KubernetesMoveSummary counts what a Kubernetes move found.
type KubernetesMoveSummary struct {
	Objects  int                  `json:"objects"`
	Refused  int                  `json:"refused"`
	Skipped  int                  `json:"skipped"`
	Flagged  int                  `json:"flagged"`
	Volumes  int                  `json:"volumes"`
	Warnings []KubernetesMoveFlag `json:"warnings"`
	Addons   []string             `json:"addons"`
}

// KubernetesMove is the API's Kubernetes move object: a plan's
// kubernetes_cluster item moved into one of your Pantech clusters.
type KubernetesMove struct {
	// `kmv_…`.
	ID           string `json:"id"`
	MigrationID  string `json:"migration_id"`
	ItemKey      string `json:"item_key"`
	ConnectionID string `json:"connection_id"`
	// The Pantech cluster the move applies to: always one of your organization's.
	ClusterID  string `json:"cluster_id"`
	SourceName string `json:"source_name"`
	// One of: exporting, ready, previewing, applying, applied, syncing, cutting_over, completed,
	// failed, cancelled.
	Status          string                 `json:"status"`
	Failure         *KubernetesMoveFailure `json:"failure"`
	Options         KubernetesMoveOptions  `json:"options"`
	Summary         KubernetesMoveSummary  `json:"summary"`
	AllowPrivileged []KubernetesMoveRef    `json:"allow_privileged"`
	PreviewedAt     *time.Time             `json:"previewed_at"`
	PreviewOK       bool                   `json:"preview_ok"`
	AppliedAt       *time.Time             `json:"applied_at"`
	CompletedAt     *time.Time             `json:"completed_at"`
	SyncGeneration  int                    `json:"sync_generation"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// KubernetesMoveStep is the API's Kubernetes move step object.
// The object's outcome in the last preview or apply; null when it has not run.
type KubernetesMoveStep struct {
	// One of: created, configured, unchanged, refused, unverified.
	Result  string  `json:"result"`
	Code    *string `json:"code"`
	Message *string `json:"message"`
}

// KubernetesMoveResource is one source object and what the move does with it
// (names only, never values).
type KubernetesMoveResource struct {
	Group      string `json:"group"`
	Kind       string `json:"kind"`
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
	APIVersion string `json:"api_version"`
	// Apply order: definitions first, custom resources last.
	Tier int `json:"tier"`
	// One of: create, refused, skipped.
	Action string               `json:"action"`
	Flags  []KubernetesMoveFlag `json:"flags"`
	// Refused for the cluster policy; name it in AllowPrivileged on preview to include it anyway.
	PolicyRefused bool                `json:"policy_refused"`
	Preview       *KubernetesMoveStep `json:"preview"`
	Apply         *KubernetesMoveStep `json:"apply"`
}

// KubernetesMoveVolume is one volume a Kubernetes move copies.
type KubernetesMoveVolume struct {
	ID                 string   `json:"id"`
	Namespace          string   `json:"namespace"`
	Claim              string   `json:"claim"`
	SizeBytes          int64    `json:"size_bytes"`
	SourceStorageClass string   `json:"source_storage_class"`
	TargetStorageClass string   `json:"target_storage_class"`
	AccessModes        []string `json:"access_modes"`
	// One of: pending, copying, synced, failed, unsupported.
	Status           string               `json:"status"`
	Flags            []KubernetesMoveFlag `json:"flags"`
	Generation       int                  `json:"generation"`
	SyncedGeneration int                  `json:"synced_generation"`
	BytesTotal       int64                `json:"bytes_total"`
	BytesDone        int64                `json:"bytes_done"`
	FilesTotal       int64                `json:"files_total"`
	// The destination's own read-back matched the source's checksum.
	ChecksumVerified bool                   `json:"checksum_verified"`
	SyncedAt         *time.Time             `json:"synced_at"`
	Failure          *KubernetesMoveFailure `json:"failure"`
}

// KubernetesMoveImage is a container image a Kubernetes move found.
type KubernetesMoveImage struct {
	Image string               `json:"image"`
	Flags []KubernetesMoveFlag `json:"flags"`
}

// KubernetesMoveHold is a workload applied scaled down (or suspended) until
// the cutover, and what the cutover returns it to.
type KubernetesMoveHold struct {
	Group     string `json:"group"`
	Kind      string `json:"kind"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Replicas  int    `json:"replicas"`
	Suspend   bool   `json:"suspend"`
}

// KubernetesMoveDetail is a Kubernetes move with its report.
type KubernetesMoveDetail struct {
	KubernetesMove
	Resources []KubernetesMoveResource `json:"resources"`
	Volumes   []KubernetesMoveVolume   `json:"volumes"`
	Images    []KubernetesMoveImage    `json:"images"`
	Holds     []KubernetesMoveHold     `json:"holds"`
}

// KubernetesMoveSourceAccess is the manifests to apply to the source cluster.
type KubernetesMoveSourceAccess struct {
	// Read-only access (get, list, watch; no Secrets).
	Reader string `json:"reader"`
	// Adds reading Secrets. Apply only to move Secrets.
	SecretsReader string `json:"secrets_reader"`
	// Lets the agent start copy pods. Bind it in each namespace whose volumes you copy.
	Copier string `json:"copier"`
}

// StartKubernetesMoveRequest is the API's start Kubernetes move request object.
type StartKubernetesMoveRequest struct {
	// A `kubernetes_cluster` item of the plan.
	ItemKey string `json:"item_key"`
	// A running Pantech cluster of your organization.
	ClusterID string `json:"cluster_id"`
	// Default false.
	IncludeSecrets *bool `json:"include_secrets,omitempty"`
	// Default false.
	CopyVolumes *bool `json:"copy_volumes,omitempty"`
	// Default true.
	PublishLoadBalancers *bool    `json:"publish_load_balancers,omitempty"`
	Namespaces           []string `json:"namespaces,omitempty"`
}

// PreviewKubernetesMoveRequest is the API's preview Kubernetes move request object.
type PreviewKubernetesMoveRequest struct {
	// Workloads the cluster policy refused that you choose to include anyway.
	AllowPrivileged []KubernetesMoveRef `json:"allow_privileged,omitempty"`
}

// CutoverKubernetesMoveRequest is the API's cutover Kubernetes move request object.
type CutoverKubernetesMoveRequest struct {
	// Must be true: writes on the source have stopped.
	SourceWritesStopped bool `json:"source_writes_stopped"`
}

// DNSZone is the API's DNS zone object.
// A hosted zone on Pantech DNS, on Pantech's own name servers. It is created
// on them once its delegation is proven; record changes then reach them
// within about a minute.
type DNSZone struct {
	// `dnsz_…`.
	ID string `json:"id"`
	// The zone's name, lower case and punycode, such as `example.com`.
	Name        string `json:"name"`
	Description string `json:"description"`
	// One of: pending, awaiting_delegation, active, suspended, deleting, deleted, failed.
	Status string `json:"status"`
	// Set when `status` is `failed`; `provider_limit` also covers no free
	// name-server pair. One of: delegation_timeout, provider_rejected, provider_limit.
	FailureCode *string `json:"failure_code"`
	// The two name servers to set at the domain's registrar, such as
	// `ana-ns.pantechdynamics.com`. Empty while `pending`.
	NameServers []string `json:"name_servers"`
	// The Pantech domain (`dom_…`) that delegates to this zone, if any.
	DomainID    *string `json:"domain_id"`
	RecordCount int     `json:"record_count"`
	// While the zone awaits delegation, when it is removed if the registrar
	// still does not delegate to its name servers; nil otherwise.
	DelegationDeadlineAt *time.Time `json:"delegation_deadline_at"`
	LastCheckedAt        *time.Time `json:"last_checked_at"`
	ActivatedAt          *time.Time `json:"activated_at"`
	SuspendedReason      *string    `json:"suspended_reason"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	Version              int        `json:"version"`
}

// CreateDNSZoneRequest is the API's create DNS zone request object.
type CreateDNSZoneRequest struct {
	// A registrable domain or a subdomain of one, such as `example.com`.
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// UpdateDNSZoneRequest is the API's update DNS zone request object.
type UpdateDNSZoneRequest struct {
	// The new description; empty clears it.
	Description string `json:"description"`
}

// DNSRecord is the API's DNS record set object: one name and type, its TTL
// and its values.
type DNSRecord struct {
	// `dnsr_…`.
	ID     string `json:"id"`
	ZoneID string `json:"zone_id"`
	// Relative to the zone: `@` for the apex, `www`, `*.dev`, `_dmarc`.
	Name string `json:"name"`
	FQDN string `json:"fqdn"`
	// One of: A, AAAA, CNAME, MX, TXT, SRV, CAA, NS, PTR, SOA.
	Type string `json:"type"`
	// Seconds, 60 to 86400.
	TTL int `json:"ttl"`
	// Host names are returned with a trailing dot.
	Values       []string `json:"values"`
	AllowPrivate bool     `json:"allow_private"`
	Comment      string   `json:"comment"`
	// The apex SOA and NS: listed, never changed.
	Managed bool `json:"managed"`
	// One of: pending, active, deleting, failed.
	Status string `json:"status"`
	// Set when `status` is `failed`, such as `invalid_change`.
	FailureCode *string   `json:"failure_code"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateDNSRecordRequest is the API's create DNS record request object.
type CreateDNSRecordRequest struct {
	// Relative to the zone (`@`, `www`); a name ending in the zone is accepted too.
	Name string `json:"name"`
	// One of: A, AAAA, CNAME, MX, TXT, SRV, CAA, NS, PTR.
	Type string `json:"type"`
	// Seconds, 60 to 86400; 300 when left out.
	TTL *int `json:"ttl,omitempty"`
	// 1 to 100 values in the type's format, such as `10 mail.example.com` for MX.
	Values []string `json:"values"`
	// Allow private and carrier-grade NAT addresses, for split-horizon use.
	AllowPrivate *bool   `json:"allow_private,omitempty"`
	Comment      *string `json:"comment,omitempty"`
}

// UpdateDNSRecordRequest is the API's update DNS record request object.
// Fields left nil are kept.
type UpdateDNSRecordRequest struct {
	TTL          *int     `json:"ttl,omitempty"`
	Values       []string `json:"values,omitempty"`
	AllowPrivate *bool    `json:"allow_private,omitempty"`
	Comment      *string  `json:"comment,omitempty"`
}

// DNSRecordInput is a record set in a batch change. An update may leave out
// Name and Type, which are the record's identity.
type DNSRecordInput struct {
	Name         string   `json:"name,omitempty"`
	Type         string   `json:"type,omitempty"`
	TTL          *int     `json:"ttl,omitempty"`
	Values       []string `json:"values,omitempty"`
	AllowPrivate *bool    `json:"allow_private,omitempty"`
	Comment      *string  `json:"comment,omitempty"`
}

// DNSRecordChange is one change of a batch.
type DNSRecordChange struct {
	// One of: create, update, delete.
	Action string `json:"action"`
	// The record to update or delete (`dnsr_…`).
	RecordID string `json:"record_id,omitempty"`
	// When set, the change fails with VERSION_CONFLICT unless the record is at this version.
	Version *int `json:"version,omitempty"`
	// Required to create or update.
	Record *DNSRecordInput `json:"record,omitempty"`
}

// DNSRecordBatchRequest is the API's DNS record batch request object.
// It is applied all or nothing.
type DNSRecordBatchRequest struct {
	// 1 to 1000 changes.
	Changes []DNSRecordChange `json:"changes"`
}

// ImportDNSZoneRequest is the API's import DNS zone request object.
type ImportDNSZoneRequest struct {
	// A BIND zone file.
	ZoneFile string `json:"zone_file"`
	// Also delete the records the file does not have (the apex SOA and NS are kept).
	Replace *bool `json:"replace,omitempty"`
}

// TLDOffer is the API's TLD offer object: what a TLD sells for, per year,
// before VAT. Registering for n years costs RegisterMinor + RenewMinor × (n-1); a
// transfer costs TransferMinor + RenewMinor × (n-1).
type TLDOffer struct {
	// Without the dot, such as `com`.
	TLD string `json:"tld"`
	// One of: NGN, USD.
	Currency      string `json:"currency"`
	RegisterMinor int64  `json:"register_minor"`
	RenewMinor    int64  `json:"renew_minor"`
	TransferMinor int64  `json:"transfer_minor"`
	Featured      bool   `json:"featured"`
	// The TLD's registry hides every registrant's contact details itself
	// (.ng and its second levels): privacy is always on and cannot be
	// turned off (422 on `privacy`, code PRIVACY_ALWAYS_ON).
	PrivacyAlwaysOn bool `json:"privacy_always_on"`
}

// RegisterPrice is what registering for years costs, in minor units.
func (o TLDOffer) RegisterPrice(years int) int64 {
	return o.RegisterMinor + o.RenewMinor*int64(max(years, 1)-1)
}

// TransferPrice is what a transfer that adds years costs, in minor units.
func (o TLDOffer) TransferPrice(years int) int64 {
	return o.TransferMinor + o.RenewMinor*int64(max(years, 1)-1)
}

// RenewPrice is what renewing for years costs, in minor units.
func (o TLDOffer) RenewPrice(years int) int64 {
	return o.RenewMinor * int64(max(years, 1))
}

// DomainSearchResult is the API's domain search result object.
type DomainSearchResult struct {
	Name string `json:"name"`
	TLD  string `json:"tld"`
	// `unknown` means it could not be checked now: try again. One of: available, unavailable, unknown.
	Status string `json:"status"`
	// Nil when the TLD is not sold.
	Price *TLDOffer `json:"price"`
}

// DomainTransferCheck is the API's domain transfer check object.
type DomainTransferCheck struct {
	Name         string `json:"name"`
	Transferable bool   `json:"transferable"`
	// A stable code when not transferable.
	Reason *string `json:"reason"`
	// Nil when the TLD is not sold.
	Price *TLDOffer `json:"price"`
}

// Domain is the API's domain object: a domain registered or transferred in
// through Pantech.
type Domain struct {
	// `dom_…`.
	ID   string `json:"id"`
	Name string `json:"name"`
	TLD  string `json:"tld"`
	// One of: pending, registering, transferring, active, expired, failed, transferred_out, suspended.
	Status      string  `json:"status"`
	FailureCode *string `json:"failure_code"`
	// Renewed by Pantech from your credit 30 days before expiry.
	AutoRenew    bool `json:"auto_renew"`
	TransferLock bool `json:"transfer_lock"`
	Privacy      bool `json:"privacy"`
	// One of: pantech, custom.
	NameServerMode string   `json:"name_server_mode"`
	NameServers    []string `json:"name_servers"`
	// The Pantech DNS zone it delegates to, in `pantech` mode.
	ZoneID       *string    `json:"zone_id"`
	RegisteredAt *time.Time `json:"registered_at"`
	ExpiresAt    *time.Time `json:"expires_at"`
	// False until the registrant clicks the registry's verification email; an
	// unverified registrant is suspended by the registry after 15 days.
	ContactVerified bool `json:"contact_verified"`
	// True while lock, privacy, name server or contact changes are being applied.
	SettingsPending bool `json:"settings_pending"`
	// True while a requested auth code is being fetched.
	AuthCodePending bool      `json:"auth_code_pending"`
	SuspendedReason *string   `json:"suspended_reason"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Version         int       `json:"version"`
}

// DomainOrder is the API's domain order object: one paid action on a domain.
// The amount is taken from credit before the registrar is called; a failed
// order returns it (Refunded). When credit does not cover an order (one you
// placed, or an automatic renewal), it waits in `awaiting_payment` for up to
// an hour. In US dollars Pantech charges the saved card the difference
// (PaymentURL is set only when the bank asks the customer to confirm); in
// naira the customer adds credit. The order goes on as soon as the money
// arrives; Domains.RetryOrderPayment tries again at once and waits another
// hour. Unpaid in time, it fails insufficient_credit and nothing is charged.
type DomainOrder struct {
	// `dor_…`.
	ID         string `json:"id"`
	DomainID   string `json:"domain_id"`
	DomainName string `json:"domain_name"`
	// One of: register, renew, transfer.
	Kind        string `json:"kind"`
	Years       int    `json:"years"`
	AmountMinor int64  `json:"amount_minor"`
	Currency    string `json:"currency"`
	// One of: pending, awaiting_payment, paid, submitted, succeeded, failed.
	// A pending order whose credit falls short becomes awaiting_payment (a
	// saved-card charge in US dollars, a deposit in naira), then paid, or
	// failed with insufficient_credit.
	Status string `json:"status"`
	// While `awaiting_payment`: the page where the bank asks the customer to
	// confirm the saved-card payment of the difference (US dollars only),
	// when it does. Nil when no confirmation is needed, while a naira order
	// waits for credit, and in every other status.
	PaymentURL *string `json:"payment_url"`
	// While `awaiting_payment`: when the order fails insufficient_credit.
	PaymentExpiresAt *time.Time `json:"payment_expires_at"`
	// One of: insufficient_credit, domain_unavailable, price_changed, registrar_rejected, transfer_rejected, invalid_auth_code, contact_rejected, tld_not_sold, billing_unavailable, domain_expired.
	FailureCode *string    `json:"failure_code"`
	Refunded    bool       `json:"refunded"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

// DomainOrderAccepted is what a registration, transfer or renewal answers
// with `202`: the domain and the order paying for it. Follow the order with
// Domains.WaitForOrder.
type DomainOrderAccepted struct {
	Domain Domain      `json:"domain"`
	Order  DomainOrder `json:"order"`
}

// DomainContact is the API's contact object: the domain's registrant. It is
// personal data: read it only to show it to its owner, and never log it.
type DomainContact struct {
	// One of: person, company.
	Type      string `json:"type"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	// Required when Type is `company`.
	OrganizationName string `json:"organization_name,omitempty"`
	Email            string `json:"email"`
	// `+CC.NUMBER`, such as `+234.8012345678`.
	Phone        string `json:"phone"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2,omitempty"`
	City         string `json:"city"`
	State        string `json:"state,omitempty"`
	// ISO 3166-1 alpha-2, such as `NG`.
	CountryCode string `json:"country_code"`
	PostalCode  string `json:"postal_code,omitempty"`
}

// String hides the contact, so printing a struct that holds one does not
// leak it into logs.
func (DomainContact) String() string { return "pantech.DomainContact{…}" }

// GoString hides the contact from %#v too.
func (DomainContact) GoString() string { return "pantech.DomainContact{…}" }

// RegisterDomainRequest is the API's register domain request object.
type RegisterDomainRequest struct {
	Name string `json:"name"`
	// 1 to 10; 1 when left out.
	Years   *int          `json:"years,omitempty"`
	Contact DomainContact `json:"contact"`
	// On when left out.
	AutoRenew *bool `json:"auto_renew,omitempty"`
	// On when left out.
	Privacy *bool `json:"privacy,omitempty"`
	// Must be true: the registrant accepted the Domain Registration Agreement
	// (https://pantechdynamics.com/legal/domain-registration).
	AcceptTerms bool `json:"accept_terms"`
}

// TransferDomainRequest is the API's transfer domain request object.
type TransferDomainRequest struct {
	Name string `json:"name"`
	// From the current registrar.
	AuthCode string `json:"auth_code"`
	// 1 to 10; a transfer adds one year when left out.
	Years     *int          `json:"years,omitempty"`
	Contact   DomainContact `json:"contact"`
	AutoRenew *bool         `json:"auto_renew,omitempty"`
	Privacy   *bool         `json:"privacy,omitempty"`
	// Must be true: the registrant accepted the Domain Registration Agreement
	// (https://pantechdynamics.com/legal/domain-registration).
	AcceptTerms bool `json:"accept_terms"`
}

// String hides the auth code and the contact.
func (r TransferDomainRequest) String() string {
	return "pantech.TransferDomainRequest{Name: " + r.Name + ", AuthCode: …}"
}

// GoString hides the auth code and the contact from %#v too.
func (r TransferDomainRequest) GoString() string { return r.String() }

// UpdateDomainRequest is the API's update domain request object. Fields left
// nil are kept.
type UpdateDomainRequest struct {
	AutoRenew    *bool `json:"auto_renew,omitempty"`
	TransferLock *bool `json:"transfer_lock,omitempty"`
	Privacy      *bool `json:"privacy,omitempty"`
}

// SetDomainNameServersRequest is the API's set domain name servers request
// object: Pantech DNS (with the zone of the same name, made if missing), or
// up to 6 custom name servers.
type SetDomainNameServersRequest struct {
	UsePantechDNS bool `json:"use_pantech_dns"`
	// With UsePantechDNS, a zone of the domain's name (`dnsz_…`).
	ZoneID *string `json:"zone_id,omitempty"`
	// Without UsePantechDNS, 2 to 6 host names.
	NameServers []string `json:"name_servers,omitempty"`
}

// UpdateDomainContactRequest is the API's update domain contact request object.
type UpdateDomainContactRequest struct {
	Contact DomainContact `json:"contact"`
}

// RenewDomainRequest is the API's renew domain request object.
type RenewDomainRequest struct {
	// 1 to 10; 1 when left out.
	Years *int `json:"years,omitempty"`
}

// DomainAuthCode is the API's domain auth code object: the code that moves
// the domain to another registrar, or Status `pending` while it is fetched.
type DomainAuthCode struct {
	// Empty while pending. Show it once; anyone holding it can move the domain away.
	AuthCode string `json:"auth_code,omitempty"`
	// When the code stops being shown.
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	// `pending` while the code is being fetched.
	Status string `json:"status,omitempty"`
}

// Pending reports whether the code is still being fetched.
func (a *DomainAuthCode) Pending() bool { return a.AuthCode == "" }

// String hides the code.
func (a DomainAuthCode) String() string { return "pantech.DomainAuthCode{…}" }

// GoString hides the code from %#v too.
func (a DomainAuthCode) GoString() string { return a.String() }

// RegistryNamespace is the API's registry namespace object: the
// organization's one namespace on the Pantech container registry. Images are
// <RegistryHost>/<Name>/<repository>:<tag>.
type RegistryNamespace struct {
	Name string `json:"name"`
	// Such as `registry.pantechdynamics.com`.
	RegistryHost string `json:"registry_host"`
	// Whether a push to a repository that does not exist creates it.
	CreateOnPush      bool      `json:"create_on_push"`
	RepositoryCount   int       `json:"repository_count"`
	StorageBytes      int64     `json:"storage_bytes"`
	StorageLimitBytes int64     `json:"storage_limit_bytes"`
	RepositoryLimit   int       `json:"repository_limit"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Version           int       `json:"version"`
}

// SetRegistryNamespaceRequest is the API's set registry namespace request object.
type SetRegistryNamespaceRequest struct {
	// 3 to 30 lower-case letters, digits and single hyphens, starting with a
	// letter. It cannot change once a repository exists.
	Name string `json:"name"`
	// On when left out.
	CreateOnPush *bool `json:"create_on_push,omitempty"`
}

// RegistryRetention is a repository's retention policy. A nil field is no
// rule. Sent in a request, it replaces the whole policy.
type RegistryRetention struct {
	// Keep the newest N tags, 1 to 1000.
	KeepLast *int `json:"keep_last"`
	// Delete untagged manifests older than this many days, 1 to 365.
	UntaggedAfterDays *int `json:"untagged_after_days"`
	// Tags matching this glob (such as `v*`), at most 128 characters, are never removed.
	KeepPattern *string `json:"keep_pattern"`
}

// RegistryRepository is the API's registry repository object.
type RegistryRepository struct {
	// `crr_…`.
	ID        string `json:"id"`
	Namespace string `json:"namespace"`
	// Such as `api/server`.
	Name string `json:"name"`
	// What to push and pull: `registry.pantechdynamics.com/acme/api/server`.
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	// One of: private, public.
	Visibility    string             `json:"visibility"`
	ImmutableTags bool               `json:"immutable_tags"`
	Retention     *RegistryRetention `json:"retention"`
	TagCount      int                `json:"tag_count"`
	StorageBytes  int64              `json:"storage_bytes"`
	LastPushedAt  *time.Time         `json:"last_pushed_at"`
	// One of: active, deleting, suspended.
	Status          string    `json:"status"`
	SuspendedReason *string   `json:"suspended_reason"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Version         int       `json:"version"`
}

// CreateRegistryRepositoryRequest is the API's create registry repository request object.
type CreateRegistryRepositoryRequest struct {
	// Lower-case path components such as `api/server`, at most 255 characters.
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	// `private` when left out. One of: private, public.
	Visibility    *string            `json:"visibility,omitempty"`
	ImmutableTags *bool              `json:"immutable_tags,omitempty"`
	Retention     *RegistryRetention `json:"retention,omitempty"`
}

// UpdateRegistryRepositoryRequest is the API's update registry repository
// request object. Fields left nil are kept; Retention replaces the whole policy.
type UpdateRegistryRepositoryRequest struct {
	Description *string `json:"description,omitempty"`
	// One of: private, public.
	Visibility    *string            `json:"visibility,omitempty"`
	ImmutableTags *bool              `json:"immutable_tags,omitempty"`
	Retention     *RegistryRetention `json:"retention,omitempty"`
}

// RegistryTag is the API's registry tag object.
type RegistryTag struct {
	Name string `json:"name"`
	// Such as `sha256:9f86…`.
	Digest    string `json:"digest"`
	MediaType string `json:"media_type"`
	SizeBytes int64  `json:"size_bytes"`
	// Such as `linux/amd64`.
	Platforms    []string   `json:"platforms"`
	PushedAt     time.Time  `json:"pushed_at"`
	LastPulledAt *time.Time `json:"last_pulled_at"`
	Immutable    bool       `json:"immutable"`
}

// RegistryToken is the API's registry token object: a credential for the
// registry only, never for the API.
type RegistryToken struct {
	// `prt_…`.
	ID   string `json:"id"`
	Name string `json:"name"`
	// `push` includes pull. One of: pull, push.
	Permission string `json:"permission"`
	// Empty means every repository in the namespace.
	Repositories []string   `json:"repositories"`
	ExpiresAt    *time.Time `json:"expires_at"`
	LastUsedAt   *time.Time `json:"last_used_at"`
	CreatedAt    time.Time  `json:"created_at"`
	RevokedAt    *time.Time `json:"revoked_at"`
	// Set only in the answer to CreateToken, and shown nowhere again: the
	// password for `docker login`. It is never printed by fmt.
	Secret string `json:"secret,omitempty"`
}

// String describes the token without its secret.
func (t RegistryToken) String() string {
	return "pantech.RegistryToken{ID: " + t.ID + ", Name: " + t.Name + ", Permission: " + t.Permission + "}"
}

// GoString describes the token without its secret for %#v too.
func (t RegistryToken) GoString() string { return t.String() }

// CreateRegistryTokenRequest is the API's create registry token request object.
type CreateRegistryTokenRequest struct {
	Name string `json:"name"`
	// One of: pull, push.
	Permission string `json:"permission"`
	// Repository names (such as `api/server`); every repository when left out.
	Repositories []string `json:"repositories,omitempty"`
	// Never expires when left out.
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// IAMAction is the API's IAM action object: one entry of the action catalogue.
type IAMAction struct {
	// `<service>:<verb>`, such as `instances:Stop`.
	Action        string   `json:"action"`
	Service       string   `json:"service"`
	Description   string   `json:"description"`
	ResourceTypes []string `json:"resource_types"`
	// Whether the organization's `read` member baseline allows it.
	Read bool `json:"read"`
}

// IAMSettings is the API's access setting object.
type IAMSettings struct {
	// What members may do with no policy: `read` (the default) or `none`. One of: read, none.
	MemberBaseline string `json:"member_baseline"`
	Version        int    `json:"version"`
	// Nil until the setting is first changed.
	UpdatedAt *time.Time `json:"updated_at"`
}

// UpdateIAMSettingsRequest is the API's update access setting request object.
type UpdateIAMSettingsRequest struct {
	// One of: read, none.
	MemberBaseline string `json:"member_baseline"`
}

// IAMTeam is the API's team object: a named group of members.
type IAMTeam struct {
	// `team_…`.
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MemberCount int    `json:"member_count"`
	// The policies attached to the team.
	PolicyIDs []string  `json:"policy_ids"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}

// CreateIAMTeamRequest is the API's create team request object.
type CreateIAMTeamRequest struct {
	// 1 to 64 characters, unique in the organization.
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// UpdateIAMTeamRequest is the API's update team request object. Fields left
// nil are kept.
type UpdateIAMTeamRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// IAMTeamMember is the API's team member object.
type IAMTeamMember struct {
	// `usr_…`.
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	// The member's organization role. One of: owner, admin, member.
	Role    string    `json:"role"`
	AddedAt time.Time `json:"added_at"`
}

// PolicyConditions must all hold for a statement to apply. A field left
// empty is no condition.
type PolicyConditions struct {
	// The request's client address is in one of these CIDRs (up to 20).
	SourceIP []string `json:"source_ip,omitempty"`
	// True: applies only when the session authenticated with MFA. False:
	// applies only when the request lacks MFA (API keys never carry it, so a
	// deny with false blocks keys too). Nil: either way.
	MFA *bool `json:"mfa,omitempty"`
}

// PolicyStatement is one statement of a policy. Build one with Allow or Deny.
type PolicyStatement struct {
	// Optional, unique within the policy, 1 to 64 of [A-Za-z0-9_-].
	Sid string `json:"sid,omitempty"`
	// One of: allow, deny.
	Effect string `json:"effect"`
	// 1 to 50 actions from Client.IAM.ListActions, with `*` wildcards:
	// `instances:Stop`, `instances:*`, `*:Read`, `*`.
	Actions []string `json:"actions"`
	// 1 to 50 of `*`, a resource id (`vm_…`), or `name:<glob>` matched
	// against the resource's name, such as `name:staging-*`.
	Resources  []string          `json:"resources"`
	Conditions *PolicyConditions `json:"conditions,omitempty"`
}

// IAMPolicy is the API's policy object.
type IAMPolicy struct {
	// `pol_…`; Pantech's managed policies are `pol_pantech_<name>`.
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// True for Pantech's policies, which cannot be changed.
	Managed         bool              `json:"managed"`
	Statements      []PolicyStatement `json:"statements"`
	AttachmentCount int               `json:"attachment_count"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Version         int               `json:"version"`
}

// CreateIAMPolicyRequest is the API's create policy request object.
type CreateIAMPolicyRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	// 1 to 20 statements.
	Statements []PolicyStatement `json:"statements"`
}

// UpdateIAMPolicyRequest is the API's update policy request object. It
// replaces the whole policy: name, description and statements.
type UpdateIAMPolicyRequest struct {
	Name        string            `json:"name"`
	Description *string           `json:"description,omitempty"`
	Statements  []PolicyStatement `json:"statements"`
}

// IAMAttachment is the API's attachment object: a policy attached to a
// member or a team.
type IAMAttachment struct {
	PolicyID string `json:"policy_id"`
	// One of: member, team.
	PrincipalType string `json:"principal_type"`
	// `usr_…` or `team_…`.
	PrincipalID string    `json:"principal_id"`
	AttachedAt  time.Time `json:"attached_at"`
}

// AttachIAMPolicyRequest is the API's attach policy request object.
type AttachIAMPolicyRequest struct {
	// One of: member, team.
	PrincipalType string `json:"principal_type"`
	PrincipalID   string `json:"principal_id"`
}

// IAMEffectivePolicy is one policy that applies to a member, and the team
// that brought it (nil when attached to the member directly).
type IAMEffectivePolicy struct {
	PolicyID string  `json:"policy_id"`
	TeamID   *string `json:"team_id"`
}

// IAMEffectiveAccess is the API's effective access object: what applies to
// a member.
type IAMEffectiveAccess struct {
	UserID string `json:"user_id"`
	// One of: owner, admin, member.
	Role           string               `json:"role"`
	Teams          []string             `json:"teams"`
	Policies       []IAMEffectivePolicy `json:"policies"`
	MemberBaseline string               `json:"member_baseline"`
}

// SimulateIAMAccessRequest is the API's simulation request object.
type SimulateIAMAccessRequest struct {
	UserID string `json:"user_id"`
	// Such as `instances:Delete`.
	Action string `json:"action"`
	// A resource id such as `vm_…`.
	Resource *string `json:"resource,omitempty"`
	// The resource's name, for `name:` globs.
	ResourceName *string `json:"resource_name,omitempty"`
	// The client address, for `source_ip` conditions.
	SourceIP *string `json:"source_ip,omitempty"`
	// Whether the session used MFA.
	MFA *bool `json:"mfa,omitempty"`
}

// IAMSimulation is the API's simulation object: the decision the evaluator
// would make.
type IAMSimulation struct {
	Allowed bool `json:"allowed"`
	// One of: owner, admin, allow, baseline_read, explicit_deny, implicit_deny.
	Decision string `json:"decision"`
	// The policy and statement that decided, for allow and explicit_deny.
	PolicyID *string `json:"policy_id"`
	Sid      *string `json:"sid"`
}

// IAMMember is the API's organization member object, as IAM sees it.
type IAMMember struct {
	// `usr_…`.
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	// One of: owner, admin, member.
	Role        string   `json:"role"`
	TeamIDs     []string `json:"team_ids"`
	PolicyCount int      `json:"policy_count"`
}

// CertificateValidation is the API's certificate validation object: how one
// name is proved for the CA.
type CertificateValidation struct {
	Name string `json:"name"`
	// `pantech_dns`: the name is in an active Pantech DNS zone and Pantech
	// writes the challenge itself. `cname`: add the CNAME CNAMEName →
	// CNAMETarget at your DNS host once. One of: pantech_dns, cname.
	Method      string  `json:"method"`
	ZoneID      *string `json:"zone_id"`
	CNAMEName   *string `json:"cname_name"`
	CNAMETarget *string `json:"cname_target"`
	// One of: pending, ready, valid, invalid.
	Status string `json:"status"`
}

// CertificateUsage is a Pantech resource that uses a certificate.
type CertificateUsage struct {
	// Such as `load_balancer`.
	Type string `json:"type"`
	ID   string `json:"id"`
}

// Certificate is the API's certificate object.
type Certificate struct {
	// `cert_…`.
	ID string `json:"id"`
	// One of: public, imported.
	Kind  string   `json:"kind"`
	Names []string `json:"names"`
	// The first name, when it is at most 64 characters.
	CommonName  *string `json:"common_name"`
	Description string  `json:"description"`
	// One of: ecdsa-p256, ecdsa-p384, rsa-2048, rsa-3072, rsa-4096.
	KeyType string `json:"key_type"`
	// Set at creation, never changed: whether Export may hand out the key.
	Exportable bool `json:"exportable"`
	// One of: pending_validation, issuing, issued, renewing, expired, revoked, failed.
	Status string `json:"status"`
	// While pending_validation, issuing or renewing. One of: awaiting_cname,
	// awaiting_dns, awaiting_ca, rate_limited, ca_unavailable, validation_failed.
	PendingReason *string `json:"pending_reason"`
	// When failed. One of: validation_timeout, validation_failed,
	// caa_forbidden, rejected_identifier, rate_limited, ca_unavailable,
	// dns_record_limit, ca_rejected.
	FailureCode *string `json:"failure_code"`
	// One per name; empty for an imported certificate.
	Validations []CertificateValidation `json:"validations"`
	// While waiting for CNAMEs on a first issue, when the request fails.
	ValidationDeadlineAt *time.Time `json:"validation_deadline_at"`
	Issuer               *string    `json:"issuer"`
	// Hexadecimal.
	SerialNumber *string `json:"serial_number"`
	// SHA-256 of the leaf's DER, lower-case hex.
	FingerprintSHA256 *string    `json:"fingerprint_sha256"`
	NotBefore         *time.Time `json:"not_before"`
	NotAfter          *time.Time `json:"not_after"`
	IssuedAt          *time.Time `json:"issued_at"`
	// Renewal starts at this time, 30 days before NotAfter.
	RenewAfter *time.Time `json:"renew_after"`
	// One of: scheduled, in_progress, failing. Nil for imported certificates.
	RenewalStatus *string `json:"renewal_status"`
	// A certificate in use cannot be revoked or deleted.
	InUseBy           []CertificateUsage `json:"in_use_by"`
	RevocationPending bool               `json:"revocation_pending"`
	RevokedAt         *time.Time         `json:"revoked_at"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
	Version           int                `json:"version"`
}

// RequestCertificateRequest is the API's request certificate request object.
type RequestCertificateRequest struct {
	// 1 to 100 DNS names; wildcards as `*.example.com`.
	Names []string `json:"names"`
	// `ecdsa-p256` when left out. One of: ecdsa-p256, rsa-2048.
	KeyType *string `json:"key_type,omitempty"`
	// Whether the private key may be exported. Fixed at creation.
	Exportable  *bool   `json:"exportable,omitempty"`
	Description *string `json:"description,omitempty"`
}

// ImportCertificateRequest is the API's import certificate request object.
// PrivateKeyPEM never prints with fmt.
type ImportCertificateRequest struct {
	CertificatePEM string `json:"certificate_pem"`
	// The intermediates, leaf excluded.
	ChainPEM *string `json:"chain_pem,omitempty"`
	// Unencrypted PKCS #1, PKCS #8 or SEC 1.
	PrivateKeyPEM string  `json:"private_key_pem"`
	Exportable    *bool   `json:"exportable,omitempty"`
	Description   *string `json:"description,omitempty"`
}

// String hides the private key.
func (r ImportCertificateRequest) String() string { return "pantech.ImportCertificateRequest{…}" }

// GoString hides the private key from %#v too.
func (r ImportCertificateRequest) GoString() string { return r.String() }

// ReimportCertificateRequest is the API's reimport certificate request
// object: the same names, a new certificate and key. PrivateKeyPEM never
// prints with fmt.
type ReimportCertificateRequest struct {
	CertificatePEM string  `json:"certificate_pem"`
	ChainPEM       *string `json:"chain_pem,omitempty"`
	PrivateKeyPEM  string  `json:"private_key_pem"`
}

// String hides the private key.
func (r ReimportCertificateRequest) String() string { return "pantech.ReimportCertificateRequest{…}" }

// GoString hides the private key from %#v too.
func (r ReimportCertificateRequest) GoString() string { return r.String() }

// UpdateCertificateRequest is the API's update certificate request object.
type UpdateCertificateRequest struct {
	Description string `json:"description"`
}

// RevokeCertificateRequest is the API's revoke certificate request object.
type RevokeCertificateRequest struct {
	// `unspecified` when left out. One of: unspecified, key_compromise.
	Reason *string `json:"reason,omitempty"`
}

// ExportedCertificate is the API's exported certificate object. It holds the
// private key: write it somewhere only its users can read, and never log it.
// PrivateKeyPEM never prints with fmt.
type ExportedCertificate struct {
	CertificatePEM string `json:"certificate_pem"`
	// The intermediates, leaf excluded; may be empty for an imported certificate.
	ChainPEM string `json:"chain_pem"`
	// PKCS #8 for certificates Pantech issued; as imported otherwise.
	PrivateKeyPEM string    `json:"private_key_pem"`
	NotAfter      time.Time `json:"not_after"`
}

// String hides the private key.
func (e ExportedCertificate) String() string {
	return "pantech.ExportedCertificate{NotAfter: " + e.NotAfter.Format(time.RFC3339) + ", PrivateKeyPEM: …}"
}

// GoString hides the private key from %#v too.
func (e ExportedCertificate) GoString() string { return e.String() }

// FullChainPEM is the leaf followed by the intermediates, the file most
// servers take as their certificate.
func (e *ExportedCertificate) FullChainPEM() string {
	if e.ChainPEM == "" {
		return e.CertificatePEM
	}
	leaf := e.CertificatePEM
	if leaf != "" && leaf[len(leaf)-1] != '\n' {
		leaf += "\n"
	}
	return leaf + e.ChainPEM
}

// DNSZoneQuote is the API's DNS zone quote object: what hosted zones cost.
// The first FreeZones live zones are free; each further one costs
// MonthlyPerZoneMinor a month, billed by the hour.
type DNSZoneQuote struct {
	FreeZones int `json:"free_zones"`
	// Zones from `awaiting_delegation` until deleted.
	LiveZones     int `json:"live_zones"`
	BillableZones int `json:"billable_zones"`
	// One of: NGN, USD.
	Currency            string `json:"currency"`
	MonthlyPerZoneMinor int64  `json:"monthly_per_zone_minor"`
	// BillableZones × MonthlyPerZoneMinor.
	MonthlyEstimateMinor int64 `json:"monthly_estimate_minor"`
	// 730.
	NominalMonthHours int `json:"nominal_month_hours"`
}

// ObjectStorageLimits are the organization's object storage quotas.
type ObjectStorageLimits struct {
	MaxBuckets      int   `json:"max_buckets"`
	MaxStorageBytes int64 `json:"max_storage_bytes"`
	MaxObjects      int64 `json:"max_objects"`
}

// ObjectStorageStatus reports whether object storage is enabled and its S3 endpoint.
type ObjectStorageStatus struct {
	Enabled             bool                 `json:"enabled"`
	Endpoint            string               `json:"endpoint"`
	Beta                bool                 `json:"beta"`
	DirectUploadEnabled bool                 `json:"direct_upload_enabled"`
	Limits              *ObjectStorageLimits `json:"limits"`
}

// ObjectStorageHourUsage is one hour of object storage transfer and request counts.
type ObjectStorageHourUsage struct {
	Hour          time.Time `json:"hour"`
	BytesSent     int64     `json:"bytes_sent"`
	BytesReceived int64     `json:"bytes_received"`
	ReadOps       int64     `json:"read_ops"`
	WriteOps      int64     `json:"write_ops"`
}

// ObjectStorageUsage is the organization's hourly transfer and request usage summary.
type ObjectStorageUsage struct {
	Range              string                   `json:"range"`
	TotalBytesSent     int64                    `json:"total_bytes_sent"`
	TotalBytesReceived int64                    `json:"total_bytes_received"`
	TotalReadOps       int64                    `json:"total_read_ops"`
	TotalWriteOps      int64                    `json:"total_write_ops"`
	Hours              []ObjectStorageHourUsage `json:"hours"`
}

// ObjectStorageKey is an S3 access key without its secret.
type ObjectStorageKey struct {
	// The key's id, starting `osk_`.
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Region      string    `json:"region"`
	AccessKeyID string    `json:"access_key_id"`
	ReadOnly    bool      `json:"read_only"`
	Buckets     []string  `json:"buckets"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreatedObjectStorageKey is a newly created S3 access key with its secret, shown once.
// SecretAccessKey never prints with fmt.
type CreatedObjectStorageKey struct {
	ObjectStorageKey
	// Shown once on creation; store it immediately.
	SecretAccessKey string `json:"secret_access_key"`
	// The S3 endpoint to use with this key.
	Endpoint string `json:"endpoint"`
}

// String hides SecretAccessKey from fmt output and logs.
func (k CreatedObjectStorageKey) String() string {
	return "pantech.CreatedObjectStorageKey{ID: " + k.ID + ", Name: " + k.Name + ", AccessKeyID: " + k.AccessKeyID + ", SecretAccessKey: …}"
}

// GoString hides SecretAccessKey from %#v too.
func (k CreatedObjectStorageKey) GoString() string { return k.String() }

// CreateObjectStorageKeyRequest creates a new S3 access key.
type CreateObjectStorageKeyRequest struct {
	// 1 to 64 letters, digits, `-`, `_` or `.`.
	Name     string    `json:"name"`
	Region   *string   `json:"region,omitempty"`
	ReadOnly *bool     `json:"read_only,omitempty"`
	Buckets  *[]string `json:"buckets,omitempty"`
}

// Bucket is an object storage bucket.
type Bucket struct {
	// The bucket's id, starting `bkt_`.
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Region         string   `json:"region"`
	Endpoint       string   `json:"endpoint"`
	QuotaBytes     int64    `json:"quota_bytes"`
	SizeBytes      int64    `json:"size_bytes"`
	ObjectCount    int64    `json:"object_count"`
	Versioning     bool     `json:"versioning"`
	ObjectLock     bool     `json:"object_lock"`
	PublicPrefixes []string `json:"public_prefixes"`
	// One of: provisioning, active, updating, deleting, failed, unknown.
	State      string     `json:"state"`
	CreatedAt  time.Time  `json:"created_at"`
	RetiredAt  *time.Time `json:"retired_at"`
	PurgeDueAt *time.Time `json:"purge_due_at"`
}

// CreatedBucket is a newly created bucket and the operation creating it.
type CreatedBucket struct {
	Bucket
	OperationID string `json:"operation_id"`
}

// UpdatedBucket is a bucket and the operation applying the change, if any.
type UpdatedBucket struct {
	Bucket
	OperationID *string `json:"operation_id"`
}

// CreateBucketRequest creates a bucket.
type CreateBucketRequest struct {
	// 3 to 63 lowercase letters, digits and hyphens, starting and ending with a letter or digit.
	Name           string    `json:"name"`
	Region         *string   `json:"region,omitempty"`
	QuotaBytes     *int64    `json:"quota_bytes,omitempty"`
	Versioning     *bool     `json:"versioning,omitempty"`
	ObjectLock     *bool     `json:"object_lock,omitempty"`
	PublicPrefixes *[]string `json:"public_prefixes,omitempty"`
}

// UpdateBucketRequest changes a bucket's quota, versioning or public prefixes.
type UpdateBucketRequest struct {
	QuotaBytes     *int64    `json:"quota_bytes,omitempty"`
	Versioning     *bool     `json:"versioning,omitempty"`
	PublicPrefixes *[]string `json:"public_prefixes,omitempty"`
}

// BucketStats holds live stored bytes and object count for a bucket.
type BucketStats struct {
	SizeBytes   int64 `json:"size_bytes"`
	ObjectCount int64 `json:"object_count"`
}

// EmptyBucketResult summarizes an empty-bucket operation.
type EmptyBucketResult struct {
	DeletedCount     int      `json:"deleted_count"`
	AbortedUploads   int      `json:"aborted_uploads"`
	LockedCount      int      `json:"locked_count"`
	SampleLockedKeys []string `json:"sample_locked_keys"`
}

// BucketObject is one object returned by listing a bucket.
type BucketObject struct {
	Key          string     `json:"key"`
	SizeBytes    int64      `json:"size_bytes"`
	ETag         string     `json:"etag"`
	ContentType  string     `json:"content_type"`
	LastModified *time.Time `json:"last_modified"`
	StorageClass string     `json:"storage_class"`
	VersionID    string     `json:"version_id"`
	IsLatest     bool       `json:"is_latest"`
	Public       bool       `json:"public"`
	ObjectPublic bool       `json:"object_public"`
	// One of: private, object, bucket_prefix.
	PublicSource string `json:"public_source"`
	PublicURL    string `json:"public_url"`
}

// BucketObjectVersion is one object version or delete marker returned when listing with Versions: true.
type BucketObjectVersion struct {
	Key            string     `json:"key"`
	VersionID      string     `json:"version_id"`
	IsLatest       bool       `json:"is_latest"`
	IsDeleteMarker bool       `json:"is_delete_marker"`
	SizeBytes      int64      `json:"size_bytes"`
	ETag           string     `json:"etag"`
	LastModified   *time.Time `json:"last_modified"`
	StorageClass   string     `json:"storage_class"`
}

// BucketObjectList is a page of folders and objects (or versions when listed with Versions: true).
type BucketObjectList struct {
	Prefix              string                `json:"prefix"`
	Folders             []string              `json:"folders"`
	Objects             []BucketObject        `json:"objects"`
	NextCursor          string                `json:"next_cursor"`
	Versions            []BucketObjectVersion `json:"versions"`
	NextKeyMarker       string                `json:"next_key_marker"`
	NextVersionIDMarker string                `json:"next_version_id_marker"`
	IsTruncated         bool                  `json:"is_truncated"`
}

// StorageTag is a key-value tag on a bucket or object.
type StorageTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// BucketObjectMeta holds properties, metadata, tags and object-lock status for one object.
type BucketObjectMeta struct {
	Key                string            `json:"key"`
	SizeBytes          int64             `json:"size_bytes"`
	ETag               string            `json:"etag"`
	ContentType        string            `json:"content_type"`
	CacheControl       string            `json:"cache_control"`
	ContentDisposition string            `json:"content_disposition"`
	LastModified       *time.Time        `json:"last_modified"`
	StorageClass       string            `json:"storage_class"`
	VersionID          string            `json:"version_id"`
	Metadata           map[string]string `json:"metadata"`
	Tags               []StorageTag      `json:"tags"`
	RetentionMode      string            `json:"retention_mode"`
	RetainUntil        *time.Time        `json:"retain_until"`
	LegalHold          bool              `json:"legal_hold"`
	Public             bool              `json:"public"`
	ObjectPublic       bool              `json:"object_public"`
	// One of: private, object, bucket_prefix.
	PublicSource string `json:"public_source"`
	PublicURL    string `json:"public_url"`
}

// PresignObjectDownloadRequest mints a short-lived direct download URL.
type PresignObjectDownloadRequest struct {
	Key              string  `json:"key"`
	VersionID        *string `json:"version_id,omitempty"`
	Inline           *bool   `json:"inline,omitempty"`
	ExpiresInSeconds *int    `json:"expires_in_seconds,omitempty"`
}

// PresignObjectUploadRequest mints a short-lived direct upload URL.
type PresignObjectUploadRequest struct {
	Key         string  `json:"key"`
	ContentType *string `json:"content_type,omitempty"`
}

// PresignedURL is a direct transfer URL or proxy indicator.
type PresignedURL struct {
	Direct           bool   `json:"direct"`
	URL              string `json:"url"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

// UploadedBucketObject summarizes an uploaded object.
type UploadedBucketObject struct {
	Key          string     `json:"key"`
	SizeBytes    int64      `json:"size_bytes"`
	ETag         string     `json:"etag"`
	ContentType  string     `json:"content_type"`
	LastModified *time.Time `json:"last_modified"`
	VersionID    string     `json:"version_id"`
	Public       bool       `json:"public"`
	ObjectPublic bool       `json:"object_public"`
	// One of: private, object, bucket_prefix.
	PublicSource string `json:"public_source"`
	PublicURL    string `json:"public_url"`
}

// CreateBucketFolderRequest creates a folder prefix marker in a bucket.
type CreateBucketFolderRequest struct {
	Prefix string `json:"prefix"`
}

// BucketFolder is the created folder prefix.
type BucketFolder struct {
	Prefix string `json:"prefix"`
}

// CopyBucketObjectRequest copies, renames or replaces metadata on an object within a bucket.
type CopyBucketObjectRequest struct {
	SourceKey          string             `json:"source_key"`
	DestinationKey     string             `json:"destination_key"`
	DeleteSource       *bool              `json:"delete_source,omitempty"`
	ReplaceMetadata    *bool              `json:"replace_metadata,omitempty"`
	ContentType        *string            `json:"content_type,omitempty"`
	CacheControl       *string            `json:"cache_control,omitempty"`
	ContentDisposition *string            `json:"content_disposition,omitempty"`
	Metadata           *map[string]string `json:"metadata,omitempty"`
}

// BucketObjectSummary summarizes a copied or completed multipart object.
type BucketObjectSummary struct {
	Key          string     `json:"key"`
	ETag         string     `json:"etag"`
	VersionID    string     `json:"version_id"`
	LastModified *time.Time `json:"last_modified"`
	Public       bool       `json:"public"`
	ObjectPublic bool       `json:"object_public"`
	// One of: private, object, bucket_prefix.
	PublicSource string `json:"public_source"`
	PublicURL    string `json:"public_url"`
}

// DeleteBucketObjectItem identifies one object or version to delete.
type DeleteBucketObjectItem struct {
	Key       string  `json:"key"`
	VersionID *string `json:"version_id,omitempty"`
}

// DeleteBucketObjectsRequest deletes objects, versions or all objects under a folder prefix.
type DeleteBucketObjectsRequest struct {
	Items           []DeleteBucketObjectItem `json:"items,omitempty"`
	RecursivePrefix *string                  `json:"recursive_prefix,omitempty"`
}

// DeleteBucketObjectFailure describes an object that could not be deleted.
type DeleteBucketObjectFailure struct {
	Key       string `json:"key"`
	VersionID string `json:"version_id"`
	Code      string `json:"code"`
	Locked    bool   `json:"locked"`
}

// DeleteBucketObjectsResult summarizes a batch delete.
type DeleteBucketObjectsResult struct {
	Deleted int                         `json:"deleted"`
	Failed  []DeleteBucketObjectFailure `json:"failed"`
}

// PutBucketObjectTagsRequest replaces key-value tags on an object or version.
type PutBucketObjectTagsRequest struct {
	Key       string       `json:"key"`
	VersionID *string      `json:"version_id,omitempty"`
	Tags      []StorageTag `json:"tags"`
}

// StorageTags wraps a list of key-value tags.
type StorageTags struct {
	Tags []StorageTag `json:"tags"`
}

// PutBucketObjectLockRequest sets retention or legal hold on an object version.
type PutBucketObjectLockRequest struct {
	Key       string  `json:"key"`
	VersionID *string `json:"version_id,omitempty"`
	// One of: GOVERNANCE, COMPLIANCE.
	RetentionMode    *string    `json:"retention_mode,omitempty"`
	RetainUntil      *time.Time `json:"retain_until,omitempty"`
	BypassGovernance *bool      `json:"bypass_governance,omitempty"`
	LegalHold        *bool      `json:"legal_hold,omitempty"`
}

// PutBucketObjectLockResult confirms an object lock update.
type PutBucketObjectLockResult struct {
	OK bool `json:"ok"`
}

// PutBucketObjectVisibilityRequest toggles public or private read access on an individual object.
type PutBucketObjectVisibilityRequest struct {
	Key       string  `json:"key"`
	VersionID *string `json:"version_id,omitempty"`
	Public    bool    `json:"public"`
}

// BucketObjectVisibility reports an object's effective visibility after an update.
type BucketObjectVisibility struct {
	Key          string `json:"key"`
	Public       bool   `json:"public"`
	ObjectPublic bool   `json:"object_public"`
	// One of: private, object, bucket_prefix.
	PublicSource string `json:"public_source"`
	PublicURL    string `json:"public_url"`
}

// StartMultipartUploadRequest initiates a multipart upload session.
type StartMultipartUploadRequest struct {
	Key                string             `json:"key"`
	ContentType        *string            `json:"content_type,omitempty"`
	CacheControl       *string            `json:"cache_control,omitempty"`
	ContentDisposition *string            `json:"content_disposition,omitempty"`
	Metadata           *map[string]string `json:"metadata,omitempty"`
	TotalSize          *int64             `json:"total_size,omitempty"`
	Public             *bool              `json:"public,omitempty"`
}

// MultipartUploadSession is returned when starting a multipart upload.
type MultipartUploadSession struct {
	UploadID      string `json:"upload_id"`
	Mode          string `json:"mode"`
	PartSizeBytes int64  `json:"part_size_bytes"`
}

// SignMultipartPartRequest signs a direct upload URL for one multipart part.
type SignMultipartPartRequest struct {
	Key        string `json:"key"`
	UploadID   string `json:"upload_id"`
	PartNumber int    `json:"part_number"`
}

// MultipartPart identifies an uploaded part by number and ETag.
type MultipartPart struct {
	PartNumber int    `json:"part_number"`
	ETag       string `json:"etag"`
}

// CompleteMultipartUploadRequest assembles uploaded parts into the final object.
type CompleteMultipartUploadRequest struct {
	Key      string          `json:"key"`
	UploadID string          `json:"upload_id"`
	Parts    []MultipartPart `json:"parts"`
}

// AbortMultipartUploadRequest cancels an in-progress multipart upload.
type AbortMultipartUploadRequest struct {
	Key      string `json:"key"`
	UploadID string `json:"upload_id"`
}

// BucketCORSRule is one cross-origin resource sharing rule on a bucket.
type BucketCORSRule struct {
	ID             string   `json:"id"`
	AllowedOrigins []string `json:"allowed_origins"`
	// Each method is one of: GET, PUT, POST, DELETE, HEAD.
	AllowedMethods []string `json:"allowed_methods"`
	AllowedHeaders []string `json:"allowed_headers"`
	ExposeHeaders  []string `json:"expose_headers"`
	MaxAgeSeconds  int      `json:"max_age_seconds"`
}

// BucketCORSConfig holds a bucket's CORS rules.
type BucketCORSConfig struct {
	Rules []BucketCORSRule `json:"rules"`
}

// BucketLifecycleRule is one expiration or cleanup rule on a bucket.
type BucketLifecycleRule struct {
	ID                           string `json:"id"`
	Prefix                       string `json:"prefix"`
	Enabled                      bool   `json:"enabled"`
	ExpirationDays               int    `json:"expiration_days"`
	NoncurrentExpirationDays     int    `json:"noncurrent_expiration_days"`
	AbortIncompleteMultipartDays int    `json:"abort_incomplete_multipart_days"`
	ExpiredObjectDeleteMarker    bool   `json:"expired_object_delete_marker"`
}

// BucketLifecycleConfig holds a bucket's lifecycle rules.
type BucketLifecycleConfig struct {
	Rules []BucketLifecycleRule `json:"rules"`
}

// BucketTagsConfig holds a bucket's key-value tags.
type BucketTagsConfig struct {
	Tags []StorageTag `json:"tags"`
}

// BucketObjectLockConfig holds default object-lock retention settings for a bucket.
type BucketObjectLockConfig struct {
	Enabled bool `json:"enabled"`
	// One of: GOVERNANCE, COMPLIANCE.
	Mode  string `json:"mode"`
	Days  int    `json:"days"`
	Years int    `json:"years"`
}
