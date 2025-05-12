package dto

import "errors"

var (
	ErrNoFieldsToUpdate = errors.New("no fields to update")
	ErrNotFound         = errors.New("no record found")
	ErrNoAccess         = errors.New("no access to this resource")
	ErrBadRequest       = errors.New("bad request")
)
