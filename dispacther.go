package main

import (
	"fmt"
	"log"
	"strings"
)

func dispatcher(commands []string, mode Mode, kv_store *store, valid_append bool) (bool, bool){

	switch strings.ToLower(commands[0]) {
	
	case "ping":
		fmt.Println("PONG")
		
	case "set":
		if err := setFunction(commands); err != nil{
			if mode == FILE{
				log.Fatalln(err)
			}else{
				fmt.Println(err.Error())
			}
			break
		}
		
		kv_store.kv_pair[commands[1]] = commands[2]
		if mode != FILE{
			fmt.Println("OK")

			if valid_append {
				if err := writeLog(commands); err != nil{
					fmt.Println(err.Error())
					fmt.Println("AOF disabled")
					return true, false 
				}
			}
			
		}
	
	case "get":
		val, present, err := getFunction(commands, kv_store.kv_pair)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		if !present {
			fmt.Println("(nil)")
			break 
		}
		fmt.Println(val)
	
	case "del":
		delKeys, err := deleteFunction(commands, kv_store.kv_pair)
		if err != nil && mode == FILE{
			log.Fatalln(err)
			break 
		}	
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		for _, val := range delKeys {
			delete(kv_store.kv_pair, val)
		}

		if mode != FILE{
			fmt.Printf("(integer) %d\n", len(delKeys))
			if valid_append {
				if err := writeLog(commands); err != nil{
					fmt.Println(err.Error())
					fmt.Println("AOF disabled")
					return true, false 
				}
			}
		}

	
	case "exist":
		count, err := existFunction(commands, kv_store.kv_pair)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Printf("(integer) %d\n", count)

	case "rename":
		val, err := renameFunction(commands, kv_store.kv_pair)
		if err != nil && mode == FILE{
			log.Fatalln(err)
			break 
		}	
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		delete(kv_store.kv_pair, commands[1])
		kv_store.kv_pair[commands[2]] = val 
		
		if mode != FILE{
			fmt.Println("OK")
			if valid_append {
				if err := writeLog(commands); err != nil{
					fmt.Println(err.Error())
					fmt.Println("AOF disabled")
					return true, false 
				}
			}
		}

	
	case "empty":
		val, err := emptyFunction(commands, kv_store.kv_pair)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Printf("(integer) %d\n", val)

	case "keys":
		matched_str, err := keysFunction(commands, kv_store.kv_pair)
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		count := 0
		for _ , key := range matched_str{
			fmt.Printf("%d) %s\n",count,key)
			count++
		}
		if (len(matched_str) == 0){
			fmt.Println("(empty array)")
		}
	
	case "expire", "pexpire":
		err := expireFunction(commands, kv_store)
		if err != nil {
			fmt.Println(err.Error())
			break
		}

	case "exit":
		fmt.Println("OK")
		return true, true  
	
	default:
	Err := fmt.Errorf("%w %s",ErrUnknownCmd, commands[0])
	fmt.Println(Err.Error())
	}

	return false, true


}
