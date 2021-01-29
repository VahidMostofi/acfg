#!/bin/bash
go build -o acfg
set -e # stop on first error
set -a # automatically export all variables
source dev.env
set +a
./acfg autoconfig --config sample-configs/teastore.yml --name teastore-test cput --indicator mean --initcpu 250 --initmem 2048 --threshold 45