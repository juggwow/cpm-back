package boq

import (
	"cpm-rad-backend/domain/request"
	"fmt"

	"github.com/shopspring/decimal"
)

type Response struct {
	SequencesNo        uint   `json:"sequencesNo"`
	ID                 uint   `json:"boqID"`
	ItemNo             string `json:"itemNo"`
	ItemName           string `json:"name"`
	ItemGroup          string `json:"group"`
	ItemQuantity       string `json:"quantity"`
	ItemAmountDelivery string `json:"delivery"`
	ItemAmountGood     string `json:"good"`
	ItemAmountBad      string `json:"bad"`
}

type ItemResponse struct {
	ID           uint            `json:"boqID" gorm:"column:ID"`
	ItemName     string          `json:"name" gorm:"column:NAME"`
	ItemQuantity decimal.Decimal `json:"quantity" gorm:"column:QUANTITY"`
	ItemUnit     string          `json:"unit" gorm:"column:UNIT"`
}

type Item struct {
	SequencesNo          uint            `gorm:"column:SEQUENCES_NO"`
	ID                   uint            `gorm:"column:ID"`
	ItemNo               string          `gorm:"column:ITEM"`
	ItemName             string          `gorm:"column:NAME"`
	ItemGroup            string          `gorm:"column:GROUPNAME"`
	ItemQuantity         decimal.Decimal `gorm:"column:QUANTITY"`
	ItemUnit             string          `gorm:"column:UNIT"`
	ItemDelivelyQuantity decimal.Decimal `gorm:"column:DELIVERY_QUANTITY"`
	ItemReceiveQuantity  decimal.Decimal `gorm:"column:RECEIVE_QUANTITY"`
	ItemDamageQuantity   decimal.Decimal `gorm:"column:DAMAGE_QUANTITY"`
}

type Items []Item

type ItemSearchSpec struct {
	request.Pagination
	SequencesNo      string
	ItemNo           string
	ItemName         string
	ItemGroup        string
	ItemQuantity     string
	ItemDelivery     string
	ItemReceive      string
	ItemDamage       string
	SortSequencesNo  string
	SortItemNo       string
	SortItemName     string
	SortItemGroup    string
	SortItemQuantity string
	SortItemDelivery string
	SortItemReceive  string
	SortItemDamage   string
}

func (Item) TableName() string {
	return "CPM.VIEW_RAD_BOQ_ITEMS"
}

func (ItemResponse) TableName() string {
	return "CPM.VIEW_RAD_BOQ_ITEMS"
}

func (item *Item) ToResponse() Response {
	res := Response{
		SequencesNo:        item.SequencesNo,
		ID:                 item.ID,
		ItemNo:             item.ItemNo,
		ItemName:           item.ItemName,
		ItemGroup:          item.ItemGroup,
		ItemQuantity:       checkZero(item.ItemQuantity, item.ItemUnit),
		ItemAmountDelivery: checkZero(item.ItemDelivelyQuantity, item.ItemUnit),
		ItemAmountGood:     checkZero(item.ItemReceiveQuantity, item.ItemUnit),
		ItemAmountBad:      checkZero(item.ItemDamageQuantity, item.ItemUnit),
	}

	return res
}

func (items *Items) ToResponse() []Response {
	res := make([]Response, len(*items))
	for i, item := range *items {
		res[i] = item.ToResponse()
	}
	return res
}

func checkZero(value decimal.Decimal, unit string) string {
	if value.IsZero() {
		return "-"
	}
	return fmt.Sprintf("%v %v", value, unit)
}
