/*
Copyright 2018 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/emicklei/go-restful"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/component-base/logs"
	"k8s.io/klog"

	"github.com/samuel-sujith/mwatch/pkg/types"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	basecmd "github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/cmd"
	fakeprov "github.com/samuel-sujith/mwatch/pkg/provider"
	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
)


//SampleAdapter encasulates the apiserver base
//Options are added to it to incorporate into metrics server
type SampleAdapter struct {
	basecmd.AdapterBase

	// Message is printed on succesful startup
	Message string
}
func (a *SampleAdapter) makeProviderOrDie(intconf types.Interimconfig, configuration types.Cfg) (provider.MetricsProvider, *restful.WebService) {
	client, err := a.DynamicClient()
	if err != nil {
		klog.Fatalf("unable to construct dynamic client: %v", err)
	}

	mapper, err := a.RESTMapper()
	if err != nil {
		klog.Fatalf("unable to construct discovery REST mapper: %v", err)
	}

	return fakeprov.NewFakeProvider(client, mapper, intconf, configuration)
}

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	w := log.NewSyncWriter(os.Stderr)
	logger := log.NewLogfmtLogger(w)

	configuration := types.Cfg{Listenaddress: "a dummy address", DesiredMetric: "a dummy metric"}
	var cert, key string
	var skipServerCertCheck bool

	a := kingpin.New(filepath.Base(os.Args[0]), "The Model metric watcher")

	a.Flag("listenaddress", "Model watcher address to watch.").Default("http://localhost:8080/metrics").StringVar(&configuration.Listenaddress)
	a.Flag("desired_metric", "Desired metric to watch out for").Default("process_open_fds").StringVar(&configuration.DesiredMetric)
	a.Flag("cert", "certificate file for client.").Default("").StringVar(&cert)
	a.Flag("key", "key for certificate file for client.").Default("").StringVar(&key)
	a.Flag("accept-invalid-cert", "Skipping cert check").Default("true").BoolVar(&skipServerCertCheck)

	cmd := &SampleAdapter{}
	cmd.Flags().StringVar(&cmd.Message, "msg", "starting adapter...", "startup message")
	/*cmd.Flags().AddGoFlagSet(flag.CommandLine) // make sure we get the klog flags
	cmd.Flags().Parse(os.Args)*/

	interimconfig := types.Interimconfig{
		Configuration: configuration,
		Logger: logger,
		Cert: cert,
		Key: key,
		SkipServerCertCheck: skipServerCertCheck,
	}

	testProvider, webService := cmd.makeProviderOrDie(interimconfig, configuration)

	
	cmd.WithCustomMetrics(testProvider)
	cmd.WithExternalMetrics(testProvider)

	level.Info(logger).Log("msg", cmd.Message)
	// Set up POST endpoint for writing fake metric values
	restful.DefaultContainer.Add(webService)
	go func() {
		// Open port for POSTing fake metrics
		//klog.Fatal(http.ListenAndServe(":8080", nil))
		level.Error(logger).Log("msg", "unable to make transport", "error", http.ListenAndServe(":8080", nil))
	}()
	if err := cmd.Run(wait.NeverStop); err != nil {
		//klog.Fatalf("unable to run custom metrics adapter: %v", err)
		level.Error(logger).Log("msg", "unable to run custom metrics adapter", "error", err)
	}
}
