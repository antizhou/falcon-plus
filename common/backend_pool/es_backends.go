// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package backend_pool

import (
	"fmt"
	"time"

	connp "github.com/toolkits/conn_pool"
	"github.com/olivere/elastic"
	"context"
)

type EsClient struct {
	cli  elastic.Client
	name string
}

func (t EsClient) Name() string {
	return t.name
}

func (t EsClient) Closed() bool {
	return true
}

func (t EsClient) Close() error {
	return nil
}

func newEsConnPool(address string, maxConns int, maxIdle int, options []elastic.ClientOptionFunc) *connp.ConnPool {
	pool := connp.NewConnPool("es", address, int32(maxConns), int32(maxIdle))

	pool.New = func(name string) (connp.NConn, error) {

		client, err := elastic.NewClient(options...)
		if err != nil {
			return nil, err
		}

		return EsClient{*client, name}, nil
	}

	return pool
}

type EsConnPoolHelper struct {
	p           *connp.ConnPool
	maxConns    int
	maxIdle     int
	connTimeout int
	callTimeout int
	address     string
}

func NewEsConnPoolHelper(address string, maxConns, maxIdle, connTimeout, callTimeout int, options []elastic.ClientOptionFunc) *EsConnPoolHelper {
	return &EsConnPoolHelper{
		p:           newEsConnPool(address, maxConns, maxIdle, options),
		maxConns:    maxConns,
		maxIdle:     maxIdle,
		connTimeout: connTimeout,
		callTimeout: callTimeout,
		address:     address,
	}
}

func (t *EsConnPoolHelper) Send(data []byte) (err error) {
	conn, err := t.p.Fetch()
	if err != nil {
		return fmt.Errorf("get connection fail: err %v. proc: %s", err, t.p.Proc())
	}

	cli := conn.(EsClient).cli

	done := make(chan error, 1)
	go func() {
		_, err := cli.Index().
			Index("monitor-metric-" + time.Now().Format("20060102")).
			Type("m").
			BodyString(string(data)).
			Do(context.TODO())
		done <- err
	}()

	select {
	case <-time.After(time.Duration(t.callTimeout) * time.Millisecond):
		t.p.ForceClose(conn)
		return fmt.Errorf("%s, call timeout", t.address)
	case err = <-done:
		if err != nil {
			t.p.ForceClose(conn)
			err = fmt.Errorf("%s, call failed, err %v. proc: %s", t.address, err, t.p.Proc())
		} else {
			t.p.Release(conn)
		}
		return err
	}
}

func (t *EsConnPoolHelper) Destroy() {
	if t.p != nil {
		t.p.Destroy()
	}
}
