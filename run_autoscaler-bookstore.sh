#!/bin/bash
go build -o acfg

export LOG_LEVEL=info
set -e # stop on first error
set -a # automatically export all variables
source dev.env
set +a

./acfg autoscaling --config $PWD/sample-configs/bookstore.yml hybrid --hpat 50 --interval 20 --usecount 0