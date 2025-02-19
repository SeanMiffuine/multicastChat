package multicast

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"local/lib/nodeset"
	"local/lib/transport/types"
	"sync"
)

type TimeVector []uint32

func NewGroup(nodeset *nodeset.Group) *Group {
	g := &Group{
		name:             nodeset.Name,
		Nodeset:          nodeset,
		VectorTimestamps: make(TimeVector, nodeset.Size()),
	}

	Groups[nodeset.Name] = g
	g.cv = sync.NewCond(&g.Mx)
	g.Nodeset.AddChangeListener(func() {
		g.Mx.Lock()
		defer g.Mx.Unlock()
		fmt.Println("Vector timestamps before update:", g.VectorTimestamps)
		if g.VectorTimestamps.Size() < g.Nodeset.Size() {
			g.VectorTimestamps = append(g.VectorTimestamps, 0)
			fmt.Println("Vector timestamps after update:", g.VectorTimestamps)
		}
		g.cv.Broadcast()
	})

	fmt.Println("Nodes:", g.Nodeset.Size())
	fmt.Println("Vector timestamps new:", g.VectorTimestamps)
	return g
}

func (vec TimeVector) Size() int {
	return len(vec)
}

func ScheduleDelivery(handleCall func(buf *bytes.Buffer, fromAddr types.NetAddr) []byte, buf *bytes.Buffer, fromAddr types.NetAddr) []byte {
	// to differentiate between regular call and multicast message, check for "MSG" as a key in buffer
	load := bytes.NewBuffer(buf.Bytes())

	var isMulticast bool
	if err := binary.Read(load, binary.NativeEndian, &isMulticast); err != nil {
		panic(err)
	}

	if !isMulticast {
		return handleCall(load, fromAddr)
	}

	var vectorTimestampsLen uint32
	if err := binary.Read(load, binary.NativeEndian, &vectorTimestampsLen); err != nil {
		panic(err)
	}

	vectorTimestamps := make([]uint32, vectorTimestampsLen)
	if err := binary.Read(load, binary.NativeEndian, &vectorTimestamps); err != nil {
		panic(err)
	}

	var passedNodeID uint32
	if err := binary.Read(load, binary.NativeEndian, &passedNodeID); err != nil {
		panic(err)
	}

	var currentGroup = Groups["chat"]

	mutex.Lock()
	if !shoudlDelay(vectorTimestamps, currentGroup.VectorTimestamps, passedNodeID) {

		updatedTimestamps := updateTimestamps(currentGroup.VectorTimestamps, vectorTimestamps)

		currentGroup.VectorTimestamps = updatedTimestamps
		result := handleCall(load, fromAddr)
		processMessageQueue(handleCall, currentGroup)
		mutex.Unlock()
		return result
	} else {
		currentGroup.messageQueue = append(currentGroup.messageQueue, Message{
			payload:    load,
			timeStampe: vectorTimestamps,
			nodeId:     passedNodeID,
			fromAddr:   fromAddr,
		})
		mutex.Unlock()
		return nil
	}
}

func updateTimestamps(A TimeVector, B TimeVector) TimeVector {

	var updated TimeVector = make(TimeVector, A.Size())
	for i := 0; i < A.Size(); i++ {
		updated[i] = max(A[i], B[i])
	}
	return updated
}

func shoudlDelay(recieve TimeVector, current TimeVector, index uint32) bool {
	/*
		- events are delayed until these two conditions are met
		1. ts(m)[i] = VCj[i] + 1
		2. ts(m)[k] ≤ VCj [k] for all k ̸= i
	*/
	for i := 0; i < len(current); i++ {
		if i == int(index) {
			// cond 1.
			if recieve[i] != current[i]+1 {
				return true
			}
		} else {
			if recieve[i] > current[i] {
				// cond 2.
				return true
			}
		}
	}
	return false
}

func processMessageQueue(handleCall func(buf *bytes.Buffer, fromAddr types.NetAddr) []byte, g *Group) {
	for i := 0; i < len(g.messageQueue); i++ {
		msg := g.messageQueue[i]
		if !shoudlDelay(msg.timeStampe, g.VectorTimestamps, msg.nodeId) {
			updatedTimestamps := updateTimestamps(g.VectorTimestamps, msg.timeStampe)
			g.VectorTimestamps = updatedTimestamps
			handleCall(msg.payload, msg.fromAddr)
			if i == len(g.messageQueue)-1 {
				g.messageQueue = g.messageQueue[:i]
			} else {
				g.messageQueue = append(g.messageQueue[i:], g.messageQueue[:i+1]...)
			}
			i-- // Adjust index after removal
			processMessageQueue(handleCall, g)
			break
		}
	}
}

type Message struct {
	payload    *bytes.Buffer
	timeStampe TimeVector
	nodeId     uint32
	fromAddr   types.NetAddr
}

type Group struct {
	name             string
	Nodeset          *nodeset.Group
	timestampIndex   uint32     // index of timestamp for this node
	VectorTimestamps TimeVector // list of timestamps per node, should be associated with each inddex of Nodeset.nodeset
	Mx               sync.Mutex
	cv               *sync.Cond
	messageQueue     []Message
}

var Groups = make(map[string]*Group)
var mutex sync.Mutex
