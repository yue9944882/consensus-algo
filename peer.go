package consensus_algo

import "time"

type peer struct {
	id int
	nt *network
}

func newPeer(id int, nt *network) *peer {
	return &peer{id, nt}
}

func (p *peer) send(id int, msg interface{}) {
	p.nt.traffic[id] <- packet{
		from: p.id,
		to:   id,
		msg:  msg,
	}
}

func (p *peer) recv(timeout time.Duration) (id int, msg interface{}, ok bool) {
	pkt, valid := <-p.nt.traffic[id]
	if valid && pkt.to == p.id {
		id = pkt.from
		msg = pkt.msg
		ok = valid
	}
	return
}
