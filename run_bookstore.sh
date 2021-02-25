#!/bin/bash
go build -o acfg
set -e # stop on first error
set -a # automatically export all variables
source dev.env
set +a
# ./acfg autoconfig --config sample-configs/bookstore.yml --name bookstore-test cput --indicator mean --initcpu 500 --initmem 256 --threshold 45

./acfg autoconfig --config sample-configs/bookstore.yml --name bookstore-bnv2 bnv2 --initialdelta 2000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

# MOBO
PYTHON_PATH=/home/vahid/.pyenv/versions/3.6.13/bin/python3
SCRIPT_PATH=$PWD/pythons/mobo-bookstore.py
echo $SCRIPT_PATH
echo $PYTHON_PATH
./acfg autoconfig --config sample-configs/bookstore.yml --name bookstore-mobo stdin --initcpu 500 --initmem 256 --pythonpath $PYTHON_PATH --scriptpath $SCRIPT_PATH