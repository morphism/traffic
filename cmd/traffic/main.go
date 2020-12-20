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
		warn           = flag.Duration("warn", time.Second, "Log warning if lagging for more than this duration")

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
			ms := time.Duration(rand.Intn(100)) * time.Millisecond
			time.Sleep(ms)
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

	tick := time.Now().UTC()

	wait := func(n int) {
		// Crudely wait for about *interval/n.
		//
		// ToDo: Use a better algorithm and distribution.
		var (
			max    = 2 * float64(*interval) / float64(n)
			sample = rand.Float64()
			d      = time.Duration(max * sample)
		)
		time.Sleep(d)
	}

LOOP:
	for {
		if 0 < *limit && *limit <= uint64(t) {
			break
		}
		emitted := 0
		n, _ := s.Counts(t)
		for i := 0; i < int(n); i++ {
			line, err := in.ReadString('\n')
			if 0 < len(line) && err != io.EOF {
				wait(int(n))
				if *ts {
					fmt.Printf("%s %s", time.Now().UTC().Format(time.RFC3339Nano), line)
				} else {
					fmt.Printf("%s", line)
				}
			}
			emitted++
			if err == io.EOF {
				break LOOP
			}
		}
		if s.Log {
			log.Printf("traffic %d total %d", t, emitted)
		}

		t++

		// Compute how to to pause then pause for that long.
		var (
			target = tick.Add(*interval)
			now    = time.Now().UTC()
			delta  = target.Sub(now)
		)
		tick = now
		if delta < -*warn {
			log.Printf("warning: lag: %v", -delta)
		}
		if delta < 0 {
			delta = 0
		}
		time.Sleep(delta)
	}

	return nil
}
