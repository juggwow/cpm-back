package form

import (
	"context"
	"cpm-rad-backend/domain/connection"
	"fmt"

	"gorm.io/gorm"
)

func Create(db *connection.DBConnection) createFunc {
	return func(ctx context.Context, req Request, createdBy string) (uint, error) {
		var formID uint

		// getEmpFromDB := employee.GetByID(db)
		// getEmpByID := func(empID string) (employee.Employee, error) {
		// 	return getEmpFromDB(ctx, empID)
		// }

		err := db.CPM.Transaction(func(tx *gorm.DB) error {
			radForm := req.ToModel()
			radForm.RadNo = fmt.Sprintf("rad-%d", radForm.ItemID)
			// now := time.Now()
			// radForm.CreatedDate = &now
			radForm.CreateBy = createdBy
			// radForm.UpdatedDate = &now
			// radForm.UpdatedBy = createdBy

			// if err := tx.Omit("EmployeeJobs", "JobStatus").Create(&radForm).Error; err != nil {
			// 	return err
			// }
			err := tx.Create(&radForm).Error
			formID = radForm.ID

			return err

			// if err := updateJobEmployees(tx, formID, &req, createdBy, getEmpByID); err != nil {
			// 	return err
			// }

			// deptChangeCode := ""
			// if err := tx.Model(&employee_job.EmployeeJob{}).Joins("Employee").Joins("EmployeeRole").Where(
			// 	"CMDC_JOB_ID = ? AND EmployeeRole.ROLE_NAME_ENG = ?", formID, employee_role.SUPERVISOR,
			// ).Pluck("Employee.DEPT_CHANGE_CODE", &deptChangeCode).Error; err != nil {
			// 	return err
			// }
		})

		return formID, err
	}
}
