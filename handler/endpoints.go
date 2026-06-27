package handler

import (
	"fmt"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)
func (s *Server) GetHello(ctx echo.Context, params generated.GetHelloParams) error {
	var resp generated.HelloResponse
	resp.Message = fmt.Sprintf("Hello User %d", params.Id)
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) PostEstate(ctx echo.Context) error {
	var req generated.PostEstateJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "invalid request body"})
	}

	if err := validateEstateInput(req.Length, req.Width); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
	}

	id := uuid.NewString()
	if err := s.Repository.CreateEstate(ctx.Request().Context(), repository.CreateEstateInput{ID: id, Length: req.Length, Width: req.Width}); err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "failed to create estate"})
	}

	resp := generated.CreateEstateResponse{Id: id}
	return ctx.JSON(http.StatusCreated, resp)
}

func validateEstateInput(length, width int) error {
	if length <= 0 || width <= 0 {
		return fmt.Errorf("length and width must be greater than zero")
	}
	if length > 50000 || width > 50000 {
		return fmt.Errorf("length and width must be between 1 and 50000")
	}
	return nil
}
