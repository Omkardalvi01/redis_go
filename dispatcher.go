package main

import (
	"fmt"
	"log"
	"strings"
)

func dispatcher(commands []string, mode Mode, kv_store *store, valid_append bool) (payload){
	r := payload{}
	r.run = true 
	r.aof_mode = true

	switch strings.ToLower(commands[0]) {
	case "ping":
		r.resp = fmt.Sprintln("PONG")
		
	case "set":
		if err := setFunction(commands); err != nil{
			if mode == FILE{
				log.Fatalln(err)
			}else{
				r.err = err
			}
			break
		}
		r.resp = fmt.Sprintln("OK")
		//data manipulation
		kv_store.mu.Lock()
		kv_store.kv_pair[commands[1]] = commands[2]
		kv_store.mu.Unlock()

		if valid_append{
			r.err = Append(commands)
		}
	
	case "get":
		val, present, err := getFunction(commands, kv_store)
		if err != nil {
			r.err = err 
			break
		}
		if !present {
			r.resp = fmt.Sprintln("(nil)")
			break 
		}
		r.resp = fmt.Sprintln(val)
	
	case "del":
		delKeys, err := deleteFunction(commands, kv_store)
		if err != nil && mode == FILE{
			log.Fatalln(err)
			break 
		}	
		if err != nil {
			r.err = err 
			break
		}
		
		kv_store.mu.Lock()
		for _, val := range delKeys {
			delete(kv_store.kv_pair, val)
		}
		kv_store.mu.Unlock()

		r.resp = fmt.Sprintf("(integer) %d\n", len(delKeys))
		if valid_append{
			r.err = Append(commands)
		}
		
	case "exist":
		count, err := existFunction(commands, kv_store)
		if err != nil {
			r.err = err 
			break
		}
		r.resp = fmt.Sprintf("(integer) %d\n", count)

	case "rename":
		val, err := renameFunction(commands, kv_store)
		if err != nil && mode == FILE{
			log.Fatalln(err)
			break 
		}	
		if err != nil {
			r.err = err 
			break
		}
		
		kv_store.mu.Lock()
		delete(kv_store.kv_pair, commands[1])
		kv_store.kv_pair[commands[2]] = val 
		kv_store.mu.Unlock()

		r.resp = fmt.Sprintln("OK")
		if valid_append {
			r.err = Append(commands)
		}

	
	case "empty":
		val, err := emptyFunction(commands, kv_store)
		if err != nil {
			r.err = err 
			break
		}
		r.resp = fmt.Sprintf("(integer) %d\n", val)

	case "keys":
		matched_str, err := keysFunction(commands, kv_store)
		if err != nil {
			r.err = err 
			break
		}

		count := 0
		for _ , key := range matched_str{
			r.resp += fmt.Sprintf("%d) %s\n",count,key)
			count++
		}
		if (len(matched_str) == 0){
			r.resp = fmt.Sprintln("(empty array)")
		}
	
	case "expire", "pexpire":
		err := expireFunction(commands, kv_store)
		if err != nil {
			r.err = err 
			break
		}

	case "exit":
		r.resp = fmt.Sprintln("OK")
		r.run = false
		return r 
	
	default:
	Err := fmt.Errorf("%w %s",ErrUnknownCmd, commands[0])
	r.err = Err 
	}

	return r
}
