# mwatch
Watching the metrics of a machine learning model

mwatch is the main module 

It calls runmetrics in the generate package

--> Currently generate package calls 3 distributions, samples them and exposes them via prometheus client to the /metrics endpoint on the localhost. This could be replaced with your own metrics sampler which samples metrics from your machine learning model or application

Once the http server which serves the metrics is up, mwatch calls the watch targetwatcher to scrape the metrics off the endpoint which is given via the configuraton listenaddress parameter. This listenaddress must be the same as the address of your metrics endpoint

Main mwatch uses pool.New to keep the scraped metrics in a buffer. You could customise this to append to the buffer and add them to a data store for further analysis. You could also use the buffer to analyse and provide statistics on concept drift if required.

Compilation and Installation
->git clone https://github.com/samuel-sujith/mwatch
->cd mwatch
->go mod tidy
->go install ./mwatch/main.go

Running the watcher
-> With default parameters
main.exe
-> WIth custom endpoint
main.exe --listenaddress=<customendpoint>

Please post issues in case of any assistance.
