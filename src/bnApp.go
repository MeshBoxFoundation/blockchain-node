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
