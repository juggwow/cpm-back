package employee

import (
	"context"
	"cpm-rad-backend/domain/connection"
	"cpm-rad-backend/domain/user_profile"
	"errors"
	"time"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

func Create(db *connection.DBConnection) func(context.Context, Employee) (uint, error) {
	return func(ctx context.Context, employee Employee) (uint, error) {
		err := db.CPM.Create(&employee).Error
		return employee.ID, err
	}
}

func GetByID(db *connection.DBConnection) func(context.Context, string) (Employee, error) {
	return func(ctx context.Context, employeeID string) (Employee, error) {
		employee := &Employee{}
		profile := &user_profile.UserProfile{}
		err := db.CPM.First(profile, "EMP_ID LIKE ?", "%"+employeeID+"%").Error
		if err != nil {
			// TODO: remove ignore error when EMPLOYEE ID not found in USER_PROFILE
			// employee.EmployeeID = employeeID
			// employee.BA = "DEV_EMPLOYEE_BA"
			// employee.DeptChangeCode = "50550002000"
			// employee.FirstName = "DEV_EMPLOYEE_FN_" + employeeID
			// employee.LastName = "DEV_EMPLOYEE_LN_" + employeeID
			// employee.Position = "DEV_EMPLOYEE_POS_" + employeeID
			// employee.Title = "Mr."
			// return *employee, nil
			return *employee, errors.New("Employee not found")
		}

		employee.Parse(profile)
		return *employee, nil
	}
}

func UpdateByID(db *connection.DBConnection) func(context.Context, uint, Employee) (uint, error) {
	return func(ctx context.Context, ID uint, employee Employee) (uint, error) {
		err := db.CPM.Model(&Employee{}).Where("ID = ?", ID).Updates(employee).Error
		return ID, err
	}
}

func DeleteByID(db *connection.DBConnection) func(context.Context, uint) error {
	return func(ctx context.Context, ID uint) error {
		result := db.CPM.Delete(&Employee{}, ID)
		if result.RowsAffected > 0 {
			return nil
		}
		return errors.New("error: deteted")
	}
}

func GetAndCreateIfNotExist(db *connection.DBConnection) func(context.Context, Employee) (Employee, error) {
	return func(ctx context.Context, employee Employee) (Employee, error) {
		getEmpByID := func(empID string) (Employee, error) {
			return GetByID(db)(ctx, empID)
		}

		emps := Employees{employee}

		if err := emps.CreateIfNotExist(db.CPM, employee.EmployeeID, getEmpByID); err != nil {
			return Employee{}, err
		}

		return emps[0], nil
	}
}

func (employees Employees) CreateIfNotExist(
	db *gorm.DB,
	createdBy string,
	getUserByID func(empID string) (Employee, error),
) error {
	var employeeIDs []string
	for i := range employees {
		emp := &employees[i]

		employeeIDs = append(employeeIDs, emp.EmployeeID)

		profile, err := getUserByID(emp.EmployeeID)
		if err != nil {
			return err
		}
		employees[i] = profile
	}

	var oldEmployees []Employee
	db.Model(&Employee{}).Find(&oldEmployees, "EMPLOYEE_ID IN ?", employeeIDs)

	if len(oldEmployees) == 0 {
		return employees.Create(db, createdBy)
	}

	now := time.Now()
	for i := range employees {
		emp := &employees[i]

		oldEmp, found := lo.Find(oldEmployees, func(old Employee) bool {
			return old.EmployeeID == employees[i].EmployeeID
		})

		if found {
			employees[i].ID = oldEmp.ID
			if isEmployeeChange(*emp, oldEmp) {
				emp.UpdatedDate = &now
				emp.UpdatedBy = createdBy
				db.Where("ID = ?", oldEmp.ID).Select(
					"FirstName",
					"LastName",
					"BA",
					"DeptChangeCode",
					"Position",
					"Title",
					"UpdatedDate",
					"UpdatedBy",
				).Updates(emp)
			}
			continue
		}

		emp.CreatedDate = &now
		emp.CreatedBy = createdBy
		emp.UpdatedDate = &now
		emp.UpdatedBy = createdBy
		err := db.Create(emp).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (employees Employees) Create(db *gorm.DB, createdBy string) error {
	if len(employees) == 0 {
		return nil
	}
	now := time.Now()
	for i := range employees {
		employees[i].CreatedDate = &now
		employees[i].CreatedBy = createdBy
		employees[i].UpdatedDate = &now
		employees[i].UpdatedBy = createdBy
	}
	if err := db.Create(&employees).Error; err != nil {
		return err
	}
	return nil
}

func isEmployeeChange(profile, oldEmp Employee) bool {
	return profile.BA != oldEmp.BA ||
		profile.DeptChangeCode != oldEmp.DeptChangeCode ||
		profile.FirstName != oldEmp.FirstName ||
		profile.LastName != oldEmp.LastName ||
		profile.Position != oldEmp.Position ||
		profile.Title != oldEmp.Title
}
