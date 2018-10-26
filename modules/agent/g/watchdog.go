package g

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	log "github.com/open-falcon/falcon-plus/logger"
)

func getApp() string {
	app := ""

	for _, i := range rand.Perm(len(Config().WatchDog.Addrs)) {
		addr := Config().WatchDog.Addrs[i]

		resp, err := http.Get(addr + "/monitor/v1/api/getAppByIp?ip=" + IP());
		if err != nil {
			continue
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
			log.Infof("[g.watchdog.getApp] init app name error")
			continue
		}

		app = string(body)
		break
	}

	if len(app) == 0 {
		app = IP()
		log.Infof("[g.watchdog.getApp] init app name error, please check watchdog addr!")
	}
	return app
}
