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
	// `provisioning` (workers being built), `ready` (every worker has joined), `updating` (a scale in progress), `deleting`, or `unavailable` (a pool made while node pools were a preview, with no node size: delete it and create it again). One of: provisioning, ready, updating, deleting, unavailable.
	Status string `json:"status"`
	// The pool's workers that have joined the cluster.
	ObservedWorkers int `json:"observed_workers"`
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
