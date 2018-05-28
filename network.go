package consensus_algo

type network struct {
	traffic map[int]chan packet
}

type packet struct {
	from int
	to   int
	msg  interface{}
}
