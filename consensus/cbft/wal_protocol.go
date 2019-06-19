package cbft

const (
	SendPrepareBlockMsg    = 0x64
	SendViewChangeMsg      = 0x65
	ConfirmedViewChangeMsg = 0x66
)

type sendPrepareBlock struct {
	PrepareBlock *prepareBlock
}

type sendViewChange struct {
	ViewChange *viewChange
	//viewChangeVotes ViewChangeVotes
	Master bool
}

type confirmedViewChange struct {
	ViewChange      *viewChange
	ViewChangeResp  *viewChangeVote
	ViewChangeVotes []*viewChangeVote
	Master          bool
}

var (
	wal_messages = []interface{}{
		prepareBlock{},
		prepareVote{},
		viewChange{},
		viewChangeVote{},
		confirmedPrepareBlock{},
		getPrepareVote{},
		prepareVotes{},
		getPrepareBlock{},
		getHighestPrepareBlock{},
		highestPrepareBlock{},
		cbftStatusData{},
		prepareBlockHash{},
		sendPrepareBlock{},
		sendViewChange{},
		confirmedViewChange{},
	}
)

func WalMessageType(msg interface{}) uint64 {
	switch msg.(type) {
	case *sendPrepareBlock:
		return SendPrepareBlockMsg
	case *sendViewChange:
		return SendViewChangeMsg
	case *confirmedViewChange:
		return ConfirmedViewChangeMsg
	}
	return MessageType(msg)
}
