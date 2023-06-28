package builtin

import (
	"testing"

	"github.com/CluEleSsUK/go-state-types/abi"
	"github.com/stretchr/testify/require"
)

func TestGenerateMethodNum(t *testing.T) {

	methodNum, err := GenerateFRCMethodNum("Receive")
	require.NoError(t, err)
	require.Equal(t, methodNum, abi.MethodNum(3726118371))
}
