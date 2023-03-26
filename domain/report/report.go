package report

import (
	"cpm-rad-backend/domain/utils"
	"mime/multipart"
	"time"

	"github.com/shopspring/decimal"
)

type Report struct {
	ID           uint        `json:"ID"`
	ItemID       string      `json:"itemID"`
	Arrival      string      `json:"arrival"`
	Inspection   string      `json:"inspection"`
	TaskMaster   string      `json:"taskMaster"`
	Invoice      string      `json:"invoice"`
	Quantity     string      `json:"quantity"`
	Country      string      `json:"country"`
	Manufacturer string      `json:"manufacturer"`
	Model        string      `json:"model"`
	Serial       string      `json:"serial"`
	PeaNo        string      `json:"peano"`
	CreateBy     string      `json:"createby"`
	Status       string      `json:"status"`
	AttachFiles  AttachFiles `json:"filesAttach"`
}

type AttachFile struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	ObjectName  string `json:"objectName"`
	DisplaySize string `json:"displaySize"`
	Size        string `json:"size"`
	Unit        string `json:"unit"`
	Path        string `json:"filePath"`
	FileType    string `json:"fileType"`
	DocType     uint   `json:"docType"`
}

type AttachFiles []AttachFile

type ReportDB struct {
	ID           uint       `gorm:"column:ID"`
	ItemID       uint       `gorm:"column:BOQ_ID"`
	RadNo        string     `gorm:"column:RAD_NO"`
	Arrival      time.Time  `gorm:"column:ARRIVAL_DATE_AT_SITE"`
	Inspection   time.Time  `gorm:"column:INSPECTION_DATE"`
	TaskMaster   string     `gorm:"column:NAME_OF_TASKMASTER"`
	Invoice      string     `gorm:"column:CONTRACTOR_INV_NO"`
	Quantity     uint       `gorm:"column:QUANTITY"`
	Country      string     `gorm:"column:COUNTRY"`
	Manufacturer string     `gorm:"column:MANUFACTURER"`
	Model        string     `gorm:"column:MODEL"`
	Serial       string     `gorm:"column:SERIAL_NO"`
	PeaNo        string     `gorm:"column:PEA_NO"`
	CreateBy     string     `gorm:"column:CREATED_BY"`
	UpdateBy     string     `gorm:"column:UPDATED_BY"`
	UpdateDate   *time.Time `gorm:"column:UPDATED_DATE"`
	Status       int        `gorm:"column:STATE_ID"`
	DelFlag      string     `gorm:"column:DEL_FLAG"`
}

func (ReportDB) TableName() string {
	return "CPM.CPM_WORK_CONTRACT_RAD"
}

func (r *Report) ToModel() ReportDB {
	return ReportDB{
		ItemID:       utils.StringToUint(r.ItemID),
		Arrival:      utils.StringToTime(r.Arrival),
		Inspection:   utils.StringToTime(r.Inspection),
		TaskMaster:   r.TaskMaster,
		Invoice:      r.Invoice,
		Quantity:     utils.StringToUint(r.Quantity),
		Country:      r.Country,
		Manufacturer: r.Manufacturer,
		Model:        r.Model,
		Serial:       r.Serial,
		PeaNo:        r.PeaNo,
		Status:       utils.StringToInt(r.Status),
	}
}

type File struct {
	Info []*multipart.FileHeader
	Type []string
}

type AttachFileDB struct {
	ID       uint            `gorm:"column:ID"`
	ReportID uint            `gorm:"column:RAD_ID"`
	DocType  uint            `gorm:"column:RAD_DOC_TYPE_ID"`
	Name     string          `gorm:"column:FILE_NAME"`
	Size     decimal.Decimal `gorm:"column:FILE_SIZE"`
	Unit     string          `gorm:"column:FILE_UNIT"`
	Path     string          `gorm:"column:FILE_PATH"`
	CreateBy string          `gorm:"column:CREATED_BY"`
}

type AttachFilesDB []AttachFileDB

func (AttachFileDB) TableName() string {
	return "CPM.CPM_WORK_CONTRACT_RAD_FILE"
}

func (f *AttachFile) ToModel(r ReportDB) AttachFileDB {
	size, _ := decimal.NewFromString(f.Size)
	file := AttachFileDB{
		ReportID: r.ID,
		DocType:  f.DocType,
		Name:     f.Name,
		Size:     size,
		Unit:     f.Unit,
		Path:     f.Path,
		CreateBy: r.CreateBy,
	}

	return file
}
