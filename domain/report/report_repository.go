package report

import (
	"context"
	"cpm-rad-backend/domain/connection"
)

func Create(db *connection.DBConnection) createFunc {
	return func(ctx context.Context, r Report, createdBy string) (uint, error) {
		var formID uint

		// getEmpFromDB := employee.GetByID(db)
		// getEmpByID := func(empID string) (employee.Employee, error) {
		// 	return getEmpFromDB(ctx, empID)
		// }

		// err := db.CPM.Transaction(func(tx *gorm.DB) error {
		// 	radForm := r.ToModel()
		// 	filesAttach := r.FilesAttach

		// 	radForm.RadNo = fmt.Sprintf("rad-%d", radForm.ItemID)
		// 	// now := time.Now()
		// 	// radForm.CreatedDate = &now
		// 	radForm.CreateBy = createdBy
		// 	// radForm.UpdatedDate = &now
		// 	// radForm.UpdatedBy = createdBy

		// 	// if err := tx.Omit("EmployeeJobs", "JobStatus").Create(&radForm).Error; err != nil {
		// 	// 	return err
		// 	// }
		// 	if err := tx.Omit("UpdateBy", "UpdateDate").Create(&radForm).Error; err != nil {
		// 		return err
		// 	}
		// 	formID = radForm.ID

		// 	for _, fileAttach := range filesAttach {
		// 		file := fileAttach.ToModel(formID, createdBy)
		// 		if err := tx.Create(&file).Error; err != nil {
		// 			return err
		// 		}
		// 	}

		// 	// if err := updateJobEmployees(tx, formID, &req, createdBy, getEmpByID); err != nil {
		// 	// 	return err
		// 	// }

		// 	// deptChangeCode := ""
		// 	// if err := tx.Model(&employee_job.EmployeeJob{}).Joins("Employee").Joins("EmployeeRole").Where(
		// 	// 	"CMDC_JOB_ID = ? AND EmployeeRole.ROLE_NAME_ENG = ?", formID, employee_role.SUPERVISOR,
		// 	// ).Pluck("Employee.DEPT_CHANGE_CODE", &deptChangeCode).Error; err != nil {
		// 	// 	return err
		// 	// }
		// 	return nil
		// })

		return formID, nil
	}
}
