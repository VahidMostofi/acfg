#!/bin/bash
set -e # stop on first error

go build -o acfg
PATH_TO_ACFG=/home/vahid/workspace/acfg/acfg
PATH_TO_CONFIG_FILE=/home/vahid/workspace/acfg/sample-configs/bookstore.yml

export LOG_LEVEL=info

set -a # automatically export all variables
source dev.env


export ACFG_LOADGENERATOR_ARGS_ARGSVUS="112"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.235"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.105"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.66"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="0.5"
export ACFG_TESTNAME=predefined_configs_00

$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 2000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

set +a