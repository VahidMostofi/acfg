#!/bin/bash
go build -o acfg
set -e # stop on first error
set -a # automatically export all variables
source dev.env
set +a
# ./acfg autoconfig --config sample-configs/bookstore.yml --name bookstore-test cput --indicator mean --initcpu 500 --initmem 512 --threshold 45

./acfg autoconfig --config sample-configs/bookstore.yml --name bookstore-test bnv2 --initialdelta 2000 --initcpu 500 --initmem 512 --maxcpuperreplica 500 --mincpu 500 --mindelta 500