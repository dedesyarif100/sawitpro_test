// This file will run automated tests for API.
package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const ApiUrl = "http://localhost:8080"

func TestApi(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip API tests")
	}

	// Check API is running before executing all test cases.
	resp, err := http.Get(ApiUrl + "/hello?id=1")
	if err != nil {
		t.Fatalf("API is not running at %s: %v", ApiUrl, err)
	}
	resp.Body.Close()

	client := &http.Client{}
	ctx := context.Background()

	for _, tc := range getTestCases() {
		tc := tc

		t.Run(tc.Name, func(t *testing.T) {
			for i := range tc.Steps {
				step := &tc.Steps[i]

				req, err := step.Request(t, ctx, &tc)
				require.NoError(t, err)

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Accept", "application/json")

				resp, err := client.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				ReadJsonResult(t, resp, step)

				step.Expect(t, ctx, &tc, resp, step.Result)
			}
		})
	}
}

func getTestCases() []TestCase {
	return []TestCase{
		{
			Name: "Test Hello",
			Steps: []TestCaseStep{
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						return http.NewRequest("GET", ApiUrl+"/hello", nil)
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
						require.Equal(t, http.StatusBadRequest, resp.StatusCode)
					},
				},
			},
		},
		{
			Name: "Test Hello with name",
			Steps: []TestCaseStep{
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						return http.NewRequest("GET", ApiUrl+"/hello?id=123", nil)
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
						require.Equal(t, http.StatusOK, resp.StatusCode)
						require.Equal(t, "Hello User 123", data["message"])
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						return http.NewRequest("GET", ApiUrl+"/hello?id=456", nil)
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
						require.Equal(t, http.StatusOK, resp.StatusCode)
						step1 := tc.Steps[0]
						require.Equal(t, "Hello User 123", step1.Result["message"])
					},
				},
			},
		},
		//----- Test for API
		{
			Name: "Test Error 1",
			Steps: []TestCaseStep{
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						return http.NewRequest("POST", ApiUrl+"/estate", nil)
					},
					Expect: ExpectBadRequest(),
				},
			},
		},
		{
			Name: "Test Error 2: Invalid Format",
			Steps: []TestCaseStep{
				{
					Request: SendRequestNewEstate(-1, -5),
					Expect:  ExpectBadRequest(),
				},
			},
		},
		{
			Name: "Test Error: Create Tree Out of Bound",
			Steps: []TestCaseStep{
				{
					Request: SendRequestNewEstate(10, 20),
					Expect:  ExpectNewEstateOk(),
				},
				{
					Request: SendRequestNewTree(5, 0, 0),
					Expect:  ExpectBadRequest(),
				},
			},
		},
		CreateNormalTestCase("Normal 1", []any{
			[]any{CreateEstate, 10, 20},
			[]any{CreateTree, 10, 5, 5},
			[]any{CreateTree, 20, 6, 5},
		}),
		CreateNormalTestCase("Normal 2", []any{
			[]any{CreateEstate, 5, 1},
			[]any{CreateTree, 10, 2, 1},
			[]any{CreateTree, 20, 3, 1},
			[]any{CreateTree, 10, 4, 1},
			[]any{GetStats, 3, 10, 20, 10},
			[]any{GetDronePlan, 0, 0},
		}),
	}
}

type TestCase struct {
	Name  string
	Steps []TestCaseStep
}

type RequestFunc func(*testing.T, context.Context, *TestCase) (*http.Request, error)
type ExpectFunc func(*testing.T, context.Context, *TestCase, *http.Response, map[string]any)

type TestCaseStep struct {
	Request RequestFunc
	Expect  ExpectFunc
	Result  map[string]any
}

func ResponseContains(t *testing.T, resp *http.Response, text string) {
	body, err := io.ReadAll(resp.Body)
	bodyStr := string(body)
	require.NoError(t, err)
	require.Contains(t, bodyStr, text)
}

func ReadJsonResult(t *testing.T, resp *http.Response, step *TestCaseStep) {
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	step.Result = make(map[string]any)

	if len(body) == 0 {
		return
	}

	err = json.Unmarshal(body, &step.Result)
	require.NoError(t, err, "response body: %s", string(body))
}

func RequireIsUUID(t *testing.T, value string) {
	_, err := uuid.Parse(value)
	require.NoError(t, err)
}

const (
	CreateEstate = iota
	CreateTree
	GetStats
	GetDronePlan
)

