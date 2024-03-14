package common

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

func GoPprofServe(port string) {
	// http://127.0.0.1:{port}/debug/pprof
	go http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", port), nil)
}
