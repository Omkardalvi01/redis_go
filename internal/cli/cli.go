package cli

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Omkardalvi01/redis_go.git/internal/engine"
)

func Run(scanner *bufio.Scanner, eng *engine.Engine, out io.Writer, errOut io.Writer) {
	for scanner.Scan() {
		result := eng.ExecuteLine(scanner.Text(), true)
		if result.Response != "" {
			_, _ = fmt.Fprint(out, result.Response)
		}
		if result.Err != nil {
			_, _ = fmt.Fprintln(errOut, result.Err)
		}
		if !result.Run {
			return
		}
	}

	if err := scanner.Err(); err != nil {
		_, _ = fmt.Fprintln(errOut, err)
	}
}
