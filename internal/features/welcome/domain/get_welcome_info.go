package domain

import "context"

// GetWelcomeInfo is the focused use case for the welcome endpoint.
type GetWelcomeInfo struct {
	repo Repository
}

// NewGetWelcomeInfo constructs the use case.
func NewGetWelcomeInfo(repo Repository) *GetWelcomeInfo {
	return &GetWelcomeInfo{repo: repo}
}

// Execute returns welcome information.
func (uc *GetWelcomeInfo) Execute(ctx context.Context) (Info, error) {
	return uc.repo.GetInfo(ctx)
}
