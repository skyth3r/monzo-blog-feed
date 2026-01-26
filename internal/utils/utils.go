package utils

import (
	"fmt"
	"os"
	"strings"
)

func FormatText(text string) string {
	text = strings.NewReplacer(
		"’", "'",
	).Replace(text)

	text = strings.NewReplacer(
		"–", "-",
	).Replace(text)
	return text
}

func MoveFile(fileName, filePath string) error {
	if err := os.Rename(fileName, filePath); err != nil {
		return fmt.Errorf("unable to move %s to '%s': %w", fileName, filePath, err)
	}
	return nil
}
