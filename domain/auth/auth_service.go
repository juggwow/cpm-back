package auth

import (
	"context"
	"cpm-rad-backend/domain/auth/employee"
	"cpm-rad-backend/domain/config"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type getEmployeeFunc func(context.Context, employee.Employee) (employee.Employee, error)
type createAuthLogFunc func(context.Context, AuthLog) error

func (fn getEmployeeFunc) GetAuthorizedEmployee(ctx context.Context, empFromToken employee.Employee) (employee.Employee, error) {
	return fn(ctx, empFromToken)
}

func (fn createAuthLogFunc) CreateAuthLog(ctx context.Context, authLog AuthLog) error {
	return fn(ctx, authLog)
}

func NewAuthenticator() (*Authenticator, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, config.AuthURL)
	if err != nil {
		return nil, err
	}

	redirectURL, err := url.JoinPath(config.AppURL, "/auth/callback")
	if err != nil {
		return nil, err
	}

	cacheConfig := bigcache.DefaultConfig(24 * time.Hour)
	cacheConfig.Verbose = true
	idTokenMap, err := bigcache.New(context.Background(), cacheConfig)
	if err != nil {
		return nil, err
	}

	config := &oauth2.Config{
		ClientID:     config.AuthClientID,
		ClientSecret: config.AuthClientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &Authenticator{
		provider:     provider,
		clientConfig: config,
		idTokenMap:   idTokenMap,
		ctx:          ctx,
	}, nil
}

func (a *Authenticator) getAuthURL() string {
	return a.clientConfig.AuthCodeURL(config.AuthState)
}

func getSecretKey(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || !t.Valid {
		return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
	}
	return []byte(config.AuthJWTSecret), nil
}

// func getClaims(token string) *JwtEmployeeClaims {
// 	parser := &jwt.Parser{
// 		SkipClaimsValidation: true,
// 	}

// 	claims := &JwtEmployeeClaims{}
// 	parser.ParseWithClaims(token, claims, getSecretKey)
// 	return claims
// }

func (a *Authenticator) getCallbackURL(token string, err error) string {
	callbackURL := config.AuthCallbackURL
	if token == "" {
		return callbackURL
	}

	// claims := getClaims(token)
	// deptChangeCode := claims.DeptChangeCode
	// if deptChangeCode != "" {
	// 	callbackURL = callbackURL + "/admin"
	// }

	urlObject, urlErr := url.Parse(callbackURL)
	if urlErr != nil {
		zap.L().Error("Invalid callback url: " + urlErr.Error())
		return ""
	}

	q := urlObject.Query()
	if token != "" {
		q.Set("token", token)
	}
	if err != nil {
		q.Set("error", err.Error())
	}
	urlObject.RawQuery = q.Encode()
	return urlObject.String()
}

func (a *Authenticator) getCallbackToken(c *echo.Context, svc getEmployeeFunc) (*JwtEmployeeClaims, string, error) {
	callbackResult, err := a.clientConfig.Exchange(a.ctx, (*c).QueryParam("code"))
	if err != nil {
		return nil, "", err
	}

	rawIDToken, ok := callbackResult.Extra("id_token").(string)
	if !ok {
		return nil, "", errors.New("no id_token field in oauth2 token")
	}

	idToken, err := a.provider.Verifier(&oidc.Config{ClientID: a.clientConfig.ClientID}).Verify(a.ctx, rawIDToken)
	if err != nil {
		return nil, "", err
	}

	rawClaims := keyClockClaims{}
	idToken.Claims(&rawClaims)

	// Check is PEA Employee
	// if rawClaims.Sub[:2] == "f:" {
	// 	emp, _ := svc.GetAuthorizedEmployee(context.TODO(), rawClaims.toEmployee())
	// 	claims := rawClaims.toEmployeeClaims(*emp.ToResponse(), rawIDToken)
	// 	return &claims, rawIDToken, err
	// }

	emp, _ := svc.GetAuthorizedEmployee(context.TODO(), rawClaims.toEmployee())
	claims := rawClaims.toEmployeeClaims(*emp.ToResponse(), rawIDToken)
	return &claims, rawIDToken, err
}

func (claims JwtEmployeeClaims) getToken(expiredDuration time.Duration) (string, error) {
	claims.ID = xid.New().String()
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(expiredDuration))
	return jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&claims,
	).SignedString([]byte(config.AuthJWTSecret))
}

// func (claims *JwtEmployeeClaims) getToken(expiredDuration time.Duration) (string, error) {
// 	claims.Id = xid.New().String()
// 	claims.IssuedAt = time.Now().Unix()
// 	claims.ExpiresAt = time.Now().Add(expiredDuration).Unix()
// 	return jwt.NewWithClaims(
// 		jwt.SigningMethodHS256,
// 		claims,
// 	).SignedString([]byte(config.AuthJWTSecret))
// }

func GetAuthorizedClaims(c echo.Context) (JwtEmployeeClaims, error) {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return JwtEmployeeClaims{}, errors.New("invalid token")
	}

	claims, ok := user.Claims.(*JwtEmployeeClaims)
	if !ok {
		return JwtEmployeeClaims{}, errors.New("invalid claimed token")
	}
	return *claims, nil
}

func (a Authenticator) getLogoutURL(claims *JwtEmployeeClaims) (string, error) {
	path, err := url.JoinPath(config.AuthURL, "/protocol/openid-connect/logout")
	if err != nil {
		return "", err
	}

	urlInfo, err := url.Parse(path)
	if err != nil {
		return "", err
	}

	// #SSO2 comment This
	idToken, err := a.idTokenMap.Get(claims.EmployeeID)
	if err != nil {
		return "", err
	}
	// #

	q := urlInfo.Query()
	q.Set("post_logout_redirect_uri", config.AppURL+"/auth")
	// #SSO2 comment This
	q.Set("id_token_hint", string(idToken))
	// #
	urlInfo.RawQuery = q.Encode()
	return urlInfo.String(), nil
}
