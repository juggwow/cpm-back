package auth

import (
	"context"
	"cpm-rad-backend/domain/config"
	"cpm-rad-backend/domain/employee"
	"cpm-rad-backend/domain/logger"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type getEmployeeFunc func(context.Context, employee.Employee) (employee.Employee, error)

func (fn getEmployeeFunc) GetAuthorizedEmployee(ctx context.Context, empFromToken employee.Employee) (employee.Employee, error) {
	return fn(ctx, empFromToken)
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
		ctx:          ctx,
	}, nil
}

func (a *Authenticator) getCallbackToken(c *echo.Context, svc getEmployeeFunc) (string, *JwtEmployeeClaims, error) {
	token, err := a.clientConfig.Exchange(a.ctx, (*c).QueryParam("code"))
	if err != nil {
		return "", nil, err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return "", nil, errors.New("No id_token field in oauth2 token.")
	}

	idToken, err := a.provider.Verifier(&oidc.Config{ClientID: a.clientConfig.ClientID}).Verify(a.ctx, rawIDToken)
	if err != nil {
		return "", nil, err
	}

	rawClaims := keyClockClaims{}
	idToken.Claims(&rawClaims)

	emp, _ := svc.GetAuthorizedEmployee(nil, rawClaims.toEmployee())
	claims := rawClaims.toEmployeeClaims(*emp.ToResponse(), rawIDToken)

	signed, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&claims,
	).SignedString([]byte(config.AuthJWTSecret))

	return signed, &claims, err
}

func (a *Authenticator) getAuthURL() string {
	return a.clientConfig.AuthCodeURL(config.AuthState)
}

func (a *Authenticator) getCallbackURL(token string, err error) string {
	callbackURL := config.AuthCallbackURL
	if token == "" {
		return callbackURL
	}

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

func (a *Authenticator) AuthenHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.Unwrap(c)

		tokenString, err := c.Cookie(config.AuthJWTKey)
		if err != nil || tokenString == nil {
			log.Error(err.Error())
			c.Redirect(http.StatusTemporaryRedirect, a.getAuthURL())
			return nil
		}

		token, err := jwt.Parse(tokenString.Value, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || !t.Valid {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(config.AuthJWTSecret), nil
		})

		if err != nil {
			log.Error(err.Error())
			c.Redirect(http.StatusTemporaryRedirect, a.getAuthURL())
			return nil
		}

		if token == nil || !token.Valid {
			c.Redirect(http.StatusTemporaryRedirect, a.getAuthURL())
			return nil
		}

		_, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Error(err.Error())
			c.Redirect(http.StatusTemporaryRedirect, a.getAuthURL())
			return nil
		}

		c.Redirect(http.StatusTemporaryRedirect, a.getCallbackURL("", nil))
		return nil
	}
}

func (a *Authenticator) AuthenCallbackHandler(svc getEmployeeFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.Unwrap(c)

		token, _, err := a.getCallbackToken(&c, svc)
		if err != nil {
			log.Error(err.Error())
			c.Redirect(http.StatusTemporaryRedirect, a.getCallbackURL("", err))
			return err
		}

		// newCookie := new(http.Cookie)
		// newCookie.Name = config.AuthJWTKey
		// newCookie.Value = token
		// newCookie.Secure = true
		// newCookie.HttpOnly = true
		// newCookie.MaxAge = 0
		// newCookie.Domain = config.AppURL
		// c.SetCookie(newCookie)

		return c.Redirect(http.StatusTemporaryRedirect, a.getCallbackURL(token, nil))
	}
}

func GetAuthorizedClaims(c echo.Context) (JwtEmployeeClaims, error) {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return JwtEmployeeClaims{}, errors.New("Invalid token")
	}

	claims, ok := user.Claims.(*JwtEmployeeClaims)
	if !ok {
		return JwtEmployeeClaims{}, errors.New("Invalid claimed token")
	}
	return *claims, nil
}

// GetCurrent godoc
// @Summary get Current Employee
// @Tags Employees
// @Accept json
// @Produce json
// @Success 200 {object} employeeClaims
// @Router /api/v1/employees/me [get]
// @Security ApiKeyAuth
func GetCurrentHandler(c echo.Context) error {
	logger.Unwrap(c)

	claims, err := GetAuthorizedClaims(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, claims)
}
