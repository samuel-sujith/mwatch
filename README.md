# mwatch
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fsamuel-sujith%2Fmwatch.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fsamuel-sujith%2Fmwatch?ref=badge_shield)


The repository components set up a custom metrics server which scrapes prometheus using prom2json package and exposing as custom metrics.

deploy manifests are included which will help you install this in your cluster

Clone the repo and run the following command in the cloned repo
docker built -t <customname>:<tagname> .

Once the image is built, substitute into your manifests.

Create your own secrets in the namespace where you deploy the manifests. Name of secret is cm-adapter-serving-certs.

Please post issues in case of any assistance.


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fsamuel-sujith%2Fmwatch.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fsamuel-sujith%2Fmwatch?ref=badge_large)