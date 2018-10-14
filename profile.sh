#!/bin/bash

rm -rf /tmp/profile* # todo maybe dont do this

go build -o ./gj cmd/gj/main.go
profile=`./gj -s 2>&1 | tail -1 | awk '{print $7}'`
echo profile: ${profile}

timestamp=`date +%s`
pprof=${timestamp}_pprof.pdf

go tool pprof --pdf ./gj ${profile} > $pprof
rm -f ./gj

atril $pprof & disown
