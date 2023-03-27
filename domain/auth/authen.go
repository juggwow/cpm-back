package auth

import (
	"context"
	"cpm-rad-backend/domain/config"
	"cpm-rad-backend/domain/employee"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt"
	"github.com/rs/xid"
	"golang.org/x/oauth2"
)

type AuthLog struct {
	ID         string    `gorm:"column:ID;primarykey"`
	IP         string    `gorm:"column:IP"`
	EmployeeID string    `gorm:"column:EMPLOYEE_ID"`
	LoginAt    time.Time `gorm:"column:LOGIN_AT"`
	IDToken    string    `gorm:"column:ID_TOKEN"`
}

func (AuthLog) TableName() string {
	return "CMDC_AUTH_LOG"
}

type Authenticator struct {
	provider     *oidc.Provider
	clientConfig *oauth2.Config
	verifier     *oidc.IDTokenVerifier
	ctx          context.Context
}

type JwtEmployeeClaims struct {
	employee.EmployeeResponse
	jwt.StandardClaims
}

type keyClockClaims struct {
	Exp               int64  `json:"exp"`
	Iat               int64  `json:"iat"`
	AuthTime          int64  `json:"auth_time"`
	Jti               string `json:"jti"`
	Iss               string `json:"iss"`
	Aud               string `json:"aud"`
	Sub               string `json:"sub"`
	Typ               string `json:"typ"`
	Azp               string `json:"azp"`
	SessionState      string `json:"session_state"`
	AtHash            string `json:"at_hash"`
	Acr               string `json:"acr"`
	EmailVerified     bool   `json:"email_verified"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
	Locale            string `json:"locale"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	Email             string `json:"email"`
}

type refreshTokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

func (a keyClockClaims) toEmployeeClaims(emp employee.EmployeeResponse, token string) JwtEmployeeClaims {
	return JwtEmployeeClaims{
		EmployeeResponse: emp,
		StandardClaims: jwt.StandardClaims{
			Audience: a.Aud,
			Id:       xid.New().String(),
			IssuedAt: a.Iat,
			Subject:  emp.EmployeeID,
			Issuer:   config.AppURL,
		},
	}
}

func (a keyClockClaims) toEmployee() employee.Employee {
	return employee.Employee{
		EmployeeID: a.PreferredUsername,
		FirstName:  a.GivenName,
		LastName:   a.FamilyName,
	}
}

func (claims JwtEmployeeClaims) getToken(expiredDuration time.Duration) (string, error) {
	tokenClaims := claims
	tokenClaims.IssuedAt = time.Now().Unix()
	tokenClaims.ExpiresAt = time.Now().Add(expiredDuration).Unix()
	return jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&tokenClaims,
	).SignedString([]byte(config.AuthJWTSecret))
}
