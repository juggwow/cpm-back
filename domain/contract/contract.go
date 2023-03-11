package contract

type ContractRequest struct {
	ContractID uint `json:"contractID"`
}

type ContractResponse struct {
	ContractID uint   `json:"contractID"`
	WorkID     uint   `json:"workID"`
	Name       string `json:"name"`
}

type Contract struct {
	ContractID uint   `gorm:"column:WORK_CONTRACT_ID"`
	WorkID     uint   `gorm:"column:WORK_ID"`
	WorkName   string `gorm:"column:WORK_NAME"`
	WorkType   string `gorm:"column:WORK_TYPE"`
}

func (Contract) TableName() string {
	return "CPM.VIEW_RAD_WORK_NAME"
}

func (contract *Contract) ToResponse() ContractResponse {
	responses := ContractResponse{
		ContractID: contract.ContractID,
		WorkID:     contract.WorkID,
		Name:       contract.WorkType + " " + contract.WorkName,
	}
	return responses
}

type AllItem struct {
	Amount     uint
	Complete   uint
	Incomplete uint
}

type CheckItem struct {
	Amount uint
	Good   uint
	Waste  uint
}

type ProgressItem struct {
	Amount uint
}

type CardResponse struct {
	AllItem      AllItem      `json:"all"`
	ProgressItem ProgressItem `json:"progress"`
	CheckItem    CheckItem    `json:"check"`
}

type NumberOfItems struct {
	All           uint `gorm:"column:ALL"`
	AllComplete   uint `gorm:"column:ALL_COMPLETE"`
	AllIncomplete uint `gorm:"column:ALL_INCOMPLETE"`
	Check         uint `gorm:"column:CHECK"`
	CheckGood     uint `gorm:"column:CHECK_GOOD"`
	CheckWaste    uint `gorm:"column:CHECK_WASTE"`
	Progress      uint `gorm:"column:PROGRESS"`
}

func (NumberOfItems) TableName() string {
	return "CPM.VIEW_RAD_NUMBER_OF_ITEMS"
}

func (n *NumberOfItems) ToResponse() CardResponse {
	responses := CardResponse{
		AllItem: AllItem{
			Amount:     n.All,
			Complete:   n.AllComplete,
			Incomplete: n.AllIncomplete,
		},
		ProgressItem: ProgressItem{
			Amount: n.Progress,
		},
		CheckItem: CheckItem{
			Amount: n.Check,
			Good:   n.CheckGood,
			Waste:  n.CheckWaste,
		},
	}
	return responses
}
