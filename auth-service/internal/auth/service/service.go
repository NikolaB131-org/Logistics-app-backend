package service

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/Nerzal/gocloak/v13/pkg/jwx"
	"github.com/NikolaB131-org/logistics-backend/auth-service/otlp"
)

type Auth struct {
	gocloak             *gocloak.GoCloak
	gocloakClientID     string
	gocloakClientSecret string
	gocloakRealm        string
}

type AccessTokenCustomClaims struct {
	jwx.Claims
}

type TokenClaims struct {
	UserID     string
	Username   string
	RealmRoles []string
}

var (
	ErrUserAlreadyExists = errors.New("user with this email already exists")
	ErrInvalidRoleName   = errors.New("invalid role name")
)

func New(gocloak *gocloak.GoCloak, gocloakClientID string, gocloakClientSecret string, gocloakRealm string) *Auth {
	return &Auth{
		gocloak:             gocloak,
		gocloakClientID:     gocloakClientID,
		gocloakClientSecret: gocloakClientSecret,
		gocloakRealm:        gocloakRealm,
	}
}

func (a *Auth) Login(ctx context.Context, email string, password string) (string, error) {
	tracerCtx, span := otlp.Tracer.Start(ctx, "Keycloak login")
	defer span.End()

	token, err := a.gocloak.Login(tracerCtx, a.gocloakClientID, a.gocloakClientSecret, a.gocloakRealm, email, password)
	if err != nil {
		log.Printf("Login failed: %s\n", err.Error())
		return "", err
	}

	return token.AccessToken, nil
}

func (a *Auth) Register(ctx context.Context, email string, password string, role string) (string, error) {
	tracerCtx, span := otlp.Tracer.Start(ctx, "Keycloak register")
	defer span.End()

	token, err := a.gocloak.LoginClient(tracerCtx, a.gocloakClientID, a.gocloakClientSecret, a.gocloakRealm)
	if err != nil {
		log.Fatal(err)
	}

	var roleKeycloak *gocloak.Role
	if role != "" {
		var roleNameLowerCase = strings.ToLower(role)
		roleKeycloak, err = a.gocloak.GetRealmRole(tracerCtx, token.AccessToken, a.gocloakRealm, roleNameLowerCase)
		if err != nil {
			slog.Error("Unable to get role by name: %s\n", slog.String("error", err.Error()))
			return "", ErrInvalidRoleName
		}
	}

	user := gocloak.User{
		Email:    gocloak.StringP(email),
		Username: gocloak.StringP(email),
		Enabled:  gocloak.BoolP(true),
		Credentials: &[]gocloak.CredentialRepresentation{{
			Value: gocloak.StringP(password),
		}},
	}

	userId, err := a.gocloak.CreateUser(tracerCtx, token.AccessToken, a.gocloakRealm, user)
	if err != nil {
		log.Printf("Failed to create user: %s\n", err.Error())
		return "", ErrUserAlreadyExists
	}

	if role != "" {
		err = a.gocloak.AddRealmRoleToUser(tracerCtx, token.AccessToken, a.gocloakRealm, userId, []gocloak.Role{*roleKeycloak})
		if err != nil {
			slog.Error("unable to add a realm role to user:", slog.String("error", err.Error()))
			return "", err
		}
	}

	return userId, nil
}

func (a *Auth) CheckToken(ctx context.Context, token string) bool {
	rptResult, err := a.gocloak.RetrospectToken(ctx, token, a.gocloakClientID, a.gocloakClientSecret, a.gocloakRealm)
	if err != nil {
		slog.Error("Token check error", slog.String("error", err.Error()))
		return false
	}

	return *rptResult.Active
}

func (a *Auth) ParseClaims(ctx context.Context, token string) (TokenClaims, error) {
	customClaims := &AccessTokenCustomClaims{}
	_, err := a.gocloak.DecodeAccessTokenCustomClaims(
		ctx,
		token,
		a.gocloakRealm,
		customClaims,
	)
	if err != nil {
		return TokenClaims{}, err
	}

	return TokenClaims{
		UserID:     customClaims.Subject,
		Username:   customClaims.PreferredUsername,
		RealmRoles: customClaims.RealmAccess.Roles,
	}, nil
}
