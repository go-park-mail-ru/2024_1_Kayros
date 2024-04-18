package functions

import (
	"mime/multipart"
	"net/http"
	"strings"
)

// GetFileMimeType - method returns mime-type of file
func GetFileMimeType(file multipart.File) (string, error) {
	// only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}
	// content-type is "application/octet-stream" if no others seemed to match.
	return http.DetectContentType(buffer), nil
}

// GetFileExtension - method returns file extension from mime-type
func GetFileExtension(mimeType string) string {
	parts := strings.Split(mimeType, "/")
	return parts[len(parts)-1]
}
