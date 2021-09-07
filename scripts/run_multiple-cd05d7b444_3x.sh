#!/bin/bash
set -e # stop on first error

PATH_TO_ACFG=/home/vahidm/workspace/acfg/acfg
PATH_TO_CONFIG_FILE=/home/vahidm/workspace/acfg/sample-configs/bookstore.yml

export LOG_LEVEL=debug

set -a # automatically export all variables
source ../dev.env

#####################   00   #######################
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="24"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="1.0"
export ACFG_TESTNAME=predefined_configs_azure_app_cd05d7b444_00

$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 1000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 50

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 30

#####################   01   #######################
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="168"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="1.0"
export ACFG_TESTNAME=predefined_configs_azure_app_cd05d7b444_01

$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 1000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 50

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 30

#####################   02   #######################
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="192"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="1.0"
export ACFG_TESTNAME=predefined_configs_azure_app_cd05d7b444_02

$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 1000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 50

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 30

#####################   03   #######################
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="120"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="1.0"
export ACFG_TESTNAME=predefined_configs_azure_app_cd05d7b444_03

$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 1000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 50

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 30

#####################   04   #######################
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="144"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="1.0"
export ACFG_TESTNAME=predefined_configs_azure_app_cd05d7b444_04

$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 1000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 50

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 30


set +a