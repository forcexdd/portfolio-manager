package repository

import "errors"

var (
	ErrAssetAlreadyExists = errors.New("asset already exists")
	ErrAssetNotFound      = errors.New("asset not found")

	ErrIndexAlreadyExists = errors.New("index already exists")
	ErrIndexNotFound      = errors.New("index not found")

	ErrPortfolioAlreadyExists = errors.New("portfolio already exists")
	ErrPortfolioNotFound      = errors.New("portfolio not found")
)
