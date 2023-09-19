package contract

type ContractRequest struct {
	ContractID uint `json:"contractID"`
}

type ContractResponse struct {
	ContractID       uint   `json:"contractID"`
	WorkID           uint   `json:"workID"`
	WorkName         string `json:"workName"`
	WorkType         string `json:"workType"`
	ProjectName      string `json:"projectName"`
	ProjectShortName string `json:"projectShortName"`
}

type Contract struct {
	ContractID       uint   `gorm:"column:WORK_CONTRACT_ID"`
	WorkID           uint   `gorm:"column:WORK_ID"`
	WorkName         string `gorm:"column:WORK_NAME"`
	WorkType         string `gorm:"column:WORK_TYPE"`
	ProjectName      string `gorm:"column:PROJECT_NAME"`
	ProjectShortName string `gorm:"column:PROJECT_SHORT_NAME"`
}

func (Contract) TableName() string {
	return "CPM.VIEW_RAD_WORK_NAME"
}

func (contract *Contract) ToResponse() ContractResponse {
	responses := ContractResponse{
		ContractID:       contract.ContractID,
		WorkID:           contract.WorkID,
		WorkName:         contract.WorkName,
		WorkType:         contract.WorkType,
		ProjectName:      contract.ProjectName,
		ProjectShortName: contract.ProjectShortName,
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
	All           uint `gorm:"column:ALL_ITEM"`
	AllComplete   uint `gorm:"column:COMPLETE"`
	AllIncomplete uint `gorm:"column:INCOMPLETE"`
	Check         uint `gorm:"column:ALL_CHECK"`
	CheckGood     uint `gorm:"column:CHECK_GOOD"`
	CheckWaste    uint `gorm:"column:CHECK_WASTE"`
	Progress      uint `gorm:"column:PROGRESS"`
}

func (NumberOfItems) TableName() string {
	return "CPM.SUM_NUMBER_OF_ITEMS"
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
