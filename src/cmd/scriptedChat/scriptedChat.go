package main

import (
	"context"
	"fmt"
	"local/chat"
	"local/chat/rpc/serverStub"
	"local/lib/finalizer"
	"local/lib/nodeset/control"
	"local/lib/rpc"
	"local/multicast"
	"local/multicast/transport"
	"os"
)

// func init() {
// 	numNodes = 3
// 	ws(0, "", "a", "b", "c")
// 	ws(1, "c", "d", "e", "f")
// 	slowChannels = []channel{{0, 2}}
// }

// test case for interleaving of messages between two nodes
// func init() {
// 	numNodes = 3
// 	ws(0, "", "a", "b", "c")
// 	ws(1, "", "e", "f", "g")
// 	slowChannels = []channel{{0, 2}}
// }

// Alice, Bob, and Charlie test case
func init() {
	numNodes = 3
	ws(0, "", "[A:1] Hi I'm Alice!", "[A:2] How is the weather")
	ws(1, "[A:1] Hi I'm Alice!", "[B:1] Hi Alice, I'm Bob!")
	ws(2, "", "[C:1] Hi I'm Charlie")
	slowChannels = []channel{{0, 1}, {1, 2}}
}

// this test case from piazza @195 by Kyle Jones
// https://piazza.com/class/m502jr8kaui37v/post/195
// func init() {
// 	numNodes = 5
// 	ws(0, "", "a", "b", "c", "d", "e", "f", "g", "h")
// 	ws(1, "h", "i", "j", "k", "l", "m", "n", "o", "p", "q")
// 	ws(2, "q", "r", "s", "t", "u", "v", "w", "x", "y", "z")
// 	ws(3, "a", "1", "2", "3", "4", "5", "6", "7", "8", "9")

// 	slowChannels = []channel{{0, 4}}
// }

func main() {
	displayCh := make(chan string)
	ctx, cancel := finalizer.WithCancel(context.Background())
	defer func() { cancel(); <-ctx.Done() }()
	control.Bind(os.Args[1])
	serverStub.Register()
	rpc.Start(ctx)

	g := multicast.NewGroup(control.Add(ctx, "chat"))
	// creates nodeset group called chat
	// each group has a nodeID to know which node itself is

	chat.SetDisplayFunc(func(msg string) {
		fmt.Println(msg)
		displayCh <- msg
	})

	control.AwaitSize(g.Nodeset, numNodes)
	// wait until group has #numNodes nodes
	for _, c := range slowChannels {
		if c.from == g.Nodeset.NodeId() {
			transport.SetSlowNode(g.Nodeset.Nodeset()[c.to].Addr)
		}
	}
	// sets all to slownodes ?

	for _, s := range script[g.Nodeset.NodeId()] {
		if s.waitFor != "" {
			for {
				if <-displayCh == s.waitFor {
					break
				}
			}
		}
		for _, i := range s.issues {
			fmt.Println(i)
			chat.Post(g, i)
		}
	}
	for {
		<-displayCh
	}
	// executes the script to send test messages
}

type step struct {
	waitFor string
	issues  []string
}
type channel struct {
	from uint32
	to   uint32
}

var script = make(map[uint32][]step)
var slowChannels []channel
var numNodes int

func ws(nodeId uint32, waitFor string, issue ...string) {
	script[nodeId] = append(script[nodeId], step{waitFor, issue})
}
