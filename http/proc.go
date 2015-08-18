package http

import (
	"encoding/json"
	"github.com/open-falcon/hbs/cache"
	"github.com/open-falcon/hbs/db"
	"github.com/open-falcon/hbs/g"
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
		if host.Ip == "" {
			RenderMsgJson(w, "Not param")
			return
		}
		target_ip := host.Ip
		if !isPrivateIP(target_ip) {
			//转化成内网IP
			target_ip = PrivateIP(host.Ip, g.Config().Nat)
		}
		host.Endpoint, _ = db.QueryEndpoint(target_ip)
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

		for i, _ := range body.Items {
			if body.Items[i].Ip == "" {
				continue
			}
			target_ip := body.Items[i].Ip
			if !isPrivateIP(target_ip) {
				//转化成内网IP
				target_ip = PrivateIP(body.Items[i].Ip, g.Config().Nat)
			}
			body.Items[i].Endpoint, _ = db.QueryEndpoint(target_ip)
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
