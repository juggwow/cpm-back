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
