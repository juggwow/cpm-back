package contract

import (
	"context"
	"cpm-rad-backend/domain/connection"
)

func GetByID(db *connection.DBConnection) getByIDFunc {
	return func(ctx context.Context, ID int) (Contract, error) {
		result := Contract{}
		cpm := db.CPM
		//fmt.Printf("Fatal id : %d \n", ID)
		// cpm = cpm.Preload("EmployeeJobs.Employee")
		// cpm = cpm.Preload("EmployeeJobs.EmployeeRole")
		//db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&results)
		err := cpm.Table("CPM.CPM_WORK_CONTRACT as cwc").
			Select("cwc.ID as WORK_CONTRACT_ID , cwc.WORK_ID ,vawdd.WORK_NAME ,vawdd.WORK_TYPE_DESCRIPTION").
			Joins("LEFT OUTER JOIN CPM.VIEW_ASSIGN_WORK_DESIGN_DETAIL as vawdd ON cwc.WORK_ID = vawdd.WORK_ID").
			Where("cwc.ID = ?", ID).
			Scan(&result).
			Error
		if err != nil {
			return result, err
		}

		return result, err
	}
}
