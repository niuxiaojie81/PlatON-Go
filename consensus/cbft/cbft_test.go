package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/stretchr/testify/assert"
	"math/big"
	"os"
	"testing"
	"time"
)

func TestNewViewChange(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, backend, validators := randomCBFT(path, 4)

	priA := validators.neibor[0]
	//addrA := crypto.PubkeyToAddress(*priA.publicKey)

	gen := backend.chain.Genesis()

	var blocks []*types.Block
	blocks = append(blocks, gen)
	for i := uint64(1); i < 10; i++ {
		blocks = append(blocks, createBlock(priA.privateKey, blocks[i-1].Hash(), blocks[i-1].NumberU64()+1))
		t.Log(blocks[i].NumberU64(), blocks[i].Hash().TerminalString(), blocks[i].ParentHash().TerminalString())
	}

	node := nodeIndexNow(validators, engine.startTimeOfEpoch)
	viewChange := makeViewChange(node.privateKey, uint64(time.Now().UnixNano()/1e6), 0, gen.Hash(), uint32(node.index), node.address, nil)

	err := engine.OnViewChange(node.nodeID, viewChange)
	assert.Nil(t, err)

}

//voteA := makeViewChangeVote(priA.privateKey, 0, 5, common.BytesToHash([]byte{1}), 0, addrA, uint32(2), addrA)
//
//err := engine.OnViewChangeVote(priA.nodeID, voteA)
//t.Log(err)

func TestNewPrepareBlock(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, backend, validators := randomCBFT(path, 4)
	owner := validators.owner

	gen := backend.chain.Genesis()
	block := createBlock(owner.privateKey, gen.Hash(), gen.NumberU64()+1)

	// test Cache
	p := buildPrepareBlock(block, owner, nil, nil)
	if err := engine.OnNewPrepareBlock(owner.nodeID, p, true); err != nil {
		t.Fatalf("test Cache error: %v", err)
	}

	view, err := engine.newViewChange()
	viewChangeVotes := buildViewChangeVote(view, validators)

	// test errFutileBlock
	if err != nil {
		t.Fatalf("newViewChange error: %v", err)
	}
	t.Log(view.BaseBlockNum, view.BaseBlockHash.Hex(), view.ProposalIndex, view.ProposalAddr, view.Timestamp)
	if err := engine.OnNewPrepareBlock(owner.nodeID, p, true); err != errFutileBlock {
		t.Fatalf("test errFutileBlock error: %v", err)
	}

	// test VerifyHeader
	header := &types.Header{Number: big.NewInt(int64(gen.NumberU64() + 1)), ParentHash: gen.Hash()}
	sign, _ := crypto.Sign(header.SealHash().Bytes(), owner.privateKey)
	header.Extra = make([]byte, 32)
	copy(header.Extra, sign[0:32])
	block = types.NewBlockWithHeader(header)
	p = buildPrepareBlock(block, owner, view, viewChangeVotes)
	if err := engine.OnNewPrepareBlock(owner.nodeID, p, true); err != errMissingSignature {
		t.Fatalf("test VerifyHeader signature error: %v", err)
	}

	// test errInvalidatorCandidateAddress
	header = &types.Header{Number: big.NewInt(int64(gen.NumberU64() + 1)), ParentHash: gen.Hash()}
	sign, _ = crypto.Sign(header.SealHash().Bytes(), owner.privateKey)
	header.Extra = make([]byte, 32+65)
	copy(header.Extra, sign)
	block = types.NewBlockWithHeader(header)
	p = buildPrepareBlock(block, owner, view, viewChangeVotes)
	p.ProposalAddr = common.HexToAddress("0x27f7e1d4b9caab9d5b13803cff6da714c51de34e")
	if err := engine.OnNewPrepareBlock(owner.nodeID, p, true); err != errInvalidatorCandidateAddress {
		t.Fatalf("test errInvalidatorCandidateAddress error: %v", err)
	}

	// test errTwoThirdViewchangeVotes
	p.ProposalAddr = owner.address
	p.ViewChangeVotes = p.ViewChangeVotes[0:1]
	if err := engine.OnNewPrepareBlock(owner.nodeID, p, true); err != errTwoThirdViewchangeVotes {
		t.Fatalf("test errTwoThirdViewchangeVotes error: %v", err)
	}

	// test errInvalidViewChangeVote
	p.ViewChangeVotes = viewChangeVotes
	if err := engine.OnNewPrepareBlock(owner.nodeID, p, true); err != errTwoThirdViewchangeVotes {
		t.Fatalf("test errTwoThirdViewchangeVotes error: %v", err)
	}
}
