package types

import (
	"testing"

	gridiron "github.com/gridiron-zone/gridiron/types"
	"github.com/stretchr/testify/require"
)

func TestDefaultParams(t *testing.T) {
	params := DefaultParams()
	require.Equal(t, gridiron.BaseDenom, params.LockDenom)
}
