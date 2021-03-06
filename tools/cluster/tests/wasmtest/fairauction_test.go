// +build ignore

package wasmtest

import (
	"fmt"
	"github.com/iotaledger/wasp/client/chainclient"
	"github.com/iotaledger/wasp/packages/coretypes"
	"testing"
	"time"

	"github.com/iotaledger/wasp/packages/subscribe"
	"github.com/stretchr/testify/assert"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	waspapi "github.com/iotaledger/wasp/packages/apilib"
	"github.com/iotaledger/wasp/packages/testutil"
	"github.com/iotaledger/wasp/packages/vm/examples/fairauction"
	"github.com/iotaledger/wasp/packages/vm/examples/fairauction/faclient"
	"github.com/iotaledger/wasp/packages/vm/examples/tokenregistry"
	"github.com/iotaledger/wasp/packages/vm/examples/tokenregistry/trclient"
)

const faWasmPath = "wasm/fairauction"
const faDescription = "Fair auction, a PoC smart contract"

func TestLoadTrAndFaAndThenRunTrMint(t *testing.T) {
	wasps := setup(t, "TestLoadTrAndFaAndThenRunTrMint")

	err := loadWasmIntoWasps(wasps, trWasmPath, trDescription)
	check(err, t)
	trProgramHash := programHash

	err = loadWasmIntoWasps(wasps, faWasmPath, faDescription)
	check(err, t)
	faProgramHash := programHash

	err = requestFunds(wasps, scOwnerAddr, "sc owner")
	check(err, t)

	auctionOwner := wallet.WithIndex(1)
	auctionOwnerAddr := auctionOwner.Address()
	err = requestFunds(wasps, auctionOwnerAddr, "auction owner")
	check(err, t)

	programHash = trProgramHash
	scTRChain, scTRAddr, scTRColor, err := startSmartContract(wasps, tokenregistry.ProgramHash, trDescription)
	checkSuccess(err, t, "TokenRegistry has been created and activated")

	programHash = faProgramHash
	scFAChain, scFAAddr, scFAColor, err := startSmartContract(wasps, fairauction.ProgramHash, faDescription)
	checkSuccess(err, t, "FairAuction has been created and activated")
	_ = scFAChain

	chidTR := (coretypes.ChainID)(*scTRAddr)
	succ := waspapi.CheckDeployment(wasps.ApiHosts(), &chidTR)
	assert.True(t, succ)

	chidFA := (coretypes.ChainID)(*scFAAddr)
	succ = waspapi.CheckDeployment(wasps.ApiHosts(), &chidFA)
	assert.True(t, succ)

	tc := trclient.NewClient(chainclient.New(
		wasps.Level1Client,
		wasps.WaspClient(0),
		scTRChain,
		auctionOwner.SigScheme(),
		20*time.Second,
	))

	// minting 1 token with TokenRegistry
	tx, err := tc.MintAndRegister(trclient.MintAndRegisterParams{
		Supply:      1,
		MintTarget:  *auctionOwnerAddr,
		Description: "Non-fungible coin 1. Very expensive",
	})
	checkSuccess(err, t, "token minted")

	mintedColor := balance.Color(tx.ID())

	if !wasps.VerifyAddressBalances(scFAAddr, 1, map[balance.Color]int64{
		*scFAColor: 1, // sc token
	}, "SC FairAuction address in the beginning") {
		t.Fail()
		return
	}
	if !wasps.VerifyAddressBalances(scTRAddr, 1, map[balance.Color]int64{
		*scTRColor: 1, // sc token
	}, "SC TokenRegistry address in the beginning") {
		t.Fail()
		return
	}
	if !wasps.VerifyAddressBalances(scOwnerAddr, testutil.RequestFundsAmount-2, map[balance.Color]int64{
		balance.ColorIOTA: testutil.RequestFundsAmount - 2,
	}, "SC owner in the beginning") {
		t.Fail()
		return
	}
	if !wasps.VerifyAddressBalances(auctionOwnerAddr, testutil.RequestFundsAmount, map[balance.Color]int64{
		balance.ColorIOTA: testutil.RequestFundsAmount - 1,
		mintedColor:       1,
	}, "Auction owner in the beginning") {
		t.Fail()
		return
	}
}

