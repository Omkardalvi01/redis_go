package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type store struct{
	mu sync.RWMutex
	kv_pair map[string]string
}
var KV_store = store{}

type Mode int 
const (
	CLI Mode = iota
	FILE 
	SERVER 
)

type payload struct{
	run bool
	aof_mode bool 
	resp string
	err error 
}

func main(){
	scanner := bufio.NewScanner(os.Stdin)
	mappings := make(map[string]string)
	KV_store.kv_pair = mappings
	var r payload 
	r.aof_mode = true
	r.run = true 

	fmt.Println("Server Running")
	ingestionFunc(&KV_store)
	http.HandleFunc("/", handler)

	go func(){
			for(r.run){

				r.err = nil 
				r.resp = ""
				commands, err := takeInput(scanner)
				if err != nil{
					log.Fatal(err.Error())
				}

				start := time.Now()

				r = dispatcher(commands, CLI, &KV_store, r.aof_mode)
				if r.resp != ""{
					fmt.Print(r.resp)
				}
				if r.err != nil{
					fmt.Print(r.err.Error())
				}
				
				end := time.Now()
				elapsed := end.Sub(start)
				fmt.Printf("%v microseconds \n",elapsed.Microseconds())
		}
	}()
	
	if err := http.ListenAndServe(":8000", nil); err != nil{
		log.Fatal(err)
	}

	
}
