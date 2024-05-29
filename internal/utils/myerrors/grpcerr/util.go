package grpcerr

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Define a Error type with a message and a status code
type Error struct {
	Message string     `json:"message"`
	Status  codes.Code `json:"-"`
}

// Implement the inbuilt error interface
func (e Error) Error() string {
	return fmt.Sprintf("Error: %s", e.Message)
}

// This member function is used by grpc when converting an error into a status
func (e Error) GRPCStatus() *status.Status {
	return status.New(e.Status, e.Error())
}

func NewError(status codes.Code, msg string) error {
	return &Error{
		Status:  status,
		Message: msg,
	}
}

func NewResponse(status codes.Code, msg string) error {
	return &Error{
		Status:  status,
		Message: msg,
	}
}

func removeErrorPrefix(grpcStatusMsg string) string {
	return strings.Replace(grpcStatusMsg, "Error: ", "", 1)
}

func Is(responseErr error, statusCode codes.Code, errClient error) bool {
	grpcErr, ok := status.FromError(responseErr)
	if ok {
		return grpcErr.Code() == statusCode && removeErrorPrefix(grpcErr.Message()) == errClient.Error()
	}
	return false
}
