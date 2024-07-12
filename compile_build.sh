#!/usr/bin/bash

docker run --rm -i -t \
    -v $PWD:/go/src/github.com/coredns/coredns -w /go/src/github.com/coredns/coredns \
        golang:1.21 sh -c 'GOFLAGS="-buildvcs=false" make gen && GOFLAGS="-buildvcs=false" make'

docker build -t coredns-dpt .
docker tag coredns-dpt hjsjhn/coredns-dpt