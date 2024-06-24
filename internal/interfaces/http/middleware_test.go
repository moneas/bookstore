package http

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/moneas/bookstore/internal/application"
	"github.com/moneas/bookstore/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindAll() ([]user.User, error) {
	args := m.Called()
	return args.Get(0).([]user.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*user.User, error) {
	args := m.Called(id)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) Save(u *user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	args := m.Called(email)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(username string) (*user.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

type MockUserService struct {
	*application.UserService
}

func NewMockUserService(repo user.Repository) *MockUserService {
	return &MockUserService{
		UserService: application.NewUserService(repo),
	}
}

func hashMD5(password string) string {
	hasher := md5.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func TestBasicAuthMiddleware(t *testing.T) {
	r := gin.New()

	mockUserRepo := new(MockUserRepository)
	mockUser := &user.User{ID: 1, Username: "user1", Password: hashMD5("password")}
	mockUserRepo.On("FindByUsername", "user1").Return(mockUser, nil)
	mockUserRepo.On("FindByUsername", "invalid").Return(nil, errors.New("user not found"))

	userService := NewMockUserService(mockUserRepo)

	r.Use(BasicAuthMiddleware(userService.UserService))

	r.GET("/protected", func(c *gin.Context) {
		c.String(http.StatusOK, "Access granted")
	})

	encodedCredentials := base64.StdEncoding.EncodeToString([]byte("user1:password"))
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Basic "+encodedCredentials)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Access granted", w.Body.String())

	encodedCredentials = base64.StdEncoding.EncodeToString([]byte("invalid:credentials"))
	req, _ = http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Basic "+encodedCredentials)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid username or password")

	req, _ = http.NewRequest(http.MethodGet, "/protected", nil)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header is required")
}
