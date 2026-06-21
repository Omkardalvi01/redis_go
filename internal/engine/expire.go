package engine

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func parseDuration(raw string) (int, error) {
	duration, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%w because %w", ErrInvalidDuration, err)
	}
	return duration, nil
}

func (e *Engine) scheduleDelete(commands []string) error {
	duration, err := parseDuration(commands[2])
	if err != nil {
		return err
	}

	key := commands[1]
	if _, ok := e.store.Get(key); !ok {
		return fmt.Errorf("(error) no such key")
	}

	// The timer is isolated so the command returns immediately while deletion happens later.
	delay := time.Second
	if strings.EqualFold(commands[0], "pexpire") {
		delay = time.Millisecond
	}

	go func() {
		timer := time.NewTimer(time.Duration(duration) * delay)
		defer timer.Stop()
		<-timer.C
		e.store.Delete(key)
	}()

	return nil
}
