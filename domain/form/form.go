package form

import (
	"cpm-rad-backend/domain/boq"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type Request struct {
	ItemID       uint        `json:"itemID"`
	Arrival      RadTime     `json:"arrival"`
	Inspection   RadTime     `json:"inspection"`
	TaskMaster   string      `json:"taskMaster"`
	Invoice      string      `json:"invoice"`
	Quantity     uint        `json:"quantity"`
	Country      string      `json:"country"`
	Manufacturer string      `json:"manufacturer"`
	Model        string      `json:"model"`
	Serial       string      `json:"serial"`
	PeaNo        string      `json:"peano"`
	CreateBy     string      `json:"createby"`
	Status       int         `json:"status"`
	FilesAttach  FilesAttach `json:"filesAttach"`
}

type Response struct {
	ID           uint            `json:"id"`
	ItemID       uint            `json:"itemID"`
	ItemName     string          `json:"itemName"`
	ItemQty      decimal.Decimal `json:"itemQty"`
	ItemUnit     string          `json:"itemUnit"`
	Arrival      RadTime         `json:"arrival"`
	Inspection   RadTime         `json:"inspection"`
	TaskMaster   string          `json:"taskMaster"`
	Invoice      string          `json:"invoice"`
	Qty          uint            `json:"quantity"`
	Country      string          `json:"country"`
	Manufacturer string          `json:"manufacturer"`
	Model        string          `json:"model"`
	Serial       string          `json:"serial"`
	PeaNo        string          `json:"peano"`
	Files        Files           `json:"filesAttach"`
}

func (form *Form) ToResponse(file Files, item boq.ItemResponse) Response {
	res := Response{
		ID:           form.ID,
		ItemID:       form.ItemID,
		ItemName:     item.ItemName,
		ItemQty:      item.ItemQuantity,
		ItemUnit:     item.ItemUnit,
		Arrival:      RadTime(form.Arrival),
		Inspection:   RadTime(form.Inspection),
		TaskMaster:   form.TaskMaster,
		Invoice:      form.Invoice,
		Qty:          form.Quantity,
		Country:      form.Country,
		Manufacturer: form.Manufacturer,
		Model:        form.Model,
		Serial:       form.Serial,
		PeaNo:        form.PeaNo,
		Files:        file,
	}
	return res
}

type RadTime time.Time

type Country struct {
	ID   uint   `json:"id" gorm:"column:ID"`
	Code string `json:"code" gorm:"column:CODE"`
	Name string `json:"name" gorm:"column:NAME"`
}

type Countrys []Country

func (Country) TableName() string {
	return `[WDDEVDB\WORKD].CPM.CPM.COUNTRY`
}

type Form struct {
	ID           uint      `gorm:"column:ID"`
	ItemID       uint      `gorm:"column:BOQ_ID"`
	RadNo        string    `gorm:"column:RAD_NO"`
	Arrival      time.Time `gorm:"column:ARRIVAL_DATE_AT_SITE"`
	Inspection   time.Time `gorm:"column:INSPECTION_DATE"`
	TaskMaster   string    `gorm:"column:NAME_OF_TASKMASTER"`
	Invoice      string    `gorm:"column:CONTRACTOR_INV_NO"`
	Quantity     uint      `gorm:"column:QUANTITY"`
	Country      string    `gorm:"column:COUNTRY"`
	Manufacturer string    `gorm:"column:MANUFACTURER"`
	Model        string    `gorm:"column:MODEL"`
	Serial       string    `gorm:"column:SERIAL_NO"`
	PeaNo        string    `gorm:"column:PEA_NO"`
	CreateBy     string    `gorm:"column:CREATED_BY"`
	Status       int       `gorm:"column:STATE_ID"`
}

func (Form) TableName() string {
	return "CPM.CPM_WORK_CONTRACT_RAD"
}

func (req *Request) ToModel() Form {
	form := Form{
		ItemID:       req.ItemID,
		Arrival:      time.Time(req.Arrival),
		Inspection:   time.Time(req.Inspection),
		TaskMaster:   req.TaskMaster,
		Invoice:      req.Invoice,
		Quantity:     req.Quantity,
		Country:      req.Country,
		Manufacturer: req.Manufacturer,
		Model:        req.Model,
		Serial:       req.Serial,
		PeaNo:        req.PeaNo,
		Status:       req.Status,
	}

	return form
}

func (c *RadTime) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`) //get rid of "
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse("2006-01-02", value) //parse time
	if err != nil {
		return err
	}
	*c = RadTime(t) //set result using the pointer
	return nil
}

func (c RadTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(c).Format("2006-01-02") + `"`), nil
}

type FileUploadResponse struct {
	Name        string `json:"name"`
	ObjectName  string `json:"objectName"`
	DisplaySize string `json:"displaySize"`
	Size        string `json:"size"`
	Unit        string `json:"unit"`
	FileType    string `json:"fileType"`
	FilePath    string `json:"filePath"`
}

type FileUploadResponses []FileUploadResponse

type FileAttach struct {
	Name string `json:"name"`
	Size string `json:"size"`
	Unit string `json:"unit"`
	Path string `json:"filePath"`
	Type uint   `json:"type"`
}

type FilesAttach []FileAttach

type File struct {
	ID      uint            `gorm:"column:ID"`
	RadID   uint            `gorm:"column:RAD_ID"`
	DocType uint            `gorm:"column:RAD_DOC_TYPE_ID"`
	Name    string          `gorm:"column:FILE_NAME"`
	Size    decimal.Decimal `gorm:"column:FILE_SIZE"`
	Unit    string          `gorm:"column:FILE_UNIT"`
	Path    string          `gorm:"column:FILE_PATH"`
}

type Files []File

func (File) TableName() string {
	return "CPM.CPM_WORK_CONTRACT_RAD_FILE"
}

type FileCreate struct {
	RadID    uint            `gorm:"column:RAD_ID"`
	DocType  uint            `gorm:"column:RAD_DOC_TYPE_ID"`
	Name     string          `gorm:"column:FILE_NAME"`
	Size     decimal.Decimal `gorm:"column:FILE_SIZE"`
	Unit     string          `gorm:"column:FILE_UNIT"`
	Path     string          `gorm:"column:FILE_PATH"`
	CreateBy string          `gorm:"column:CREATED_BY"`
}

func (FileCreate) TableName() string {
	return "CPM.CPM_WORK_CONTRACT_RAD_FILE"
}

func (f *FileAttach) ToModel(radID uint, createBy string) FileCreate {
	size, _ := decimal.NewFromString(f.Size)
	file := FileCreate{
		RadID:    radID,
		DocType:  f.Type,
		Name:     f.Name,
		Size:     size,
		Unit:     f.Unit,
		Path:     f.Path,
		CreateBy: createBy,
	}

	return file
}

type DocType struct {
	ID   uint   `json:"id" gorm:"column:ID"`
	Name string `json:"name" gorm:"column:DESCRIPTION"`
}
type DocTypes []DocType

func (DocType) TableName() string {
	return "CPM.CPM_RAD_DOC_TYPE"
}
