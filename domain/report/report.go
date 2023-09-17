package report

import (
	"cpm-rad-backend/domain/config"
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

type RequestReportUpdate struct {
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
	UpdateBy   string `json:"updateBy"`
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

func (r *RequestReportUpdate) ToModel() ReportDB {
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
		UpdateBy:   r.UpdateBy,
		UpdateDate: &time.Time{},
		Status:     utils.StringToInt(r.Status),
	}
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
		CreateBy:   r.CreateBy,
		UpdateDate: &time.Time{},
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
	FileID uint `json:"fileId"`
	TypeID uint `json:"typeId"`
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
	Name     string `json:"name,omitempty" gorm:"column:FILE_NAME"`
	Size     string `json:"size,omitempty" gorm:"column:FILE_SIZE"`
	Unit     string `json:"unit,omitempty" gorm:"column:FILE_UNIT"`
	TypeID   uint   `json:"typeID,omitempty" gorm:"column:TYPE_ID"`
	TypeName string `json:"typeName,omitempty" gorm:"column:TYPE_NAME"`
}
type ResponseAttachFiles []ResponseAttachFile

type ReportDetailDB struct {
	ID                 uint      `gorm:"column:ID"`
	ItemID             uint      `gorm:"column:ITEM_ID"`
	ItemName           string    `gorm:"column:ITEM_NAME"`
	ItemUnit           string    `gorm:"column:ITEM_UNIT"`
	RadNo              string    `gorm:"column:RAD_NO"`
	Arrival            time.Time `gorm:"column:ARRIVAL_DATE_AT_SITE"`
	Inspection         time.Time `gorm:"column:INSPECTION_DATE"`
	TaskMaster         string    `gorm:"column:NAME_OF_TASKMASTER"`
	Invoice            string    `gorm:"column:CONTRACTOR_INV_NO"`
	Quantity           uint      `gorm:"column:QUANTITY"`
	Country            string    `gorm:"column:COUNTRY"`
	Brand              string    `gorm:"column:MANUFACTURER"`
	Model              string    `gorm:"column:MODEL"`
	Serial             string    `gorm:"column:SERIAL_NO"`
	PeaNo              string    `gorm:"column:PEA_NO"`
	StateID            string    `gorm:"column:STATE_ID"`
	StateName          string    `gorm:"column:STATE_NAME"`
	Remark             string    `gorm:"column:REMARK"`
	ReceiveQuantity    uint      `gorm:"column:RECEIVE_QUANTITY"`
	DefectQuantity     uint      `gorm:"column:DEFECT_QUANTITY"`
	InCompleteQuantity uint      `gorm:"column:INCOMPLETE_QUANTITY"`
	MismatchQuantity   uint      `gorm:"column:MISMATCH_QUANTITY"`
}
type MultiReportDetailDB []ReportDetailDB

type ResponseReportDetail struct {
	ID          uint                `json:"id"`
	ItemID      uint                `json:"itemID,omitempty"`
	ItemName    string              `json:"itemName,omitempty"`
	ItemUnit    string              `json:"itemUnit,omitempty"`
	RadNo       string              `json:"RadNo,omitempty"`
	Arrival     utils.Time          `json:"arrival,omitempty"`
	Inspection  utils.Time          `json:"inspection,omitempty"`
	TaskMaster  string              `json:"taskMaster,omitempty"`
	Invoice     string              `json:"invoice,omitempty"`
	Quantity    uint                `json:"quantity,omitempty"`
	Country     string              `json:"country,omitempty"`
	Brand       string              `json:"brand,omitempty"`
	Model       string              `json:"model,omitempty"`
	Serial      string              `json:"serial,omitempty"`
	PeaNo       string              `json:"peano,omitempty"`
	StateID     string              `json:"stateID,omitempty"`
	StateName   string              `json:"stateName,omitempty"`
	AttachFiles ResponseAttachFiles `json:"attachFiles,omitempty"`
	RadFiles    ResponseAttachFile  `json:"radFile,omitempty"`
	Comment     ResponseComment     `json:"comment,omitempty"`
}

func (r *ReportDetailDB) ToResponse(attachFiles ResponseAttachFiles, attachImages DbAttachImages) ResponseReportDetail {
	var radFile ResponseAttachFile
	var attach ResponseAttachFiles
	for _, file := range attachFiles {
		if file.TypeID == 0 {
			radFile = file
		} else {
			attach = append(attach, file)
		}
	}

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
		StateID:     r.StateID,
		StateName:   r.StateName,
		AttachFiles: attach,
		RadFiles:    radFile,
		Comment:     r.ToResponseComment(attachImages),
	}
	return res
}

func (r *ReportDetailDB) ToResponseComment(attachImages DbAttachImages) ResponseComment {
	res := ResponseComment{
		OverAll: utils.IfThenElse((r.DefectQuantity+r.InCompleteQuantity+r.MismatchQuantity) > 0, "พบปัญหา", "รับของเรียบร้อย").(string),
		Problem: Problem{
			Defect:     r.DefectQuantity,
			Incomplate: r.InCompleteQuantity,
			Mismatch:   r.MismatchQuantity,
		},
		NoProblem:    r.ReceiveQuantity,
		Images:       attachImages.ToResponseSrc(),
		ExtraComment: r.Remark,
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

type FileResponse struct {
	Obj  []byte
	Ext  string
	Name string
}

type DbAttachImage struct {
	ID         uint            `gorm:"column:ID"`
	ReportID   uint            `gorm:"column:RAD_ID"`
	Name       string          `gorm:"column:FILE_NAME"`
	Size       decimal.Decimal `gorm:"column:FILE_SIZE"`
	Unit       string          `gorm:"column:FILE_UNIT"`
	Path       string          `gorm:"column:FILE_PATH"`
	CreateBy   string          `gorm:"column:CREATED_BY"`
	UpdateBy   string          `gorm:"column:UPDATED_BY"`
	UpdateDate *time.Time      `gorm:"column:UPDATED_DATE"`
	DelFlag    string          `gorm:"column:DEL_FLAG"`
}

type DbAttachImages []DbAttachImage

func (DbAttachImage) TableName() string {
	return "CPM.RAD_ATTACH_PHOTO"
}

func (item *DbAttachImage) ToResponseSrc() Image {
	res := Image{
		Src: config.AppURL + "/image/" + item.Path,
	}
	return res
}

func (items *DbAttachImages) ToResponseSrc() []Image {
	res := make([]Image, len(*items))
	for i, item := range *items {
		res[i] = item.ToResponseSrc()
	}
	return res
}

type ResponseComment struct {
	OverAll      string  `json:"overAll,omitempty"`
	Problem      Problem `json:"problem,omitempty"`
	NoProblem    uint    `json:"noProblem,omitempty"`
	Images       []Image `json:"images,omitempty"`
	ExtraComment string  `json:"extraComment,omitempty"`
}

type Problem struct {
	Defect     uint `json:"defect,omitempty"`
	Incomplate uint `json:"incomplate,omitempty"`
	Mismatch   uint `json:"mismatch,omitempty"`
}

type Image struct {
	Src string `json:"src,omitempty"`
}
