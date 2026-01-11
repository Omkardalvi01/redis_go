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

func main(){
	scanner := bufio.NewScanner(os.Stdin)
	mappings := make(map[string]string)
	kv_store := store{}
	kv_store.kv_pair = mappings
	exit_val := false
	


	fmt.Println("Server Running")
	

	
	for(!exit_val){
		
		commands, err := takeInput(scanner)
		if err != nil{
			log.Fatal(err.Error())
		}

		start := time.Now()
		exit_val = dispatcher(commands, &kv_store)
		end := time.Now()
		elapsed := end.Sub(start)
		fmt.Printf("%v microseconds \n",elapsed.Microseconds())
	}
	
}
