# pcf - wrapper CLI for interacting with OM, BOSH and others for PCF environments
[![Build Status](https://travis-ci.com/pivotal/pcf-cli.svg?token=jUqzM7hyJNi7CRu5xyLL&branch=master)](https://travis-ci.com/pivotal/pcf-cli)

## Installation

You can build `pcf` from the source if you have Go:
```bash
$ git clone git@github.com:pivotal/pcf-cli.git && cd pcf-cli && go install github.com/pivotal/pcf-cli/cmd/pcf
```

## Development

Unit and integration tests can be run if you have [Ginkgo](https://github.com/onsi/ginkgo) installed:
```bash
$ ginkgo -r .
```

