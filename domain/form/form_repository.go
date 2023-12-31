package form

import (
	"context"
	"cpm-rad-backend/domain/connection"
	"cpm-rad-backend/domain/minio"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/inhies/go-bytesize"
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
			filesAttach := req.FilesAttach

			radForm.RadNo = fmt.Sprintf("rad-%d", radForm.ItemID)
			// now := time.Now()
			// radForm.CreatedDate = &now
			radForm.CreateBy = createdBy
			// radForm.UpdatedDate = &now
			// radForm.UpdatedBy = createdBy

			// if err := tx.Omit("EmployeeJobs", "JobStatus").Create(&radForm).Error; err != nil {
			// 	return err
			// }
			if err := tx.Omit("UpdateBy", "UpdateDate").Create(&radForm).Error; err != nil {
				return err
			}
			formID = radForm.ID

			for _, fileAttach := range filesAttach {
				file := fileAttach.ToModel(formID, createdBy)
				if err := tx.Create(&file).Error; err != nil {
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

		return formID, err
	}
}

func Update(db *connection.DBConnection) updateFunc {
	return func(ctx context.Context, req UpdateRequest) error {

		// getEmpFromDB := employee.GetByID(db)
		// getEmpByID := func(empID string) (employee.Employee, error) {
		// 	return getEmpFromDB(ctx, empID)
		// }

		err := db.CPM.Transaction(func(tx *gorm.DB) error {
			radForm := req.ToModel()
			radForm.UpdateBy = "คนแก้ เอกสาร"
			now := time.Now()
			radForm.UpdateDate = &now

			// if err := tx.Omit("EmployeeJobs", "JobStatus").Create(&radForm).Error; err != nil {
			// 	return err
			// }
			// if err := tx.Where("ID = ?", req.ID).Updates(&radForm).Error; err != nil {
			if err := tx.Updates(&radForm).Error; err != nil {
				return err
			}
			// log.Info(fmt.Sprintf("%v", req.File))
			files := req.File
			for _, f := range files {
				if f.ID > 0 {
					file := FileUpdate{
						ID:         f.ID,
						DocType:    f.Type,
						UpdateBy:   radForm.UpdateBy,
						UpdateDate: radForm.UpdateDate,
						DelFlag:    "",
					}

					if err := tx.Omit("DEL_FLAG").Updates(&file).Error; err != nil {
						return err
					}
					continue
				}
				file := f.ToModel(req.ID, radForm.UpdateBy)
				if err := tx.Create(&file).Error; err != nil {
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

		return err
	}
}

func Delete(db *connection.DBConnection) dateteFunc {
	return func(ctx context.Context, id uint) error {

		// getEmpFromDB := employee.GetByID(db)
		// getEmpByID := func(empID string) (employee.Employee, error) {
		// 	return getEmpFromDB(ctx, empID)
		// }
		var delForm deleteForm
		err := db.CPM.Transaction(func(tx *gorm.DB) error {
			delForm.ID = id
			delForm.UpdateBy = "คนลบ เอกสาร"
			now := time.Now()
			delForm.UpdateDate = &now
			delForm.DelFlag = "Y"

			form := delForm.ToModel()
			if err := tx.Updates(&form).Error; err != nil {
				return err
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

		return err
	}
}

func Get(db *connection.DBConnection) getFunc {
	return func(ctx context.Context, id uint) (Response, error) {
		var res Response

		var result FormDB
		cpm := db.CPM.Model(&result)
		err := cpm.Table("CPM.GetReportDetail(?)", id).
			Scan(&result).
			Error
		if err != nil {
			return res, err
		}

		var files AttachFiles
		cpm = db.CPM.Model(&files)
		err = cpm.Table("CPM.GetFileAttachment(?)", result.ID).Order("TYPE_ID").
			Scan(&files).
			Error
		if err != nil {
			return res, err
		}

		res = result.ToResponse(files)
		return res, err
	}
}

// func View(db *connection.DBConnection) viewFunc {
// 	return func(ctx context.Context, id uint) (ResponseView, error) {
// 		var res ResponseView

// 		var result FormViewDB
// 		cpm := db.CPM.Model(&result)
// 		err := cpm.Table("CPM.GetReportDetail(?)", id).
// 			Scan(&result).
// 			Error
// 		if err != nil {
// 			return res, err
// 		}

// 		var files ViewFilesDB
// 		cpm = db.CPM.Model(&files)
// 		err = cpm.Table("CPM.GetFileAttachment(?)", result.ID).Order("TYPE_ID").
// 			Scan(&files).
// 			Error
// 		if err != nil {
// 			return res, err
// 		}

// 		var attachFiles AttachFiles
// 		var viewFiles ViewFiles
// 		var temp uint = 0
// 		count := len(files)
// 		var typeName string

// 		for i, f := range files {
// 			fmt.Println(i)
// 			if count == 1 {
// 				attachFiles = append(attachFiles, AttachFile{
// 					TypeName: f.TypeName,
// 					Files: []ViewFile{{
// 						ID:   f.ID,
// 						Name: f.Name,
// 						Size: f.Size,
// 						Unit: f.Unit,
// 					}},
// 				})
// 			} else if count-1 == i {
// 				if f.TypeID == temp {
// 					viewFiles = append(viewFiles, ViewFile{
// 						ID:   f.ID,
// 						Name: f.Name,
// 						Size: f.Size,
// 						Unit: f.Unit,
// 					})
// 					attachFiles = append(attachFiles, AttachFile{
// 						TypeName: f.TypeName,
// 						Files:    viewFiles,
// 					})
// 				} else {
// 					if len(viewFiles) > 0 {
// 						attachFiles = append(attachFiles, AttachFile{
// 							TypeName: typeName,
// 							Files:    viewFiles,
// 						})
// 					}
// 					attachFiles = append(attachFiles, AttachFile{
// 						TypeName: f.TypeName,
// 						Files: []ViewFile{{
// 							ID:   f.ID,
// 							Name: f.Name,
// 							Size: f.Size,
// 							Unit: f.Unit,
// 						}},
// 					})
// 				}
// 			} else {
// 				if temp == 0 || f.TypeID == temp {
// 					typeName = f.TypeName
// 					viewFiles = append(viewFiles, ViewFile{
// 						ID:   f.ID,
// 						Name: f.Name,
// 						Size: f.Size,
// 						Unit: f.Unit,
// 					})
// 				} else {
// 					attachFiles = append(attachFiles, AttachFile{
// 						TypeName: typeName,
// 						Files:    viewFiles,
// 					})
// 					viewFiles = nil
// 					typeName = f.TypeName
// 					viewFiles = append(viewFiles, ViewFile{
// 						ID:   f.ID,
// 						Name: f.Name,
// 						Size: f.Size,
// 						Unit: f.Unit,
// 					})
// 				}
// 				temp = f.TypeID
// 			}
// 		}

// 		res = result.ToResponse(attachFiles)
// 		return res, err
// 	}
// }

func GetCountry(db *connection.DBConnection) getCountryFunc {
	return func(ctx context.Context, filter string) (Countrys, error) {
		var result Countrys
		cpm := db.CPM.Model(&result)
		err := cpm.Table("CPM.GetCountry(?)", filter).
			Scan(&result).
			Error
		if err != nil {
			return result, err
		}

		return result, err
	}
}

func FileUpload(db *connection.DBConnection, m minio.Client) fileUploadFunc {
	return func(ctx context.Context, file *multipart.FileHeader, itemID int) (FileUploadResponse, error) {
		// var result FileUploadResponse
		info, objectName, err := m.Upload(ctx, file, string(rune(itemID)))

		b := bytesize.New(float64(info.Size))
		displaySize := b.Format("%.2f ", "", false)
		words := strings.Fields(displaySize)

		result := FileUploadResponse{
			Name:        file.Filename,
			ObjectName:  objectName,
			DisplaySize: displaySize,
			Size:        words[0],
			Unit:        words[1],
			FileType:    strings.Replace(filepath.Ext(info.Key), ".", "", 1),
			FilePath:    info.Key,
		}
		return result, err
	}
}

func GetDocType(db *connection.DBConnection) getDocTypeFunc {
	return func(ctx context.Context) (DocTypes, error) {
		var result DocTypes
		cpm := db.CPM.Model(&result)
		err := cpm.Select("ID,DESCRIPTION").Scan(&result).Error
		if err != nil {
			return result, err
		}

		return result, err
	}
}

func FileDelete(db *connection.DBConnection, m minio.Client) fileDeleteFunc {
	return func(ctx context.Context, itemID int, objectName string) error {
		// var result FileUploadResponse
		err := m.Delete(ctx, objectName, uint(itemID))

		return err
	}
}

func FileDownload(db *connection.DBConnection, m minio.Client) fileDownloadFunc {
	return func(ctx context.Context, fileID uint) (FileResponse, error) {
		var result File
		var fileResponse FileResponse
		cpm := db.CPM.Model(&result)
		err := cpm.Where("ID = ?", fileID).Scan(&result).Error
		if err != nil {
			return fileResponse, err
		}

		obj, ext, err := m.Download(ctx, result.Path)
		fileResponse = FileResponse{
			Obj:  obj,
			Ext:  ext,
			Name: result.Name,
		}
		return fileResponse, err
	}
}
