package main

import (
	"fmt"
	"log"
	"strings"
)

func dispatcher(commands []string, mode Mode, kv_store *store) bool{

	switch strings.ToLower(commands[0]) {
	
	case "ping":
		fmt.Println("PONG")
		
	case "set":
		err := setFunction(commands)
		if mode == FILE{
			if err != nil{
				fmt.Println(err.Error())
				break 
			}	
			kv_store.kv_pair[commands[1]] = commands[2]
			return true
		}	
		if err != nil{
			fmt.Println(err.Error())
			break 
		}	
		kv_store.kv_pair[commands[1]] = commands[2]
		fmt.Println("OK")
		if err = writeLog(commands); err != nil{
			log.Fatal(err)
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
		if mode == FILE{
			if err != nil{
				fmt.Println(err.Error())
				break 
			}	
			for _, val := range delKeys {
				delete(kv_store.kv_pair, val)
			}
			return true
		}	
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		writeLog(commands)
		fmt.Printf("(integer) %d\n", len(delKeys))
	
	case "exist":
		count, err := existFunction(commands, kv_store.kv_pair)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Printf("(integer) %d\n", count)

	case "rename":
		val, err := renameFunction(commands, kv_store.kv_pair)
		if mode == FILE{
			if err != nil{
				fmt.Println(err.Error())
				break 
			}	
			delete(kv_store.kv_pair, commands[1])
			kv_store.kv_pair[commands[2]] = val 
			return true
		}	
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		writeLog(commands)
		delete(kv_store.kv_pair, commands[1])
		kv_store.kv_pair[commands[2]] = val 
		fmt.Println("OK")
	
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
		return true  
	
	default:
	Err := fmt.Errorf("%w %s",ErrUnknownCmd, commands[0])
	fmt.Println(Err.Error())
	}

	return false


}
