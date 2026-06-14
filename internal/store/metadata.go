package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Status is the lifecycle state of a client ADP workspace.
type Status string

const (
	// StatusDraft: newly created, no ADP generated yet (the default).
	StatusDraft Status = "draft"
	// StatusUpdating: a skill is iterating the knowledge base or ADP.
	StatusUpdating Status = "updating"
	// StatusReady: the ADP output is generated and current.
	StatusReady Status = "ready"
	// StatusStale: new materials arrived but the ADP has not been regenerated.
	StatusStale Status = "stale"
)

// Client is the per-client metadata persisted at <workspace>/metadata.json.
// Field keys are English by design; values may be Chinese.
type Client struct {
	Name           string    `json:"name"`
	Owner          string    `json:"owner"`
	Stage          string    `json:"stage"`
	Status         Status    `json:"status"`
	Created        time.Time `json:"created"`
	Updated        time.Time `json:"updated"`
	Model          string    `json:"model,omitempty"`
	MaterialsCount int       `json:"materials_count"`
}

// MetadataPath returns the path to a client's metadata.json.
func MetadataPath(workspace string) string {
	return filepath.Join(workspace, "metadata.json")
}

// ReadMetadata loads a client's metadata. A missing file yields a minimal
// Client (StatusDraft) so pre-metadata workspaces degrade gracefully.
func ReadMetadata(workspace string) (Client, error) {
	c := Client{
		Name:    filepath.Base(workspace),
		Status:  StatusDraft,
		Created: time.Now().UTC(),
		Updated: time.Now().UTC(),
	}
	data, err := os.ReadFile(MetadataPath(workspace))
	if err != nil {
		if os.IsNotExist(err) {
			return c, nil
		}
		return c, err
	}
	if err := json.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("parse %s: %w", MetadataPath(workspace), err)
	}
	if c.Status == "" {
		c.Status = StatusDraft
	}
	if c.Name == "" {
		c.Name = filepath.Base(workspace)
	}
	return c, nil
}

// WriteMetadata persists a client's metadata, refreshing Updated. Created is
// preserved from the value passed in.
func WriteMetadata(workspace string, c Client) error {
	c.Updated = time.Now().UTC()
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		return err
	}
	return os.WriteFile(MetadataPath(workspace), data, 0o644)
}
