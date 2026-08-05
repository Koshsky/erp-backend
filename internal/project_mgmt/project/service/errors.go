package service

import "errors"

var (
	// ErrForbidden возвращается, когда у пользователя нет прав на операцию с проектом.
	ErrForbidden = errors.New("forbidden")
	// ErrNotFound возвращается, когда проект не существует или недоступен пользователю.
	ErrNotFound = errors.New("project not found")
)
