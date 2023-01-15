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

type Role string

const (
	ADMIN Role = "ADMIN"
	STAFF Role = "STAFF"
)

type Authenticator struct {
	provider     *oidc.Provider
	clientConfig *oauth2.Config
	verifier     *oidc.IDTokenVerifier
	ctx          context.Context
}

type JwtEmployeeClaims struct {
	employeeClaims
	jwt.StandardClaims
}

type employeeClaims struct {
	employee.EmployeeResponse
	Role  Role   `json:"role"`
	Token string `json:"token"`
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

func (a keyClockClaims) toEmployeeClaims(emp employee.EmployeeResponse, token string) JwtEmployeeClaims {
	exp := time.Now().Add(time.Hour * time.Duration(8)).Unix()
	return JwtEmployeeClaims{
		employeeClaims: employeeClaims{
			Role:             STAFF,
			EmployeeResponse: emp,
			Token:            token,
		},
		StandardClaims: jwt.StandardClaims{
			Audience:  a.Aud,
			ExpiresAt: exp,
			Id:        xid.New().String(),
			IssuedAt:  a.Iat,
			Subject:   emp.EmployeeID,
			Issuer:    config.AppURL,
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
