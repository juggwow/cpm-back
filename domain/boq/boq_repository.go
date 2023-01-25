package boq

import (
	"context"
	"cpm-rad-backend/domain/connection"

	"gorm.io/gorm"
)

func Get(db *connection.DBConnection) getFunc {
	return func(ctx context.Context, spec ItemSearchSpec, id uint) (Items, int64, error) {
		jobs, count, err := spec.searchBoQItems(db.CPM, &id)
		if err != nil {
			return jobs, count, err
		}

		return jobs, count, err
	}
}

func (spec *ItemSearchSpec) searchBoQItems(db *gorm.DB, id *uint) (Items, int64, error) {
	var Items Items

	// jobIDs := spec.getSearchCriteriaIDs(db)
	// if jobIDs != nil && len(*jobIDs) == 0 {
	// 	return boqItems, 0, nil
	// }

	countDB := db.Model(&Items)
	countDB, err := spec.buildSearchCriteria(countDB, id)
	if err != nil {
		return Items, 0, err
	}
	var count int64
	if err := countDB.Count(&count).Error; err != nil {
		return Items, 0, err
	}

	findDB := db.Model(&Items)
	//findDB = findDB.Preload("EmployeeJobs.Employee").Preload("EmployeeJobs.EmployeeRole")
	findDB = findDB.Offset(spec.Offset()).Limit(spec.GetLimit())
	findDB, err = spec.buildSearchCriteria(findDB, id)
	if err != nil {
		return Items, 0, err
	}

	return Items, count, findDB.Find(&Items).Error
}

func (spec *ItemSearchSpec) buildSearchCriteria(db *gorm.DB, id *uint) (*gorm.DB, error) {

	db = db.Select("ROW_NUMBER() OVER(ORDER BY ID ASC) AS SEQUENCES_NO,ID,ITEM,NAME,GROUPNAME,QUANTITY,UNIT").
		Where("WORK_CONTRACT_ID = ?", id).Order("ID ASC")

	return db, db.Error
}

func GetItemByID(db *connection.DBConnection) getItemByIDFunc {
	return func(ctx context.Context, id uint) (ItemResponse, error) {
		var result ItemResponse
		cpm := db.CPM.Model(&result)
		err := cpm.Where("ID = ?", id).Scan(&result).Error

		if err != nil {
			return result, err
		}

		return result, err
	}
}
