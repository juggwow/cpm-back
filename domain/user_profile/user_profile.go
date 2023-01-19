package user_profile

type UserProfile struct {
	ID             uint   `gorm:"column:ID;primary"`
	EmpID          string `gorm:"column:EMP_ID"`
	BusinessArea   string `gorm:"column:BUSINESS_AREA"`
	FirstName      string `gorm:"column:FIRST_NAME"`
	LastName       string `gorm:"column:LAST_NAME"`
	DeptChangeCode string `gorm:"column:DEPT_CHANGE_CODE"`
	Title          string `gorm:"column:TITLE_S_DESC"`
	Position       string `gorm:"column:PLANS_TEXT_SHORT"`
}

type UserProfiles []UserProfile

func (UserProfile) TableName() string {
	return "CPM_USER_PROFILE"
}
