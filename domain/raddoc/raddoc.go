package raddoc

import (
	"cpm-rad-backend/domain/rad"
	"cpm-rad-backend/domain/request"
)

type Response struct {
	Item      Item      `json:"item"`
	ListOfDoc listOfDoc `json:"data"`
}

type Item struct {
	ID          uint   `json:"id" gorm:"column:ID"`
	Name        string `json:"name" gorm:"column:NAME"`
	ContractQTY string `json:"contractQTY" gorm:"column:CONTRACT_QTY"`
	ReceiveQTY  string `json:"receiveQTY" gorm:"column:RECEIVE_QTY"`
}

func (Item) TableName() string {
	return "CPM.VIEW_RAD_BOQ_ITEMS"
}

type RadDoc struct {
	ID         uint     `json:"id" gorm:"column:ID"`
	Seq        uint     `json:"seq" gorm:"column:SEQ_NO"`
	InvNo      string   `json:"invNo" gorm:"column:CONTRACTOR_INV_NO"`
	Qty        string   `json:"qty" gorm:"column:DELIVERED_QTY"`
	Arrival    rad.Time `json:"arrival" gorm:"column:ARRIVAL_DATE_AT_SITE"`
	Inspection rad.Time `json:"inspection" gorm:"column:INSPECTION_DATE"`
	CreateBy   string   `json:"createBy" gorm:"column:CREATED_BY"`
	StateID    uint     `json:"stateID" gorm:"column:STATE_ID"`
	StateName  string   `json:"stateName" gorm:"column:STATE_NAME"`
}
type listOfDoc []RadDoc

func (RadDoc) TableName() string {
	return "CPM.VIEW_RAD_LIST_DOC"
}

type SearchSpec struct {
	request.Pagination
	SeqNo          string
	InvNo          string
	Qty            string
	Arrival        string
	Inspection     string
	CreateBy       string
	StateName      string
	SortSeqNo      string
	SortInvNo      string
	SortQty        string
	SortArrival    string
	SortInspection string
	SortCreateBy   string
	SortStateName  string
}
