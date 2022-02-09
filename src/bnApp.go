package main

import (
	"flag"
	"fmt"

	"github.com/ergo-services/ergo"
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
	"github.com/ergo-services/ergo/node"
)

var (
	NodeName         string
	Cookie           string
	err              error
	ListenRangeBegin int
	ListenRangeEnd   int = 35000
	Listen           string
	ListenEPMD       int

	EnableRPC bool
)

type bnApp struct {
	gen.Application
}

func (da *bnApp) Load(args ...etf.Term) (gen.ApplicationSpec, error) {
	return gen.ApplicationSpec{
		Name:        "bnApp",
		Description: "bn Applicatoin",
		Version:     "v.1.0",
		Environment: map[string]interface{}{
			"envName1": 123,
			"envName2": "Hello world",
		},
		Children: []gen.ApplicationChildSpec{
			gen.ApplicationChildSpec{
				Child: &bnSup{},
				Name:  "bnSup",
			},
			gen.ApplicationChildSpec{
				Child: &bnGenServ{},
				Name:  "justbnGS",
			},
		},
	}, nil
}

func (da *bnApp) Start(process gen.Process, args ...etf.Term) {
	fmt.Println("Application started!")
}

type bnSup struct {
	gen.Supervisor
}

func (ds *bnSup) Init(args ...etf.Term) (gen.SupervisorSpec, error) {
	spec := gen.SupervisorSpec{
		Name: "bnAppSup",
		Children: []gen.SupervisorChildSpec{
			gen.SupervisorChildSpec{
				Name:  "bnServer01",
				Child: &bnGenServ{},
			},
			gen.SupervisorChildSpec{
				Name:  "bnServer02",
				Child: &bnGenServ{},
				Args:  []etf.Term{12345},
			},
			gen.SupervisorChildSpec{
				Name:  "bnServer03",
				Child: &bnGenServ{},
				Args:  []etf.Term{"abc", 67890},
			},
		},
		Strategy: gen.SupervisorStrategy{
			Type: gen.SupervisorStrategyOneForAll,
			// Type:      gen.SupervisorStrategyRestForOne,
			// Type:      gen.SupervisorStrategyOneForOne,
			Intensity: 2,
			Period:    5,
			Restart:   gen.SupervisorStrategyRestartTemporary,
			// Restart: gen.SupervisorStrategyRestartTransient,
			// Restart: gen.SupervisorStrategyRestartPermanent,
		},
	}
	return spec, nil
}

// gen.Server implementation structure
type bnGenServ struct {
	gen.Server
}

func (dgs *bnGenServ) HandleCast(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	fmt.Printf("HandleCast (%s): %v\n", process.Name(), message)
	switch message {
	case etf.Atom("stop"):
		return gen.ServerStatusStopWithReason("stop they said")
	}
	return gen.ServerStatusOK
}

func (dgs *bnGenServ) HandleCall(process *gen.ServerProcess, from gen.ServerFrom, message etf.Term) (etf.Term, gen.ServerStatus) {
	fmt.Printf("HandleCall (%s): %v, From: %v\n", process.Name(), message, from)

	switch message {
	case etf.Atom("hello"):
		return etf.Atom("hi"), gen.ServerStatusOK
	}

	reply := etf.Tuple{etf.Atom("error"), etf.Atom("unknown_request")}
	return reply, gen.ServerStatusOK
}

func init() {
	flag.IntVar(&ListenRangeBegin, "listen_begin", 15151, "listen port range")
	flag.IntVar(&ListenRangeEnd, "listen_end", 25151, "listen port range")
	flag.StringVar(&NodeName, "name", "bn@127.0.0.1", "node name")
	flag.IntVar(&ListenEPMD, "epmd", 4369, "EPMD port")
	flag.StringVar(&Cookie, "cookie", "123", "cookie for interaction with erlang cluster")
}

func main() {
	flag.Parse()

	opts := node.Options{
		ListenRangeBegin: uint16(ListenRangeBegin),
		ListenRangeEnd:   uint16(ListenRangeEnd),
		EPMDPort:         uint16(ListenEPMD),
	}

	// Initialize new node with given name, cookie, listening port range and epmd port
	bnNode, _ := ergo.StartNode(NodeName, Cookie, opts)

	// start application
	if _, err := bnNode.ApplicationLoad(&bnApp{}); err != nil {
		panic(err)
	}

	appProcess, _ := bnNode.ApplicationStart("bnApp")
	fmt.Println("Run erl shell:")
	fmt.Printf("erl -name %s -setcookie %s\n", "erl-"+bnNode.Name(), Cookie)

	fmt.Println("-----Examples that can be tried from 'erl'-shell")
	fmt.Printf("gen_server:cast({%s,'%s'}, stop).\n", "bnServer01", bnNode.Name())
	fmt.Printf("gen_server:call({%s,'%s'}, hello).\n", "bnServer01", bnNode.Name())

	appProcess.Wait()
	bnNode.Stop()
}
