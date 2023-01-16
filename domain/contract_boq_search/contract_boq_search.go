package contract_boq_search

import "cpm-rad-backend/domain/request"

type BoQSearchSpec struct {
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
