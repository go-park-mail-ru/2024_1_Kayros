package auth

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"2024_1_kayros/internal/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignUp(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("already signed up", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()

		request := httptest.NewRequest(http.MethodGet, "/api/v1/signup", nil)
		ctx = request.Context()
		ctx = context.WithValue(ctx, "email", "aaa@aa.aa")
		respWriter := httptest.NewRecorder()
		s.handler.SignUp(respWriter, request.WithContext(ctx))
		resp := respWriter.Result()

		_, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		err = resp.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("can not read request body", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()

		request := httptest.NewRequest(http.MethodGet, "/api/v1/signup", &errorReader{})
		respWriter := httptest.NewRecorder()
		s.handler.SignUp(respWriter, request)
		resp := respWriter.Result()

		_, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		err = resp.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("unmarshall error", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()

		request := httptest.NewRequest(http.MethodGet, "/api/v1/signup", strings.NewReader("{"))
		respWriter := httptest.NewRecorder()
		s.handler.SignUp(respWriter, request)
		resp := respWriter.Result()

		_, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		err = resp.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("validation error", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()

		request := httptest.NewRequest(http.MethodGet, "/api/v1/signup", strings.NewReader("{}"))
		respWriter := httptest.NewRecorder()
		s.handler.SignUp(respWriter, request)
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

		request := httptest.NewRequest(http.MethodGet, "/api/v1/signup", strings.NewReader(`{"name": "aaa", "email": "aaa@aa.aa", "password": "qqqqqqq1"}`))
		respWriter := httptest.NewRecorder()
		s.mockUcAuth.EXPECT().SignUp(request.Context(), gomock.Any()).Return(nil, fmt.Errorf("error"))
		s.handler.SignUp(respWriter, request)
		resp := respWriter.Result()

		_, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		err = resp.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("ok", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()

		request := httptest.NewRequest(http.MethodGet, "/api/v1/signup", strings.NewReader(`{"name": "aaa", "email": "aaa@aa.aa", "password": "qqqqqqq1"}`))
		respWriter := httptest.NewRecorder()
		s.mockUcAuth.EXPECT().SignUp(request.Context(), gomock.Any()).Return(&entity.User{}, nil)
		s.mockUcSess.EXPECT().SetValue(request.Context(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
		s.handler.SignUp(respWriter, request)
		resp := respWriter.Result()

		_, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		err = resp.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestSignIn(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("already signed up", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()

		request := httptest.NewRequest(http.MethodGet, "/api/v1/signin", nil)
		ctx = request.Context()
		ctx = context.WithValue(ctx, "email", "aaa@aa.aa")
		respWriter := httptest.NewRecorder()
		s.handler.SignIn(respWriter, request.WithContext(ctx))
		resp := respWriter.Result()

		_, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		err = resp.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("validation error", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()

		request := httptest.NewRequest(http.MethodGet, "/api/v1/signin", strings.NewReader("{}"))
		respWriter := httptest.NewRecorder()
		s.handler.SignIn(respWriter, request)
		resp := respWriter.Result()

		_, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		err = resp.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ok", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()

		request := httptest.NewRequest(http.MethodGet, "/api/v1/signin", strings.NewReader(`{"email": "aaa@aa.aa", "password": "qqqqqqq1"}`))
		respWriter := httptest.NewRecorder()
		s.mockUcAuth.EXPECT().SignIn(request.Context(), gomock.Any(), gomock.Any()).Return(&entity.User{}, nil)
		s.mockUcSess.EXPECT().SetValue(request.Context(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
		s.handler.SignIn(respWriter, request)
		resp := respWriter.Result()

		_, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		err = resp.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
