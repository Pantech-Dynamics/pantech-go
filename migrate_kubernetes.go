package pantech

import (
	"context"
	"net/url"
)

// Kubernetes move statuses. Completed, failed and cancelled are final.
const (
	KubernetesMoveExporting   = "exporting"
	KubernetesMoveReady       = "ready"
	KubernetesMovePreviewing  = "previewing"
	KubernetesMoveApplying    = "applying"
	KubernetesMoveApplied     = "applied"
	KubernetesMoveSyncing     = "syncing"
	KubernetesMoveCuttingOver = "cutting_over"
	KubernetesMoveCompleted   = "completed"
	KubernetesMoveFailed      = "failed"
	KubernetesMoveCancelled   = "cancelled"
)

// StartKubernetesMove starts moving a plan's kubernetes_cluster item into one
// of your running Pantech clusters: POST /migrations/{migration_id}/kubernetes-moves.
// The Pantech Migrate agent exports the source cluster with the read-only
// access from GetKubernetesMoveSourceAccess. Secrets move only with
// IncludeSecrets, volume data only with CopyVolumes. The move is free; the
// cluster bills as usual.
// Scope: write.
func (s *MigrateService) StartKubernetesMove(ctx context.Context, migrationID string, req StartKubernetesMoveRequest, opts ...RequestOption) (*KubernetesMove, error) {
	var out KubernetesMove
	if err := s.c.Do(ctx, "POST", "/migrations/"+url.PathEscape(migrationID)+"/kubernetes-moves", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListKubernetesMoves lists the organization's moves, newest first, without
// their resources and volumes: GET /kubernetes-moves.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *MigrateService) ListKubernetesMoves(ctx context.Context, opts ...RequestOption) ([]KubernetesMove, error) {
	return listAll[KubernetesMove](ctx, s.c, "/kubernetes-moves", nil, opts...)
}

// GetKubernetesMoveSourceAccess returns the manifests to apply to the source
// cluster: a read-only reader, an optional Secrets reader and an optional
// copier role: GET /kubernetes-moves/source-access.
// Scope: read.
func (s *MigrateService) GetKubernetesMoveSourceAccess(ctx context.Context, opts ...RequestOption) (*KubernetesMoveSourceAccess, error) {
	var out KubernetesMoveSourceAccess
	if err := s.c.Do(ctx, "GET", "/kubernetes-moves/source-access", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetKubernetesMove gets a move with its report: every source object and what
// the move does with it (names only, never values), the volumes and their
// copy progress, the images and the workloads held until the cutover:
// GET /kubernetes-moves/{move_id}.
// Scope: read.
func (s *MigrateService) GetKubernetesMove(ctx context.Context, moveID string, opts ...RequestOption) (*KubernetesMoveDetail, error) {
	var out KubernetesMoveDetail
	if err := s.c.Do(ctx, "GET", "/kubernetes-moves/"+url.PathEscape(moveID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// PreviewKubernetesMove queues a server-side dry run of every object against
// the target cluster: POST /kubernetes-moves/{move_id}/preview. The move's
// namespaces are created (empty) so objects inside them can be checked;
// nothing else is created.
// Scope: write.
func (s *MigrateService) PreviewKubernetesMove(ctx context.Context, moveID string, req PreviewKubernetesMoveRequest, opts ...RequestOption) (*KubernetesMove, error) {
	return s.actOnKubernetesMove(ctx, moveID, "preview", &req, opts)
}

// ApplyKubernetesMove queues the apply of what the last successful preview
// checked: POST /kubernetes-moves/{move_id}/apply. Workloads that mount a
// copied volume start scaled down until the cutover.
// Scope: write.
func (s *MigrateService) ApplyKubernetesMove(ctx context.Context, moveID string, opts ...RequestOption) (*KubernetesMove, error) {
	return s.actOnKubernetesMove(ctx, moveID, "apply", nil, opts)
}

// SyncKubernetesMoveVolumes queues the next copy of every volume; only what
// changed since the last copy is sent: POST /kubernetes-moves/{move_id}/sync-volumes.
// Scope: write.
func (s *MigrateService) SyncKubernetesMoveVolumes(ctx context.Context, moveID string, opts ...RequestOption) (*KubernetesMove, error) {
	return s.actOnKubernetesMove(ctx, moveID, "sync-volumes", nil, opts)
}

// CutOverKubernetesMove takes the final copy of every volume, then scales the
// held workloads up to their source replicas: POST /kubernetes-moves/{move_id}/cutover.
// Stop writes on the source first and set req.SourceWritesStopped.
// Scope: write.
func (s *MigrateService) CutOverKubernetesMove(ctx context.Context, moveID string, req CutoverKubernetesMoveRequest, opts ...RequestOption) (*KubernetesMove, error) {
	return s.actOnKubernetesMove(ctx, moveID, "cutover", &req, opts)
}

// CancelKubernetesMove stops a move and removes its copy pods, staged data
// and export: POST /kubernetes-moves/{move_id}/cancel. Objects already
// applied to your Pantech cluster stay.
// Scope: write.
func (s *MigrateService) CancelKubernetesMove(ctx context.Context, moveID string, opts ...RequestOption) (*KubernetesMove, error) {
	return s.actOnKubernetesMove(ctx, moveID, "cancel", nil, opts)
}

// DownloadKubernetesMoveManifests returns the transformed manifests as one
// YAML file, for GitOps: GET /kubernetes-moves/{move_id}/manifests. Secret
// values are never included.
// Scope: read.
func (s *MigrateService) DownloadKubernetesMoveManifests(ctx context.Context, moveID string, opts ...RequestOption) ([]byte, error) {
	var out []byte
	if err := s.c.Do(ctx, "GET", "/kubernetes-moves/"+url.PathEscape(moveID)+"/manifests", nil, nil, &out, append([]RequestOption{WithRequestHeader("Accept", "application/yaml")}, opts...)...); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *MigrateService) actOnKubernetesMove(ctx context.Context, moveID, action string, body any, opts []RequestOption) (*KubernetesMove, error) {
	var out KubernetesMove
	if err := s.c.Do(ctx, "POST", "/kubernetes-moves/"+url.PathEscape(moveID)+"/"+action, nil, body, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
