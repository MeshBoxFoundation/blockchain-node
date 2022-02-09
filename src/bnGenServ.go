package src

import (
	"fmt"
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
)

// gen.Server implementation structure
type BnGenServ struct {
	gen.Server
}

func (dgs *BnGenServ) HandleCast(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	fmt.Printf("HandleCast (%s): %v\n", process.Name(), message)
	switch message {
	case etf.Atom("stop"):
		return gen.ServerStatusStopWithReason("stop they said")
	}
	return gen.ServerStatusOK
}

func (dgs *BnGenServ) HandleCall(process *gen.ServerProcess, from gen.ServerFrom, message etf.Term) (etf.Term, gen.ServerStatus) {
	fmt.Printf("HandleCall (%s): %v, From: %v\n", process.Name(), message, from)

	switch message {
	case etf.Atom("hello"):
		return etf.Atom("hi"), gen.ServerStatusOK
	}

	reply := etf.Tuple{etf.Atom("error"), etf.Atom("unknown_request")}
	return reply, gen.ServerStatusOK
}
