package usecase

import (
	"context"
	"fmt"
	"testing"

	"2024_1_kayros/gen/go/rest"
	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("internal error", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()

		s.mockRepo.EXPECT().GetAll(ctx).Return(nil, fmt.Errorf("error"))
		rests, err := s.layer.GetAll(ctx, &rest.Empty{})
		assert.Equal(t, &rest.RestList{}, rests)
		assert.Error(t, err)
	})
}

func TestGetById(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("internal error", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()
		restId := &rest.RestId{Id: 1}

		s.mockRepo.EXPECT().GetById(ctx, restId).Return(nil, fmt.Errorf("error"))
		r, err := s.layer.GetById(ctx, restId)
		assert.Equal(t, &rest.Rest{}, r)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()
		var id uint64 = 1
		restId := &rest.RestId{Id: id}
		restExpected := &rest.Rest{
			Id:               1,
			Name:             "a",
			ShortDescription: "a",
			LongDescription:  "a",
			Address:          "a",
			ImgUrl:           "a",
			Rating:           5,
			CommentCount:     5,
		}

		s.mockRepo.EXPECT().GetById(ctx, restId).Return(restExpected, nil)
		r, err := s.layer.GetById(ctx, restId)
		assert.Equal(t, restExpected, r)
		assert.NoError(t, err)
	})
}

func TestGetByFilter(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("internal error", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()
		restId := &rest.Id{Id: 1}

		s.mockRepo.EXPECT().GetByFilter(ctx, restId).Return(nil, fmt.Errorf("error"))
		rests, err := s.layer.GetByFilter(ctx, restId)
		assert.Equal(t, &rest.RestList{}, rests)
		assert.Error(t, err)
	})
}
