package report

type Report struct {
	ItemID       string      `json:"itemID"`
	Arrival      string      `json:"arrival"`
	Inspection   string      `json:"inspection"`
	TaskMaster   string      `json:"taskMaster"`
	Invoice      string      `json:"invoice"`
	Quantity     string      `json:"quantity"`
	Country      string      `json:"country"`
	Manufacturer string      `json:"manufacturer"`
	Model        string      `json:"model"`
	Serial       string      `json:"serial"`
	PeaNo        string      `json:"peano"`
	CreateBy     string      `json:"createby"`
	Status       string      `json:"status"`
	AttachFiles  AttachFiles `json:"filesAttach"`
}

type AttachFile struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Size string `json:"size"`
	Unit string `json:"unit"`
	Path string `json:"filePath"`
	Type uint   `json:"docType"`
}

type AttachFiles []AttachFile
