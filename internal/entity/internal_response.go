package entity

import "errors"

type DataType any
type ErrorType any

// RaiseError возвращает объект ошибки с сообщением об ошибке
func RaiseError(message string) (DataType, ErrorType) {
	return nil, errors.New(message)
}

// GenerateResponse возвращает объект с полезными данными
func GenerateResponse(dataResponse any) (DataType, ErrorType) {
	return dataResponse, nil
}
