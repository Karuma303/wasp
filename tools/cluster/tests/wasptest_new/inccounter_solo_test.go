package wasptest

import (
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/stretchr/testify/require"
	"testing"
)

var incFile = util.LocateFile("inccounter_bg.wasm", "contracts/wasm")

const incName = "inccounter"
const incDescription = "IncCounter, a PoC smart contract"

var incHname = coretypes.Hn(incName)

const varCounter = "counter"
const varNumRepeats = "numRepeats"

func TestIncSoloInc(t *testing.T) {
	al := solo.New(t, false, false)
	chain := al.NewChain(nil, "chain1")
	err := chain.DeployWasmContract(nil, incName, incFile)
	require.NoError(t, err)
	req := solo.NewCallParams(incName, "increment").
		WithTransfer(balance.ColorIOTA, 1)
	_, err = chain.PostRequest(req, nil)
	require.NoError(t, err)
	ret, err := chain.CallView(incName, "getCounter")
	require.NoError(t, err)
	counter, _, err := codec.DecodeInt64(ret.MustGet(varCounter))
	require.NoError(t, err)
	require.EqualValues(t, 1, counter)
}

func TestIncSoloRepeatMany(t *testing.T) {
	al := solo.New(t, false, false)
	chain := al.NewChain(nil, "chain1")
	err := chain.DeployWasmContract(nil, incName, incFile)
	require.NoError(t, err)
	req := solo.NewCallParams(incName, "repeatMany", varNumRepeats, 2).
		WithTransfer(balance.ColorIOTA, 1)
	_, err = chain.PostRequest(req, nil)
	require.NoError(t, err)
	chain.WaitForEmptyBacklog()
	ret, err := chain.CallView(incName, "getCounter")
	require.NoError(t, err)
	counter, _, err := codec.DecodeInt64(ret.MustGet(varCounter))
	require.NoError(t, err)
	require.EqualValues(t, 3, counter)
}

func TestIncSoloResultsTest(t *testing.T) {
	al := solo.New(t, false, false)
	chain := al.NewChain(nil, "chain1")
	err := chain.DeployWasmContract(nil, incName, incFile)
	require.NoError(t, err)
	req := solo.NewCallParams(incName, "resultsTest").
		WithTransfer(balance.ColorIOTA, 1)
	ret, err := chain.PostRequest(req, nil)
	require.NoError(t, err)
	//ret, err = chain.CallView(incName, "results_check")
	//require.NoError(t, err)
	require.EqualValues(t, 8, len(ret))
}

func TestIncSoloStateTest(t *testing.T) {
	al := solo.New(t, false, false)
	chain := al.NewChain(nil, "chain1")
	err := chain.DeployWasmContract(nil, incName, incFile)
	require.NoError(t, err)
	req := solo.NewCallParams(incName, "stateTest").
		WithTransfer(balance.ColorIOTA, 1)
	ret, err := chain.PostRequest(req, nil)
	require.NoError(t, err)
	ret, err = chain.CallView(incName, "stateCheck")
	require.NoError(t, err)
	require.EqualValues(t, 0, len(ret))
}
