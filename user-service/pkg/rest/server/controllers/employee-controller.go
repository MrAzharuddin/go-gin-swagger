package controllers

import (
	"errors"
	"github.com/MrAzharuddin/swagger-test/user-service/pkg/rest/server/daos/clients/sqls"
	"github.com/MrAzharuddin/swagger-test/user-service/pkg/rest/server/models"
	"github.com/MrAzharuddin/swagger-test/user-service/pkg/rest/server/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"os"
	"strconv"
)

type EmployeeController struct {
	employeeService *services.EmployeeService
}

func NewEmployeeController() (*EmployeeController, error) {
	employeeService, err := services.NewEmployeeService()
	if err != nil {
		return nil, err
	}
	return &EmployeeController{
		employeeService: employeeService,
	}, nil
}

// CreateEmployee creates a new employee for the employee service
//	@Summary		Creates a new employee
//	@Description	Creates a new employee
//	@Tags			employees
//	@Accept			json
//	@Produce		json
//	@Param			employee	body		models.Employee	true	"Create employee"
//	@Success		201			{object}	models.Employee
//	@Failure		422			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/employees [post]
func (employeeController *EmployeeController) CreateEmployee(context *gin.Context) {
	// validate input
	var input models.Employee
	if err := context.ShouldBindJSON(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// trigger employee creation
	employeeCreated, err := employeeController.employeeService.CreateEmployee(&input)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, employeeCreated)
}

// FetchEmployee fetches a single employee for the employee service
//	@Summary		Fetches a single employee
//	@Description	Fetches a single employee
//	@Tags			employees
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"id"
//	@Success		200	{object}	models.Employee
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/employees/{id} [get]
func (employeeController *EmployeeController) FetchEmployee(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// trigger employee fetching
	employee, err := employeeController.employeeService.GetEmployee(id)
	if err != nil {
		log.Error(err)
		if errors.Is(err, sqls.ErrNotExists) {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	serviceName := os.Getenv("SERVICE_NAME")
	collectorURL := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if len(serviceName) > 0 && len(collectorURL) > 0 {
		// get the current span by the request context
		currentSpan := trace.SpanFromContext(context.Request.Context())
		currentSpan.SetAttributes(attribute.String("employee.id", strconv.FormatInt(employee.Id, 10)))
	}

	context.JSON(http.StatusOK, employee)
}

// UpdateEmployee updates a single employee for the employee service
//	@Summary		Updates a single employee
//	@Description	Updates a single employee
//	@Tags			employees
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int				true	"id"
//	@Param			employee	body		models.Employee	true	"Update employee"
//	@Success		204			{object}	interface{}
//	@Failure		404			{object}	ErrorResponse
//	@Failure		422			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/employees/{id} [put]
func (employeeController *EmployeeController) UpdateEmployee(context *gin.Context) {
	// validate input
	var input models.Employee
	if err := context.ShouldBindJSON(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// trigger employee update
	if _, err := employeeController.employeeService.UpdateEmployee(id, &input); err != nil {
		log.Error(err)
		if errors.Is(err, sqls.ErrNotExists) {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusNoContent, gin.H{})
}

// DeleteEmployee deletes a single employee for the employee service
//	@Summary		Deletes a single employee
//	@Description	Deletes a single employee
//	@Tags			employees
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"id"
//	@Success		204	{object}	interface{}
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/employees/{id} [delete]
func (employeeController *EmployeeController) DeleteEmployee(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// trigger employee deletion
	if err := employeeController.employeeService.DeleteEmployee(id); err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusNoContent, gin.H{})
}
