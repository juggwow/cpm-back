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
	var result listOfDoc

	countDB := db.Model(&result)
	countDB = countDB.Where("BOQ_ID = ?", id) //.Order("ID DESC")
	countDB, err := spec.buildSearch(countDB)
	if err != nil {
		return result, 0, err
	}
	var count int64
	if err := countDB.Count(&count).Error; err != nil {
		return result, 0, err
	}

	countDB = countDB.Offset(spec.Offset()).Limit(spec.GetLimit())
	countDB, err = spec.buildOrder(countDB)
	if err != nil {
		return result, 0, err
	}

	return result, count, countDB.Find(&result).Error
}

func (spec *SearchSpec) buildSearch(db *gorm.DB) (*gorm.DB, error) {

	if spec.SeqNo != "" {
		db = db.Where("SEQ_NO LIKE ?", "%"+spec.SeqNo+"%")
	}

	if spec.InvNo != "" {
		db = db.Where("CONTRACTOR_INV_NO LIKE ?", "%"+spec.InvNo+"%")
	}

	if spec.Qty != "" {
		db = db.Where("DELIVERED_QTY LIKE ?", "%"+spec.Qty+"%")
	}

	if spec.Arrival != "" {
		db = db.Where("ARRIVAL_DATE_AT_SITE LIKE ?", "%"+spec.Arrival+"%")
	}

	if spec.Inspection != "" {
		db = db.Where("INSPECTION_DATE LIKE ?", "%"+spec.Inspection+"%")
	}

	if spec.CreateBy != "" {
		db = db.Where("CREATED_BY LIKE ?", "%"+spec.CreateBy+"%")
	}

	if spec.StateName != "" {
		db = db.Where("STATE_NAME LIKE ?", "%"+spec.StateName+"%")
	}

	return db, db.Error
}

func (spec *SearchSpec) buildOrder(db *gorm.DB) (*gorm.DB, error) {

	if spec.SortSeqNo != "" {
		db = db.Order("SEQ_NO " + spec.SortSeqNo)
	}

	if spec.SortInvNo != "" {
		db = db.Order("CONTRACTOR_INV_NO " + spec.SortInvNo)
	}

	if spec.SortQty != "" {
		db = db.Order("DELIVERED_QTY " + spec.SortQty)
	}

	if spec.SortArrival != "" {
		db = db.Order("ARRIVAL_DATE_AT_SITE " + spec.SortArrival)
	}

	if spec.SortInspection != "" {
		db = db.Order("INSPECTION_DATE " + spec.SortInspection)
	}

	if spec.SortCreateBy != "" {
		db = db.Order("CREATED_BY " + spec.SortCreateBy)
	}

	if spec.SortStateName != "" {
		db = db.Order("STATE_NAME " + spec.SortStateName)
	}

	return db, db.Error
}
