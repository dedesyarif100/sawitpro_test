package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestGetHello(t *testing.T) {
	t.Run("returns hello message", func(t *testing.T) {
		srv := NewServer(NewServerOptions{})

		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/hello?id=123", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		params := generated.GetHelloParams{
			Id: 123,
		}

		err := srv.GetHello(ctx, params)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var resp generated.HelloResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		expected := "Hello User 123"
		if resp.Message != expected {
			t.Fatalf("expected message %q, got %q", expected, resp.Message)
		}
	})
}

func TestPostEstate(t *testing.T) {
	t.Run("returns 201 when estate is created", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().
			CreateEstate(
				gomock.Any(),
				gomock.Any(),
			).
			DoAndReturn(func(_ context.Context, input repository.CreateEstateInput) error {
				if input.ID == "" {
					t.Fatal("expected generated id")
				}
				if input.Name == "" {
					t.Fatal("expected generated name")
				}
				if input.Length != 10 {
					t.Fatalf("expected length 10, got %d", input.Length)
				}
				if input.Width != 20 {
					t.Fatalf("expected width 20, got %d", input.Width)
				}
				return nil
			})

		srv := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()

		body := `{"length":10,"width":20}`
		req := httptest.NewRequest(http.MethodPost, "/estate", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		err := srv.PostEstate(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
		}
	})

	t.Run("returns 400 for invalid request body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		srv := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/estate", strings.NewReader(`{`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		err := srv.PostEstate(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("returns 400 for invalid estate input", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		srv := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()

		body := `{"length":0,"width":20}`
		req := httptest.NewRequest(http.MethodPost, "/estate", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		err := srv.PostEstate(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("returns 500 when repository returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().
			CreateEstate(
				gomock.Any(),
				gomock.Any(),
			).
			Return(errors.New("database error"))

		srv := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()

		body := `{"length":10,"width":20}`
		req := httptest.NewRequest(http.MethodPost, "/estate", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		err := srv.PostEstate(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestPostEstateIdTree(t *testing.T) {
	t.Run("returns 201 when tree is created", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		mockRepo.EXPECT().
			EstateExists(gomock.Any(), "estate-1").
			Return(true, nil)

		mockRepo.EXPECT().
			TreeExists(gomock.Any(), "estate-1", 1, 2).
			Return(false, nil)

		mockRepo.EXPECT().
			CreateTree(
				gomock.Any(),
				gomock.Any(),
			).
			DoAndReturn(func(_ context.Context, input repository.CreateTreeInput) error {
				if input.ID == "" {
					t.Fatal("expected generated tree id")
				}
				if input.EstateID != "estate-1" {
					t.Fatalf("expected estate id estate-1, got %s", input.EstateID)
				}
				if input.X != 1 {
					t.Fatalf("expected X=1, got %d", input.X)
				}
				if input.Y != 2 {
					t.Fatalf("expected Y=2, got %d", input.Y)
				}
				if input.Height != 5 {
					t.Fatalf("expected height=5, got %d", input.Height)
				}
				return nil
			})

		srv := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()

		body := `{"x":1,"y":2,"height":5}`
		req := httptest.NewRequest(http.MethodPost, "/estate/estate-1/tree", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("estate-1")

		err := srv.PostEstateIdTree(ctx, "estate-1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
		}
	})

	t.Run("returns 500 when EstateExists fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		mockRepo.EXPECT().
			EstateExists(gomock.Any(), "estate-1").
			Times(1).
			Return(false, errors.New("database error"))

		srv := NewServer(NewServerOptions{
			Repository: mockRepo,
		})

		e := echo.New()

		req := httptest.NewRequest(
			http.MethodPost,
			"/estate/estate-1/tree",
			strings.NewReader(`{
			"x":1,
			"y":1,
			"height":10
		}`),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("estate-1")

		err := srv.PostEstateIdTree(ctx, "estate-1")
		require.NoError(t, err)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("returns 400 for invalid request body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		srv := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/estate/estate-1/tree", strings.NewReader(`{`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("estate-1")

		err := srv.PostEstateIdTree(ctx, "estate-1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("returns 400 for invalid tree input", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		srv := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()

		body := `{"x":0,"y":2,"height":5}`
		req := httptest.NewRequest(http.MethodPost, "/estate/estate-1/tree", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("estate-1")

		err := srv.PostEstateIdTree(ctx, "estate-1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("returns 500 when EstateExists fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		mockRepo.EXPECT().
			EstateExists(gomock.Any(), "estate-1").
			Return(false, errors.New("db error"))

		srv := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()

		body := `{"x":1,"y":2,"height":5}`
		req := httptest.NewRequest(http.MethodPost, "/estate/estate-1/tree", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("estate-1")

		_ = srv.PostEstateIdTree(ctx, "estate-1")

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("returns 404 when estate does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		mockRepo.EXPECT().
			EstateExists(gomock.Any(), "estate-1").
			Return(false, nil)

		srv := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()

		body := `{"x":1,"y":2,"height":5}`
		req := httptest.NewRequest(http.MethodPost, "/estate/estate-1/tree", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("estate-1")

		_ = srv.PostEstateIdTree(ctx, "estate-1")

		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("returns 500 when TreeExists fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		mockRepo.EXPECT().
			EstateExists(gomock.Any(), "estate-1").
			Return(true, nil)

		mockRepo.EXPECT().
			TreeExists(gomock.Any(), "estate-1", 1, 2).
			Return(false, errors.New("db error"))

		srv := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()

		body := `{"x":1,"y":2,"height":5}`
		req := httptest.NewRequest(http.MethodPost, "/estate/estate-1/tree", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("estate-1")

		_ = srv.PostEstateIdTree(ctx, "estate-1")

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("returns 400 when tree already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		mockRepo.EXPECT().
			EstateExists(gomock.Any(), "estate-1").
			Return(true, nil)

		mockRepo.EXPECT().
			TreeExists(gomock.Any(), "estate-1", 1, 2).
			Return(true, nil)

		srv := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()

		body := `{"x":1,"y":2,"height":5}`
		req := httptest.NewRequest(http.MethodPost, "/estate/estate-1/tree", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("estate-1")

		_ = srv.PostEstateIdTree(ctx, "estate-1")

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("returns 500 when CreateTree fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		mockRepo.EXPECT().
			EstateExists(gomock.Any(), "estate-1").
			Return(true, nil)

		mockRepo.EXPECT().
			TreeExists(gomock.Any(), "estate-1", 1, 2).
			Return(false, nil)

		mockRepo.EXPECT().
			CreateTree(gomock.Any(), gomock.Any()).
			Return(errors.New("db error"))

		srv := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()

		body := `{"x":1,"y":2,"height":5}`
		req := httptest.NewRequest(http.MethodPost, "/estate/estate-1/tree", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("estate-1")

		_ = srv.PostEstateIdTree(ctx, "estate-1")

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestGetEstateStats(t *testing.T) {
	t.Run("returns stats for existing estate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().EstateExists(gomock.Any(), "estate-1").Return(true, nil)
		mockRepo.EXPECT().GetTreeStats(gomock.Any(), "estate-1").Return(repository.TreeStats{Count: 3, Max: 15, Min: 10, Median: 12}, nil)

		srv := NewServer(NewServerOptions{Repository: mockRepo})
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/estate/estate-1/stats", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("estate-1")

		err := srv.GetEstateIdStats(ctx, "estate-1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("returns 500 when EstateExists fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		mockRepo.EXPECT().
			EstateExists(gomock.Any(), "estate-1").
			Return(false, errors.New("database error"))

		srv := NewServer(NewServerOptions{
			Repository: mockRepo,
		})

		e := echo.New()

		req := httptest.NewRequest(
			http.MethodGet,
			"/estate/estate-1/stats",
			nil,
		)

		rec := httptest.NewRecorder()

		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("estate-1")

		err := srv.GetEstateIdStats(ctx, "estate-1")
		require.NoError(t, err)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("returns 500 when GetTreeStats fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		mockRepo.EXPECT().
			EstateExists(gomock.Any(), "estate-1").
			Return(true, nil)

		mockRepo.EXPECT().
			GetTreeStats(gomock.Any(), "estate-1").
			Return(repository.TreeStats{}, errors.New("database error"))

		srv := NewServer(NewServerOptions{
			Repository: mockRepo,
		})

		e := echo.New()

		req := httptest.NewRequest(
			http.MethodGet,
			"/estate/estate-1/stats",
			nil,
		)

		rec := httptest.NewRecorder()

		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("estate-1")

		err := srv.GetEstateIdStats(ctx, "estate-1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("returns 404 for missing estate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().EstateExists(gomock.Any(), "missing").Return(false, nil)

		srv := NewServer(NewServerOptions{Repository: mockRepo})
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/estate/missing/stats", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("missing")

		err := srv.GetEstateIdStats(ctx, "missing")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})
}

func TestGetEstateDronePlan(t *testing.T) {
	t.Run("returns distance and rest for existing estate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().EstateExists(gomock.Any(), "estate-1").Return(true, nil)
		mockRepo.EXPECT().GetDronePlanSummary(gomock.Any(), "estate-1").Return(repository.DronePlanSummary{Distance: 100, HorizontalDistance: 60, VerticalDistance: 40}, nil)

		srv := NewServer(NewServerOptions{Repository: mockRepo})
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/estate/estate-1/drone-plan?max_distance=70", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("estate-1")

		params := generated.GetEstateIdDronePlanParams{MaxDistance: intPtr(70)}
		err := srv.GetEstateIdDronePlan(ctx, "estate-1", params)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("returns 500 when EstateExists fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		mockRepo.EXPECT().
			EstateExists(gomock.Any(), "estate-1").
			Times(1).
			Return(false, errors.New("database error"))

		srv := NewServer(NewServerOptions{
			Repository: mockRepo,
		})

		e := echo.New()

		req := httptest.NewRequest(
			http.MethodGet,
			"/estate/estate-1/drone-plan",
			nil,
		)

		rec := httptest.NewRecorder()

		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("estate-1")

		params := generated.GetEstateIdDronePlanParams{}

		err := srv.GetEstateIdDronePlan(ctx, "estate-1", params)
		require.NoError(t, err)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("returns 400 for invalid max_distance", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().EstateExists(gomock.Any(), "estate-1").Return(true, nil)

		srv := NewServer(NewServerOptions{Repository: mockRepo})
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/estate/estate-1/drone-plan?max_distance=0", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("estate-1")

		params := generated.GetEstateIdDronePlanParams{MaxDistance: intPtr(0)}
		err := srv.GetEstateIdDronePlan(ctx, "estate-1", params)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("returns 404 for missing estate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().EstateExists(gomock.Any(), "missing").Return(false, nil)

		srv := NewServer(NewServerOptions{Repository: mockRepo})
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/estate/missing/drone-plan", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("missing")

		params := generated.GetEstateIdDronePlanParams{}
		err := srv.GetEstateIdDronePlan(ctx, "missing", params)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("returns 500 when GetDronePlanSummary fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		mockRepo.EXPECT().
			EstateExists(gomock.Any(), "estate-1").
			Return(true, nil)

		mockRepo.EXPECT().
			GetDronePlanSummary(gomock.Any(), "estate-1").
			Return(repository.DronePlanSummary{}, errors.New("database error"))

		srv := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/estate/estate-1/drone-plan", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("estate-1")

		params := generated.GetEstateIdDronePlanParams{}

		err := srv.GetEstateIdDronePlan(ctx, "estate-1", params)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected %d got %d", http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("returns rest when distance is within max distance", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		mockRepo.EXPECT().
			EstateExists(gomock.Any(), "estate-1").
			Return(true, nil)

		mockRepo.EXPECT().
			GetDronePlanSummary(gomock.Any(), "estate-1").
			Return(repository.DronePlanSummary{
				Distance:           50,
				HorizontalDistance: 30,
				VerticalDistance:   20,
			}, nil)

		srv := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/estate/estate-1/drone-plan?max_distance=100", nil)
		rec := httptest.NewRecorder()

		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("estate-1")

		maxDistance := 100
		params := generated.GetEstateIdDronePlanParams{
			MaxDistance: &maxDistance,
		}

		err := srv.GetEstateIdDronePlan(ctx, "estate-1", params)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("expected %d got %d", http.StatusOK, rec.Code)
		}

		var resp generated.DronePlanResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}

		if resp.Rest == nil {
			t.Fatal("expected rest")
		}

		if resp.Rest.X != 30 {
			t.Fatalf("expected x=30 got %d", resp.Rest.X)
		}

		if resp.Rest.Y != 20 {
			t.Fatalf("expected y=20 got %d", resp.Rest.Y)
		}
	})
}

func TestGenerateEstateName(t *testing.T) {
	name := generateEstateName("123e4567-e89b-12d3-a456-426614174000")
	if name == "" {
		t.Fatal("expected a generated estate name")
	}
	if name == "123e4567-e89b-12d3-a456-426614174000" {
		t.Fatal("expected generated name to be derived from the id")
	}
}

func TestValidateEstateInput(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		width   int
		wantErr bool
	}{
		{name: "valid", length: 10, width: 5, wantErr: false},
		{name: "zero length", length: 0, width: 5, wantErr: true},
		{name: "negative width", length: 10, width: -1, wantErr: true},
		{name: "too large", length: 50001, width: 1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEstateInput(tt.length, tt.width)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error for length=%d width=%d", tt.length, tt.width)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func intPtr(v int) *int {
	return &v
}

func TestValidateTreeInput(t *testing.T) {
	tests := []struct {
		name    string
		x       int
		y       int
		height  int
		wantErr bool
	}{
		{name: "valid", x: 2, y: 3, height: 10, wantErr: false},
		{name: "zero x", x: 0, y: 3, height: 10, wantErr: true},
		{name: "height too high", x: 2, y: 3, height: 31, wantErr: true},
		{name: "too large coordinate", x: 50001, y: 1, height: 10, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTreeInput(tt.x, tt.y, tt.height)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error for x=%d y=%d height=%d", tt.x, tt.y, tt.height)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}
