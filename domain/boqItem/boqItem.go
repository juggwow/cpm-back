package boqItem

import "cpm-rad-backend/domain/request"

type BoqItemListDB struct {
	RowNo            uint   `gorm:"column:ROW_NO"`
	ID               uint   `gorm:"column:ID"`
	WorkContractID   uint   `gorm:"column:WORK_CONTRACT_ID"`
	Number           string `gorm:"column:NUMBER"`
	GroupNumber      string `gorm:"column:GROUP_NUMBER"`
	GroupName        string `gorm:"column:GROUP_NAME"`
	Name             string `gorm:"column:NAME"`
	Quantity         string `gorm:"column:QUANTITY"`
	Unit             string `gorm:"column:UNIT"`
	DeliveryQuantity string `gorm:"column:DELIVERY_QUANTITY"`
	ReceiveQuantity  string `gorm:"column:RECEIVE_QUANTITY"`
	DamageQuantity   string `gorm:"column:DAMAGE_QUANTITY"`
	ReceiveStatus    string `gorm:"column:RECEIVE_STATUS"`
}
type BoqItemListDBs []BoqItemListDB

type BoqItemList struct {
	RowNo uint `json:"rowNo"`
	ID    uint `json:"id"`
	// WorkContractID   uint   `json:"workContractID"`
	Number string `json:"number"`
	// GroupNumber      string `json:"groupNumber"`
	GroupName        string `json:"groupName"`
	Name             string `json:"name"`
	Quantity         string `json:"quantity"`
	DeliveryQuantity string `json:"deliveryQty"`
	ReceiveQuantity  string `json:"receiveQty"`
	DamageQuantity   string `json:"damageQty"`
	ReceiveStatus    string `json:"receiveStatus"`
}
type BoqItemLists []BoqItemList

func (item *BoqItemListDB) ToResponse() BoqItemList {
	return BoqItemList{
		RowNo:            item.RowNo,
		ID:               item.ID,
		Number:           item.Number,
		GroupName:        item.GroupName,
		Name:             item.Name,
		Quantity:         item.Quantity,
		DeliveryQuantity: item.DeliveryQuantity,
		ReceiveQuantity:  item.ReceiveQuantity,
		DamageQuantity:   item.DamageQuantity,
		ReceiveStatus:    item.ReceiveStatus,
	}
}

func (items *BoqItemListDBs) ToResponse() BoqItemLists {
	res := make(BoqItemLists, len(*items))
	for i, item := range *items {
		res[i] = item.ToResponse()
	}
	return res
}

type SearchSpec struct {
	request.Pagination
	SearchRowNo       string
	SearchNumber      string
	SearchGroupName   string
	SearchName        string
	SearchQuantity    string
	SearchDeliveryQty string
	SearchReceiveQty  string
	SearchDamageQty   string
	SortRowNo         string
	SortNumber        string
	SortGroupName     string
	SortName          string
	SortQuantity      string
	SortDeliveryQty   string
	SortReceiveQty    string
	SortDamageQty     string
}
