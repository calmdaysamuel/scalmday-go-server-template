package permissions

import (
	"context"
	"encoding/json"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/palantir/pkg/bearertoken"
	werror "github.com/palantir/witchcraft-go-error"
	"scalmday-go-server-template/config"
)

type User struct {
	Aal string `json:"aal"`
	Amr []struct {
		Method    string `json:"method"`
		Timestamp int    `json:"timestamp"`
	} `json:"amr"`
	AppMetadata struct {
		Provider  string   `json:"provider"`
		Providers []string `json:"providers"`
	} `json:"app_metadata"`
	Aud          string `json:"aud"`
	Email        string `json:"email"`
	Exp          int    `json:"exp"`
	Iat          int    `json:"iat"`
	IsAnonymous  bool   `json:"is_anonymous"`
	Iss          string `json:"iss"`
	Phone        string `json:"phone"`
	Role         string `json:"role"`
	SessionID    string `json:"session_id"`
	Sub          string `json:"sub"`
	UserMetadata struct {
		EmailVerified bool `json:"email_verified"`
	} `json:"user_metadata"`
	claims jwt.Claims
}

func (u *User) ID() string {
	return u.Sub
}

func (u *User) IsAnon() bool {
	return u.IsAnonymous
}

type UserProvider interface {
	GetUser(ctx context.Context, token bearertoken.Token) (User, error)
}

func NewUserProvider(ctx context.Context, config config.InstallConfig) (UserProvider, error) {
	k, err := keyfunc.NewDefaultCtx(ctx, config.TrustedPublicKeyURLs) // Context is used to end the refresh goroutine.
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "Failed to create a keyfunc.Keyfunc from the server's URL.")
	}
	return &userProvider{keyFunc: k}, nil
}

type userProvider struct {
	keyFunc keyfunc.Keyfunc
}

func (u *userProvider) GetUser(ctx context.Context, token bearertoken.Token) (User, error) {
	// Parse the JWT.
	parsed, err := jwt.Parse(string(token), u.keyFunc.Keyfunc)
	if err != nil {
		return User{}, werror.WrapWithContextParams(ctx, err, "failed to verify jwt")
	}
	claimsBytes, err := json.Marshal(parsed.Claims)
	if err != nil {
		return User{}, werror.WrapWithContextParams(ctx, err, "failed to unmarshal claims bytes")
	}
	var user User
	if err := json.Unmarshal(claimsBytes, &user); err != nil {
		return User{}, werror.WrapWithContextParams(ctx, err, "failed to unmarshal claims bytes to a verified user")
	}
	user.claims = parsed.Claims

	return user, nil
}
