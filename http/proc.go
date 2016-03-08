package http

import (
	"encoding/json"
	"github.com/ZeaLoVe/hbs/cache"
	"github.com/ZeaLoVe/hbs/db"
	"github.com/ZeaLoVe/hbs/g"
	"github.com/open-falcon/common/model"
	"net/http"
	"strings"
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

	//get ,API of all hosts, use in agent alive check.
	http.HandleFunc("/all/hosts", func(w http.ResponseWriter, r *http.Request) {
		var hosts []ResponseHost
		var host ResponseHost
		cache.HostMap.Lock()
		//cache中的map的key就是hostname，也就是endpoint；value是hostid没用
		for key, _ := range cache.HostMap.M {
			host.Endpoint = key
			host.Ip = cache.HostMap.M2[key] //通过hostname找IP
			if strings.EqualFold(host.Ip, "0.0.0.0") {
				continue
			}
			hosts = append(hosts, host)
		}
		cache.HostMap.Unlock()
		RenderJson(w, hosts)
	})

	//get ,API of all virtual hosts.ip=0.0.0.0
	http.HandleFunc("/all/vhosts", func(w http.ResponseWriter, r *http.Request) {
		var hosts []ResponseHost
		var host ResponseHost
		cache.HostMap.Lock()
		//cache中的map的key就是hostname，也就是endpoint；value是hostid没用
		for key, _ := range cache.HostMap.M {
			host.Endpoint = key
			host.Ip = cache.HostMap.M2[key] //通过hostname找IP
			if strings.EqualFold(host.Ip, "0.0.0.0") {
				hosts = append(hosts, host)
			}
		}
		cache.HostMap.Unlock()
		RenderJson(w, hosts)
	})

	//API add virtual host:GET
	http.HandleFunc("/vhost/add", func(w http.ResponseWriter, r *http.Request) {
		var args model.AgentReportRequest
		args.Hostname = r.FormValue("name")
		args.IP = "0.0.0.0"
		args.AgentVersion = "0.0.0"
		args.PluginVersion = "0.0.0"
		if args.Hostname == "" {
			RenderMsgJson(w, "require host name")
			return
		}
		cache.Agents.Put(&args)
		RenderMsgJson(w, "add virtual host done.")
	})

}

type Endpoint struct {
	Endpoint string `json:"endpoint,omitempty"`
}

type ResponseHost struct {
	Ip       string `json:"ip,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
}

type ResponseEndpoints struct {
	Items []ResponseHost `json:"items,omitempty"`
}
