package pantech

import (
	"context"
	"net/url"
)

// SSHKeysService manages the SSH keys installed on new instances.
type SSHKeysService struct{ c *Client }

// List lists SSH keys: GET /ssh-keys.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *SSHKeysService) List(ctx context.Context, opts ...RequestOption) ([]SSHKey, error) {
	return listAll[SSHKey](ctx, s.c, "/ssh-keys", nil, opts...)
}

// Create adds or generates an SSH key: POST /ssh-keys.
// Scope: write.
func (s *SSHKeysService) Create(ctx context.Context, req CreateSSHKeyRequest, opts ...RequestOption) (*GeneratedSSHKey, error) {
	var out GeneratedSSHKey
	if err := s.c.Do(ctx, "POST", "/ssh-keys", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get gets an SSH key: GET /ssh-keys/{ssh_key_id}.
// Scope: read.
func (s *SSHKeysService) Get(ctx context.Context, sshKeyID string, opts ...RequestOption) (*SSHKey, error) {
	var out SSHKey
	if err := s.c.Do(ctx, "GET", "/ssh-keys/"+url.PathEscape(sshKeyID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes an SSH key: DELETE /ssh-keys/{ssh_key_id}.
// Scope: write.
func (s *SSHKeysService) Delete(ctx context.Context, sshKeyID string, opts ...RequestOption) error {
	return s.c.Do(ctx, "DELETE", "/ssh-keys/"+url.PathEscape(sshKeyID), nil, nil, nil, opts...)
}
