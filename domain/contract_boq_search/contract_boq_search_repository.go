package contract_boq_search

import (
	"context"
	"cpm-rad-backend/domain/connection"
	"cpm-rad-backend/domain/contract_boq"

	"gorm.io/gorm"
)

func Get(db *connection.DBConnection) getFunc {
	return func(ctx context.Context, spec BoQSearchSpec, id uint) (contract_boq.BoQItems, int64, error) {
		jobs, count, err := spec.searchBoQItems(db.RAD, &id)
		if err != nil {
			return jobs, count, err
		}

		return jobs, count, err
	}
}

func (spec *BoQSearchSpec) searchBoQItems(db *gorm.DB, id *uint) (contract_boq.BoQItems, int64, error) {
	var boqItems contract_boq.BoQItems

	// jobIDs := spec.getSearchCriteriaIDs(db)
	// if jobIDs != nil && len(*jobIDs) == 0 {
	// 	return boqItems, 0, nil
	// }

	countDB := db.Model(&boqItems)
	countDB, err := spec.buildSearchCriteria(countDB, id)
	if err != nil {
		return boqItems, 0, err
	}
	var count int64
	if err := countDB.Count(&count).Error; err != nil {
		return boqItems, 0, err
	}

	findDB := db.Model(&boqItems)
	//findDB = findDB.Preload("EmployeeJobs.Employee").Preload("EmployeeJobs.EmployeeRole")
	findDB = findDB.Offset(spec.Offset()).Limit(spec.GetLimit())
	findDB, err = spec.buildSearchCriteria(findDB, id)
	if err != nil {
		return boqItems, 0, err
	}

	return boqItems, count, findDB.Find(&boqItems).Error
}

func (spec *BoQSearchSpec) buildSearchCriteria(db *gorm.DB, id *uint) (*gorm.DB, error) {

	db = db.Select("ROW_NUMBER() OVER(ORDER BY ID ASC) AS SEQUENCES_NO,ID,ITEM,NAME,GROUPNAME,CONCAT(QUANTITY,' ',UNIT)  AS QUANTITY").
		Where("WORK_CONTRACT_ID = ?", id).Order("ID ASC")

	return db, db.Error
}
