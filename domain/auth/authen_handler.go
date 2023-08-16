package auth

import (
	"cpm-rad-backend/domain/config"
	"cpm-rad-backend/domain/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// type getEmployeeFunc func(context.Context, employee.Employee) (employee.Employee, error)
// type getIDTokenFunc func(context.Context, string) (string, error)
// type createAuthLogFunc func(context.Context, AuthLog) error

// func (fn getEmployeeFunc) GetAuthorizedEmployee(ctx context.Context, empFromToken employee.Employee) (employee.Employee, error) {
// 	return fn(ctx, empFromToken)
// }

// func (fn getIDTokenFunc) GetIDToken(ctx context.Context, ID string) (string, error) {
// 	return fn(ctx, ID)
// }

// func (fn createAuthLogFunc) CreateAuthLog(ctx context.Context, authLog AuthLog) error {
// 	return fn(ctx, authLog)
// }

// func NewAuthenticator() (*Authenticator, error) {
// 	ctx := context.Background()

// 	provider, err := oidc.NewProvider(ctx, config.AuthURL)
// 	if err != nil {
// 		return nil, err
// 	}

// 	redirectURL, err := url.JoinPath(config.AppURL, "/auth/callback")
// 	if err != nil {
// 		return nil, err
// 	}

// 	config := &oauth2.Config{
// 		ClientID:     config.AuthClientID,
// 		ClientSecret: config.AuthClientSecret,
// 		RedirectURL:  redirectURL,
// 		Endpoint:     provider.Endpoint(),
// 		Scopes:       []string{oidc.ScopeOpenID, "profile"},
// 	}

// 	return &Authenticator{
// 		provider:     provider,
// 		clientConfig: config,
// 		ctx:          ctx,
// 	}, nil
// }

// func (a *Authenticator) getCallbackToken(c *echo.Context, svc getEmployeeFunc) (*JwtEmployeeClaims, string, error) {
// 	callbackResult, err := a.clientConfig.Exchange(a.ctx, (*c).QueryParam("code"))
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	rawIDToken, ok := callbackResult.Extra("id_token").(string)
// 	if !ok {
// 		return nil, "", errors.New("No id_token field in oauth2 token.")
// 	}

// 	idToken, err := a.provider.Verifier(&oidc.Config{ClientID: a.clientConfig.ClientID}).Verify(a.ctx, rawIDToken)
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	rawClaims := keyClockClaims{}
// 	idToken.Claims(&rawClaims)

// 	// log := logger.Unwrap(*c)
// 	// log.Info(fmt.Sprintln(rawClaims))
// 	// fmt.Printf("%+v\n", rawClaims)

// 	var empRes employee.EmployeeResponse
// 	if rawClaims.Sub[:2] == "f:" {
// 		emp, _ := svc.GetAuthorizedEmployee(context.TODO(), rawClaims.toEmployee())
// 		empRes = *emp.ToResponse()
// 	} else {
// 		empRes = employee.EmployeeResponse{
// 			Id:             0,
// 			EmployeeID:     rawClaims.PreferredUsername,
// 			Title:          "",
// 			FirstName:      rawClaims.GivenName,
// 			LastName:       rawClaims.FamilyName,
// 			Position:       "",
// 			BA:             "",
// 			DeptChangeCode: "",
// 		}
// 	}

// 	claims := rawClaims.toEmployeeClaims(empRes, rawIDToken)
// 	return &claims, rawIDToken, err
// }

// func (a *Authenticator) getAuthURL() string {
// 	return a.clientConfig.AuthCodeURL(config.AuthState)
// }

// func (a *Authenticator) getCallbackURL(token string, err error) string {
// 	callbackURL := config.AuthCallbackURL
// 	if token == "" {
// 		return callbackURL
// 	}

// 	urlObject, urlErr := url.Parse(callbackURL)
// 	if urlErr != nil {
// 		zap.L().Error("Invalid callback url: " + urlErr.Error())
// 		return ""
// 	}

// 	q := urlObject.Query()
// 	if token != "" {
// 		q.Set("token", token)
// 	}
// 	if err != nil {
// 		q.Set("error", err.Error())
// 	}
// 	urlObject.RawQuery = q.Encode()
// 	return urlObject.String()
// }

func (a *Authenticator) AuthenHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		logger.Unwrap(c)
		webRedirectURL := c.QueryParam("page")
		a.clientConfig.RedirectURL = config.AppURL + "/auth/callback" + "?page=" + webRedirectURL

		tokenString, err := c.Cookie(config.AuthJWTKey)
		if err != nil || tokenString == nil {
			fmt.Println("round1")
			c.Redirect(http.StatusTemporaryRedirect, a.getAuthURL())
			return nil
		}

		token, err := jwt.Parse(tokenString.Value, getSecretKey)

		if err != nil {
			fmt.Println("round2")
			c.Redirect(http.StatusTemporaryRedirect, a.getAuthURL())
			return nil
		}

		if token == nil || !token.Valid {
			fmt.Println("round3")
			c.Redirect(http.StatusTemporaryRedirect, a.getAuthURL())
			return nil
		}

		_, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("round4")
			c.Redirect(http.StatusTemporaryRedirect, a.getAuthURL())
			return nil
		}

		c.Redirect(http.StatusTemporaryRedirect, a.getCallbackURL(webRedirectURL, "", err))
		return nil
	}
}

