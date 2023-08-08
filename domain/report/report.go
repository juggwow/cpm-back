package report

import (
	"cpm-rad-backend/domain/utils"
	"mime/multipart"
	"time"

	"github.com/shopspring/decimal"
)

type RequestReportCreate struct {
	ID         uint
	ItemID     string `json:"itemID"`
	Arrival    string `json:"arrival"`
	Inspection string `json:"inspection"`
	TaskMaster string `json:"taskMaster"`
	Invoice    string `json:"invoice"`
	Quantity   string `json:"quantity"`
	Country    string `json:"country"`
	Brand      string `json:"brand"`
	Model      string `json:"model"`
	Serial     string `json:"serial"`
	PeaNo      string `json:"peano"`
	CreateBy   string `json:"createby"`
	Status     string `json:"status"`
}

type FileDetail struct {
	ID          uint
	Name        string
	ObjectName  string
	DisplaySize string
	Size        string
	Unit        string
	Path        string
	FileType    string
	DocType     uint
}

// type AttachFiles []AttachFile

type ReportDB struct {
	ID         uint       `gorm:"column:ID"`
	ItemID     uint       `gorm:"column:BOQ_ID"`
	RadNo      string     `gorm:"column:RAD_NO"`
	Arrival    time.Time  `gorm:"column:ARRIVAL_DATE_AT_SITE"`
	Inspection time.Time  `gorm:"column:INSPECTION_DATE"`
	TaskMaster string     `gorm:"column:NAME_OF_TASKMASTER"`
	Invoice    string     `gorm:"column:CONTRACTOR_INV_NO"`
	Quantity   uint       `gorm:"column:QUANTITY"`
	Country    string     `gorm:"column:COUNTRY"`
	Brand      string     `gorm:"column:MANUFACTURER"`
	Model      string     `gorm:"column:MODEL"`
	Serial     string     `gorm:"column:SERIAL_NO"`
	PeaNo      string     `gorm:"column:PEA_NO"`
	CreateBy   string     `gorm:"column:CREATED_BY"`
	UpdateBy   string     `gorm:"column:UPDATED_BY"`
	UpdateDate *time.Time `gorm:"column:UPDATED_DATE"`
	Status     int        `gorm:"column:STATE_ID"`
	DelFlag    string     `gorm:"column:DEL_FLAG"`
}

func (ReportDB) TableName() string {
	return "CPM.CPM_WORK_CONTRACT_RAD"
}

func (r *RequestReportCreate) ToModel() ReportDB {
	return ReportDB{
		ID:         r.ID,
		ItemID:     utils.StringToUint(r.ItemID),
		Arrival:    utils.StringToTime(r.Arrival),
		Inspection: utils.StringToTime(r.Inspection),
		TaskMaster: r.TaskMaster,
		Invoice:    r.Invoice,
		Quantity:   utils.StringToUint(r.Quantity),
		Country:    r.Country,
		Brand:      r.Brand,
		Model:      r.Model,
		Serial:     r.Serial,
		PeaNo:      r.PeaNo,
		Status:     utils.StringToInt(r.Status),
	}
}

type File struct {
	Info   []*multipart.FileHeader
	Type   []string
	Update []UpdateFile
	Delete []string
}
type UpdateFile struct {
	FileID  uint `json:"id"`
	DocType uint `json:"docType"`
}

type AttachFileDB struct {
	ID         uint            `gorm:"column:ID"`
	ReportID   uint            `gorm:"column:RAD_ID"`
	DocType    uint            `gorm:"column:RAD_DOC_TYPE_ID"`
	Name       string          `gorm:"column:FILE_NAME"`
	Size       decimal.Decimal `gorm:"column:FILE_SIZE"`
	Unit       string          `gorm:"column:FILE_UNIT"`
	Path       string          `gorm:"column:FILE_PATH"`
	CreateBy   string          `gorm:"column:CREATED_BY"`
	UpdateBy   string          `gorm:"column:UPDATED_BY"`
	UpdateDate *time.Time      `gorm:"column:UPDATED_DATE"`
	DelFlag    string          `gorm:"column:DEL_FLAG"`
}

type AttachFilesDB []AttachFileDB

func (AttachFileDB) TableName() string {
	return "CPM.CPM_WORK_CONTRACT_RAD_FILE"
}

