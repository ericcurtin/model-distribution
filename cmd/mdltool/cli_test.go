package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCLIDefaultRegistry(t *testing.T) {
	// Build the mdltool binary for testing
	buildCmd := exec.Command("go", "build", "-o", "test-mdltool", ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build test binary: %v", err)
	}
	defer os.Remove("test-mdltool")

	tests := []struct {
		name     string
		args     []string
		contains string
	}{
		{
			name:     "help shows default-registry option",
			args:     []string{"--help"},
			contains: "-default-registry",
		},
		{
			name:     "version works",
			args:     []string{"--version"},
			contains: "model-distribution-tool version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./test-mdltool", tt.args...)
			output, err := cmd.CombinedOutput()
			
			// For help and version, these should exit successfully
			if tt.name == "help shows default-registry option" || tt.name == "version works" {
				if err != nil {
					t.Errorf("Command failed: %v, output: %s", err, output)
				}
			}

			if !strings.Contains(string(output), tt.contains) {
				t.Errorf("Expected output to contain %q, got: %s", tt.contains, output)
			}
		})
	}
}