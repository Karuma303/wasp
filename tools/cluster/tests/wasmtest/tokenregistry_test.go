// +build ignore

package wasmtest

import (
	"fmt"
	"github.com/iotaledger/wasp/client/chainclient"
	"github.com/iotaledger/wasp/packages/coretypes"
	"testing"
	"time"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/testutil"
	"github.com/iotaledger/wasp/packages/vm/examples/tokenregistry"
	"github.com/iotaledger/wasp/packages/vm/examples/tokenregistry/trclient"
	"github.com/iotaledger/wasp/packages/vm/vmconst"
)

const trWasmPath = "wasm/tokenregistry"
const trDescription = "Token registry, a PoC smart contract"

func TestTrMintSupply(t *testing.T) {
	wasps := setup(t, "TestTrMintSupply")

	err := loadWasmIntoWasps(wasps, trWasmPath, trDescription)
	check(err, t)

	err = requestFunds(wasps, scOwnerAddr, "sc owner")
	check(err, t)

	minter := wallet.WithIndex(1)
	minterAddr := minter.Address()
	err = requestFunds(wasps, minterAddr, "minter")
	check(err, t)

	scChain, scAddr, scColor, err := startSmartContract(wasps, tokenregistry.ProgramHash, trDescription)
	checkSuccess(err, t, "smart contract has been created and activated")

	if !wasps.VerifyAddressBalances(scAddr, 1, map[balance.Color]int64{
		*scColor: 1, // sc token
	}, "SC address in the beginning") {
		t.Fail()
		return
	}
	if !wasps.VerifyAddressBalances(scOwnerAddr, testutil.RequestFundsAmount-1, map[balance.Color]int64{
		balance.ColorIOTA: testutil.RequestFundsAmount - 1,
	}, "owner in the beginning") {
		t.Fail()
		return
	}

	tc := trclient.NewClient(chainclient.New(
		wasps.Level1Client,
		wasps.WaspClient(0),
		scChain,
		minter.SigScheme(),
		15*time.Second,
	))

	tx1, err := tc.MintAndRegister(trclient.MintAndRegisterParams{
		Supply:      1,
		MintTarget:  *minterAddr,
		Description: "Non-fungible coin 1",
	})
	checkSuccess(err, t, "token minted and registered successfully")

	for {
		// the sleep 1 second is usually enough
		time.Sleep(time.Second)
		reqId := coretypes.NewRequestID(tx1.ID(), 0)
		r, err := wasps.WaspClient(0).RequestStatus(scAddr, &reqId)
		check(err, t)
		if r.IsProcessed {
			break
		}
		fmt.Println("Busy waiting for transaction to be processed")
	}

	mintedColor1 := balance.Color(tx1.ID())

	if !wasps.VerifyAddressBalances(scAddr, 1, map[balance.Color]int64{
		balance.ColorIOTA: 0,
		*scColor:          1,
	}, "SC address in the end") {
		t.Fail()
	}

	if !wasps.VerifyAddressBalances(minterAddr, testutil.RequestFundsAmount, map[balance.Color]int64{
		mintedColor1:      1,
		balance.ColorIOTA: testutil.RequestFundsAmount - 1,
	}, "minter1 in the end") {
		t.Fail()
		return
	}

	if !wasps.VerifySCStateVariables2(scAddr, map[kv.Key]interface{}{
		vmconst.VarNameOwnerAddress:      scOwnerAddr[:],
		vmconst.VarNameProgramData:       programHash[:],
		tokenregistry.VarStateListColors: []byte(mintedColor1.String()),
		vmconst.VarNameDescription:       trDescription,
	}) {
		t.Fail()
	}
}
