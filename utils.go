package main

import (
	"bufio"
	"context"
	"os"
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

	return commands, nil 

}

func writeLog(commands []string) (error) {
	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	cmd := strings.Join(commands, " ")
	_, err = f.WriteString(cmd + "\n")
	if err != nil {
		return err
	}
	return nil 
}

func ingestionfunc(kv_store *store) (error) {
	f, err := os.OpenFile("log.txt", os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	AOF_mode := true
	scanner := bufio.NewScanner(f)
	for scanner.Scan(){
		cmd_line := scanner.Text()
		commands := strings.Split(cmd_line, " ")
		_, AOF_mode = dispatcher(commands, FILE, kv_store, AOF_mode)
	}
	
	return nil
}
