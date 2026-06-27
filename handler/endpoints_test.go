package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
	gomock "go.uber.org/mock/gomock"
)

func TestGetHello(t *testing.T) {

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
