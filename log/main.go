package log

import (
	"bytes"
	"fmt"
	"github.com/egymgmbh/go-prefix-writer/prefixer"
	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"runtime"
	"time"
)

func ColoredLinePrefixWriter(w io.Writer, group string, colorV color.Attribute) io.Writer {
	var buf bytes.Buffer
	color.New(colorV).Fprint(&buf, group)
	extender := "="
	for i := 10 - len(group); i > 0; i-- {
		extender += "="
	}

	prefix := fmt.Sprintf("[%s] %s> ", buf.String(), extender)
	return prefixer.New(w, func() string {
		return prefix
	})
}

func InitializeGlobalLogger(myName string, prettyInsteadOfJSON bool) {
	// make sure we dont log reqs/resp in grpc
	output := ColoredLinePrefixWriter(os.Stdout, myName, color.FgGreen)

	go func() {
		for {
			<-time.After(time.Second * 30)
			log.Debug().Msgf("Current goroutine count %d", runtime.NumGoroutine())
		}
	}()

	if prettyInsteadOfJSON {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: output}).With().Timestamp().Logger()
	} else {
		log.Logger = zerolog.New(output).With().Timestamp().Logger()
	}
}
