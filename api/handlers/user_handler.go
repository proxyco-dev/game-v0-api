package handlers

import (
	"game-v0-api/api/presenter"
	common "game-v0-api/pkg/common"
	entities "game-v0-api/pkg/entities"
	repository "game-v0-api/pkg/user"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepo  repository.UserRepository
	validator *validator.Validate
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo:  userRepo,
		validator: validator.New(),
	}
}

// getMe godoc
// @Summary Get current user
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} presenter.MeResponse
// @Failure 401 {object} presenter.ErrorResponse
// @Router /api/user/me [get]
func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	return c.Status(200).JSON(presenter.MeResponse{Email: email})
}

// signIn godoc
// @Summary Sign in
// @Tags User
// @Accept json
// @Produce json
// @Param user body presenter.SignInRequest true "User credentials"
// @Success 200 {object} presenter.SignInResponse
// @Failure 400 {object} presenter.ErrorResponse
// @Failure 401 {object} presenter.ErrorResponse
// @Failure 500 {object} presenter.ErrorResponse
// @Router /api/user/sign-in [post]
func (h *UserHandler) SignIn(c *fiber.Ctx) error {
	var request presenter.SignInRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(presenter.ErrorResponse{Error: "დაფიქსირდა შეცდომა"})
	}

	user, err := h.userRepo.FindByEmail(request.Email)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(presenter.ErrorResponse{Error: "შეყვანილი მონაცემები არ არის სწორი"})
	}

	if !checkPasswordHash(request.Password, user.Password) {
		return c.Status(401).JSON(presenter.ErrorResponse{Error: "შეყვანილი მონაცემები არ არის სწორი"})
	}

	token, err := generateToken(user.ID.String(), user.Email)
	if err != nil {
		return c.Status(500).JSON(presenter.ErrorResponse{Error: "ვერ მოხერხდა ტოკენის გენერაცია"})
	}

	return c.Status(200).JSON(presenter.SignInResponse{Token: token})
}

// signUp godoc
// @Summary Sign up
// @Tags User
// @Accept json
// @Produce json
// @Param user body presenter.SignUpRequest true "User credentials"
// @Success 200 {object} presenter.SignUpResponse "Returns message and user object with excluded fields"
// @Failure 400 {object} presenter.ErrorResponse
// @Failure 500 {object} presenter.ErrorResponse
// @Router /api/user/sign-up [post]
func (h *UserHandler) SignUp(c *fiber.Ctx) error {
	var request presenter.SignUpRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(presenter.ErrorResponse{Error: "მონაცემები ვერ დამუშავდა"})
	}

	if err := h.validator.Struct(request); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(400).JSON(presenter.ErrorResponse{Error: validationErrors.Error()})
	}

	_, err := h.userRepo.FindByEmail(request.Email)
	if err == nil {
		return c.Status(400).JSON(presenter.ErrorResponse{Error: "მომხმარებელი უკვე არსებობს"})
	}

	hashedPassword, err := hashPassword(request.Password)
	if err != nil {
		return c.Status(500).JSON(presenter.ErrorResponse{Error: "ვერ მოხერხდა პაროლის ჰეშირება"})
	}

	newUser := &entities.User{
		Username: request.Username,
		Email:    request.Email,
		Password: hashedPassword,
	}

	if err := h.userRepo.Create(newUser); err != nil {
		return c.Status(500).JSON(presenter.ErrorResponse{Error: "ვერ მოხერხდა მომხმარებლის დამატება"})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "მომხმარებელი წარმატებით დაემატა",
		"user":    common.Exclude(newUser, []string{"password", "refreshToken", "created_at", "avatar", "updated_at"}),
	})
}

type FindUserResponse struct {
	Message string `json:"message"`
	Users   []struct {
		Id       string `json:"id"`
		Username string `json:"username"`
	} `json:"users"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func generateToken(id string, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
