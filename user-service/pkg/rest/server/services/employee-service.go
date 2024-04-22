package services

import (
	"github.com/MrAzharuddin/swagger-test/user-service/pkg/rest/server/daos"
	"github.com/MrAzharuddin/swagger-test/user-service/pkg/rest/server/models"
)

type EmployeeService struct {
	employeeDao *daos.EmployeeDao
}

func NewEmployeeService() (*EmployeeService, error) {
	employeeDao, err := daos.NewEmployeeDao()
	if err != nil {
		return nil, err
	}
	return &EmployeeService{
		employeeDao: employeeDao,
	}, nil
}

func (employeeService *EmployeeService) CreateEmployee(employee *models.Employee) (*models.Employee, error) {
	return employeeService.employeeDao.CreateEmployee(employee)
}

func (employeeService *EmployeeService) GetEmployee(id int64) (*models.Employee, error) {
	return employeeService.employeeDao.GetEmployee(id)
}

func (employeeService *EmployeeService) UpdateEmployee(id int64, employee *models.Employee) (*models.Employee, error) {
	return employeeService.employeeDao.UpdateEmployee(id, employee)
}

func (employeeService *EmployeeService) DeleteEmployee(id int64) error {
	return employeeService.employeeDao.DeleteEmployee(id)
}
