#!/bin/bash
go build -o acfg
set -e # stop on first error
set -a # automatically export all variables
source dev.env
set +a

./acfg autoscaling --config $PWD/sample-configs/bookstore.yml hybrid --hpat 50 --interval 20 --usecount 2 --allocationsFile $PWD/pythons/workload_ranges_top_10_bin_25_with_replica_train_sorted.json