// scenario with 2 smart contracts
func TestTrMintAndFaAuctionWith2Bids(t *testing.T) {
	wasps := setup(t, "TestTrMintAndFaAuctionWith2Bids")

	err := loadWasmIntoWasps(wasps, trWasmPath, trDescription)
	check(err, t)
	trProgramHash := programHash

	err = loadWasmIntoWasps(wasps, faWasmPath, faDescription)
	check(err, t)
	faProgramHash := programHash

	err = requestFunds(wasps, scOwnerAddr, "sc owner")
	check(err, t)

	auctionOwner := wallet.WithIndex(1)
	auctionOwnerAddr := auctionOwner.Address()
	err = requestFunds(wasps, auctionOwnerAddr, "auction owner")
	check(err, t)

	bidder1 := wallet.WithIndex(2)
	bidder1Addr := bidder1.Address()
	err = requestFunds(wasps, bidder1Addr, "bidder 1")
	check(err, t)

	bidder2 := wallet.WithIndex(3)
	bidder2Addr := bidder2.Address()
	err = requestFunds(wasps, bidder2Addr, "bidder 2")
	check(err, t)

	programHash = trProgramHash
	scTRChain, scTRAddr, scTRColor, err := startSmartContract(wasps, tokenregistry.ProgramHash, trDescription)
	checkSuccess(err, t, "TokenRegistry has been created and activated")

	programHash = faProgramHash
	scFAChain, scFAAddr, scFAColor, err := startSmartContract(wasps, fairauction.ProgramHash, faDescription)
	checkSuccess(err, t, "FairAuction has been created and activated")

	tc := trclient.NewClient(chainclient.New(
		wasps.Level1Client,
		wasps.WaspClient(0),
		scTRChain,
		auctionOwner.SigScheme(),
		20*time.Second,
	))

	// minting 1 token with TokenRegistry
	tx, err := tc.MintAndRegister(trclient.MintAndRegisterParams{
		Supply:      1,
		MintTarget:  *auctionOwnerAddr,
		Description: "Non-fungible coin 1. Very expensive",
	})
	checkSuccess(err, t, "token minted")

	mintedColor := balance.Color(tx.ID())

	if !wasps.VerifyAddressBalances(scFAAddr, 1, map[balance.Color]int64{
		*scFAColor: 1, // sc token
	}, "SC FairAuction address in the beginning") {
		t.Fail()
		return
	}
	if !wasps.VerifyAddressBalances(scTRAddr, 1, map[balance.Color]int64{
		*scTRColor: 1, // sc token
	}, "SC TokenRegistry address in the beginning") {
		t.Fail()
		return
	}
	if !wasps.VerifyAddressBalances(scOwnerAddr, testutil.RequestFundsAmount-2, map[balance.Color]int64{
		balance.ColorIOTA: testutil.RequestFundsAmount - 2,
	}, "SC owner in the beginning") {
		t.Fail()
		return
	}
	if !wasps.VerifyAddressBalances(auctionOwnerAddr, testutil.RequestFundsAmount, map[balance.Color]int64{
		balance.ColorIOTA: testutil.RequestFundsAmount - 1,
		mintedColor:       1,
	}, "Auction owner in the beginning") {
		t.Fail()
		return
	}
	if !wasps.VerifyAddressBalances(bidder1Addr, testutil.RequestFundsAmount, map[balance.Color]int64{
		balance.ColorIOTA: testutil.RequestFundsAmount,
	}, "Bidder1 in the beginning") {
		t.Fail()
		return
	}
	if !wasps.VerifyAddressBalances(bidder2Addr, testutil.RequestFundsAmount, map[balance.Color]int64{
		balance.ColorIOTA: testutil.RequestFundsAmount,
	}, "Bidder2 in the beginning") {
		t.Fail()
		return
	}

	faclientOwner := faclient.NewClient(chainclient.New(
		wasps.Level1Client,
		wasps.WaspClient(0),
		scFAChain,
		auctionOwner.SigScheme(),
		20*time.Second,
	), 0)

	_, err = faclientOwner.StartAuction("selling my only token", &mintedColor, 1, 100, 1)
	checkSuccess(err, t, "StartAuction created")

	faclientBidder1 := faclient.NewClient(chainclient.New(wasps.Level1Client, wasps.WaspClient(0), scFAChain, bidder1.SigScheme()), 0)
	faclientBidder2 := faclient.NewClient(chainclient.New(wasps.Level1Client, wasps.WaspClient(0), scFAChain, bidder2.SigScheme()), 0)

	subs, err := subscribe.SubscribeMulti(wasps.PublisherHosts(), "request_out")
	check(err, t)

	tx1, err := faclientBidder1.PlaceBid(&mintedColor, 110)
	check(err, t)
	tx2, err := faclientBidder2.PlaceBid(&mintedColor, 110)
	check(err, t)

	patterns := [][]string{
		{"request_out", scFAAddr.String(), tx1.ID().String(), "0"},
		{"request_out", scFAAddr.String(), tx2.ID().String(), "0"},
	}
	err = nil
	if !subs.WaitForPatterns(patterns, 40*time.Second) {
		err = fmt.Errorf("didn't receive completion message in time")
	}
	checkSuccess(err, t, "2 bids have been placed")

	// wait for auction to finish
	time.Sleep(70 * time.Second)

	if !wasps.VerifyAddressBalances(scFAAddr, 2, map[balance.Color]int64{
		*scFAColor:        1, // sc token
		balance.ColorIOTA: 1, // 1 i for sending request to itself
	}, "SC FairAuction address in the end") {
		t.Fail()
	}
	if !wasps.VerifyAddressBalances(scTRAddr, 1, map[balance.Color]int64{
		*scTRColor: 1, // sc token
	}, "SC TokenRegistry address in the end") {
		t.Fail()
	}
	if !wasps.VerifyAddressBalances(scOwnerAddr, testutil.RequestFundsAmount-2+4, map[balance.Color]int64{
		balance.ColorIOTA: testutil.RequestFundsAmount - 2 + 4,
	}, "SC owner in the end") {
		t.Fail()
	}
	if !wasps.VerifyAddressBalances(auctionOwnerAddr, testutil.RequestFundsAmount-2+110-4, map[balance.Color]int64{
		balance.ColorIOTA: testutil.RequestFundsAmount - 2 + 110 - 4,
	}, "Auction owner in the end") {
		t.Fail()
	}
	if !wasps.VerifyAddressBalances(bidder1Addr, testutil.RequestFundsAmount-110+1, map[balance.Color]int64{
		balance.ColorIOTA: testutil.RequestFundsAmount - 110,
		mintedColor:       1,
	}, "Bidder1 in the end") {
		t.Fail()
	}
	if !wasps.VerifyAddressBalances(bidder2Addr, testutil.RequestFundsAmount, map[balance.Color]int64{
		balance.ColorIOTA: testutil.RequestFundsAmount,
	}, "Bidder2 in the end") {
		t.Fail()
	}
}
