package main

import (
	"bufio"
	"errors"
	"fmt"
	"maps"
	"os"
	"regexp"
	"strings"
)

var (
	ErrNumberArguments = errors.New("(error) Wrong number of arguments")
	ErrUnkownCmd = errors.New("(error) Unkown command ")
	ErrNokey = errors.New("(error) no such key")
	ErrRegex = errors.New("(error) in the regex operation")
)
func main(){
	scanner := bufio.NewScanner(os.Stdin)
	kv_store := make(map[string]string)
	fmt.Println("Server Running")
	exit_val := false

	for(!exit_val){
		scanner.Scan()
		if err := scanner.Err(); err != nil{
			fmt.Print("Error while reading input")
		}

		line_command := string(scanner.Bytes())
		commands := strings.Split(line_command, " ")

		
		switch  strings.ToLower(commands[0]) {
		case "ping":
			fmt.Println("PONG")
		
		case "get":
			if (len(commands) < 2 || commands[1] == ""){
				Err := fmt.Errorf("%w for %s command",ErrNumberArguments, commands[0])
				fmt.Println(Err.Error())
				continue
			}
			if _, ok := kv_store[commands[1]]; ok{
				fmt.Println(kv_store[commands[1]])
			}else{
				fmt.Println("(nil)")
			}
			
		case "set":
			if len(commands) < 3 {
				Err := fmt.Errorf("%w for %s command",ErrNumberArguments, commands[0])
				fmt.Println(Err.Error())
				continue
			}
			kv_store[commands[1]] = commands[2]
			fmt.Println("OK")
		
		case "del":
			if len(commands) < 2{
				Err := fmt.Errorf("%w for %s command",ErrNumberArguments, commands[0])
				fmt.Println(Err.Error())
				continue
			}

			count := 0
			for i := 1; i < len(commands); i++ {
				if _ , ok := kv_store[commands[i]]; ok{
					delete(kv_store, commands[i])
					count++
				}
			}
			fmt.Printf("(integer) %d\n",count)
		
		case "exists":
			if len(commands) != 2 {
				Err := fmt.Errorf("%w for %s command",ErrNumberArguments, commands[0])
				fmt.Println(Err.Error())
				continue
			}

			if len(commands) < 2 {
				Err := fmt.Errorf("%w for %s command",ErrNumberArguments, commands[0])
				fmt.Println(Err.Error())
				continue
			}

			count := 0
			for i := 1; i < len(commands); i++ {
				if _ , ok := kv_store[commands[i]]; ok {
					count++
				}
			}
			fmt.Printf("(integer) %d\n",count)
		
		case "empty":
			if len(commands) != 1 {
				Err := fmt.Errorf("%w for %s command",ErrNumberArguments, commands[0])
				fmt.Println(Err.Error())
				continue
			}
			fmt.Printf("(integer) %d\n", len(kv_store))
		
		case "keys":
			if len(commands) != 2 {
				Err := fmt.Errorf("%w for %s command",ErrNumberArguments, commands[0])
				fmt.Println(Err.Error())
				continue
			}

			count := 1
			if commands[1] == "*"{
				keys := maps.Keys(kv_store)
				for key := range keys{
					fmt.Printf("%d) %s\n",count,key)
					count++
				}
				continue
			}

			if strings.Count(commands[1], "*") == 0{
				if _, ok := kv_store[commands[1]]; ok{
					fmt.Printf("%d) %s\n",count,commands[1])
				}
				continue
			}

			regex_str := commands[1]
			if regex_str[len(regex_str)-1] == '*' {
				regex_str = regex_str[:len(regex_str)-1]
				regex_str += "([a-z]+)"
			}

			regex_str = strings.ReplaceAll(regex_str, "*", "([a-z])")
			keys := maps.Keys(kv_store)
			matched_str := make([]string,0)
			for key := range keys{
				match, err := regexp.Match(regex_str, []byte(key))
				if err != nil{
					Err := fmt.Errorf("%w caused due to %w while finding %s", ErrRegex, err, key)
					fmt.Println(Err)
					continue
				}
				if match {
					matched_str = append(matched_str, key)
				}
			}

			for _ , key := range matched_str{
				fmt.Printf("%d) %s\n",count,key)
				count++
			}

			if (len(matched_str) == 0){
				fmt.Println("(empty array)")
			}

		
		case "rename":
			if len(commands) != 3{
				Err := fmt.Errorf("%w for %s command",ErrNumberArguments, commands[0])
				fmt.Println(Err.Error())
				continue
			}

			val , ok := kv_store[commands[1]]
			if !ok {
				fmt.Println(ErrNokey.Error())
				continue
			}

			delete(kv_store, commands[1])
			kv_store[commands[2]] = val 
			fmt.Println("OK")

		case "exit":
			fmt.Printf("OK")
			exit_val=true

		default:
			Err := fmt.Errorf("%w %s",ErrUnkownCmd, commands[0])
			fmt.Println(Err.Error())
		}
	}
	
}