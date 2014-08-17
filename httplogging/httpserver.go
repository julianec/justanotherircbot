package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

func logRequest(rw http.ResponseWriter, req *http.Request) {
	var err error
	var body []byte

	log.Print(req.Method, " ", req.RequestURI, " ", req.Proto)
	for headername, headervalues := range req.Header {
		for _, headervalue := range headervalues {
			log.Print(headername, ": ", headervalue)
		}
	}
	body, err = ioutil.ReadAll(req.Body)
	if err != nil {
		log.Print("Error reading body: ", err)
		return
	}
	log.Print(string(body))
}

func main() {
	var err error
	var bindto *string = flag.String("bind-to", ":8080", "host:port")
        flag.Parse()

	http.HandleFunc("/", logRequest)
	err = http.ListenAndServe(*bindto, nil)
	if err != nil {
		log.Print("Error setting up http-server: ", err)
	}
}
