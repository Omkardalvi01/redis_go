package engine

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Omkardalvi01/redis_go.git/internal/aof"
	"github.com/Omkardalvi01/redis_go.git/internal/kv"
)

var (
	ErrEmptyInput      = errors.New("(error) empty command")
	ErrNumberArguments = errors.New("(error) wrong number of arguments")
	ErrUnknownCmd      = errors.New("(error) unknown command")
	ErrInvalidDuration = errors.New("(error) duration is invalid")
)

type Config struct {
	LogPath   string
	EnableAOF bool
}

type Result struct {
	Response string
	Err      error
	Run      bool
}

type Engine struct {
	store     *kv.Store
	log       *aof.Log
	enableAOF bool
}

func New(cfg Config) (*Engine, error) {
	eng := &Engine{
		store:     kv.NewStore(),
		enableAOF: cfg.EnableAOF,
	}

	if cfg.EnableAOF {
		eng.log = aof.New(cfg.LogPath)
		if err := eng.log.Replay(func(commands []string) error {
			return eng.execute(commands, false).Err
		}); err != nil {
			return nil, err
		}
	}

	return eng, nil
}

func (e *Engine) ExecuteLine(line string, persist bool) Result {
	line = strings.TrimSpace(line)
	if line == "" {
		return Result{Err: ErrEmptyInput, Run: true}
	}

	commands := strings.Fields(line)
	return e.execute(commands, persist)
}

func (e *Engine) execute(commands []string, persist bool) Result {
	if len(commands) == 0 {
		return Result{Err: ErrEmptyInput, Run: true}
	}

	result := Result{Run: true}
	switch strings.ToLower(commands[0]) {
	case "ping":
		result.Response = "PONG\n"

	case "set":
		if err := validateArity(commands, 3); err != nil {
			result.Err = err
			break
		}
		e.store.Set(commands[1], commands[2])
		result.Response = "OK\n"
		result.Err = e.appendIfNeeded(commands, persist)

	case "get":
		if err := validateArity(commands, 2); err != nil {
			result.Err = err
			break
		}
		value, ok := e.store.Get(commands[1])
		if !ok {
			result.Response = "(nil)\n"
			break
		}
		result.Response = value + "\n"

	case "del":
		if err := validateMinArity(commands, 2); err != nil {
			result.Err = err
			break
		}
		deleted := e.store.Delete(commands[1:]...)
		result.Response = fmt.Sprintf("(integer) %d\n", deleted)
		result.Err = e.appendIfNeeded(commands, persist)

	case "exist":
		if err := validateMinArity(commands, 2); err != nil {
			result.Err = err
			break
		}
		count := e.store.Exists(commands[1:]...)
		result.Response = fmt.Sprintf("(integer) %d\n", count)

	case "rename":
		if err := validateArity(commands, 3); err != nil {
			result.Err = err
			break
		}
		if _, err := e.store.Rename(commands[1], commands[2]); err != nil {
			result.Err = err
			break
		}
		result.Response = "OK\n"
		result.Err = e.appendIfNeeded(commands, persist)

	case "empty":
		if err := validateArity(commands, 1); err != nil {
			result.Err = err
			break
		}
		result.Response = fmt.Sprintf("(integer) %d\n", e.store.Len())

	case "keys":
		if err := validateArity(commands, 2); err != nil {
			result.Err = err
			break
		}
		matched, err := e.store.Keys(commands[1])
		if err != nil {
			result.Err = err
			break
		}
		if len(matched) == 0 {
			result.Response = "(empty array)\n"
			break
		}
		var b strings.Builder
		for i, key := range matched {
			fmt.Fprintf(&b, "%d) %s\n", i, key)
		}
		result.Response = b.String()

	case "expire", "pexpire":
		if err := validateArity(commands, 3); err != nil {
			result.Err = err
			break
		}
		result.Err = e.scheduleDelete(commands)
		if result.Err == nil {
			result.Response = "OK\n"
			result.Err = e.appendIfNeeded(commands, persist)
		}

	case "exit":
		result.Response = "OK\n"
		result.Run = false

	default:
		result.Err = fmt.Errorf("%w %s", ErrUnknownCmd, commands[0])
	}

	return result
}

func (e *Engine) appendIfNeeded(commands []string, persist bool) error {
	if !persist || !e.enableAOF || e.log == nil {
		return nil
	}
	return e.log.Append(commands)
}

func validateArity(commands []string, expected int) error {
	if len(commands) != expected {
		return fmt.Errorf("%w for %s command", ErrNumberArguments, commands[0])
	}
	return nil
}

func validateMinArity(commands []string, min int) error {
	if len(commands) < min {
		return fmt.Errorf("%w for %s command", ErrNumberArguments, commands[0])
	}
	return nil
}
