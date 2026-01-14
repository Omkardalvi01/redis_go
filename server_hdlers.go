package main

import "net/http"

func handler(w http.ResponseWriter, r *http.Request){
	commands := r.URL.Query()["cmd"]
	p := dispatcher(commands, SERVER, &KV_store, true)
	if p.err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(p.err.Error()))
	}else{
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(p.resp))
	}
	
}