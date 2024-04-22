package daos

import (
	"errors"
	"github.com/MrAzharuddin/swagger-test/user-service/pkg/rest/server/models"
)

var employees = make(map[int64]*models.Employee)

type EmployeeDao struct {
}

func NewEmployeeDao() (*EmployeeDao, error) {
	return &EmployeeDao{}, nil
}

func (employeeDao *EmployeeDao) CreateEmployee(employee *models.Employee) (*models.Employee, error) {
	employees[employee.Id] = employee

	return employee, nil
}

func (employeeDao *EmployeeDao) GetEmployee(id int64) (*models.Employee, error) {
	if employee, ok := employees[id]; ok {
		return employee, nil
	}

	return &models.Employee{}, errors.New("employee not found")
}

func (employeeDao *EmployeeDao) UpdateEmployee(id int64, employee *models.Employee) (*models.Employee, error) {
	if id != employee.Id {
		return nil, errors.New("id and payload don't match")
	}
	employees[employee.Id] = employee

	return employee, nil
}

func (employeeDao *EmployeeDao) DeleteEmployee(id int64) error {
	if _, ok := employees[id]; ok {
		delete(employees, id)
		return nil
	}

	return errors.New("employee not found")
}
