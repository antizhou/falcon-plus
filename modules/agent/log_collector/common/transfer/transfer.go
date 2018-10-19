package transfer

import (
	"time"
	"math/rand"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/log_collector/common/dlog"
	"fmt"
	"net/rpc/jsonrpc"
	"log"
)

var Cache = make(chan LogPoint, 1024*100)

func init() {
	go func() {
		for {
			select {
			case p := <-Cache:
				points := make([]LogPoint, 0)
				points = append(points, p)
			DONE:
				for {
					select {
					case point := <-Cache:
						points = append(points, point)
						continue
					default:
						break DONE
					}
				}
				//开一个协程，异步发送至odin-agent
				go saveLog(points)
			}
			time.Sleep(10 * time.Second)
		}
	}()
}

type LogPoint struct {
	App     string
	Content string
	Time    time.Time
}

type Req struct {
	Data []LogPoint
}

func saveLog(points []LogPoint) string {
	req := Req{
		Data: points,
	}
	resp := ""
	var e error

	fmt.Println(req)
	for _, i := range rand.Perm(len(g.Config().Transfer.Addrs)) {
		addr := g.Config().Transfer.Addrs[i]

		conn, err := jsonrpc.Dial("tcp", addr)
		if err != nil {
			log.Println("dialing:", err)
			continue
		}

		err = conn.Call("Transfer.SaveLog2ES", req, &resp)
		if err != nil {
			e = err
			continue
		}
		break
	}
	if resp == "" {
		dlog.Error("save log error: %v", e)
	}
	return resp
}