func (a *Authenticator) AuthenCallbackHandler(svc getEmployeeFunc, svc2 createAuthLogFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.Unwrap(c)
		webRedirectURL := c.QueryParam("page")
		claims, idToken, err := a.getCallbackToken(&c, svc)
		if err != nil {
			log.Error("Error authen callback", zap.Error(err))
			c.Redirect(http.StatusTemporaryRedirect, a.getCallbackURL(webRedirectURL, "", err))
			return err
		}

		token, err := claims.getToken(config.AuthJWTExpiredDuration)
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, a.getCallbackURL(webRedirectURL, "", err))
			return err
		}

		a.idTokenMap.Set(claims.EmployeeID, []byte(idToken))

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

		return c.Redirect(http.StatusTemporaryRedirect, a.getCallbackURL(webRedirectURL, token, nil))
	}
}

// func GetAuthorizedClaims(c echo.Context) (JwtEmployeeClaims, error) {
// 	user, ok := c.Get("user").(*jwt.Token)
// 	if !ok {
// 		return JwtEmployeeClaims{}, errors.New("Invalid token")
// 	}

// 	claims, ok := user.Claims.(*JwtEmployeeClaims)
// 	if !ok {
// 		return JwtEmployeeClaims{}, errors.New("Invalid claimed token")
// 	}
// 	return *claims, nil
// }

// GetRefreshTokenHandler godoc
// @Summary Refresh Token
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} refreshTokenResponse
// @Router /auth/refreshToken [get]
// @Security ApiKeyAuth
func (a Authenticator) GetRefreshTokenHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.Unwrap(c)

		claims, err := GetAuthorizedClaims(c)
		if err != nil {
			log.Warn("Refresh token: no claims", zap.Error(err))
			return c.NoContent(http.StatusUnauthorized)
		}

		idToken, err := a.idTokenMap.Get(claims.EmployeeID)
		if err != nil {
			log.Warn("Refresh token: no idToken", zap.Error(err))
			return c.NoContent(http.StatusUnauthorized)
		}

		if _, err := a.provider.Verifier(&oidc.Config{ClientID: a.clientConfig.ClientID}).Verify(a.ctx, string(idToken)); err != nil {
			log.Warn("Refresh token: unable to verify", zap.Error(err))
			return c.NoContent(http.StatusUnauthorized)
		}

		token, err := claims.getToken(config.AuthJWTExpiredDuration)
		if err != nil {
			log.Warn("Refresh token: generate token failed", zap.Error(err))
			return c.NoContent(http.StatusUnauthorized)
		}
		a.idTokenMap.Set(claims.EmployeeID, idToken)

		refreshToken, err := claims.getToken(config.AuthJWTExpiredRefreshDuration)
		if err != nil {
			log.Warn("Refresh token: generate refresh token failed", zap.Error(err))
			return c.NoContent(http.StatusUnauthorized)
		}
		a.idTokenMap.Set(claims.EmployeeID, idToken)

		return c.JSON(http.StatusOK, refreshTokenResponse{
			Token:        token,
			RefreshToken: refreshToken,
		})
	}
}

func (a Authenticator) LogoutHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		logger.Unwrap(c)

		// parser := &jwt.Parser{
		// 	SkipClaimsValidation: true,
		// }
		parser := jwt.NewParser(jwt.WithoutClaimsValidation())
		// claims := JwtEmployeeClaims{}

		claims := &JwtEmployeeClaims{}
		parser.ParseWithClaims(c.Param("token"), claims, getSecretKey)

		logoutURL, err := a.getLogoutURL(claims)
		if err != nil {
			return err
		}

		a.idTokenMap.Delete(claims.EmployeeID)

		return c.Redirect(http.StatusTemporaryRedirect, logoutURL)
	}
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
	fmt.Sprintln("GetCurrentHandler")
	claims, err := GetAuthorizedClaims(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, claims)
}
