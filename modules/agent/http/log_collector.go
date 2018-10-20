package http

import (
	"net/http"
	"github.com/open-falcon/falcon-plus/modules/agent/collector/strategy"
	"github.com/open-falcon/falcon-plus/modules/agent/collector/worker"
)

func configLogCollector() {
	//http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
	//	RenderDataJson(w, "ok")
	//})

	http.HandleFunc("/strategy", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, strategy.GetListAll())
	})

	http.HandleFunc("/cached", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, worker.GetCachedAll())
	})

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		//RenderDataJson(w, CheckLogByStrategy(log))
	})
}
