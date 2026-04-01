package handoff

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Handoff represents a task handoff between agents.
type Handoff struct {
	From         string   `json:"from"`
	To           string   `json:"to"`
	Status       string   `json:"status"`
	Feature      string   `json:"feature"`
	Summary      string   `json:"summary"`
	Artifacts    []string `json:"artifacts"`
	Instructions string   `json:"instructions"`
	Timestamp    time.Time `json:"timestamp"`
}

// Validate returns an error if any required field is empty.
func (h Handoff) Validate() error {
	if h.From == "" {
		return fmt.Errorf("handoff missing required field: From")
	}
	if h.To == "" {
		return fmt.Errorf("handoff missing required field: To")
	}
	if h.Status == "" {
		return fmt.Errorf("handoff missing required field: Status")
	}
	if h.Feature == "" {
		return fmt.Errorf("handoff missing required field: Feature")
	}
	return nil
}

// Write serializes a Handoff to the given file path.
func Write(path string, h Handoff) error {
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal handoff: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// Read deserializes a Handoff from the given file path.
func Read(path string) (Handoff, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Handoff{}, fmt.Errorf("read handoff: %w", err)
	}
	var h Handoff
	if err := json.Unmarshal(data, &h); err != nil {
		return Handoff{}, fmt.Errorf("unmarshal handoff: %w", err)
	}
	return h, nil
}
