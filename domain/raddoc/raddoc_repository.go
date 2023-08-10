package raddoc

import (
	"context"
	"cpm-rad-backend/domain/connection"

	"gorm.io/gorm"
)

func GetByItem(db *connection.DBConnection) getByItemFunc {
	return func(ctx context.Context, spec SearchSpec, id uint) (Response, int64, error) {
		var result Response
		lod, count, err := spec.search(db.CPM, &id)
		if err != nil {
			return result, count, err
		}
		var item Item

		cpm := db.CPM.Model(&item)
		err = cpm.Select("ID,NAME,CONCAT(QUANTITY,' ',UNIT) AS CONTRACT_QTY,CONCAT(RECEIVE_QUANTITY,' ',UNIT) AS RECEIVE_QTY").
			Where("ID = ?", id).Scan(&item).Error

		result = Response{
			Item:      item,
			ListOfDoc: lod,
		}

		return result, count, err
	}

}

func (spec *SearchSpec) search(db *gorm.DB, id *uint) (listOfDoc, int64, error) {
	var lod listOfDoc

	// jobIDs := spec.getSearchCriteriaIDs(db)
	// if jobIDs != nil && len(*jobIDs) == 0 {
	// 	return boqItems, 0, nil
	// }

	countDB := db.Model(&lod)
	countDB, err := spec.buildSearchCriteria(countDB, id)
	if err != nil {
		return lod, 0, err
	}
	var count int64
	if err := countDB.Count(&count).Error; err != nil {
		return lod, 0, err
	}

	findDB := db.Model(&lod)
	//findDB = findDB.Preload("EmployeeJobs.Employee").Preload("EmployeeJobs.EmployeeRole")
	findDB = findDB.Offset(spec.Offset()).Limit(spec.GetLimit())
	findDB, err = spec.buildSearchCriteria(findDB, id)
	if err != nil {
		return lod, 0, err
	}

	return lod, count, findDB.Find(&lod).Error
}

func (spec *SearchSpec) buildSearchCriteria(db *gorm.DB, id *uint) (*gorm.DB, error) {

	db = db.Where("BOQ_ID = ?", id).Order("ID DESC")

	return db, db.Error
}
