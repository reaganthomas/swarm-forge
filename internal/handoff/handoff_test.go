package handoff_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/swarm-forge/swarm-forge/internal/handoff"
)

func TestWriteAndRead(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "handoff.json")

	ts := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)
	h := handoff.Handoff{
		From:         "architect",
		To:           "e2e-interpreter",
		Status:       "pending",
		Feature:      "features/login.feature",
		Summary:      "Login feature spec",
		Artifacts:    []string{"features/login.feature"},
		Instructions: "Generate acceptance tests",
		Timestamp:    ts,
	}

	if err := handoff.Write(path, h); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	got, err := handoff.Read(path)
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	if got.From != h.From {
		t.Errorf("From = %q, want %q", got.From, h.From)
	}
	if got.To != h.To {
		t.Errorf("To = %q, want %q", got.To, h.To)
	}
	if got.Status != h.Status {
		t.Errorf("Status = %q, want %q", got.Status, h.Status)
	}
	if got.Feature != h.Feature {
		t.Errorf("Feature = %q, want %q", got.Feature, h.Feature)
	}
	if got.Summary != h.Summary {
		t.Errorf("Summary = %q, want %q", got.Summary, h.Summary)
	}
	if len(got.Artifacts) != 1 || got.Artifacts[0] != h.Artifacts[0] {
		t.Errorf("Artifacts = %v, want %v", got.Artifacts, h.Artifacts)
	}
	if got.Instructions != h.Instructions {
		t.Errorf("Instructions = %q, want %q", got.Instructions, h.Instructions)
	}
	if !got.Timestamp.Equal(h.Timestamp) {
		t.Errorf("Timestamp = %v, want %v", got.Timestamp, h.Timestamp)
	}
}

func TestReadNonexistentFile(t *testing.T) {
	_, err := handoff.Read("/nonexistent/path/handoff.json")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestWriteInvalidPath(t *testing.T) {
	h := handoff.Handoff{From: "test"}
	err := handoff.Write("/nonexistent/dir/handoff.json", h)
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestReadCorruptedJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0o644)

	_, err := handoff.Read(path)
	if err == nil {
		t.Fatal("expected error for corrupted JSON")
	}
}

func TestValidateValidHandoff(t *testing.T) {
	h := handoff.Handoff{
		From:    "architect",
		To:      "coder",
		Status:  "pending",
		Feature: "features/test.feature",
	}
	if err := h.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateEmptyFrom(t *testing.T) {
	h := handoff.Handoff{
		From:    "",
		To:      "coder",
		Status:  "pending",
		Feature: "features/test.feature",
	}
	err := h.Validate()
	if err == nil {
		t.Fatal("expected error for empty From")
	}
}

func TestValidateEmptyTo(t *testing.T) {
	h := handoff.Handoff{
		From:    "architect",
		To:      "",
		Status:  "pending",
		Feature: "features/test.feature",
	}
	err := h.Validate()
	if err == nil {
		t.Fatal("expected error for empty To")
	}
}

func TestValidateEmptyStatus(t *testing.T) {
	h := handoff.Handoff{
		From:    "architect",
		To:      "coder",
		Status:  "",
		Feature: "features/test.feature",
	}
	err := h.Validate()
	if err == nil {
		t.Fatal("expected error for empty Status")
	}
}

func TestValidateEmptyFeature(t *testing.T) {
	h := handoff.Handoff{
		From:    "architect",
		To:      "coder",
		Status:  "pending",
		Feature: "",
	}
	err := h.Validate()
	if err == nil {
		t.Fatal("expected error for empty Feature")
	}
}

func TestValidateAllEmpty(t *testing.T) {
	h := handoff.Handoff{}
	err := h.Validate()
	if err == nil {
		t.Fatal("expected error for all-empty handoff")
	}
}
