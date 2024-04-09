package functions

import (
	"path/filepath"
	"strings"
)

func GetFileExtension(fileName string) string {
	ext := filepath.Ext(fileName)
	// Удаляем точку из расширения, если она есть
	ext = strings.TrimPrefix(ext, ".")
	return ext
}
