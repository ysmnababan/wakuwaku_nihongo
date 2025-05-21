package example_feat

import (
	"fmt"
	"wakuwaku_nihongo/internals/factory"
	"wakuwaku_nihongo/internals/utils/response"

	"github.com/labstack/echo/v4"
)

type IUserService interface {
	Get(ctx echo.Context) (out []*UserResponse, err error)
	Create(ctx echo.Context, in *UserCreateRequest) (err error)
	Login(ctx echo.Context, req *UserLoginRequest) (out *UserLoginResponse, err error)
}

type handler struct {
	service IUserService
}

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		service: NewService(f),
	}
}

// @Summary Get List of User
// @Description Get list of User
// @Tags user
// @Produce json
// @Success 200 {object} response.Success{data=[]UserResponse}
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Param Authorization header string true "Bearer Token"
// @Router /api/v1/users [get]
func (h *handler) GetUsers(c echo.Context) error {
	fmt.Println(c.Get("user_id"))
	res, err := h.service.Get(c)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	return response.SuccessResponse(res).Send(c)
}

// @Summary Create User
// @Description Create new User
// @Tags user
// @Accept json
// @Produce json
// @Param payload body UserCreateRequest true "Payload"
// @Success 200 {object} response.Success{data=string}
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /api/v1/users [post]
func (h *handler) CreateUser(c echo.Context) error {
	req := &UserCreateRequest{}
	err := c.Bind(req)
	if err != nil {
		return response.ErrorWrap(response.ErrUnprocessableEntity, err).Send(c)
	}

	err = c.Validate(req)
	if err != nil {
		return response.ErrorWrap(response.ErrValidation, err).Send(c)
	}
	err = h.service.Create(c, req)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse("mantap").Send(c)
}

// @Summary Login
// @Description User Login
// @Tags user
// @Accept json
// @Produce json
// @Param payload body UserLoginRequest true "Payload"
// @Success 200 {object} response.Success{data=UserLoginResponse}
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /api/v1/users/auth [post]
func (h *handler) Login(c echo.Context) error {
	req := &UserLoginRequest{}
	err := c.Bind(req)
	if err != nil {
		return response.ErrorWrap(response.ErrUnprocessableEntity, err).Send(c)
	}

	err = c.Validate(req)
	if err != nil {
		return response.ErrorWrap(response.ErrValidation, err).Send(c)
	}

	res, err := h.service.Login(c, req)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(res).Send(c)
}
