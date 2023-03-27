package auth

import (
	"context"
	"cpm-rad-backend/domain/connection"
)

func GetIDToken(db *connection.DBConnection) getIDTokenFunc {
	return func(ctx context.Context, ID string) (string, error) {
		idToken := ""
		if err := db.CPM.Model(&AuthLog{}).Where("ID = ?", ID).
			Pluck("ID_TOKEN", &idToken).
			Error; err != nil {
			return "", err
		}
		return idToken, nil
	}
}

func CreateLog(db *connection.DBConnection) createAuthLogFunc {
	return func(ctx context.Context, authLog AuthLog) error {
		return db.CPM.Create(&authLog).Error
	}
}
