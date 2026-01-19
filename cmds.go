package main

import (
	"errors"
	"fmt"
	"maps"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrNumberArguments = errors.New("(error) Wrong number of arguments")
	ErrUnknownCmd      = errors.New("(error) Unkown command ")
	ErrNokey           = errors.New("(error) no such key")
	ErrRegex           = errors.New("(error) in the regex operation")
	ErrInvalidDuration = errors.New("(error) duraion is invalid")
)

func setFunction(commands []string) error {
	
	if len(commands) != 3 {
		Err := fmt.Errorf("%w for %s command",ErrNumberArguments, commands[0])
		return Err
	}
	return nil
}

func getFunction(commands []string, kv_store *store) (string, bool, error) {
	kv_store.mu.RLock()
	defer kv_store.mu.RUnlock()

	if (len(commands) != 2){
		Err := fmt.Errorf("%w for %s command",ErrNumberArguments, commands[0])
		return "" , false, Err
	}
	val, ok := kv_store.kv_pair[commands[1]]
	return val, ok, nil
}

func deleteFunction(commands []string, kv_store *store) ([]string, error){
	kv_store.mu.RLock()
	defer kv_store.mu.RUnlock()

	delKeys := make([]string, 0)
	if len(commands) < 2{
				Err := fmt.Errorf("%w for %s command",ErrNumberArguments, commands[0])
				return delKeys, Err 
			}


	for i := 1; i < len(commands); i++ {
		if _ , ok := kv_store.kv_pair[commands[i]]; ok{
			delKeys = append(delKeys, commands[i])
		}
	}
	return delKeys, nil 
}

func existFunction(commands []string, kv_store *store) (int, error){
	kv_store.mu.RLock()
	defer kv_store.mu.RUnlock()
	if len(commands) < 2 {
		Err := fmt.Errorf("%w for %s command",ErrNumberArguments, commands[0])
		return 0, Err 
	}

	count := 0
	for i := 1; i < len(commands); i++ {
		if _ , ok := kv_store.kv_pair[commands[i]]; ok {
			count++
		}
	}
	return count, nil 
}

func renameFunction(commands []string, kv_store *store) (string, error){
	kv_store.mu.RLock()
	defer kv_store.mu.RUnlock()
	if len(commands) != 3{
		Err := fmt.Errorf("%w for %s command",ErrNumberArguments, commands[0])
		return "", Err 
	}

	val , ok := kv_store.kv_pair[commands[1]]
	if !ok {
		return "", ErrNokey
	}
	
	return val, nil 
}

func emptyFunction(commands []string, kv_store *store) (int, error){
	kv_store.mu.RLock()
	defer kv_store.mu.RUnlock()
	if len(commands) != 1 {
		Err := fmt.Errorf("%w for %s command",ErrNumberArguments, commands[0])
		return 0, Err
	}
	return len(kv_store.kv_pair), nil 
}

func keysFunction(commands []string, kv_store *store) ([]string, error) {
	kv_store.mu.RLock()
	defer kv_store.mu.RUnlock()
	matched_str := make([]string, 0)
	if len(commands) != 2 {
		Err := fmt.Errorf("%w for %s command",ErrNumberArguments, commands[0])
		return matched_str, Err
	}

	if strings.Count(commands[1], "*") == 0{
		if key, ok := kv_store.kv_pair[commands[1]]; ok{
			matched_str = append(matched_str, key)
		}
		return matched_str, nil 
	}

	regex_str := commands[1]
	if regex_str[len(regex_str)-1] == '*' {
		regex_str = regex_str[:len(regex_str)-1]
		regex_str += "([a-z]+)"
	}

	regex_str = strings.ReplaceAll(regex_str, "*", "([a-z])")
	keys := maps.Keys(kv_store.kv_pair)
	for key := range keys{
		match, err := regexp.Match(regex_str, []byte(key))
		if err != nil{
			Err := fmt.Errorf("%w caused due to %w while finding %s", ErrRegex, err, key)
			return matched_str, Err 
		}
		if match {
			matched_str = append(matched_str, key)
		}
	}
	return matched_str, nil
}

func expireFunction (commands []string, kv_store *store) (error){
	if len(commands) != 3{
		Err := fmt.Errorf("%w for %s command",ErrNumberArguments, commands[0])
		return Err
	}
	
	duration, err := strconv.Atoi(commands[2])
	if err != nil {
		Err := fmt.Errorf("%w because %w",ErrInvalidDuration, err)
		return Err 
	}

	if _, ok := kv_store.kv_pair[commands[1]]; !ok {
		return ErrNokey
	}
	
	if (commands[0] == "expire"){
		go delayedDel(kv_store, commands[1], duration, Second)
	}else{
		go delayedDel(kv_store, commands[1], duration, Milli)
	}

	return nil 
}
