package pantech

import (
	"context"
	"iter"
	"net/url"
)

// InstancesService manages instances (VPS) and their orders.
type InstancesService struct{ c *Client }

// List lists instances: GET /instances.
// Scope: read.
func (s *InstancesService) List(ctx context.Context, params *ListParams, opts ...RequestOption) (*Page[Instance], error) {
	var out Page[Instance]
	if err := s.c.Do(ctx, "GET", "/instances", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// All iterates over every item of List, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *InstancesService) All(ctx context.Context, params *ListParams, opts ...RequestOption) iter.Seq2[Instance, error] {
	var p ListParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[Instance], error) {
		p.Cursor = cursor
		return s.List(ctx, &p, opts...)
	})
}

// Create creates an instance: POST /instances.
// Scope: write.
func (s *InstancesService) Create(ctx context.Context, req CreateInstanceRequest, opts ...RequestOption) (*OrderAccepted, error) {
	var out OrderAccepted
	if err := s.c.Do(ctx, "POST", "/instances", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get gets an instance: GET /instances/{instance_id}.
// Scope: read.
func (s *InstancesService) Get(ctx context.Context, instanceID string, opts ...RequestOption) (*Instance, error) {
	var out Instance
	if err := s.c.Do(ctx, "GET", "/instances/"+url.PathEscape(instanceID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes an instance: DELETE /instances/{instance_id}.
// Scope: write.
func (s *InstancesService) Delete(ctx context.Context, instanceID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "DELETE", "/instances/"+url.PathEscape(instanceID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Start starts an instance: POST /instances/{instance_id}/start.
// Scope: write.
func (s *InstancesService) Start(ctx context.Context, instanceID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/instances/"+url.PathEscape(instanceID)+"/start", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Stop stops an instance: POST /instances/{instance_id}/stop.
// Scope: write.
func (s *InstancesService) Stop(ctx context.Context, instanceID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/instances/"+url.PathEscape(instanceID)+"/stop", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Reboot reboots an instance: POST /instances/{instance_id}/reboot.
// Scope: write.
func (s *InstancesService) Reboot(ctx context.Context, instanceID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/instances/"+url.PathEscape(instanceID)+"/reboot", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Rebuild rebuilds an instance: POST /instances/{instance_id}/rebuild.
// Scope: write.
func (s *InstancesService) Rebuild(ctx context.Context, instanceID string, req RebuildRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/instances/"+url.PathEscape(instanceID)+"/rebuild", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Rename renames an instance: POST /instances/{instance_id}/rename.
// Scope: write.
func (s *InstancesService) Rename(ctx context.Context, instanceID string, req RenameRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/instances/"+url.PathEscape(instanceID)+"/rename", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Resize resizes an instance: POST /instances/{instance_id}/resize.
// Scope: write.
func (s *InstancesService) Resize(ctx context.Context, instanceID string, req ResizeRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/instances/"+url.PathEscape(instanceID)+"/resize", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ChangeSSHKey changes an instance's SSH key: POST /instances/{instance_id}/change-ssh-key.
// Scope: write.
func (s *InstancesService) ChangeSSHKey(ctx context.Context, instanceID string, req ChangeSSHKeyRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/instances/"+url.PathEscape(instanceID)+"/change-ssh-key", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Restore restores an instance from a snapshot: POST /instances/{instance_id}/restore.
// Scope: write.
func (s *InstancesService) Restore(ctx context.Context, instanceID string, req RestoreRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/instances/"+url.PathEscape(instanceID)+"/restore", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ChangeSecurityGroup changes an instance's security group: PUT /instances/{instance_id}/security-group.
// Scope: write.
func (s *InstancesService) ChangeSecurityGroup(ctx context.Context, instanceID string, req ChangeSecurityGroupRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "PUT", "/instances/"+url.PathEscape(instanceID)+"/security-group", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// AttachPrivateNetwork adds a private database network interface to an instance: POST /instances/{instance_id}/private-network.
// Scope: write.
func (s *InstancesService) AttachPrivateNetwork(ctx context.Context, instanceID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/instances/"+url.PathEscape(instanceID)+"/private-network", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// DetachPrivateNetwork removes an instance's private database network interface: DELETE /instances/{instance_id}/private-network.
// Scope: write.
func (s *InstancesService) DetachPrivateNetwork(ctx context.Context, instanceID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "DELETE", "/instances/"+url.PathEscape(instanceID)+"/private-network", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// OpenConsole opens the browser console: POST /instances/{instance_id}/console.
// Scope: write.
func (s *InstancesService) OpenConsole(ctx context.Context, instanceID string, opts ...RequestOption) (*InstanceConsole, error) {
	var out InstanceConsole
	if err := s.c.Do(ctx, "POST", "/instances/"+url.PathEscape(instanceID)+"/console", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// IssueVolumeKeyCredential issues a volume key credential: POST /instances/{instance_id}/volume-key-credential.
// Scope: write.
func (s *InstancesService) IssueVolumeKeyCredential(ctx context.Context, instanceID string, opts ...RequestOption) (*VolumeKeyCredential, error) {
	var out VolumeKeyCredential
	if err := s.c.Do(ctx, "POST", "/instances/"+url.PathEscape(instanceID)+"/volume-key-credential", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMetrics gets instance metrics: GET /instances/{instance_id}/metrics.
// Scope: read.
func (s *InstancesService) GetMetrics(ctx context.Context, instanceID string, params *InstanceMetricsParams, opts ...RequestOption) (*MetricsSeries, error) {
	var out MetricsSeries
	if err := s.c.Do(ctx, "GET", "/instances/"+url.PathEscape(instanceID)+"/metrics", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetUsage gets an instance's usage: GET /instances/{instance_id}/usage.
// Scope: read.
func (s *InstancesService) GetUsage(ctx context.Context, instanceID string, opts ...RequestOption) (*UsageSummary, error) {
	var out UsageSummary
	if err := s.c.Do(ctx, "GET", "/instances/"+url.PathEscape(instanceID)+"/usage", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListVolumes lists an instance's volumes: GET /instances/{instance_id}/volumes.
// Scope: read.
func (s *InstancesService) ListVolumes(ctx context.Context, instanceID string, params *ListParams, opts ...RequestOption) (*Page[Volume], error) {
	var out Page[Volume]
	if err := s.c.Do(ctx, "GET", "/instances/"+url.PathEscape(instanceID)+"/volumes", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// AllVolumes iterates over every item of ListVolumes, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *InstancesService) AllVolumes(ctx context.Context, instanceID string, params *ListParams, opts ...RequestOption) iter.Seq2[Volume, error] {
	var p ListParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[Volume], error) {
		p.Cursor = cursor
		return s.ListVolumes(ctx, instanceID, &p, opts...)
	})
}

// ListSnapshots lists an instance's snapshots: GET /instances/{instance_id}/snapshots.
// Scope: read.
func (s *InstancesService) ListSnapshots(ctx context.Context, instanceID string, params *ListParams, opts ...RequestOption) (*Page[Snapshot], error) {
	var out Page[Snapshot]
	if err := s.c.Do(ctx, "GET", "/instances/"+url.PathEscape(instanceID)+"/snapshots", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// AllSnapshots iterates over every item of ListSnapshots, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *InstancesService) AllSnapshots(ctx context.Context, instanceID string, params *ListParams, opts ...RequestOption) iter.Seq2[Snapshot, error] {
	var p ListParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[Snapshot], error) {
		p.Cursor = cursor
		return s.ListSnapshots(ctx, instanceID, &p, opts...)
	})
}

// CreateSnapshot snapshots an instance: POST /instances/{instance_id}/snapshots.
// Scope: write.
func (s *InstancesService) CreateSnapshot(ctx context.Context, instanceID string, req CreateSnapshotRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/instances/"+url.PathEscape(instanceID)+"/snapshots", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetSnapshotSchedule gets an instance's snapshot schedule: GET /instances/{instance_id}/snapshot-schedule.
// Scope: read.
func (s *InstancesService) GetSnapshotSchedule(ctx context.Context, instanceID string, opts ...RequestOption) (*SnapshotSchedule, error) {
	var out SnapshotSchedule
	if err := s.c.Do(ctx, "GET", "/instances/"+url.PathEscape(instanceID)+"/snapshot-schedule", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// PutSnapshotSchedule sets an instance's snapshot schedule: PUT /instances/{instance_id}/snapshot-schedule.
// Scope: write.
func (s *InstancesService) PutSnapshotSchedule(ctx context.Context, instanceID string, req PutSnapshotScheduleRequest, opts ...RequestOption) (*SnapshotSchedule, error) {
	var out SnapshotSchedule
	if err := s.c.Do(ctx, "PUT", "/instances/"+url.PathEscape(instanceID)+"/snapshot-schedule", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// OpenSSHAccess grants temporary SSH access: POST /instances/{instance_id}/ssh-access.
// Scope: write.
func (s *InstancesService) OpenSSHAccess(ctx context.Context, instanceID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/instances/"+url.PathEscape(instanceID)+"/ssh-access", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetSSHAccess gets an SSH access grant: GET /instances/{instance_id}/ssh-access/{grant_id}.
// Scope: read.
func (s *InstancesService) GetSSHAccess(ctx context.Context, instanceID string, grantID string, opts ...RequestOption) (*SSHAccessGrant, error) {
	var out SSHAccessGrant
	if err := s.c.Do(ctx, "GET", "/instances/"+url.PathEscape(instanceID)+"/ssh-access/"+url.PathEscape(grantID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// RevokeSSHAccess revokes SSH access: DELETE /instances/{instance_id}/ssh-access/{grant_id}.
// Scope: write.
func (s *InstancesService) RevokeSSHAccess(ctx context.Context, instanceID string, grantID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "DELETE", "/instances/"+url.PathEscape(instanceID)+"/ssh-access/"+url.PathEscape(grantID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListOrders lists instance orders: GET /instance-orders.
// Scope: read.
func (s *InstancesService) ListOrders(ctx context.Context, params *PageParams, opts ...RequestOption) (*Page[InstanceOrder], error) {
	var out Page[InstanceOrder]
	if err := s.c.Do(ctx, "GET", "/instance-orders", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// AllOrders iterates over every item of ListOrders, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *InstancesService) AllOrders(ctx context.Context, params *PageParams, opts ...RequestOption) iter.Seq2[InstanceOrder, error] {
	var p PageParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[InstanceOrder], error) {
		p.Cursor = cursor
		return s.ListOrders(ctx, &p, opts...)
	})
}

// GetOrder gets an instance order: GET /instance-orders/{order_id}.
// Scope: read.
func (s *InstancesService) GetOrder(ctx context.Context, orderID string, opts ...RequestOption) (*InstanceOrder, error) {
	var out InstanceOrder
	if err := s.c.Do(ctx, "GET", "/instance-orders/"+url.PathEscape(orderID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
