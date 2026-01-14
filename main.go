package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type store struct{
	mu sync.Mutex
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

	for(!r.run){

		r.err = nil 
		r.resp = ""
		commands, err := takeInput(scanner)
		if err != nil{
			log.Fatal(err.Error())
		}

		start := time.Now()

		r = dispatcher(commands, CLI, &KV_store, r.aof_mode)
		if r.resp != ""{
			fmt.Println(r.resp)
		}
		if r.err != nil{
			fmt.Println(r.err.Error())
		}
		
		end := time.Now()
		elapsed := end.Sub(start)
		fmt.Printf("%v microseconds \n",elapsed.Microseconds())
	}
	
}
