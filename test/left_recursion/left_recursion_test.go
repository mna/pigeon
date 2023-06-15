package leftrecursion

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLeftRecursion(t *testing.T) {
	t.Parallel()
	data := "7+10/2*-4+5*3%6-8*6"
	res, err := Parse("", []byte(data))
	require.NoError(t, err)
	str, ok := res.(string)
	require.True(t, ok)
	require.Equal(t, str, "(((7+((10/2)*(-4)))+((5*3)%6))-(8*6))")
}
