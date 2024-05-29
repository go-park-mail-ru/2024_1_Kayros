package restaurants

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"2024_1_kayros/internal/entity"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRestaurantList(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("get all without filter", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()

		s.mockUcRest.EXPECT().GetAll(ctx).Return(nil, fmt.Errorf("error"))
		request := httptest.NewRequest(http.MethodGet, "/api/v1/restaurants", nil)
		respWriter := httptest.NewRecorder()
		s.handler.RestaurantList(respWriter, request)
		resp := respWriter.Result()

		_, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		err = resp.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("filter conversion error", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()

		request := httptest.NewRequest(http.MethodGet, "/api/v1/restaurants?filter=x", nil)
		respWriter := httptest.NewRecorder()
		s.handler.RestaurantList(respWriter, request)
		resp := respWriter.Result()

		_, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		err = resp.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("get with filter", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()

		s.mockUcRest.EXPECT().GetByFilter(ctx, gomock.Any()).Return(nil, fmt.Errorf("error"))
		request := httptest.NewRequest(http.MethodGet, "/api/v1/restaurants?filter=1", nil)
		respWriter := httptest.NewRecorder()
		s.handler.RestaurantList(respWriter, request)
		resp := respWriter.Result()

		_, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		err = resp.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("get all success", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()

		s.mockUcRest.EXPECT().GetAll(ctx).Return([]*entity.Restaurant{}, nil)
		request := httptest.NewRequest(http.MethodGet, "/api/v1/restaurants", nil)
		respWriter := httptest.NewRecorder()
		s.handler.RestaurantList(respWriter, request)
		resp := respWriter.Result()

		_, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		err = resp.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

}

func TestRestaurantById(t *testing.T) {
	t.Parallel()

	t.Run("bad id", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()

		request := httptest.NewRequest(http.MethodGet, "/api/v1/restaurants/x", nil)
		request = mux.SetURLVars(request, map[string]string{"id": "x"})
		respWriter := httptest.NewRecorder()
		s.handler.RestaurantById(respWriter, request)
		resp := respWriter.Result()

		_, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		err = resp.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("internal error", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()

		s.mockUcRest.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
		request := httptest.NewRequest(http.MethodGet, "/api/v1/restaurants/1", nil)
		request = mux.SetURLVars(request, map[string]string{"id": "1"})
		respWriter := httptest.NewRecorder()
		s.handler.RestaurantById(respWriter, request)
		resp := respWriter.Result()

		_, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		err = resp.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

}
