
#!/bin/bash
set -e # stop on first error

PATH_TO_ACFG=/home/vahid/workspace/acfg/acfg
PATH_TO_CONFIG_FILE=/home/vahid/workspace/acfg/sample-configs/bookstore.yml

export LOG_LEVEL=info

set -a # automatically export all variables
source ../dev.env

#####################   00   #######################
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="168"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.232142857143"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.557738095238"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.205357142857"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="0.5"
export ACFG_TESTNAME=predefined_configs_00

$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 2000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

#####################   01   #######################
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="210"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.185714285714"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.647619047619"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.164285714286"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="0.5"
export ACFG_TESTNAME=predefined_configs_01

$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 2000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

#####################   02   #######################
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="125"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.312"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.4104"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.276"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="0.5"
export ACFG_TESTNAME=predefined_configs_02

$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 2000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

#####################   03   #######################
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="252"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.154761904762"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.70753968254"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.136904761905"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="0.5"
export ACFG_TESTNAME=predefined_configs_03

$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 2000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

#####################   04   #######################
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="295"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.132203389831"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.74813559322"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.116949152542"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="0.5"
export ACFG_TESTNAME=predefined_configs_04

$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 2000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

#####################   05   #######################
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="207"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.376811594203"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.452657004831"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.166666666667"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="0.5"
export ACFG_TESTNAME=predefined_configs_05

$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 2000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

#####################   06   #######################
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="249"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.313253012048"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.546184738956"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.138554216867"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="0.5"
export ACFG_TESTNAME=predefined_configs_06

$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 2000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

#####################   07   #######################
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="291"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.268041237113"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.612714776632"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.118556701031"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="0.5"
export ACFG_TESTNAME=predefined_configs_07

$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 2000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

#####################   08   #######################
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="337"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.115727002967"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.780415430267"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.10237388724"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="0.5"
export ACFG_TESTNAME=predefined_configs_08

$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 2000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500

#####################   09   #######################
export ACFG_LOADGENERATOR_ARGS_ARGSVUS="334"
export ACFG_LOADGENERATOR_ARGS_ARGSLOGINPROB="0.233532934132"
export ACFG_LOADGENERATOR_ARGS_ARGSEDITBOOKPROB="0.660778443114"
export ACFG_LOADGENERATOR_ARGS_ARGSGETBOOKPROB="0.103293413174"
export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION="0.5"
export ACFG_TESTNAME=predefined_configs_09

$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 2000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500


set +a
