package utils

import (
	"cpm-rad-backend/domain/config"
	"errors"
	"mime/multipart"
)

func IsPdf(files []*multipart.FileHeader) error {
	// const LIMIT int64 = 52428800
	for _, file := range files {
		if (file.Header["Content-Type"][0] != "application/pdf") && (config.FILE_SIZE_LIMIT < file.Size) {
			return errors.New(":Invalid PDF File Request")
		}
	}
	return nil
}
