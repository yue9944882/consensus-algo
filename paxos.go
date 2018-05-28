package consensus_algo

import (
	"time"
)

const (
	paxosMessageTypePrepare = "prepare"
	paxosMessageTypePromise = "promise"
	paxosMessageTypePropose = "propose"
	paxosMessageTypeAccept  = "accept"
)

type paxosMessage struct {
	msgType  string
	ok       bool
	promiseN int
	acceptN  int
	acceptV  interface{}
}

type acceptor struct {
	*peer
	acceptedV interface{}
	promisedN int
	acceptedN int
}

func (a *acceptor) promise(msg *paxosMessage) (*paxosMessage, error) {
	promiseMsg := &paxosMessage{
		msgType: paxosMessageTypePromise,
		acceptN: a.acceptedN,
		acceptV: a.acceptedV,
	}
	if a.promisedN < msg.promiseN {
		promiseMsg.ok = true
	}
	return promiseMsg, nil
}

func (a *acceptor) accept(msg *paxosMessage) (*paxosMessage, error) {
	acceptMsg := &paxosMessage{
		msgType: paxosMessageTypeAccept,
	}
	if msg.acceptN == a.promisedN {
		a.acceptedV = msg.acceptV
		msg.ok = true
	}
	return acceptMsg, nil
}

type proposer struct {
	*peer
	acceptorsPromiseN map[int]int

	n int
}

func (p *proposer) prepare() {
	majority := len(p.acceptorsPromiseN)/2 + 1
	promisedCount := 0
	for {
		for id := range p.acceptorsPromiseN {
			p.send(id, &paxosMessage{
				msgType:  paxosMessageTypePrepare,
				promiseN: p.n,
			})
		}
		recvCount := 0
		for promisedCount < majority && recvCount < len(p.acceptorsPromiseN) {
			id, msg, ok := p.recv(time.Second)
			if ok {
				pmsg := msg.(*paxosMessage)
				if pmsg.ok {
					p.acceptorsPromiseN[id] = pmsg.promiseN
					promisedCount++
				} else if p.acceptorsPromiseN[id] < pmsg.acceptN {
					p.acceptorsPromiseN[id] = pmsg.acceptN
				}
			}
		}
		// TODO: propose
	}
}

func (p *proposer) propose() {

}

type learner struct {
	*peer
}
