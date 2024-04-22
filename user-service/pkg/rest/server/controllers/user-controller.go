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

type UserController struct {
	userService *services.UserService
}

func NewUserController() (*UserController, error) {
	userService, err := services.NewUserService()
	if err != nil {
		return nil, err
	}
	return &UserController{
		userService: userService,
	}, nil
}

// CreateUser creates a new user for the user service
//	@Summary		Creates a new user
//	@Description	Creates a new user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.User	true	"Create user"
//	@Success		201		{object}	models.User
//	@Failure		422		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/users [post]
func (userController *UserController) CreateUser(context *gin.Context) {
	// validate input
	var input models.User
	if err := context.ShouldBindJSON(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// trigger user creation
	userCreated, err := userController.userService.CreateUser(&input)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, userCreated)
}

// FetchUser fetches a single user for the user service
//	@Summary		Fetches a single user
//	@Description	Fetches a single user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"id"
//	@Success		200	{object}	models.User
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/users/{id} [get]
func (userController *UserController) FetchUser(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// trigger user fetching
	user, err := userController.userService.GetUser(id)
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
		currentSpan.SetAttributes(attribute.String("user.id", strconv.FormatInt(user.Id, 10)))
	}

	context.JSON(http.StatusOK, user)
}

// UpdateUser updates a single user for the user service
//	@Summary		Updates a single user
//	@Description	Updates a single user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int			true	"id"
//	@Param			user	body		models.User	true	"Update user"
//	@Success		204		{object}	interface{}
//	@Failure		404		{object}	ErrorResponse
//	@Failure		422		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/users/{id} [put]
func (userController *UserController) UpdateUser(context *gin.Context) {
	// validate input
	var input models.User
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

	// trigger user update
	if _, err := userController.userService.UpdateUser(id, &input); err != nil {
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

// DeleteUser deletes a single user for the user service
//	@Summary		Deletes a single user
//	@Description	Deletes a single user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"id"
//	@Success		204	{object}	interface{}
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/users/{id} [delete]
func (userController *UserController) DeleteUser(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// trigger user deletion
	if err := userController.userService.DeleteUser(id); err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusNoContent, gin.H{})
}
