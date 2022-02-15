
#!/bin/bash
set -e # stop on first error

PATH_TO_ACFG=/home/vahidm/workspace/acfg/acfg
PATH_TO_CONFIG_FILE=/home/vahidm/workspace/acfg/sample-configs/bookstore.yml

export LOG_LEVEL=info

set -a # automatically export all variables
source ./dev.env

# #####################   00   #######################
# export ACFG_LOADGENERATOR_ARGS_ARGSVUS="24"
# export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.33208333333333334"
# export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.33208333333333334"
# export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.33208333333333334"
# export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="1.0"
# export ACFG_TESTNAME=predefined_configs_azure_app_cd05d7b444_3x_00

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 1000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 50

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 30

# #####################   01   #######################
# export ACFG_LOADGENERATOR_ARGS_ARGSVUS="168"
# export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.33208333333333334"
# export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.33208333333333334"
# export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.33208333333333334"
# export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="1.0"
# export ACFG_TESTNAME=predefined_configs_azure_app_cd05d7b444_3x_01

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 1000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 50

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 30

# #####################   02   #######################
# export ACFG_LOADGENERATOR_ARGS_ARGSVUS="192"
# export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.33208333333333334"
# export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.33208333333333334"
# export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.33208333333333334"
# export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="1.0"
# export ACFG_TESTNAME=predefined_configs_azure_app_cd05d7b444_3x_02

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 1000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 50

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 30

# #####################   03   #######################
# export ACFG_LOADGENERATOR_ARGS_ARGSVUS="120"
# export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.33208333333333334"
# export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.33208333333333334"
# export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.33208333333333334"
# export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="1.0"
# export ACFG_TESTNAME=predefined_configs_azure_app_cd05d7b444_3x_03

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 1000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 50

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 30

# #####################   04   #######################
# export ACFG_LOADGENERATOR_ARGS_ARGSVUS="144"
# export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.33208333333333334"
# export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.33208333333333334"
# export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.33208333333333334"
# export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="1.0"
# export ACFG_TESTNAME=predefined_configs_azure_app_cd05d7b444_3x_04

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 1000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 50

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 30

#####################   05   #######################
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="216"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="1.0"
export ACFG_TESTNAME=predefined_configs_azure_app_cd05d7b444_3x_05

$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 1000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 50

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 30

#####################   06   #######################
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="96"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="1.0"
export ACFG_TESTNAME=predefined_configs_azure_app_cd05d7b444_3x_06

$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 1000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 50

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 30

#####################   07   #######################
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="240"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="1.0"
export ACFG_TESTNAME=predefined_configs_azure_app_cd05d7b444_3x_07

$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 1000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 50

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 30

#####################   08   #######################
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="264"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="1.0"
export ACFG_TESTNAME=predefined_configs_azure_app_cd05d7b444_3x_08

$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 1000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 50

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 30

#####################   09   #######################
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="72"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.33208333333333334"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="1.0"
export ACFG_TESTNAME=predefined_configs_azure_app_cd05d7b444_3x_09

$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 1000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 50

# $PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE hpa --initcpu 500 --initmem 256 --threshold 30


set +a
