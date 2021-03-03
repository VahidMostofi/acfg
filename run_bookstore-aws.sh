#!/bin/bash
go build -o acfg
set -e # stop on first error
set -a # automatically export all variables
source dev.env
set +a
# Note that you need to change the config file for differnt workloads.
# # BNV Delta=2
# ./acfg autoconfig --config sample-configs/bookstore-aws.yml --name bookstore-bnv2-2-aws-200-1.0 bnv2 --initialdelta 2000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500
# ./acfg autoconfig --config sample-configs/bookstore-aws.yml --name bookstore-bnv2-2-aws-300-1.0 bnv2 --initialdelta 2000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

# MOBO
PYTHON_PATH=/home/vahid/.virtualenvs/py37-tf1-gpu/bin/python
SCRIPT_PATH=$PWD/pythons/mobo-bookstore.py
echo $SCRIPT_PATH
echo $PYTHON_PATH
./acfg autoconfig --config sample-configs/bookstore-aws.yml --name bookstore-bnv2-2-aws-300-1.0 stdin --initcpu 500 --initmem 256 --pythonpath $PYTHON_PATH --scriptpath $SCRIPT_PATH