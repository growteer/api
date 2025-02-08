package gqlutil_test

import (
	"context"
	"errors"
	"testing"

	"github.com/growteer/api/internal/api/graphql/gqlutil"
	"github.com/growteer/api/internal/app/shared/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Recover(t *testing.T) {
	t.Run("error passed", func(t *testing.T) {
		// given
		ctx := context.Background()
		err := errors.New("some error")

		// when
		result := gqlutil.Recover(ctx, err)

		// then
		require.Error(t, result)

		var internal apperrors.Internal
		assert.ErrorAs(t, result, &internal)
		assert.Equal(t, err, internal.Wrapped)
	})

	t.Run("string passed", func(t *testing.T) {
		// given
		ctx := context.Background()
		err := "some other error"

		// when
		result := gqlutil.Recover(ctx, err)

		// then
		require.Error(t, result)

		var internal apperrors.Internal
		assert.ErrorAs(t, result, &internal)
		assert.Nil(t, internal.Wrapped)
	})
}
