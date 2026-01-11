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

type Mode int 
const (
	CLI Mode = iota
	FILE 
	SERVER 
)

func main(){
	scanner := bufio.NewScanner(os.Stdin)
	mappings := make(map[string]string)
	kv_store := store{}
	kv_store.kv_pair = mappings
	exit_val := false
	


	fmt.Println("Server Running")
	ingestionfunc(&kv_store)

	
	for(!exit_val){
		
		commands, err := takeInput(scanner)
		if err != nil{
			log.Fatal(err.Error())
		}

		start := time.Now()
		exit_val = dispatcher(commands, CLI, &kv_store)
		end := time.Now()
		elapsed := end.Sub(start)
		fmt.Printf("%v microseconds \n",elapsed.Microseconds())
	}
	
}
