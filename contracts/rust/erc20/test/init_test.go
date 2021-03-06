package erc20

import (
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	erc20name = "erc20test"
	erc20file = "../pkg/erc20_bg.wasm"
)

func TestDeployErc20(t *testing.T) {
	env := solo.New(t, false, false)
	chain := env.NewChain(nil, "chain1")

	creator := env.NewSignatureSchemeWithFunds()
	creatorAgentID := coretypes.NewAgentIDFromAddress(creator.Address())
	err := chain.DeployWasmContract(nil, erc20name, erc20file,
		PARAM_SUPPLY, 1000000,
		PARAM_CREATOR, creatorAgentID,
	)
	require.NoError(t, err)
	_, rec := chain.GetInfo()
	require.EqualValues(t, 5, len(rec))

	_, err = chain.FindContract(erc20name)
	require.NoError(t, err)

	// deploy second time
	err = chain.DeployWasmContract(nil, erc20name, erc20file,
		PARAM_SUPPLY, 1000000,
		PARAM_CREATOR, creatorAgentID,
	)
	require.Error(t, err)
	_, rec = chain.GetInfo()
	require.EqualValues(t, 5, len(rec))
}

func TestDeployErc20Fail1(t *testing.T) {
	env := solo.New(t, false, false)
	chain := env.NewChain(nil, "chain1")
	err := chain.DeployWasmContract(nil, erc20name, erc20file)
	require.Error(t, err)
	_, rec := chain.GetInfo()
	require.EqualValues(t, 4, len(rec))
}

func TestDeployErc20Fail2(t *testing.T) {
	env := solo.New(t, false, false)
	chain := env.NewChain(nil, "chain1")
	err := chain.DeployWasmContract(nil, erc20name, erc20file,
		PARAM_SUPPLY, 1000000,
	)
	require.Error(t, err)
	_, rec := chain.GetInfo()
	require.EqualValues(t, 4, len(rec))
}

func TestDeployErc20Fail3(t *testing.T) {
	env := solo.New(t, false, false)
	chain := env.NewChain(nil, "chain1")
	creator := env.NewSignatureSchemeWithFunds()
	creatorAgentID := coretypes.NewAgentIDFromAddress(creator.Address())
	err := chain.DeployWasmContract(nil, erc20name, erc20file,
		PARAM_CREATOR, creatorAgentID,
	)
	require.Error(t, err)
	_, rec := chain.GetInfo()
	require.EqualValues(t, 4, len(rec))
}

func TestDeployErc20Fail3Repeat(t *testing.T) {
	env := solo.New(t, false, false)
	chain := env.NewChain(nil, "chain1")
	creator := env.NewSignatureSchemeWithFunds()
	creatorAgentID := coretypes.NewAgentIDFromAddress(creator.Address())
	err := chain.DeployWasmContract(nil, erc20name, erc20file,
		PARAM_CREATOR, creatorAgentID,
	)
	require.Error(t, err)
	_, rec := chain.GetInfo()
	require.EqualValues(t, 4, len(rec))

	// repeat after failure
	err = chain.DeployWasmContract(nil, erc20name, erc20file,
		PARAM_SUPPLY, 1000000,
		PARAM_CREATOR, creatorAgentID,
	)
	require.NoError(t, err)
	_, rec = chain.GetInfo()
	require.EqualValues(t, 5, len(rec))

	_, err = chain.FindContract(erc20name)
	require.NoError(t, err)
}
