package functions

import (
	"path/filepath"
	"strings"
)

func GetFileExtension(fileName string) string {
	ext := filepath.Ext(fileName)
	// Delete dot from extension (.img --> img)
	ext = strings.TrimPrefix(ext, ".")
	return ext
}
