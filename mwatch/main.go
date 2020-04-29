// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// A simple example exposing fictional RPC latencies with different types of
// random distributions (uniform, normal, and exponential) as Prometheus
// metrics.
package main

//
import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/prometheus/pkg/pool"
	"github.com/samuel-sujith/mwatch/pkg/generate"
	"github.com/samuel-sujith/mwatch/pkg/types"
	"github.com/samuel-sujith/mwatch/pkg/watch"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {

	configuration := types.Cfg{Listenaddress: "a dummy address"}

	a := kingpin.New(filepath.Base(os.Args[0]), "The Model metric watcher")

	a.Flag("listenaddress", "Model watcher address to watch.").Default("http://localhost:8080/metrics").StringVar(&configuration.Listenaddress)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error parsing commandline arguments"))
		a.Usage(os.Args[1:])
		os.Exit(2)
	}

	/*conf, err := config.Loadfile(configuration.configfile)
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "couldn't load configuration (--config.file=%q)", configuration.configfile))
	}*/

	ch := make(chan bool)
	//buf := &bytes.Buffer{}
	w := log.NewSyncWriter(os.Stderr)
	logger := log.NewLogfmtLogger(w)

	//fmt.Println("About to run boiiiii")

	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(time.Second)

	go generate.Runmetrics(ch)

	<-ch

	level.Info(logger).Log("msg", "Generated metrics")
	buffers := pool.New(1e3, 1e6, 3, func(sz int) interface{} { return make([]byte, 0, sz) })
	b := buffers.Get(1024).([]byte)
	buf := bytes.NewBuffer(b)
	//fmt.Println("configuration is ", configuration)
mainLoop:
	for {

		select {
		case <-term:
			level.Info(logger).Log("msg", "exiting due to interrupt")
			break mainLoop
		case <-ticker.C:
			contenttype, watcherr := watch.Targetwatching(configuration, buf, logger)
			if watcherr == nil {
				b = buf.Bytes()
				level.Info(logger).Log("msg", "nil response from watcher", "ctype", contenttype)

				//fmt.Println("The content type from the watcher server is ", contenttype)
				//fmt.Println("The response from the watcher is ", b)
			}

			if watcherr != nil {
				//fmt.Println("There is error in watching the metrics", watcherr)
				level.Error(logger).Log("msg", "err in target watching", "err", err)
			}
			//TODO
			buffers.Put(b)
		}
	}
	//fmt.Println("put into the buffer")
	level.Info(logger).Log("msg", "See you next time!")
}
