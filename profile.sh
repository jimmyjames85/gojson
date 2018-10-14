#!/bin/bash
rm -f ./gj
go build -o ./gj cmd/gj/main.go
./gj

# echo file is at: $thefile
# now=`date +%s`
echo 'go tool pprof --pdf ./gj /tmp/profile.../cpu.pprof > pprof.pdf'
