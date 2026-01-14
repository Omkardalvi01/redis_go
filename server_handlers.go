package main

import (
	"fmt"
	"net/http"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()["cmd"][0]
	commands := strings.Split(params, " ")
	p := dispatcher(commands, SERVER, &KV_store, true)
	if p.err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, p.err.Error())
	}else{
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, p.resp)
	}
	
}