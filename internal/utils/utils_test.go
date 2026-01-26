package utils_test

import (
	"testing"

	"github.com/skyth3r/monzo-blog-feed/internal/utils"
)

func FormatText_Test(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello, world!", "Hello, world!"},
		{"It’s a test", "It's a test"},
		{"– is a dash", "- is a dash"},
	}

	for _, test := range tests {
		result := utils.FormatText(test.input)
		if result != test.expected {
			t.Errorf("expected %q, got %q", test.expected, result)
		}
	}
}

func MoveFile_Test(t *testing.T) {
	// This test is not implemented as it requires file system operations
	// and is not suitable for unit testing without a mock file system.
	t.Skip("MoveFile test is skipped because it requires file system operations.")
}
