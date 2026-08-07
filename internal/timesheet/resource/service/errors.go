package service

import "errors"

var (
	// ErrForbidden возвращается, когда пользователь не владеет ресурсом
	// (не admin и ресурс не в его собственности).
	ErrForbidden = errors.New("forbidden")
	// ErrNotFound возвращается, когда ресурс не существует или удалён.
	ErrNotFound = errors.New("resource not found")
)
