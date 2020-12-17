package traffic

import (
	"fmt"
	"log"
	"os"
	"time"
)

type writer struct{}

func (w *writer) Write(bytes []byte) (int, error) {
	return fmt.Fprintf(os.Stderr, "%s %s", time.Now().UTC().Format(time.RFC3339Nano), string(bytes))
}

func LogRFC3339Nano() {
	log.SetFlags(0)
	log.SetOutput(&writer{})
}
