package main

import (
	"flag"
	"fmt"
	"github.com/MeshBoxFoundation/blockchain-core/src"

	"github.com/ergo-services/ergo"
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
	node, err := ergo.StartNode(NodeName, Cookie, opts)
	if err != nil {
		panic(err)
	}

	// Spawn supervisor process
	process, err := node.Spawn("bn_sup", gen.ProcessOptions{}, &src.BnCoreSup{})
	if err != nil {
		panic(err)
	}

	fmt.Println("spawned bn_sup")

	follower := &src.Follower{}
	_, err = node.Spawn("follower", gen.ProcessOptions{MailboxSize: 10000}, follower)
	if err != nil {
		panic(err)
	}
	fmt.Println("OK")

	process.Wait()
	node.Stop()
	node.Wait()
}
