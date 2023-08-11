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
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type getEmployeeFunc func(context.Context, employee.Employee) (employee.Employee, error)
type getIDTokenFunc func(context.Context, string) (string, error)
type createAuthLogFunc func(context.Context, AuthLog) error

func (fn getEmployeeFunc) GetAuthorizedEmployee(ctx context.Context, empFromToken employee.Employee) (employee.Employee, error) {
	return fn(ctx, empFromToken)
}

func (fn getIDTokenFunc) GetIDToken(ctx context.Context, ID string) (string, error) {
	return fn(ctx, ID)
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

func (a *Authenticator) getCallbackToken(c *echo.Context, svc getEmployeeFunc) (*JwtEmployeeClaims, string, error) {
	callbackResult, err := a.clientConfig.Exchange(a.ctx, (*c).QueryParam("code"))
	if err != nil {
		return nil, "", err
	}

	rawIDToken, ok := callbackResult.Extra("id_token").(string)
	if !ok {
		return nil, "", errors.New("No id_token field in oauth2 token.")
	}

	idToken, err := a.provider.Verifier(&oidc.Config{ClientID: a.clientConfig.ClientID}).Verify(a.ctx, rawIDToken)
	if err != nil {
		return nil, "", err
	}

	rawClaims := keyClockClaims{}
	idToken.Claims(&rawClaims)

	log := logger.Unwrap(*c)
	log.Info(fmt.Sprintln(rawClaims))
	// Check is PEA Employee
	// if rawClaims.Sub[:2] == "f:" {
	// 	emp, _ := svc.GetAuthorizedEmployee(context.TODO(), rawClaims.toEmployee())
	// 	claims := rawClaims.toEmployeeClaims(*emp.ToResponse(), rawIDToken)
	// 	return &claims, rawIDToken, err
	// }

	//emp, _ := svc.GetAuthorizedEmployee(nil, rawClaims.toEmployee())
	emp, _ := svc.GetAuthorizedEmployee(context.TODO(), rawClaims.toEmployee())
	claims := rawClaims.toEmployeeClaims(*emp.ToResponse(), rawIDToken)

	return &claims, rawIDToken, err
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
		logger.Unwrap(c)

		tokenString, err := c.Cookie(config.AuthJWTKey)
		if err != nil || tokenString == nil {
			c.Redirect(http.StatusTemporaryRedirect, a.getAuthURL())
			return nil
		}

		token, err := jwt.Parse(tokenString.Value, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || !t.Valid {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(config.AuthJWTSecret), nil
		})

		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, a.getAuthURL())
			return nil
		}

		if token == nil || !token.Valid {
			c.Redirect(http.StatusTemporaryRedirect, a.getAuthURL())
			return nil
		}

		_, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Redirect(http.StatusTemporaryRedirect, a.getAuthURL())
			return nil
		}

		c.Redirect(http.StatusTemporaryRedirect, a.getCallbackURL("", err))
		return nil
	}
}

func (a *Authenticator) AuthenCallbackHandler(
	svc getEmployeeFunc,
	svc2 createAuthLogFunc,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.Unwrap(c)

		claims, idToken, err := a.getCallbackToken(&c, svc)
		if err != nil {
			log.Error("Error authen callback", zap.Error(err))
			c.Redirect(http.StatusTemporaryRedirect, a.getCallbackURL("", err))
			return err
		}

		token, err := claims.getToken(config.AuthJWTExpiredDuration)
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, a.getCallbackURL("", err))
			return err
		}

		authLog := AuthLog{
			ID:         claims.ID,
			IP:         c.RealIP(),
			IDToken:    idToken,
			EmployeeID: claims.EmployeeID,
			LoginAt:    time.Now(),
		}
		svc2.CreateAuthLog(c.Request().Context(), authLog)

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

// GetCurrentHandler godoc
// @Summary get Current Employee
// @Tags Employees
// @Accept json
// @Produce json
// @Success 200 {object} employee.EmployeeResponse
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

// GetRefreshTokenHandler godoc
// @Summary Refresh Token
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} refreshTokenResponse
// @Router /auth/refreshToken [get]
// @Security ApiKeyAuth
func (a Authenticator) GetRefreshTokenHandler(svc getIDTokenFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger.Unwrap(c)

		claims, err := GetAuthorizedClaims(c)
		if err != nil {
			return err
		}

		idToken, err := svc.GetIDToken(c.Request().Context(), claims.ID)
		if err != nil {
			return err
		}

		if _, err := a.provider.Verifier(&oidc.Config{ClientID: a.clientConfig.ClientID}).Verify(a.ctx, idToken); err != nil {
			return err
		}

		token, err := claims.getToken(config.AuthJWTExpiredDuration)
		if err != nil {
			return err
		}

		refreshToken, err := claims.getToken(config.AuthJWTExpiredRefreshDuration)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, refreshTokenResponse{
			Token:        token,
			RefreshToken: refreshToken,
		})
	}
}

func (a Authenticator) LogoutHandler(svc getIDTokenFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger.Unwrap(c)

		// parser := &jwt.Parser{
		// 	SkipClaimsValidation: true,
		// }
		parser := jwt.NewParser(jwt.WithoutClaimsValidation())

		claims := JwtEmployeeClaims{}
		parser.ParseWithClaims(c.Param("token"), claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AuthJWTSecret), nil
		})

		path, err := url.JoinPath(config.AuthURL, "/protocol/openid-connect/logout")
		if err != nil {
			return err
		}

		urlInfo, err := url.Parse(path)
		if err != nil {
			return err
		}

		idToken, err := svc.GetIDToken(c.Request().Context(), claims.ID)
		if err != nil {
			return err
		}

		q := urlInfo.Query()
		q.Set("post_logout_redirect_uri", config.AppURL+"/auth")
		q.Set("id_token_hint", idToken)
		urlInfo.RawQuery = q.Encode()

		return c.Redirect(http.StatusTemporaryRedirect, urlInfo.String())
	}
}
