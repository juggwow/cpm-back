package report

import (
	"cpm-rad-backend/domain/request"
	"cpm-rad-backend/domain/utils"
	"fmt"
)

type ProgressReport struct {
	Seq            int        `json:"seq"`
	ID             int        `json:"id"`
	InvoidNo       string     `json:"invNo"`
	BoqItemName    string     `json:"itemName"`
	ArrivalDate    utils.Time `json:"arrival"`
	InspectionDate utils.Time `json:"inspection"`
	StateID        int        `json:"stateID"`
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
		Seq:            int(p.SequencesNo),
		ID:             int(p.ID),
		InvoidNo:       p.Invoice,
		BoqItemName:    p.ItemName,
		ArrivalDate:    p.Arrival,
		InspectionDate: p.Inspection,
		StateID:        int(p.StateID),
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
