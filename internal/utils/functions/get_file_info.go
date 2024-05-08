package functions

import (
	"net/http"
	"path/filepath"
	"strings"
)

// GetFileMimeType - method returns mime-type of file
func GetFileMimeType(buffer []byte) (string, error) {
	bufferCopy := make([]byte, 512)
	copy(bufferCopy, buffer[0:512])
	// content-type is "application/octet-stream" if no others seemed to match.
	return http.DetectContentType(bufferCopy), nil
}

// GetFileExtension - method returns file extension from mime-type
func GetFileExtension(fileName string) string {
	ext := filepath.Ext(fileName)
	// Удаляем точку из расширения, если она есть
	ext = strings.TrimPrefix(ext, ".")
	return ext
}
