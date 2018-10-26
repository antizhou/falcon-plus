package strategy

import (
	"encoding/json"
	"fmt"
	"time"

	dlog "github.com/open-falcon/falcon-plus/logger"
	"github.com/open-falcon/falcon-plus/modules/agent/collector/common/scheme"
	"github.com/parnurzeal/gorequest"
)

func getHTTPStrategy(addrs []string, uri string, timeoutInSec int) (strategies []*scheme.Strategy, e error) {
	for _, addr := range addrs {

		url := fmt.Sprintf("%s%s", addr, uri)

		dlog.Infof("URL in get strategy : [%s]", url)

		body, err := getRequest(url, timeoutInSec)

		if err != nil {
			strategies = nil
			e = err
			continue
		}

		var strategyResp []*scheme.Strategy
		err = json.Unmarshal([]byte(body), &strategyResp)
		if err != nil {
			dlog.Errorf("json decode error when update strategy : [%v]", err)
			return nil, err
		}
		return strategyResp, nil
	}
	return
}

func getRequest(url string, timeout int) (string, error) {
	request := gorequest.New().Timeout(time.Duration(timeout) * time.Second)
	resp, body, errs := request.Get(url).End()

	if errs == nil {
		if resp.StatusCode != 200 {
			dlog.Errorf("get HTTP Request Response: [code:%d][body:%s][errs:%v]", resp.StatusCode, body, errs)
			return body, fmt.Errorf("Code is not 200")
		}
		dlog.Infof("get HTTP Request Response : [code:%d][body:%s]", 200, body)
		return body, nil
	}
	dlog.Errorf("get HTTP Request Response: [body:%s][errs:%v]", body, errs)
	return body, fmt.Errorf("%v", errs)
}
