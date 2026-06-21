package aof

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

// Log serializes file access so concurrent writes cannot interleave.
type Log struct {
	mu   sync.Mutex
	path string
}

func New(path string) *Log {
	return &Log{path: path}
}

func (l *Log) Append(commands []string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(strings.Join(commands, " ") + "\n")
	return err
}

func (l *Log) Replay(apply func([]string) error) error {
	f, err := os.Open(l.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Fields collapses repeated spaces and prevents empty commands from reaching dispatch.
		commands := strings.Fields(line)
		if err := apply(commands); err != nil {
			return fmt.Errorf("replay %q: %w", line, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
