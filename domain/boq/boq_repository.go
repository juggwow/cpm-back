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

	// countDB := db.Model(&Items)
	// countDB, err := spec.buildSearchCriteria(countDB, id)
	// if err != nil {
	// 	return Items, 0, err
	// }
	var count int64
	// if err := countDB.Count(&count).Error; err != nil {
	// 	return Items, 0, err
	// }

	subQuery := db.Model(&Items).Select("ROW_NUMBER() OVER(ORDER BY ID ASC) AS SEQUENCES_NO,ID,ITEM,NAME,GROUPNAME,QUANTITY,UNIT,DELIVERY_QUANTITY,RECEIVE_QUANTITY,DAMAGE_QUANTITY").
		Where("WORK_CONTRACT_ID = ?", id)

	// findDB := db.Model(&Items)
	//findDB = findDB.Preload("EmployeeJobs.Employee").Preload("EmployeeJobs.EmployeeRole")
	findDB, err := spec.buildSearch(db, subQuery, id)
	if err := findDB.Count(&count).Error; err != nil {
		return Items, 0, err
	}
	findDB = spec.buildOrder(findDB)
	findDB = findDB.Offset(spec.Offset()).Limit(spec.GetLimit())
	// findDB, err = spec.buildSearchCriteria(findDB, id)
	if err != nil {
		return Items, 0, err
	}

	return Items, count, findDB.Find(&Items).Error
}

func (spec *ItemSearchSpec) buildOrder(db *gorm.DB) *gorm.DB {
	if spec.SortSequencesNo != "" {
		db.Order("SEQUENCES_NO " + spec.SortSequencesNo)
	}

	if spec.SortItemNo != "" {
		db.Order("ITEM " + spec.SortItemNo)
	}

	if spec.SortItemName != "" {
		db.Order("NAME " + spec.SortItemName)
	}

	if spec.SortItemGroup != "" {
		db.Order("GROUPNAME " + spec.SortItemGroup)
	}

	if spec.SortItemQuantity != "" {
		db.Order("QUANTITY " + spec.SortItemQuantity)
	}

	if spec.SortItemDelivery != "" {
		db.Order("DELIVERY_QUANTITY " + spec.SortItemDelivery)
	}

	if spec.SortItemReceive != "" {
		db.Order("RECEIVE_QUANTITY " + spec.SortItemReceive)
	}

	if spec.SortItemDamage != "" {
		db.Order("DAMAGE_QUANTITY " + spec.SortItemDamage)
	}

	return db
}

func (spec *ItemSearchSpec) buildSearch(db *gorm.DB, subQuery *gorm.DB, id *uint) (*gorm.DB, error) {

	db = db.Table("(?) as u", subQuery)

	if spec.SequencesNo != "" {
		db = db.Where("SEQUENCES_NO LIKE ?", spec.SequencesNo+"%")
	}

	if spec.ItemNo != "" {
		db = db.Where("ITEM LIKE ?", spec.ItemNo+"%")
	}

	if spec.ItemName != "" {
		db = db.Where("NAME LIKE ?", spec.ItemName+"%")
	}

	if spec.ItemGroup != "" {
		db = db.Where("GROUPNAME LIKE ?", spec.ItemGroup+"%")
	}

	if spec.ItemQuantity != "" {
		db = db.Where("QUANTITY LIKE ?", spec.ItemQuantity+"%")
	}

	if spec.ItemDelivery != "" {
		db = db.Where("DELIVERY_QUANTITY LIKE ?", spec.ItemDelivery+"%")
	}

	if spec.ItemReceive != "" {
		db = db.Where("RECEIVE_QUANTITY LIKE ?", spec.ItemReceive+"%")
	}

	if spec.ItemDamage != "" {
		db = db.Where("DAMAGE_QUANTITY LIKE ?", spec.ItemDamage+"%")
	}

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
