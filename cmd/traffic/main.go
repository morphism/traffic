package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"gihub.com/morphism/traffic"
	"github.com/jsccast/yaml"
	"golang.org/x/exp/rand"
)

func main() {
	// script -q -c "traffic -test-source" /dev/null | traffic

	traffic.LogRFC3339Nano()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var (
		testInput = flag.Bool("test-source", false, "Just run a test input source")

		interval       = flag.Duration("interval", time.Second, "Tick duration")
		seed           = flag.Uint64("seed", 0, "Seed for the RNG (defaults to current time in nanoseconds)")
		configFilename = flag.String("config", "traffic.json", "Name of configuration file")
		limit          = flag.Uint64("limit", 0, "Number of ticks (0 means run forever)")
		ts             = flag.Bool("timestamps", false, "Prefix each line with current timestamp")
		logging        = flag.Bool("log", false, "Turn on some logging output")

		s traffic.System
		t = int64(0)

		in = bufio.NewReader(os.Stdin)
	)

	flag.Parse()

	if *testInput {
		// stdbuf -o0
		i := int64(0)
		for {
			fmt.Printf("%09d %s\n", i, time.Now().UTC().Format(time.RFC3339Nano))
			i++
		}
		return nil // Won't get here.
	}

	{
		bs, err := ioutil.ReadFile(*configFilename)
		if err != nil {
			panic(fmt.Errorf("can't read '%s': %v", *configFilename, err))
		}

		if strings.HasSuffix(*configFilename, ".yaml") {
			err = yaml.Unmarshal(bs, &s)
		} else if strings.HasSuffix(*configFilename, ".json") {
			err = json.Unmarshal(bs, &s)
		} else {
			err = fmt.Errorf("unknown config file syntax ('%s')", *configFilename)
		}
		if err != nil {
			panic(err)
		}
	}

	s.Log = *logging

	if *seed == 0 {
		*seed = uint64(time.Now().UnixNano())
	}
	src := rand.NewSource(*seed)

	if err := s.Init(src); err != nil {
		return err
	}

LOOP:
	for {
		if 0 < *limit && *limit <= uint64(t) {
			break
		}
		emitted := 0
		for _, n := range s.Counts(t) {
			for i := 0; i < int(n); i++ {
				line, err := in.ReadString('\n')
				if *ts {
					fmt.Printf("%s %s", time.Now().UTC().Format(time.RFC3339Nano), line)
				} else {
					fmt.Printf("%s", line)
				}
				emitted++
				if err == io.EOF {
					break LOOP
				}
			}
		}
		if s.Log {
			log.Printf("traffic %d total %d", t, emitted)
		}

		t++
		time.Sleep(*interval)
	}

	return nil
}
