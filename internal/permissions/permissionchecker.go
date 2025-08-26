package permissions

import (
	"context"
	"github.com/palantir/pkg/bearertoken"
)

type Manager interface {
	IsLoggedIn(ctx context.Context, token bearertoken.Token) bool
}

type permissionManager struct {
	userprovider UserProvider
}

func NewManager(ctx context.Context, provider UserProvider) Manager {
	return &permissionManager{userprovider: provider}
}

func (p *permissionManager) IsLoggedIn(ctx context.Context, token bearertoken.Token) bool {
	user, err := p.userprovider.GetUser(ctx, token)
	return err != nil && !user.IsAnon()
}
