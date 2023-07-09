package clear

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTokenizer(t *testing.T) {

	input := []string{"hello", "foo", "-t", "10"}
	tt := NewTokenizer(input...)

	got := []string{}
	peeks := []string{}

	for tt.HasNext() {
		peeks = append(peeks, tt.Peek())
		got = append(got, tt.Next())
	}

	require.Equal(t, input, got)
	require.Equal(t, input, peeks)
}

func TestTokenizerChain(t *testing.T) {

	input := []string{"hello", "foo", "-t", "-f", "-o"}
	tt1 := NewTokenizer(input[0:2]...)
	tt2 := NewTokenizer(input[2:]...)

	tt := NewTokenizerChain(tt1, tt2)

	got := []string{}
	peeks := []string{}

	for tt.HasNext() {
		peeks = append(peeks, tt.Peek())
		got = append(got, tt.Next())
	}

	require.Equal(t, input, got)
	require.Equal(t, input, peeks)
}
