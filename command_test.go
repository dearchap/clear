package clear

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommand(t *testing.T) {

	cmd := NewCmd("foo")

	tzer := NewTokenizer("bar", "2")

	require.NoError(t, cmd.Run(context.Background(), tzer))

	tzer = NewTokenizer("bar", "2")

	cmd = NewCmd("foo", CmdArgs(IntArg("f")))

	require.Error(t, cmd.Run(context.Background(), tzer))
	tzer = NewTokenizer("2", "10")
	require.NoError(t, cmd.Run(context.Background(), tzer))

	cmd = NewCmd("foo", CmdArgs(IntArg("-f"), FloatArg("bar")))

	require.NoError(t, cmd.Run(context.Background(), NewTokenizer()))
}
