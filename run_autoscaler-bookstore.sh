#!/bin/bash
go build -o acfg
set -e # stop on first error
set -a # automatically export all variables
source dev.env
set +a

./acfg autoscaling --config $PWD/sample-configs/bookstore.yml hybrid --allocationsFile a --hpat 20 --interval 20 --usecount 1