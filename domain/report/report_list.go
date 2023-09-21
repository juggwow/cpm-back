package report

import (
	"cpm-rad-backend/domain/request"
	"cpm-rad-backend/domain/utils"
	"fmt"
)

type ProgressReport struct {
	Seq            uint       `json:"seq"`
	ID             uint       `json:"id"`
	InvoidNo       string     `json:"invNo"`
	ItemID         uint       `json:"itemID"`
	ItemName       string     `json:"itemName"`
	ArrivalDate    utils.Time `json:"arrival"`
	InspectionDate utils.Time `json:"inspection"`
	StateID        uint       `json:"stateID"`
	StateName      string     `json:"stateName"`
}
type ProgressReports []ProgressReport

type ProgressReportDB struct {
	SequencesNo uint       `gorm:"column:SEQ_NO"`
	ID          uint       `gorm:"column:ID"`
	ItemID      uint       `gorm:"column:BOQ_ID"`
	ItemName    string     `gorm:"column:BOQ_ITEM_NAME"`
	Invoice     string     `gorm:"column:CONTRACTOR_INV_NO"`
	Arrival     utils.Time `gorm:"column:ARRIVAL_DATE_AT_SITE"`
	Inspection  utils.Time `gorm:"column:INSPECTION_DATE"`
	StateID     uint       `gorm:"column:STATE_ID"`
	StateName   string     `gorm:"column:STATE_NAME"`
}

type ProgressReportDBs []ProgressReportDB

func (ProgressReportDB) TableName() string {
	return "CPM.RAD_LIST_PROGRESS_DOC"
}

type ProgressReportSearch struct {
	request.Pagination
	SequencesNo     string
	ItemName        string
	Invoice         string
	Arrival         string
	Inspection      string
	StateName       string
	SortSequencesNo string
	SortItemName    string
	SortInvoice     string
	SortArrival     string
	SortInspection  string
	SortStateName   string
}

func (p *ProgressReportDB) ToResponse() ProgressReport {
	return ProgressReport{
		Seq:            p.SequencesNo,
		ID:             p.ID,
		InvoidNo:       p.Invoice,
		ItemID:         p.ItemID,
		ItemName:       p.ItemName,
		ArrivalDate:    p.Arrival,
		InspectionDate: p.Inspection,
		StateID:        p.StateID,
		StateName:      p.StateName,
	}
}

func (p *ProgressReportDBs) ToResponse() []ProgressReport {
	res := make([]ProgressReport, len(*p))
	for i, item := range *p {
		res[i] = item.ToResponse()
	}
	return res
}

type CheckReport struct {
	Seq            int        `json:"seq"`
	ID             int        `json:"id"`
	InvoidNo       string     `json:"invNo"`
	BoqItemName    string     `json:"itemName"`
	ArrivalDate    utils.Time `json:"arrival"`
	InspectionDate utils.Time `json:"inspection"`
	Amount         string     `json:"amount"`
	Good           string     `json:"good"`
	Waste          string     `json:"waste"`
}
type CheckReports []CheckReport

type CheckReportDB struct {
	SequencesNo uint       `gorm:"column:SEQ_NO"`
	ID          uint       `gorm:"column:ID"`
	ItemID      uint       `gorm:"column:BOQ_ID"`
	ItemName    string     `gorm:"column:BOQ_ITEM_NAME"`
	Invoice     string     `gorm:"column:CONTRACTOR_INV_NO"`
	Arrival     utils.Time `gorm:"column:ARRIVAL_DATE_AT_SITE"`
	Inspection  utils.Time `gorm:"column:INSPECTION_DATE"`
	Amount      string     `gorm:"column:DELIVERED_QTY"`
	Good        uint       `gorm:"column:GOOD_QTY"`
	Waste       uint       `gorm:"column:WASTE_QTY"`
	Uint        string     `gorm:"column:UNIT"`
}

type CheckReportDBs []CheckReportDB

func (CheckReportDB) TableName() string {
	return "CPM.RAD_LIST_CHECK_DOC"
}

