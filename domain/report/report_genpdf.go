package report

import (
	"cpm-rad-backend/domain/config"
	"fmt"
)

type ReportPdfdata struct {
	Inspection   string            `json:"inspection"`
	Project      string            `json:"project"`
	Station      string            `json:"station"`
	Company      string            `json:"company"`
	Arrival      string            `json:"arrival"`
	Invoice      string            `json:"invoice"`
	Manufacturer string            `json:"manufacturer"`
	Serial       string            `json:"serial"`
	Receive      string            `json:"receive"`
	Doc          FileAttachments   `json:"doc"`
	BoqItem      map[string]string `json:"boqItem"`
}

type ReportPdfDB struct {
	Inspection   string `gorm:"column:INSPECTION_DATE"`
	Project      string `gorm:"column:PROJECT_NAME"`
	Station      string `gorm:"column:WORK_NAME"`
	Company      string `gorm:"column:CONTRACTOR_NAME"`
	Arrival      string `gorm:"column:ARRIVAL_DATE_AT_SITE"`
	Invoice      string `gorm:"column:CONTRACTOR_INV_NO"`
	Manufacturer string `gorm:"column:MANUFACTURER"`
	Serial       string `gorm:"column:SERIAL_NO"`
	Receive      string `gorm:"column:ROUND_NO"`
	ItemQuantity string `gorm:"column:QUANTITY"`
	ItemName     string `gorm:"column:NAME"`
	ItemUnit     string `gorm:"column:UNIT"`
}

type FileAttachmentDB struct {
	ID   int    `gorm:"column:ID"`
	Name string `gorm:"column:FILE_NAME"`
	Type string `gorm:"column:FILE_TYPE"`
}

type FileAttachmentDBs []FileAttachmentDB

type FileAttachment struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type FileAttachments []FileAttachment

func (f *FileAttachmentDB) ToResponse() FileAttachment {
	return FileAttachment{
		Type: f.Type,
		Name: f.Name,
		Link: config.WebLoadFileURL + fmt.Sprint(f.ID),
	}
}

func (f *FileAttachmentDBs) ToResponse() []FileAttachment {
	res := make([]FileAttachment, len(*f))
	for i, item := range *f {
		res[i] = item.ToResponse()
	}
	return res
}
