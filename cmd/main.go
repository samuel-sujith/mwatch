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
	"flag"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	
	
	"github.com/samuel-sujith/mwatch/generate"


	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)


func main() {

	cfg := struct {
		configFile string
	}

	a := kingpin.New(filepath.Base(os.Args[0]), "The Model metric watcher")

	a.Flag("config.file", "Model watcher configuration file path.").Default("mwatch.yml").StringVar(&cfg.configFile)
	
	_, err := a.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error parsing commandline arguments"))
		a.Usage(os.Args[1:])
		os.Exit(2)
	}


	ch := make(chan bool)
	buf := &bytes.Buffer{}
	logger := log.NewLogfmtLogger(buf)

	var g run.Group
	{
		g.add(

			func() error {
				go generate.Runmetrics(ch)
				return nil
			},
			func(err error) {
				level.Error(logger).Log("err in generating metrics", err)
				os.exit(3)
			},
		)
		g.add(

			func() error {
				<-ch
				level.Info(logger).Log("Generated metrics")
				//TODO
				return nil
			},
			func(err error) {
				level.Error(logger).Log("err in scraping metrics", err)
				os.exit(3)
			},
		)

	}
	if err := g.Run(); err != nil {
		level.Error(logger).Log("err in running the Ok groups", err)
		os.Exit(1)
	}
	level.Info(logger).Log("msg", "See you next time!")
}
