package transport

import (
	"bytes"
	"context"
	"encoding/binary"
	"local/lib/transport"
	"local/lib/transport/types"
	"local/multicast"
	"math/rand/v2"
	"time"
)

func Multicast(payload *bytes.Buffer, g *multicast.Group) *bytes.Buffer {
	g.Mx.Lock()
	newLoad := bytes.NewBuffer(make([]byte, 0))
	g.VectorTimestamps[g.Nodeset.NodeId()]++

	binary.Write(newLoad, binary.NativeEndian, true)
	binary.Write(newLoad, binary.NativeEndian, uint32(len(g.VectorTimestamps)))
	binary.Write(newLoad, binary.NativeEndian, []uint32(g.VectorTimestamps))
	binary.Write(newLoad, binary.NativeEndian, uint32(g.Nodeset.NodeId()))
	binary.Write(newLoad, binary.NativeEndian, payload.Bytes())

	g.Mx.Unlock()
	for _, node := range g.Nodeset.Nodeset() {
		if node.Id != g.Nodeset.NodeId() {
			go func() {
				if node.Addr == slowNode {
					time.Sleep(time.Duration(rand.Uint32()%5000) * time.Millisecond)
				}
				_, e := transport.Call(newLoad, ResolveNetAddr(node.Addr))
				if e != nil {
					panic(e)
				}
			}()
		}
	}
	return bytes.NewBuffer(make([]byte, 0))
}

func Call(payload *bytes.Buffer, to types.NetAddr) (result *bytes.Buffer, err error) {
	newLoad := bytes.NewBuffer(make([]byte, 0))
	binary.Write(newLoad, binary.NativeEndian, false)
	if err := binary.Write(newLoad, binary.NativeEndian, payload.Bytes()); err != nil {
		panic(err)
	}
	result, err = transport.Call(newLoad, to)
	return
}

func Listen(context context.Context, handleCall func(msg *bytes.Buffer, from types.NetAddr) []byte) {
	transport.Listen(context, func(msg *bytes.Buffer, from types.NetAddr) []byte {
		return multicast.ScheduleDelivery(handleCall, msg, from)
	})
}

func LocalAddr() types.NetAddr {
	return transport.LocalAddr()
}

func ResolveNetAddr(addr string) types.NetAddr {
	return transport.ResolveNetAddr(addr)
}

func SetSlowNode(node string) {
	slowNode = node
}

var slowNode string
