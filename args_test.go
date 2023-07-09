package clear

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type getter interface {
	Get() any
}

func TestArgs(t *testing.T) {

	captured := []float64{}
	eachAction := func(index int, f float64) error {
		require.Equal(t, len(captured), index)
		captured = append(captured, f)
		return nil
	}

	arg := FloatArg("simple", EachValAction(eachAction))

	tt := NewTokenizer()
	require.NoError(t, arg.Consume(tt))
	require.Equal(t, []float64{}, arg.(getter).Get())

	tt = NewTokenizer("10.1", "-111")
	require.NoError(t, arg.Consume(tt))
	require.Equal(t, []float64{10.1, -111}, arg.(getter).Get())
	require.Equal(t, []float64{10.1, -111}, captured)

	arg = FloatArg("simple", Min[float64](2), Max[float64](5))

	tt = NewTokenizer()
	require.Error(t, arg.Consume(tt))
	require.Equal(t, []float64{}, arg.(getter).Get())

	tt = NewTokenizer("10.1", "-111")
	require.NoError(t, arg.Consume(tt))
	require.Equal(t, []float64{10.1, -111}, arg.(getter).Get())

	tt = NewTokenizer("10.1", "-111", "invalid")
	require.Error(t, arg.Consume(tt))
	require.Equal(t, []float64{10.1, -111}, arg.(getter).Get())

	tt = NewTokenizer("10.1", "-111", "0.1", "11", "233", "-098")
	require.NoError(t, arg.Consume(tt))
	require.Equal(t, []float64{10.1, -111, 0.1, 11, 233}, arg.(getter).Get())
}
