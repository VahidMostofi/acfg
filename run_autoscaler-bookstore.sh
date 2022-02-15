#!/bin/bash
go build -o acfg

export LOG_LEVEL=info
set -e # stop on first error
set -a # automatically export all variables
source dev.env
set +a

#./acfg autoscaling --config $PWD/sample-configs/bookstore.yml hybrid --hpat 30 --interval 30 --usecount 5 --allocationsFile $PWD/pythons/cd05d7b4445349ee645ea290586fd28c0c675a155eb1522485535c5c0329a908workload_ranges_top_5_with_replica_train_sorted_bnv2.json
#./acfg autoscaling --config $PWD/sample-configs/bookstore.yml hybrid --hpat 30 --interval 20 --usecount 0 --allocationsFile $PWD/pythons/fargate_cd05d7b4445349ee645_ranges_top_5_without_replica_train_sorted.json
#./acfg autoscaling --config $PWD/sample-configs/bookstore.yml hybrid --hpat 50 --interval 30 --usecount 10 --allocationsFile $PWD/pythons/cd05d7b4445349ee645ea290586fd28c0c675a155eb1522485535c5c0329a908workload_ranges_top_10_with_replica_train_sorted_bnv2.json
./acfg autoscaling --config $PWD/sample-configs/bookstore.yml hybrid --hpat 30 --interval 30 --usecount 10 --allocationsFile $PWD/pythons/cd05d7b4445349ee645ea290586fd28c0c675a155eb1522485535c5c0329a908workload_ranges_top_10_with_replica_train_sorted_bnv2.json
