package transfer

import (
	"testing"
	"log"
	"net/rpc/jsonrpc"
)

func Test_saveLog(t *testing.T) {

	logPoint := LogPoint{
		App:     "monitor-watchdog",
		Content: "sdfasdf",
		Time:    1539682470000,
	}

	data := make([]LogPoint, 0)
	data = append(data, logPoint)

	req := Req{
		Data: data,
	}

	addr := "127.0.0.1:8433"

	conn, err := jsonrpc.Dial("tcp", addr)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	conn.Call("Transfer.SaveLog2ES", req, nil)
}
