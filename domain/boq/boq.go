package boq

import "cpm-rad-backend/domain/request"

// type BoQRequest struct {
// 	SequencesNo           uint   `json:"sequencesNo"`
// 	BoQItemNo             string `json:"itemNo"`
// 	BoQItemName           string `json:"name"`
// 	BoQItemGroup          string `json:"group"`
// 	BoQItemQuantity       string `json:"quantity"`
// 	BoQItemAmountDelivery string `json:"delivery"`
// 	BoQItemAmountGood     string `json:"good"`
// 	BoQItemAmountBad      string `json:"bad"`
// }

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

type Item struct {
	SequencesNo  uint   `gorm:"column:SEQUENCES_NO"`
	ID           uint   `gorm:"column:ID"`
	ItemNo       string `gorm:"column:ITEM"`
	ItemName     string `gorm:"column:NAME"`
	ItemGroup    string `gorm:"column:GROUPNAME"`
	ItemQuantity string `gorm:"column:QUANTITY"`
}

type Items []Item

type ItemSearchSpec struct {
	request.Pagination
	SequencesNo      int
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

// func (req *BoQRequest) ToModel() BoQItem {

// 	boqItem := BoQItem{
// 		SequencesNo:     req.SequencesNo,
// 		BoQID:           0,
// 		BoQItemNo:       req.BoQItemNo,
// 		BoQItemName:     req.BoQItemName,
// 		BoQItemGroup:    req.BoQItemGroup,
// 		BoQItemQuantity: 0,
// 		BoQItemUnit:     "",
// 	}

// 	return boqItem
// }

func (item *Item) ToResponse() Response {
	res := Response{
		SequencesNo:        item.SequencesNo,
		ID:                 item.ID,
		ItemNo:             item.ItemNo,
		ItemName:           item.ItemName,
		ItemGroup:          item.ItemGroup,
		ItemQuantity:       item.ItemQuantity,
		ItemAmountDelivery: "",
		ItemAmountGood:     "",
		ItemAmountBad:      "",
	}

	return res
}

func (items *Items) ToResponse() []Response {
	responses := make([]Response, len(*items))
	for i, item := range *items {
		responses[i] = item.ToResponse()
	}
	return responses
}
