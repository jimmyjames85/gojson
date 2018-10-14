#!/bin/bash
rm -rf /tmp/profile*
rm -f ./gj
go build -o ./gj cmd/gj/main.go
./gj > /dev/null 2>&1

# echo file is at: $thefile
# now=`date +%s`
echo 'go tool pprof --pdf ./gj /tmp/profile.../cpu.pprof > pprof.pdf'
