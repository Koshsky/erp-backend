package service

import "errors"

var (
	// ErrForbidden возвращается, когда пользователь не управляет этим сотрудником
	// (не admin и сотрудник не в его подчинении).
	ErrForbidden = errors.New("forbidden")
	// ErrNotFound возвращается, когда сотрудник не существует или удалён.
	ErrNotFound = errors.New("employee not found")
	// ErrResourceNotFound возвращается, когда должность (ресурс) не найдена или удалена.
	ErrResourceNotFound = errors.New("resource not found")
)
