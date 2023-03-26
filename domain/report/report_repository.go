package report

import (
	"context"
	"cpm-rad-backend/domain/connection"
	"cpm-rad-backend/domain/minio"
	"cpm-rad-backend/domain/utils"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/inhies/go-bytesize"
	"gorm.io/gorm"
)

func Create(db *connection.DBConnection, m minio.Client) createFunc {
	return func(ctx context.Context, r Report, f File) (Report, error) {

		// getEmpFromDB := employee.GetByID(db)
		// getEmpByID := func(empID string) (employee.Employee, error) {
		// 	return getEmpFromDB(ctx, empID)
		// }

		err := db.CPM.Transaction(func(tx *gorm.DB) error {
			report := r.ToModel()

			report.RadNo = fmt.Sprintf("rad-%d", report.ItemID)
			// now := time.Now()
			// report.CreatedDate = &now
			report.CreateBy = "createdBy"
			// report.UpdatedDate = &now
			// report.UpdatedBy = createdBy

			// if err := tx.Omit("EmployeeJobs", "JobStatus").Create(&report).Error; err != nil {
			// 	return err
			// }

			if err := tx.Omit("UpdateBy", "UpdateDate").Create(&report).Error; err != nil {
				return err
			}

			for i, file := range f.Info {
				src, err := file.Open()
				if err != nil {
					return err
				}
				defer src.Close()

				info, objectName, err := m.Upload(ctx, file, fmt.Sprintf("%d/%d", report.ItemID, report.ID))
				if err != nil {
					return err
				}

				b := bytesize.New(float64(info.Size))
				displaySize := b.Format("%.2f ", "", false)
				words := strings.Fields(displaySize)

				f := AttachFile{
					ID:          0,
					Name:        file.Filename,
					ObjectName:  objectName,
					DisplaySize: displaySize,
					Size:        words[0],
					Unit:        words[1],
					Path:        info.Key,
					FileType:    strings.Replace(filepath.Ext(info.Key), ".", "", 1),
					DocType:     utils.StringToUint(f.Type[i]),
				}

				r.AttachFiles = append(r.AttachFiles, f)
				file := f.ToModel(report)
				if err := tx.Create(&file).Error; err != nil {
					//remove file minio
					return err
				}
			}

			// if err := updateJobEmployees(tx, formID, &req, createdBy, getEmpByID); err != nil {
			// 	return err
			// }

			// deptChangeCode := ""
			// if err := tx.Model(&employee_job.EmployeeJob{}).Joins("Employee").Joins("EmployeeRole").Where(
			// 	"CMDC_JOB_ID = ? AND EmployeeRole.ROLE_NAME_ENG = ?", formID, employee_role.SUPERVISOR,
			// ).Pluck("Employee.DEPT_CHANGE_CODE", &deptChangeCode).Error; err != nil {
			// 	return err
			// }

			return nil
		})

		return r, err
	}
}