func (f *FileDetail) ToModel(r ReportDB) AttachFileDB {
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

// type AttachFileResponse struct {
// 	ID          uint   `json:"id"`
// 	Name        string `json:"name"`
// 	DisplaySize string `json:"displaySize"`
// 	FileType    string `json:"fileType"`
// 	DocType     uint   `json:"docType"`
// }

// type AttachFilesResponse []AttachFileResponse

// func (af *AttachFile) ToResponse() AttachFileResponse {
// 	return AttachFileResponse{
// 		ID:          af.ID,
// 		Name:        af.Name,
// 		DisplaySize: af.DisplaySize,
// 		FileType:    af.FileType,
// 		DocType:     af.DocType,
// 	}
// }

type ResponseAttachFile struct {
	ID       uint   `json:"id" gorm:"column:ID"`
	Name     string `json:"name" gorm:"column:FILE_NAME"`
	Size     string `json:"size" gorm:"column:FILE_SIZE"`
	Unit     string `json:"unit" gorm:"column:FILE_UNIT"`
	TypeID   uint   `json:"typeID" gorm:"column:TYPE_ID"`
	TypeName string `json:"typeName" gorm:"column:TYPE_NAME"`
}
type ResponseAttachFiles []ResponseAttachFile

type ReportDetailDB struct {
	ID         uint      `gorm:"column:ID"`
	ItemID     uint      `gorm:"column:ITEM_ID"`
	ItemName   string    `gorm:"column:ITEM_NAME"`
	ItemUnit   string    `gorm:"column:ITEM_UNIT"`
	RadNo      string    `gorm:"column:RAD_NO"`
	Arrival    time.Time `gorm:"column:ARRIVAL_DATE_AT_SITE"`
	Inspection time.Time `gorm:"column:INSPECTION_DATE"`
	TaskMaster string    `gorm:"column:NAME_OF_TASKMASTER"`
	Invoice    string    `gorm:"column:CONTRACTOR_INV_NO"`
	Quantity   uint      `gorm:"column:QUANTITY"`
	Country    string    `gorm:"column:COUNTRY"`
	Brand      string    `gorm:"column:MANUFACTURER"`
	Model      string    `gorm:"column:MODEL"`
	Serial     string    `gorm:"column:SERIAL_NO"`
	PeaNo      string    `gorm:"column:PEA_NO"`
	StateName  string    `gorm:"column:STATE_NAME"`
}
type MultiReportDetailDB []ReportDetailDB

type ResponseReportDetail struct {
	ID          uint                `json:"id"`
	ItemID      uint                `json:"itemID"`
	ItemName    string              `json:"itemName"`
	ItemUnit    string              `json:"itemUnit"`
	RadNo       string              `json:"RadNo"`
	Arrival     utils.Time          `json:"arrival"`
	Inspection  utils.Time          `json:"inspection"`
	TaskMaster  string              `json:"taskMaster"`
	Invoice     string              `json:"invoice"`
	Quantity    uint                `json:"quantity"`
	Country     string              `json:"country"`
	Brand       string              `json:"brand"`
	Model       string              `json:"model"`
	Serial      string              `json:"serial"`
	PeaNo       string              `json:"peano"`
	StateName   string              `json:"stateName"`
	AttachFiles ResponseAttachFiles `json:"attachFiles"`
}

func (r *ReportDetailDB) ToResponse(attachFiles ResponseAttachFiles) ResponseReportDetail {
	res := ResponseReportDetail{
		ID:          r.ID,
		ItemID:      r.ItemID,
		ItemName:    r.ItemName,
		ItemUnit:    r.ItemUnit,
		RadNo:       r.RadNo,
		Arrival:     utils.Time(r.Arrival),
		Inspection:  utils.Time(r.Inspection),
		TaskMaster:  r.TaskMaster,
		Invoice:     r.Invoice,
		Quantity:    r.Quantity,
		Country:     r.Country,
		Brand:       r.Brand,
		Model:       r.Model,
		Serial:      r.Serial,
		PeaNo:       r.PeaNo,
		StateName:   r.StateName,
		AttachFiles: attachFiles,
	}
	return res
}

func (r *ReportDB) ToResponse(attachFiles ResponseAttachFiles) ResponseReportDetail {
	res := ResponseReportDetail{
		ID:          r.ID,
		ItemID:      r.ItemID,
		RadNo:       r.RadNo,
		Arrival:     utils.Time(r.Arrival),
		Inspection:  utils.Time(r.Inspection),
		TaskMaster:  r.TaskMaster,
		Invoice:     r.Invoice,
		Quantity:    r.Quantity,
		Country:     r.Country,
		Brand:       r.Brand,
		Model:       r.Model,
		Serial:      r.Serial,
		PeaNo:       r.PeaNo,
		AttachFiles: attachFiles,
	}
	return res
}
