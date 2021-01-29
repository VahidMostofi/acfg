#!/bin/bash
source dev.env
go run main.go --config sample-configs/teastore.yml --name teastore-test cput --indicator mean --initcpu 250 --initmem 2048 --threshold 45