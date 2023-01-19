package contract

import (
	"context"
	"cpm-rad-backend/domain/connection"
)

func GetByID(db *connection.DBConnection) getByIDFunc {
	return func(ctx context.Context, ID int) (Contract, error) {
		var result Contract
		rad := db.RAD.Model(&result)
		//fmt.Printf("Fatal id : %d \n", ID)
		// cpm = cpm.Preload("EmployeeJobs.Employee")
		// cpm = cpm.Preload("EmployeeJobs.EmployeeRole")
		//db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&results)
		err := rad.Where("WORK_CONTRACT_ID = ?", ID).
			Scan(&result).
			Error
		if err != nil {
			return result, err
		}

		return result, err
	}
}
