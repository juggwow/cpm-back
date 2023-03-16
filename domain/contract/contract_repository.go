package contract

import (
	"context"
	"cpm-rad-backend/domain/connection"
)

func GetByID(db *connection.DBConnection) getByIDFunc {
	return func(ctx context.Context, ID int) (Contract, error) {
		var result Contract
		cpm := db.CPM.Model(&result)
		//fmt.Printf("Fatal id : %d \n", ID)
		// cpm = cpm.Preload("EmployeeJobs.Employee")
		// cpm = cpm.Preload("EmployeeJobs.EmployeeRole")
		//db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&results)
		err := cpm.Where("WORK_CONTRACT_ID = ?", ID).
			Scan(&result).
			Error
		if err != nil {
			return result, err
		}

		return result, err
	}
}

func GetNumberOfItem(db *connection.DBConnection) getNumberOfItemFunc {
	return func(ctx context.Context, ID int) (NumberOfItems, error) {
		var result NumberOfItems
		// result = NumberOfItems{
		// 	All:           50,
		// 	AllComplete:   1,
		// 	AllIncomplete: 49,
		// 	Check:         4,
		// 	CheckGood:     3,
		// 	CheckWaste:    1,
		// 	Progress:      4,
		// }
		cpm := db.CPM.Model(&result)
		err := cpm.Where("WORK_CONTRACT_ID = ?", ID).
			Scan(&result).
			Error
		if err != nil {
			return result, err
		}
		return result, nil
	}
}
