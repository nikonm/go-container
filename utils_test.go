package go_container

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type test struct{}

func TestGetPkgPath(t *testing.T) {
	s := GetPkgPath(&test{})
	require.Equal(t, "go_container.test", s)
}
