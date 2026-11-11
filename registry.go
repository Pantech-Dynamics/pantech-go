package pantech

import (
	"context"
	"iter"
	"net/url"
)

// DefaultRegistryHost is the Pantech container registry. Log in to it with
// `docker login registry.pantechdynamics.com -u pantech --password-stdin`,
// giving an API key or a registry token's Secret on stdin.
const DefaultRegistryHost = "registry.pantechdynamics.com"

// Registry repository visibilities.
const (
	RegistryPrivate = "private"
	RegistryPublic  = "public"
)

// Registry repository statuses.
const (
	RegistryRepositoryActive    = "active"
	RegistryRepositoryDeleting  = "deleting"
	RegistryRepositorySuspended = "suspended"
)

// Registry token permissions. Push includes pull.
const (
	RegistryPull = "pull"
	RegistryPush = "push"
)

// RegistryService manages the organization's private container registry:
// its namespace, repositories, tags and the registry tokens CI and clusters
// log in with. Images are pushed and pulled with docker or any OCI client,
// not through this API.
type RegistryService struct{ c *Client }

// GetNamespace gets the organization's namespace: GET /registry. It answers
// REGISTRY_NAMESPACE_NOT_SET (404) until one is claimed with SetNamespace.
// Scope: read.
func (s *RegistryService) GetNamespace(ctx context.Context, opts ...RequestOption) (*RegistryNamespace, error) {
	var out RegistryNamespace
	if err := s.c.Do(ctx, "GET", "/registry", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetNamespace claims the organization's namespace or changes its settings:
// PUT /registry. The name cannot change once a repository exists.
// Scope: write.
func (s *RegistryService) SetNamespace(ctx context.Context, req SetRegistryNamespaceRequest, opts ...RequestOption) (*RegistryNamespace, error) {
	var out RegistryNamespace
	if err := s.c.Do(ctx, "PUT", "/registry", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListRepositories lists repositories: GET /registry/repositories.
// Scope: read.
func (s *RegistryService) ListRepositories(ctx context.Context, params *ListRegistryParams, opts ...RequestOption) (*Page[RegistryRepository], error) {
	var out Page[RegistryRepository]
	if err := s.c.Do(ctx, "GET", "/registry/repositories", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// AllRepositories iterates over every item of ListRepositories, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *RegistryService) AllRepositories(ctx context.Context, params *ListRegistryParams, opts ...RequestOption) iter.Seq2[RegistryRepository, error] {
	var p ListRegistryParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[RegistryRepository], error) {
		p.Cursor = cursor
		return s.ListRepositories(ctx, &p, opts...)
	})
}

// CreateRepository creates a repository: POST /registry/repositories. With
// the namespace's CreateOnPush on, a first push creates one too.
// Scope: write.
func (s *RegistryService) CreateRepository(ctx context.Context, req CreateRegistryRepositoryRequest, opts ...RequestOption) (*RegistryRepository, error) {
	var out RegistryRepository
	if err := s.c.Do(ctx, "POST", "/registry/repositories", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetRepository gets a repository: GET /registry/repositories/{repository_id}.
// Scope: read.
func (s *RegistryService) GetRepository(ctx context.Context, repositoryID string, opts ...RequestOption) (*RegistryRepository, error) {
	var out RegistryRepository
	if err := s.c.Do(ctx, "GET", repositoryPath(repositoryID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateRepository changes a repository's description, visibility, tag
// immutability or retention: PATCH /registry/repositories/{repository_id}.
// Scope: write.
func (s *RegistryService) UpdateRepository(ctx context.Context, repositoryID string, req UpdateRegistryRepositoryRequest, opts ...RequestOption) (*RegistryRepository, error) {
	var out RegistryRepository
	if err := s.c.Do(ctx, "PATCH", repositoryPath(repositoryID), nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteRepository deletes a repository and every image in it:
// DELETE /registry/repositories/{repository_id}. It is returned deleting.
// Scope: delete.
func (s *RegistryService) DeleteRepository(ctx context.Context, repositoryID string, opts ...RequestOption) (*RegistryRepository, error) {
	var out RegistryRepository
	if err := s.c.Do(ctx, "DELETE", repositoryPath(repositoryID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListTags lists a repository's tags, newest push first:
// GET /registry/repositories/{repository_id}/tags.
// Scope: read.
func (s *RegistryService) ListTags(ctx context.Context, repositoryID string, params *ListRegistryParams, opts ...RequestOption) (*Page[RegistryTag], error) {
	var out Page[RegistryTag]
	if err := s.c.Do(ctx, "GET", repositoryPath(repositoryID)+"/tags", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// AllTags iterates over every item of ListTags, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *RegistryService) AllTags(ctx context.Context, repositoryID string, params *ListRegistryParams, opts ...RequestOption) iter.Seq2[RegistryTag, error] {
	var p ListRegistryParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[RegistryTag], error) {
		p.Cursor = cursor
		return s.ListTags(ctx, repositoryID, &p, opts...)
	})
}

// GetTag gets a tag: GET /registry/repositories/{repository_id}/tags/{tag}.
// Scope: read.
func (s *RegistryService) GetTag(ctx context.Context, repositoryID, tag string, opts ...RequestOption) (*RegistryTag, error) {
	var out RegistryTag
	if err := s.c.Do(ctx, "GET", tagPath(repositoryID, tag), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteTag deletes a tag: DELETE /registry/repositories/{repository_id}/tags/{tag}.
// A repository with immutable tags answers TAG_IMMUTABLE. The space is
// reclaimed by the registry's next garbage collection.
// Scope: delete.
func (s *RegistryService) DeleteTag(ctx context.Context, repositoryID, tag string, opts ...RequestOption) error {
	return s.c.Do(ctx, "DELETE", tagPath(repositoryID, tag), nil, nil, nil, opts...)
}

// ListTokens lists registry tokens, without their secrets: GET /registry/tokens.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *RegistryService) ListTokens(ctx context.Context, opts ...RequestOption) ([]RegistryToken, error) {
	return listAll[RegistryToken](ctx, s.c, "/registry/tokens", nil, opts...)
}

// CreateToken creates a registry token: POST /registry/tokens. Its Secret is
// in this answer only; it is the password for `docker login`.
// Scope: write.
func (s *RegistryService) CreateToken(ctx context.Context, req CreateRegistryTokenRequest, opts ...RequestOption) (*RegistryToken, error) {
	var out RegistryToken
	if err := s.c.Do(ctx, "POST", "/registry/tokens", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteToken revokes a registry token: DELETE /registry/tokens/{token_id}.
// Logins with it stop within 5 minutes.
// Scope: delete.
func (s *RegistryService) DeleteToken(ctx context.Context, tokenID string, opts ...RequestOption) error {
	return s.c.Do(ctx, "DELETE", "/registry/tokens/"+url.PathEscape(tokenID), nil, nil, nil, opts...)
}

func repositoryPath(repositoryID string) string {
	return "/registry/repositories/" + url.PathEscape(repositoryID)
}

func tagPath(repositoryID, tag string) string {
	return repositoryPath(repositoryID) + "/tags/" + url.PathEscape(tag)
}
