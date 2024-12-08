package http

import (
	"net/http"
	user "product-recommendation/internal/application/usecases"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

type UserHandler struct {
	registerUseCase  *user.RegisterUserUseCase
	ListUsersUseCase *user.ListUsersUseCase
}

type createUserRequest struct {
	Name     string `json:"name" validate:"required,min=4"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func NewUserHandler(registerUseCase *user.RegisterUserUseCase, listUseCase *user.ListUsersUseCase) *UserHandler {
	return &UserHandler{
		registerUseCase:  registerUseCase,
		ListUsersUseCase: listUseCase,
	}
}

func (handler *UserHandler) Register(context *gin.Context) {
	var request createUserRequest

	id := uuid.New().String()

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := validate.Struct(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := handler.registerUseCase.Execute(id, request.Name, request.Email, request.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Usuário registrado com sucesso"})
}

func (handler *UserHandler) List(context *gin.Context) {
	users, err := handler.ListUsersUseCase.Execute()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	response := make([]gin.H, len(users))

	for i, user := range users {
		response[i] = gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		}
	}

	context.JSON(http.StatusOK, response)
}