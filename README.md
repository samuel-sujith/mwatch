# mwatch

The repository components set up a custom metrics server which scrapes prometheus using prom2json package and exposing as custom metrics.

deploy manifests are included which will help you install this in your cluster

Clone the repo and run the following command in the cloned repo
docker built -t <customname>:<tagname> .

Once the image is built, substitute into your manifests.

Please post issues in case of any assistance.
