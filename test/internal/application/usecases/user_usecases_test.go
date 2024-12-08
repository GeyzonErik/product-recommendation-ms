package user

import (
	"errors"
	usecases "product-recommendation/internal/application/usecases"
	"product-recommendation/internal/domain/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockRepository struct {
	data map[string]*user.User
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		data: make(map[string]*user.User),
	}
}

func (mockRepository *MockRepository) Save(user *user.User) error {
	if user.ID == "exists" {
		return errors.New("Usuário já cadastrado!")
	}
	mockRepository.data[user.ID] = user
	return nil
}

func (mockRepo *MockRepository) FindAll() ([]*user.User, error) {
	result := []*user.User{}
	for _, user := range mockRepo.data {
		result = append(result, user)
	}
	return result, nil
}

func TestRegisterUserUseCase_Execute(test *testing.T) {
	mockRepo := NewMockRepository()
	useCase := usecases.NewRegisterUserUseCase(mockRepo)

	err := useCase.Execute("1", "João das Neves", "mail@mail.com", "12345")
	assert.NoError(test, err, "Não deve retornar erro com registro de usuário")

	err = useCase.Execute("exists", "João das Neves", "mail@mail.com", "123456")
	assert.EqualError(test, err, "Usuário já cadastrado!", "Não deve criar com ID existente")
}
