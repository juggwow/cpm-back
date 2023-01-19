package contract_boq

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

type BoQResponse struct {
	SequencesNo           uint   `json:"sequencesNo"`
	BoQID                 uint   `json:"boqID"`
	BoQItemNo             string `json:"itemNo"`
	BoQItemName           string `json:"name"`
	BoQItemGroup          string `json:"group"`
	BoQItemQuantity       string `json:"quantity"`
	BoQItemAmountDelivery string `json:"delivery"`
	BoQItemAmountGood     string `json:"good"`
	BoQItemAmountBad      string `json:"bad"`
}

type BoQItem struct {
	SequencesNo     uint   `gorm:"column:SEQUENCES_NO"`
	BoQID           uint   `gorm:"column:ID"`
	BoQItemNo       string `gorm:"column:ITEM"`
	BoQItemName     string `gorm:"column:NAME"`
	BoQItemGroup    string `gorm:"column:GROUPNAME"`
	BoQItemQuantity string `gorm:"column:QUANTITY"`
}

type BoQItems []BoQItem

func (BoQItem) TableName() string {
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

func (boqItem *BoQItem) ToResponse() BoQResponse {
	res := BoQResponse{
		SequencesNo:           boqItem.SequencesNo,
		BoQID:                 boqItem.BoQID,
		BoQItemNo:             boqItem.BoQItemNo,
		BoQItemName:           boqItem.BoQItemName,
		BoQItemGroup:          boqItem.BoQItemGroup,
		BoQItemQuantity:       boqItem.BoQItemQuantity,
		BoQItemAmountDelivery: "",
		BoQItemAmountGood:     "",
		BoQItemAmountBad:      "",
	}

	return res
}

func (boqItems *BoQItems) ToResponse() []BoQResponse {
	responses := make([]BoQResponse, len(*boqItems))
	for i, boqItem := range *boqItems {
		responses[i] = boqItem.ToResponse()
	}
	return responses
}
