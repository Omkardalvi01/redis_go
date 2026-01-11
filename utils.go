package main

import (
	"bufio"
	"context"
	"strings"
	"time"
)

type TimeFormat int 
const (
	Milli TimeFormat = iota 
	Second
	UX_TF 
	UX_TF_Milli
)

func delayedDel(kv_store *store, key string, t int, tf TimeFormat) {
	ctx := context.Background()
	var cancel context.CancelFunc

	switch tf {
	case Second:
		ctx, cancel = context.WithTimeout(ctx, time.Duration(t)*time.Second)
		defer cancel()

	case Milli:
		ctx, cancel = context.WithTimeout(ctx, time.Duration(t)*time.Millisecond)
		defer cancel()

	}
	
	<-ctx.Done()
	kv_store.mu.Lock()
	delete(kv_store.kv_pair, key)
	kv_store.mu.Unlock()
}

func takeInput(scanner *bufio.Scanner) ([]string, error){
	scanner.Scan()
	if err := scanner.Err(); err != nil{
		return nil, err 
	}
	
	
	line_command := string(scanner.Bytes())
	line_command_trimmed := strings.TrimRight(line_command, " ")
	commands := strings.Split(line_command_trimmed, " ")

	// for i, str := range commands{
	// 	commands[i] = strings.TrimSpace(str)
	// }

	return commands, nil 

}
