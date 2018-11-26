package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/corbamico/cmdserver/utils"
)

type webCmdShell struct{}

func newwebCmdShell() *webCmdShell {
	return new(webCmdShell)
}

func (w *webCmdShell) ServeHTTP(r http.ResponseWriter, rq *http.Request) {
	switch rq.URL.Path {
	case "/status":
		stat := utils.StatProc()
		result := fmt.Sprintf("{\"result\":%d}", stat)
		r.Write([]byte(result))
	case "/restart":
		utils.RestartProc()
		r.Write([]byte(`{"result":0}`))
	default:
		r.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	err := http.ListenAndServe(":8080", newwebCmdShell())
	if err != nil {
		log.Fatalln("CmdShell server failed: ", err.Error())
	}
}
