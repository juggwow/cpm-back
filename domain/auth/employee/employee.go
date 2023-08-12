package employee

import (
	"cpm-rad-backend/domain/auth/user_profile"
	"cpm-rad-backend/domain/connection"
)

type Employee struct {
	connection.BaseModel
	EmployeeID     string `gorm:"column:EMPLOYEE_ID"`
	Token          string `gorm:"column:TOKEN"`
	Title          string `gorm:"column:TITLE"`
	FirstName      string `gorm:"column:FIRSTNAME"`
	LastName       string `gorm:"column:LASTNAME"`
	BA             string `gorm:"column:BA"`
	DeptChangeCode string `gorm:"column:DEPT_CHANGE_CODE"`
	Position       string `gorm:"column:POSITION"`
}

func (Employee) TableName() string {
	return "EMPLOYEE"
}

type Employees []Employee

type EmployeeRequest struct {
	EmployeeID string `json:"employeeId"`
	Title      string `json:"title"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Position   string `json:"position"`
}

type EmployeeResponse struct {
	ID             uint   `json:"id"`
	EmployeeID     string `json:"employeeId"`
	Title          string `json:"title"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Position       string `json:"position"`
	BA             string `json:"businessArea"`
	DeptChangeCode string `json:"deptChangeCode"`
}

func (employee *Employee) ToResponse() *EmployeeResponse {
	return &EmployeeResponse{
		ID:             employee.ID,
		EmployeeID:     employee.EmployeeID,
		Title:          employee.Title,
		FirstName:      employee.FirstName,
		LastName:       employee.LastName,
		Position:       employee.Position,
		BA:             employee.BA,
		DeptChangeCode: employee.DeptChangeCode,
	}
}

func (req *EmployeeRequest) ToModel() *Employee {
	return &Employee{
		EmployeeID: req.EmployeeID,
		Title:      req.Title,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Position:   req.Position,
	}
}

func (employee *Employee) Parse(u *user_profile.UserProfile) {
	employee.EmployeeID = u.EmpID
	employee.FirstName = u.FirstName
	employee.LastName = u.LastName
	employee.Title = u.Title
	employee.Position = u.Position
	employee.BA = u.BusinessArea
	employee.DeptChangeCode = u.DeptChangeCode
}
