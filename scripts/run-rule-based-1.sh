#!/bin/bash

set -e # stop on first error

PATH_TO_ACFG="$PWD/../acfg"
PATH_TO_CONFIG_FILE="$PWD/../sample-configs/bookstore.yml"

export LOG_LEVEL=info

set -a # automatically export all variables
source ../dev.env
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="75"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.33"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.33"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.34"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="1.0"

export ACFG_TESTNAME=rule_based_1_75_threshold_30
$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE cput --indicator mean --initcpu 500 --initmem 512 --threshold 30

export ACFG_TESTNAME=rule_based_1_75_threshold_50
$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE cput --indicator mean --initcpu 500 --initmem 512 --threshold 50

export ACFG_TESTNAME=rule_based_1_75_threshold_70
$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE cput --indicator mean --initcpu 500 --initmem 512 --threshold 70

##############
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="100"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.33"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.33"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.34"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="1.0"

export ACFG_TESTNAME=rule_based_1_100_threshold_30
$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE cput --indicator mean --initcpu 500 --initmem 512 --threshold 30

export ACFG_TESTNAME=rule_based_1_100_threshold_50
$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE cput --indicator mean --initcpu 500 --initmem 512 --threshold 50

export ACFG_TESTNAME=rule_based_1_100_threshold_70
$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE cput --indicator mean --initcpu 500 --initmem 512 --threshold 70

################
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="125"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.33"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.33"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.34"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="1.0"

export ACFG_TESTNAME=rule_based_1_125_threshold_30
$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE cput --indicator mean --initcpu 500 --initmem 512 --threshold 30

export ACFG_TESTNAME=rule_based_1_125_threshold_50
$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE cput --indicator mean --initcpu 500 --initmem 512 --threshold 50

export ACFG_TESTNAME=rule_based_1_125_threshold_70
$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE cput --indicator mean --initcpu 500 --initmem 512 --threshold 70

################
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="150"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.33"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.33"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.34"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="1.0"

export ACFG_TESTNAME=rule_based_1_150_threshold_30
$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE cput --indicator mean --initcpu 500 --initmem 512 --threshold 30

export ACFG_TESTNAME=rule_based_1_150_threshold_50
$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE cput --indicator mean --initcpu 500 --initmem 512 --threshold 50

export ACFG_TESTNAME=rule_based_1_150_threshold_70
$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE cput --indicator mean --initcpu 500 --initmem 512 --threshold 70
set +a