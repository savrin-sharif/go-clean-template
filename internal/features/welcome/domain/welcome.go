// Package domain holds welcome business types and contracts.
package domain

import "context"

// Info is the welcome payload shown by the landing endpoint.
type Info struct {
	AppName     string
	Environment string
	Message     string
}

// Repository loads welcome information.
type Repository interface {
	GetInfo(ctx context.Context) (Info, error)
}
