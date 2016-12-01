package g

import (
	"log"
	"runtime"
)

// change log:
// 1.0.7: code refactor for open source
// 1.0.8: bugfix loop init cache
// 1.0.9: update host table anyway
// 1.1.0: remove Checksum when query plugins
// 1.1.0.sdpv001 add api to get endpoint by ip
// 1.1.0.sdpv002 add function to get endpoint by both public ip and private ip
// 1.1.0.sdpv003 add function to get hosts information with <ip,endpoint> pair
// 1.1.0.sdpv004 modify import from open-falcon to ZeaLoVe ,fix run_time of metrics
// 1.1.0.sdpv005 add /hosts/id get host_id by name
const (
	VERSION = "1.1.0.sdpv005"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