func CreateNormalTestCase(name string, a []any) TestCase {
	tc := TestCase{}
	tc.Name = name

	for _, step := range a {
		switch step.([]any)[0].(int) {
		case CreateEstate:
			tc.Steps = append(tc.Steps, TestCaseStep{
				Request: SendRequestNewEstate(step.([]any)[1].(int), step.([]any)[2].(int)),
				Expect:  ExpectNewEstateOk(),
			})
		case CreateTree:
			tc.Steps = append(tc.Steps, TestCaseStep{
				Request: SendRequestNewTree(step.([]any)[1].(int), step.([]any)[2].(int), step.([]any)[3].(int)),
				Expect:  ExpectNewTreeOk(),
			})
		case GetStats:
			tc.Steps = append(tc.Steps, TestCaseStep{
				Request: SendRequestGetStats(),
				Expect:  ExpectGetStatsOk(step.([]any)[1].(int), step.([]any)[2].(int), step.([]any)[3].(int), step.([]any)[4].(int)),
			})
		case GetDronePlan:
			tc.Steps = append(tc.Steps, TestCaseStep{
				Request: SendRequestGetDronePlan(step.([]any)[1].(int)),
				Expect:  ExpectGetDronePlanOk(step.([]any)[2].(int)),
			})
		}

	}
	return tc
}

func SendRequestNewEstate(length, width int) RequestFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
		req := map[string]int{
			"length": length,
			"width":  width,
		}
		body, err := json.Marshal(req)
		require.NoError(t, err)
		return http.NewRequest("POST", ApiUrl+"/estate", bytes.NewReader(body))
	}
}

func ExpectNewEstateOk() ExpectFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
		RequireReturnIsUUID(t, resp, data)
	}
}

func SendRequestNewTree(height, x, y int) RequestFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {

		id, ok := tc.Steps[0].Result["id"].(string)
		require.True(t, ok, "estate id not found")

		req := map[string]int{
			"height": height,
			"x":      x,
			"y":      y,
		}

		body, err := json.Marshal(req)
		require.NoError(t, err)

		return http.NewRequest(
			http.MethodPost,
			ApiUrl+"/estate/"+id+"/tree",
			bytes.NewReader(body),
		)
	}
}

func ExpectNewTreeOk() ExpectFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
		RequireReturnIsUUID(t, resp, data)
	}
}

func SendRequestGetStats() RequestFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {

		id, ok := tc.Steps[0].Result["id"].(string)
		require.True(t, ok, "estate id not found")

		return http.NewRequest(
			http.MethodGet,
			ApiUrl+"/estate/"+id+"/stats",
			nil,
		)
	}
}

func ExpectGetStatsOk(count, min, max, median int) ExpectFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
		RequireStats(t, resp, data, count, min, max, median)
	}
}

func SendRequestGetDronePlan(distance int) RequestFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {

		id, ok := tc.Steps[0].Result["id"].(string)
		require.True(t, ok, "estate id not found")

		url := fmt.Sprintf("%s/estate/%s/drone-plan", ApiUrl, id)

		if distance > 0 {
			url += fmt.Sprintf("?max_distance=%d", distance)
		}

		return http.NewRequest(http.MethodGet, url, nil)
	}
}

func ExpectGetDronePlanOk(distance int) ExpectFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
		RequireDistance(t, resp, data, distance)
	}
}

func RequireReturnIsUUID(t *testing.T, resp *http.Response, data map[string]any) {
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	id, ok := data["id"].(string)
	require.True(t, ok, "response does not contain id")

	RequireIsUUID(t, id)
}

func RequireStats(t *testing.T, resp *http.Response, data map[string]any, count, min, max, median int) {
	require.Equal(t, http.StatusOK, resp.StatusCode)

	require.Equal(t, float64(count), data["count"])
	require.Equal(t, float64(min), data["min"])
	require.Equal(t, float64(max), data["max"])
	require.Equal(t, float64(median), data["median"])
}

func RequireDistance(t *testing.T, resp *http.Response, data map[string]any, distance int) {
	require.Equal(t, http.StatusOK, resp.StatusCode)

	require.Equal(t, float64(distance), data["distance"])
}

func ExpectBadRequest() ExpectFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
		body, _ := io.ReadAll(resp.Body)

		require.Equal(
			t,
			http.StatusBadRequest,
			resp.StatusCode,
			string(body),
		)
	}
}
