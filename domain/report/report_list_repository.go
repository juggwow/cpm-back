package report

import (
	"context"
	"cpm-rad-backend/domain/connection"

	"gorm.io/gorm"
)

func GetProgressReport(db *connection.DBConnection) getProgressReportFunc {
	return func(ctx context.Context, spec ProgressReportSearch, id uint) (ProgressReports, int64, error) {
		var result ProgressReports
		lod, count, err := spec.search(db.CPM, &id)

		if err != nil {
			return result, count, err
		}
		result = lod.ToResponse()
		return result, count, err
	}

}

func (spec *ProgressReportSearch) search(db *gorm.DB, id *uint) (ProgressReportDBs, int64, error) {
	var lod ProgressReportDBs

	countDB := db.Model(&lod)
	countDB = countDB.Where("CONTRACT_ID = ?", id)
	countDB, err := spec.buildSearch(countDB)
	if err != nil {
		return lod, 0, err
	}
	var count int64
	if err := countDB.Count(&count).Error; err != nil {
		return lod, 0, err
	}

	countDB = countDB.Offset(spec.Offset()).Limit(spec.GetLimit())
	countDB, err = spec.buildOrder(countDB)
	if err != nil {
		return lod, 0, err
	}

	return lod, count, countDB.Find(&lod).Error
}

func (spec *ProgressReportSearch) buildSearch(db *gorm.DB) (*gorm.DB, error) {

	// db = db.Where("BOQ_ID = ?", id).Order("ID DESC")
	if spec.SequencesNo != "" {
		db = db.Where("SEQ_NO LIKE ?", "%"+spec.SequencesNo+"%")
	}

	if spec.Invoice != "" {
		db = db.Where("CONTRACTOR_INV_NO LIKE ?", "%"+spec.Invoice+"%")
	}

	if spec.ItemName != "" {
		db = db.Where("BOQ_ITEM_NAME LIKE ?", "%"+spec.ItemName+"%")
	}

	if spec.Arrival != "" {
		db = db.Where("ARRIVAL_DATE_AT_SITE LIKE ?", "%"+spec.Arrival+"%")
	}

	if spec.Inspection != "" {
		db = db.Where("INSPECTION_DATE LIKE ?", "%"+spec.Inspection+"%")
	}

	if spec.StateName != "" {
		db = db.Where("STATE_NAME LIKE ?", "%"+spec.StateName+"%")
	}

	return db, db.Error
}

func (spec *ProgressReportSearch) buildOrder(db *gorm.DB) (*gorm.DB, error) {

	// db = db.Where("BOQ_ID = ?", id).Order("ID DESC")
	if spec.SortSequencesNo != "" {
		db = db.Order("SEQ_NO " + spec.SortSequencesNo)
	}

	if spec.SortInvoice != "" {
		db = db.Order("CONTRACTOR_INV_NO " + spec.SortInvoice)
	}

	if spec.SortItemName != "" {
		db = db.Order("BOQ_ITEM_NAME " + spec.SortItemName)
	}

	if spec.SortArrival != "" {
		db = db.Order("ARRIVAL_DATE_AT_SITE " + spec.SortArrival)
	}

	if spec.SortInspection != "" {
		db = db.Order("INSPECTION_DATE " + spec.SortInspection)
	}

	if spec.SortStateName != "" {
		db = db.Order("STATE_NAME " + spec.SortStateName)
	}

	return db, db.Error
}

func GetCheckReport(db *connection.DBConnection) getCheckReportFunc {
	return func(ctx context.Context, spec CheckReportSearch, id uint) (CheckReports, int64, error) {
		var result CheckReports
		data, count, err := spec.search(db.CPM, &id)

		if err != nil {
			return result, count, err
		}
		result = data.ToResponse()
		return result, count, err
	}

}

func (spec *CheckReportSearch) search(db *gorm.DB, id *uint) (CheckReportDBs, int64, error) {
	var result CheckReportDBs

	countDB := db.Model(&result)
	countDB = countDB.Where("CONTRACT_ID = ?", id)
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

func (spec *CheckReportSearch) buildSearch(db *gorm.DB) (*gorm.DB, error) {

	if spec.SequencesNo != "" {
		db = db.Where("SEQ_NO LIKE ?", "%"+spec.SequencesNo+"%")
	}

	if spec.Invoice != "" {
		db = db.Where("CONTRACTOR_INV_NO LIKE ?", "%"+spec.Invoice+"%")
	}

	if spec.ItemName != "" {
		db = db.Where("BOQ_ITEM_NAME LIKE ?", "%"+spec.ItemName+"%")
	}

	if spec.Arrival != "" {
		db = db.Where("ARRIVAL_DATE_AT_SITE LIKE ?", "%"+spec.Arrival+"%")
	}

	if spec.Inspection != "" {
		db = db.Where("INSPECTION_DATE LIKE ?", "%"+spec.Inspection+"%")
	}

	if spec.Amount != "" {
		db = db.Where("DELIVERED_QTY LIKE ?", "%"+spec.Amount+"%")
	}

	if spec.Good != "" {
		db = db.Where("GOOD_QTY LIKE ?", "%"+spec.Good+"%")
	}

	if spec.Waste != "" {
		db = db.Where("WASTE_QTY LIKE ?", "%"+spec.Waste+"%")
	}

	return db, db.Error
}

func (spec *CheckReportSearch) buildOrder(db *gorm.DB) (*gorm.DB, error) {

	if spec.SortSequencesNo != "" {
		db = db.Order("SEQ_NO " + spec.SortSequencesNo)
	}

	if spec.SortInvoice != "" {
		db = db.Order("CONTRACTOR_INV_NO " + spec.SortInvoice)
	}

	if spec.SortItemName != "" {
		db = db.Order("BOQ_ITEM_NAME " + spec.SortItemName)
	}

	if spec.SortArrival != "" {
		db = db.Order("ARRIVAL_DATE_AT_SITE " + spec.SortArrival)
	}

	if spec.SortInspection != "" {
		db = db.Order("INSPECTION_DATE " + spec.SortInspection)
	}

	if spec.SortAmount != "" {
		db = db.Order("DELIVERED_QTY " + spec.SortAmount)
	}

	if spec.SortGood != "" {
		db = db.Order("GOOD_QTY " + spec.SortGood)
	}

	if spec.SortWaste != "" {
		db = db.Order("WASTE_QTY " + spec.SortWaste)
	}

	return db, db.Error
}
