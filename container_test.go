package go_container

import (
	berr "errors"
	"github.com/nikonm/go-container/errors"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

type mockService1 struct{}

func New1(s *mockService2) (*mockService1, error) {
	return &mockService1{}, nil
}

type mockService2 struct{}

func New2(s *mockService3) (*mockService2, error) {
	return &mockService2{}, nil
}

type mockService3 struct{}

func New3() (*mockService3, error) {
	return &mockService3{}, nil
}

func TestNew(t *testing.T) {
	t.Parallel()
	c, err := New(&Options{
		BasePkg: "go_container",
		Services: []interface{}{
			New2,
			New3,
			New1,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	require.NotNil(t, c)
	s, err := c.GetService("mockService3")
	_, ok := s.(*mockService3)
	require.True(t, ok)
}

type mockService4 struct{}

func New4(s *mockService5) (*mockService4, error) {
	return &mockService4{}, nil
}

type mockService5 struct{}

func New5(s *mockService4) (*mockService5, error) {
	return &mockService5{}, nil
}

func TestCircleError(t *testing.T) {
	t.Parallel()
	_, err := New(&Options{
		BasePkg: "go_container",
		Services: []interface{}{
			New4,
			New5,
		},
	})
	require.Equal(t, errors.CircleLoop.GetMsg(), err.(errors.Error).GetMsg())
}

func TestErrorSign(t *testing.T) {
	t.Parallel()
	f := func() {}
	_, err := New(&Options{
		BasePkg: "go_container",
		Services: []interface{}{
			f,
		},
	})
	require.Equal(t, errors.InvalidInitSign.GetMsg(), err.(errors.Error).GetMsg())
}

func TestErrorServiceNotFound(t *testing.T) {
	t.Parallel()
	c, _ := New(&Options{})
	_, err := c.GetService("service")
	require.Equal(t, errors.ServiceNotFound.GetMsg(), err.(errors.Error).GetMsg())
}

func TestErrorInvalidOutArgumentsCount(t *testing.T) {
	t.Parallel()
	f := func() interface{} {
		return ""
	}
	_, err := New(&Options{
		BasePkg: "go_container",
		Services: []interface{}{
			f,
		},
	})
	require.Equal(t, errors.InvalidOutArgumentsCount.GetMsg(), err.(errors.Error).GetMsg())
}

func TestErrorInvalidOutSign(t *testing.T) {
	t.Parallel()
	f := func() (interface{}, int) {
		return "", 0
	}
	_, err := New(&Options{
		BasePkg: "go_container",
		Services: []interface{}{
			f,
		},
	})
	require.Equal(t, errors.InvalidOutSign.GetMsg(), err.(errors.Error).GetMsg())
}

func TestErrorNotResolveDep(t *testing.T) {
	t.Parallel()
	f := func(s mockService1) (interface{}, error) {
		return "", nil
	}
	_, err := New(&Options{
		BasePkg: "go_container",
		Services: []interface{}{
			f,
		},
	})
	require.Equal(t, errors.NotResolveDep.GetMsg(), err.(errors.Error).GetMsg())
}

func TestErrorServiceInitError(t *testing.T) {
	t.Parallel()
	f := func() (interface{}, error) {
		return "", berr.New("test")
	}
	_, err := New(&Options{
		BasePkg: "go_container",
		Services: []interface{}{
			f,
		},
	})
	require.True(t, strings.HasPrefix(err.(errors.Error).GetMsg(), "service init error"))
}

type mockService6 struct{}

func (s mockService6) Stop() error { return berr.New("test") }
func New6() (*mockService6, error) {
	return &mockService6{}, nil
}

type mockService7 struct{}

func (s *mockService7) Stop() error { return nil }
func New7(s *mockService6) (*mockService7, error) {
	return &mockService7{}, nil
}

func TestStopper(t *testing.T) {
	t.Parallel()
	c, err := New(&Options{
		BasePkg: "go_container",
		Services: []interface{}{
			New6,
			New7,
		},
	})
	require.NoError(t, err)
	err = c.StopAll()
	require.EqualError(t, err, "test")
}