type CheckReportSearch struct {
	request.Pagination
	SequencesNo     string
	ItemName        string
	Invoice         string
	Arrival         string
	Inspection      string
	Amount          string
	Good            string
	Waste           string
	SortSequencesNo string
	SortItemName    string
	SortInvoice     string
	SortArrival     string
	SortInspection  string
	SortAmount      string
	SortGood        string
	SortWaste       string
}

func (p *CheckReportDB) ToResponse() CheckReport {
	return CheckReport{
		Seq:            int(p.SequencesNo),
		ID:             int(p.ID),
		InvoidNo:       p.Invoice,
		BoqItemName:    p.ItemName,
		ArrivalDate:    p.Arrival,
		InspectionDate: p.Inspection,
		Amount:         p.Amount,
		Good:           checkZero(p.Good, p.Uint),
		Waste:          checkZero(p.Waste, p.Uint),
	}
}

func (p *CheckReportDBs) ToResponse() []CheckReport {
	res := make([]CheckReport, len(*p))
	for i, item := range *p {
		res[i] = item.ToResponse()
	}
	return res
}

func checkZero(value uint, unit string) string {
	if value == 0 {
		return ""
	}
	return fmt.Sprintf("%v %v", value, unit)
}

type ResponseReport struct {
	Seq              uint       `json:"seq"`
	ID               uint       `json:"id"`
	DeliveryNumber   string     `json:"deliveryNumber,omitempty"`
	WorkName         string     `json:"workName,omitempty"`
	ProjectShortName string     `json:"projectShortName,omitempty"`
	ItemID           uint       `json:"itemID,omitempty"`
	ItemName         string     `json:"itemName,omitempty"`
	ArrivalDate      utils.Time `json:"arrival,omitempty"`
	InspectionDate   utils.Time `json:"inspection,omitempty"`
	StateID          uint       `json:"stateID,omitempty"`
	StateName        string     `json:"stateName,omitempty"`
}
type ResponseReportList []ResponseReport

type DbWaitForApprovReport struct {
	SequencesNo      uint       `gorm:"column:SEQ_NO"`
	ID               uint       `gorm:"column:ID"`
	ContractID       uint       `gorm:"column:CONTRACT_ID"`
	WorkName         string     `gorm:"column:WORK_NAME"`
	ProjectShortName string     `gorm:"column:PROJECT_SHORT_NAME"`
	ItemID           uint       `gorm:"column:ITEM_ID"`
	ItemName         string     `gorm:"column:ITEM_NAME"`
	DeliveryNumber   string     `gorm:"column:DELIVERY_NUMBER"`
	Arrival          utils.Time `gorm:"column:ARRIVAL_DATE"`
	Inspection       utils.Time `gorm:"column:INSPECTION_DATE"`
}

type DbWaitForApprovReports []DbWaitForApprovReport

// func (ProgressReportDB) TableName() string {
// 	return "CPM.RAD_LIST_PROGRESS_DOC"
// }

type SearchSortWaitForApprovReport struct {
	request.Pagination
	SequencesNo          string
	DeliveryNumber       string
	ItemName             string
	WorkName             string
	ProjectShortName     string
	Arrival              string
	Inspection           string
	SortSequencesNo      string
	SortDeliveryNumber   string
	SortItemName         string
	SortWorkName         string
	SortProjectShortName string
	SortArrival          string
	SortInspection       string
}

func (db *DbWaitForApprovReport) ToResponse() ResponseReport {
	return ResponseReport{
		Seq:              db.SequencesNo,
		ID:               db.ID,
		DeliveryNumber:   db.DeliveryNumber,
		WorkName:         db.WorkName,
		ProjectShortName: db.ProjectShortName,
		ItemID:           db.ItemID,
		ItemName:         db.ItemName,
		ArrivalDate:      db.Arrival,
		InspectionDate:   db.Inspection,
	}
}

func (p *DbWaitForApprovReports) ToResponse() []ResponseReport {
	res := make([]ResponseReport, len(*p))
	for i, item := range *p {
		res[i] = item.ToResponse()
	}
	return res
}
