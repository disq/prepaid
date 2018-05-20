![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)
![Tag](https://img.shields.io/github/tag/disq/prepaid.svg)
[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/disq/prepaid)
[![Go Report](https://goreportcard.com/badge/github.com/disq/prepaid)](https://goreportcard.com/report/github.com/disq/prepaid)

# Prepaid

Prepaid card PoC

## Build / Deploy

    # Clone
    git clone https://github.com/disq/prepaid.git
    cd prepaid

    # Install serverless
    npm install -g serverless

    # Deploy
    make deploy

## Configuration

See [serverless.yml](./serverless.yml) for configuration.

## Deploying to Production

    make dep sls-build
    serverless deploy -s production

# License

MIT.
