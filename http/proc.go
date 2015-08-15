package http

import (
	"github.com/open-falcon/hbs/cache"
	"github.com/open-falcon/hbs/db"
	"net/http"
)

func configProcRoutes() {
	http.HandleFunc("/expressions", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, cache.ExpressionCache.Get())
	})

	http.HandleFunc("/plugins/", func(w http.ResponseWriter, r *http.Request) {
		hostname := r.URL.Path[len("/plugins/"):]
		RenderDataJson(w, cache.GetPlugins(hostname))
	})

	//API get endpoint by name
	http.HandleFunc("/endpoint", func(w http.ResponseWriter, r *http.Request) {
		var res ResponseEndpoints
		var host ResponseHost
		host.Ip = r.FormValue("ip")
		host.Endpoint, _ = db.QueryEndpoint(host.Ip)
		res.Items = append(res.Items, host)
		RenderJson(w, res)
	})

	http.HandleFunc("/endpoints", func(w http.ResponseWriter, r *http.Request) {
		//body also response
		var body ResponseEndpoints
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&body)

		if err != nil {
			RenderMsgJson(w, "Not param, may be with wrong format")
			return
		}

		for i, _ := range body {
			if body.Items[i].Ip == "" {
				continue
			}
			body.Items[i].Endpoint, _ = db.QueryEndpoint(body.Items[i].Ip)
		}
		RenderJson(w, body)
	})

}

type ResponseHost struct {
	Ip       string `json:"ip,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
}

type ResponseEndpoints struct {
	Items []ResponseHost `json:"items,omitempty"`
}
