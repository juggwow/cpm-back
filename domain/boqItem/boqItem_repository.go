package boqItem

import (
	"context"
	"cpm-rad-backend/domain/connection"

	"gorm.io/gorm"
)

func Get(db *connection.DBConnection) getFunc {
	return func(ctx context.Context, spec SearchSpec, id uint) (BoqItemLists, int64, error) {

		// var result FormDB
		// cpm := db.CPM.Model(&result)
		// err := cpm.Table("CPM.GetReportDetail(?)", id).
		// 	Scan(&result).
		// 	Error
		// if err != nil {
		// 	return res, err
		// }

		data, count, err := spec.searchBoQItems(db.CPM, &id)
		if err != nil {
			return data, count, err
		}

		return data, count, err
	}
}

func (spec *SearchSpec) searchBoQItems(db *gorm.DB, id *uint) (BoqItemLists, int64, error) {
	var result BoqItemListDBs
	var count int64

	findcount := db.Table("CPM.GetBoqItems(?)", id)
	//findDB = findDB.Preload("EmployeeJobs.Employee").Preload("EmployeeJobs.EmployeeRole")
	findcount = spec.buildSearch(findcount)
	if err := findcount.Distinct("ID").Count(&count).Error; err != nil {
		return result.ToResponse(), 0, err
	}

	findDB := db.Model(&result).Table("CPM.GetBoqItems(?)", id)
	findDB = spec.buildSearch(findDB)
	findDB = spec.buildOrder(findDB)
	findDB = findDB.Offset(spec.Offset()).Limit(spec.GetLimit())
	if err := findDB.Find(&result).Error; err != nil {
		return result.ToResponse(), 0, err
	}

	return result.ToResponse(), count, nil
}

func (spec *SearchSpec) buildSearch(db *gorm.DB) *gorm.DB {

	// db = db.Table("(?) as u", subQuery)

	if spec.SearchRowNo != "" {
		db = db.Where("ROW_NO LIKE ?", "%"+spec.SearchRowNo+"%")
	}

	if spec.SearchNumber != "" {
		db = db.Where("NUMBER LIKE ?", "%"+spec.SearchNumber+"%")
	}

	if spec.SearchGroupName != "" {
		db = db.Where("GROUP_NAME LIKE ?", "%"+spec.SearchGroupName+"%")
	}

	if spec.SearchName != "" {
		db = db.Where("NAME LIKE ?", "%"+spec.SearchName+"%")
	}

	if spec.SearchQuantity != "" {
		db = db.Where("QUANTITY LIKE ?", "%"+spec.SearchQuantity+"%")
	}

	if spec.SearchDeliveryQty != "" {
		db = db.Where("DELIVERY_QUANTITY LIKE ?", "%"+spec.SearchDeliveryQty+"%")
	}

	if spec.SearchReceiveQty != "" {
		db = db.Where("RECEIVE_QUANTITY LIKE ?", "%"+spec.SearchReceiveQty+"%")
	}

	if spec.SearchDamageQty != "" {
		db = db.Where("DAMAGE_QUANTITY LIKE ?", "%"+spec.SearchReceiveQty+"%")
	}

	return db
}

func (spec *SearchSpec) buildOrder(db *gorm.DB) *gorm.DB {
	if spec.SortRowNo != "" {
		db.Order("ROW_NO " + spec.SortRowNo)
	}

	if spec.SortNumber != "" {
		db.Order("NUMBER " + spec.SortNumber)
	}

	if spec.SortGroupName != "" {
		db.Order("GROUP_NAME " + spec.SortGroupName)
	}

	if spec.SortName != "" {
		db.Order("NAME " + spec.SortName)
	}

	if spec.SortQuantity != "" {
		db.Order("QUANTITY " + spec.SortQuantity)
	}

	if spec.SortDeliveryQty != "" {
		db.Order("DELIVERY_QUANTITY " + spec.SortDeliveryQty)
	}

	if spec.SortReceiveQty != "" {
		db.Order("RECEIVE_QUANTITY " + spec.SortReceiveQty)
	}

	if spec.SortDamageQty != "" {
		db.Order("DAMAGE_QUANTITY " + spec.SortDamageQty)
	}

	return db
}
