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
const (
	VERSION = "1.1.0.sdpv003"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
