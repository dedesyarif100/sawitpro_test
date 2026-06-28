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
	name := generateEstateName(id)
	if err := s.Repository.CreateEstate(ctx.Request().Context(), repository.CreateEstateInput{ID: id, Name: name, Length: req.Length, Width: req.Width}); err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "failed to create estate"})
	}

	resp := generated.CreateEstateResponse{Id: id}
	return ctx.JSON(http.StatusCreated, resp)
}

func (s *Server) PostEstateIdTree(ctx echo.Context, id string) error {
	var req generated.PostEstateIdTreeJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "invalid request body"})
	}

	if err := validateTreeInput(req.X, req.Y, req.Height); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
	}

	exists, err := s.Repository.EstateExists(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "failed to check estate"})
	}
	if !exists {
		return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{Message: "estate not found"})
	}

	treeExists, err := s.Repository.TreeExists(ctx.Request().Context(), id, req.X, req.Y)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "failed to check tree"})
	}
	if treeExists {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "plot already has tree"})
	}

	treeID := uuid.NewString()
	if err := s.Repository.CreateTree(ctx.Request().Context(), repository.CreateTreeInput{ID: treeID, EstateID: id, X: req.X, Y: req.Y, Height: req.Height}); err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "failed to create tree"})
	}

	resp := generated.CreateTreeResponse{Id: treeID}
	return ctx.JSON(http.StatusCreated, resp)
}

func (s *Server) GetEstateIdStats(ctx echo.Context, id string) error {
	exists, err := s.Repository.EstateExists(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "failed to check estate"})
	}
	if !exists {
		return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{Message: "estate not found"})
	}

	stats, err := s.Repository.GetTreeStats(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "failed to get tree stats"})
	}

	resp := generated.TreeStatsResponse{Count: stats.Count, Max: stats.Max, Min: stats.Min, Median: stats.Median}
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) GetEstateIdDronePlan(ctx echo.Context, id string, params generated.GetEstateIdDronePlanParams) error {
	exists, err := s.Repository.EstateExists(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "failed to check estate"})
	}
	if !exists {
		return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{Message: "estate not found"})
	}

	if params.MaxDistance != nil && *params.MaxDistance <= 0 {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "max_distance must be greater than zero"})
	}

	summary, err := s.Repository.GetDronePlanSummary(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "failed to get drone plan distance"})
	}

	resp := generated.DronePlanResponse{Distance: summary.Distance}
	if params.MaxDistance != nil {
		resp.Rest = &generated.DronePlanRest{}
		if summary.Distance <= int64(*params.MaxDistance) {
			resp.Rest.X = summary.HorizontalDistance
			resp.Rest.Y = summary.VerticalDistance
		} else {
			// Return the last point coordinates for the full path when the maximum distance is exceeded.
			resp.Rest.X = summary.HorizontalDistance
			resp.Rest.Y = summary.VerticalDistance
		}
	}
	return ctx.JSON(http.StatusOK, resp)
}

func generateEstateName(id string) string {
	return fmt.Sprintf("Estate-%s", id[:8])
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

func validateTreeInput(x, y, height int) error {
	if x <= 0 || y <= 0 || height <= 0 {
		return fmt.Errorf("x, y, and height must be greater than zero")
	}
	if x > 50000 || y > 50000 {
		return fmt.Errorf("x and y must be between 1 and 50000")
	}
	if height > 30 {
		return fmt.Errorf("height must be between 1 and 30")
	}
	return nil
}
