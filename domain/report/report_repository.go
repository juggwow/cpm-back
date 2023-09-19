package report

import (
	"context"
	"cpm-rad-backend/domain/connection"
	"cpm-rad-backend/domain/minio"
	"cpm-rad-backend/domain/utils"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/inhies/go-bytesize"
	"gorm.io/gorm"
)

func Create(db *connection.DBConnection, m minio.Client) createFunc {
	return func(ctx context.Context, r RequestReportCreate, f File) (ResponseReportDetail, error) {
		var res ResponseReportDetail
		var attachFiles ResponseAttachFiles
		// getEmpFromDB := employee.GetByID(db)
		// getEmpByID := func(empID string) (employee.Employee, error) {
		// 	return getEmpFromDB(ctx, empID)
		// }
		err := db.CPM.Transaction(func(tx *gorm.DB) error {
			report := r.ToModel()

			// report.RadNo = fmt.Sprintf("rad-%d", report.ItemID)
			// report.CreateBy = "createdBy"

			if err := tx.Omit("UpdateBy", "UpdateDate", "DelFlag", "RadNo").Create(&report).Error; err != nil {
				return err
			}
			// res.ID = report.ID

			for i, file := range f.Info {
				src, err := file.Open()
				if err != nil {
					return err
				}
				defer src.Close()

				fileDetail, shouldReturn, returnValue := uploadFileToMinio(m, ctx, file, report, f.Type[i])
				if shouldReturn {
					return returnValue
				}

				file := fileDetail.ToModel(report)
				if err := tx.Omit("UpdateBy", "UpdateDate", "DelFlag").Create(&file).Error; err != nil {
					//remove file minio
					return err
				}
				resFile := ResponseAttachFile{
					ID:     file.ID,
					Name:   file.Name,
					Size:   file.Size.String(),
					Unit:   file.Unit,
					TypeID: file.DocType,
				}
				attachFiles = append(attachFiles, resFile)
			}

			res = report.ToResponse(attachFiles)

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

		return res, err
	}
}

func Get(db *connection.DBConnection) getFunc {
	return func(ctx context.Context, id uint) (ResponseReportDetail, error) {
		var res ResponseReportDetail
		var result ReportDetailDB

		cpm := db.CPM.Model(&result)
		err := cpm.Table("CPM.GET_REPORT_DETAIL(?)", id).
			Scan(&result).
			Error
		if err != nil {
			return res, err
		}

		var files ResponseAttachFiles
		cpm = db.CPM.Model(&files)
		err = cpm.Table("CPM.GetFileAttachment(?)", result.ID).Order("TYPE_ID").
			Scan(&files).
			Error
		if err != nil {
			return res, err
		}

		var images DbAttachImages
		err = db.CPM.Select("Path").Where("RAD_ID = ? AND DEL_FLAG = ?", result.ID, "N").Find(&images).Error
		if err != nil {
			return res, err
		}

		res = result.ToResponse(files, images)
		return res, err
	}
}

func Update(db *connection.DBConnection, m minio.Client) updateFunc {
	return func(ctx context.Context, r RequestReportUpdate, f File) (ResponseReportDetail, error) {
		var res ResponseReportDetail
		var attachFiles ResponseAttachFiles
		// getEmpFromDB := employee.GetByID(db)
		// getEmpByID := func(empID string) (employee.Employee, error) {
		// 	return getEmpFromDB(ctx, empID)
		// }

		err := db.CPM.Transaction(func(tx *gorm.DB) error {
			report := r.ToModel()
			// report.UpdateBy = "คนแก้ เอกสาร"
			now := time.Now()
			report.UpdateDate = &now

			if err := tx.Omit("ItemID", "CreateBy", "DelFlag", "RadNo").Updates(&report).Error; err != nil {
				return err
			}

			for i, file := range f.Info {
				src, err := file.Open()
				if err != nil {
					return err
				}
				defer src.Close()

				fileDetail, shouldReturn, returnValue := uploadFileToMinio(m, ctx, file, report, f.Type[i])
				if shouldReturn {
					return returnValue
				}
				report.CreateBy = "createdBy"
				file := fileDetail.ToModel(report)
				if err := tx.Omit("UpdateBy", "UpdateDate", "DelFlag").Create(&file).Error; err != nil {
					//remove file minio
					return err
				}
				resFile := ResponseAttachFile{
					ID:     file.ID,
					Name:   file.Name,
					Size:   file.Size.String(),
					Unit:   file.Unit,
					TypeID: file.DocType,
				}
				attachFiles = append(attachFiles, resFile)
			}

			for _, updateFile := range f.Update {
				file := AttachFileDB{
					ID:         updateFile.FileID,
					DocType:    updateFile.TypeID,
					UpdateBy:   report.UpdateBy,
					UpdateDate: report.UpdateDate,
				}

				if err := tx.Select("DocType", "UpdateBy", "UpdateDate").Updates(&file).Error; err != nil {
					return err
				}
			}

			for _, delFile := range f.Delete {
				file := AttachFileDB{
					ID:         utils.StringToUint(delFile),
					UpdateBy:   report.UpdateBy,
					UpdateDate: report.UpdateDate,
					DelFlag:    "Y",
				}

				if err := tx.Select("DelFlag", "UpdateBy", "UpdateDate").Updates(&file).Error; err != nil {
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
			res = report.ToResponse(attachFiles)
			return nil
		})

		return res, err
	}
}

func UpdateBasicDetails(db *connection.DBConnection) updateBasicDetailsFunc {
	return func(ctx context.Context, req RequestReportUpdate) (ResponseReportDetail, error) {
		var res ResponseReportDetail
		data := req.ToModel()
		if err := db.CPM.Select("Arrival", "Inspection", "TaskMaster").Updates(&data).Error; err != nil {
			return res, err
		}
		res = data.ToResponse(ResponseAttachFiles{})
		return res, nil
	}
}

func uploadFileToMinio(m minio.Client, ctx context.Context, file *multipart.FileHeader, report ReportDB, docType string) (FileDetail, bool, error) {
	info, objectName, err := m.Upload(ctx, file, fmt.Sprintf("%d/%d", report.ItemID, report.ID))
	if err != nil {
		return FileDetail{}, true, err
	}

	b := bytesize.New(float64(info.Size))
	displaySize := b.Format("%.2f ", "", false)
	words := strings.Fields(displaySize)

	attachFile := FileDetail{
		Name:        file.Filename,
		ObjectName:  objectName,
		DisplaySize: displaySize,
		Size:        words[0],
		Unit:        words[1],
		Path:        info.Key,
		FileType:    strings.Replace(filepath.Ext(info.Key), ".", "", 1),
		DocType:     utils.StringToUint(docType),
	}
	return attachFile, false, nil
}
