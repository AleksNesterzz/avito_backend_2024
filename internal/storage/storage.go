package storage

import (
	"errors"
)

var (
	ErrBannerNotFound  = errors.New("banner not found")
	ErrBannerNotExists = errors.New("banner not exists")
	ErrBannerExists    = errors.New("banner already exists")
	ErrUserNotFound    = errors.New("user not found")
)
