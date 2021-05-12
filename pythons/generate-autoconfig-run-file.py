import os
import sys
import json
import math
endpoints = ["login", "getbook", "editbook"]

if len(sys.argv) < 3: 
  print('you must pass 2 args, path of the json file and a prefix name of the workloads.')
  sys.exit(1)

workload_range_conditions_file = sys.argv[1]
workload_name=sys.argv[2]

with open(workload_range_conditions_file) as f:
  workload_ranges=json.load(f)

command = """
#!/bin/bash
set -e # stop on first error

PATH_TO_ACFG=/home/vahid/workspace/acfg/acfg
PATH_TO_CONFIG_FILE=/home/vahid/workspace/acfg/sample-configs/bookstore.yml

export LOG_LEVEL=info

set -a # automatically export all variables
source ../dev.env

"""

for idx, item in enumerate(workload_ranges):
  sub_name = str(idx)
  if len(sub_name) == 1:
    sub_name = "0" + sub_name
  
  command += "#####################   " + sub_name + "   #######################\n"
  
  workload_range = item["workload-range"]
  total = 0
  for endpoint in workload_range:
    total += workload_range[endpoint]["high"]
  total = int(math.ceil(total))
  envs = "export ACFG_LOADGENERATOR_ARGS_ARGSVUS=\""+str(total)+"\""
  for endpoint in workload_range:
    envs += "\n"
    envs += "export ACFG_LOADGENERATOR_ARGS_ARGS"+endpoint.upper()+"PROB=\""+str(workload_range[endpoint]["high"]/total)+"\""
  envs += "\n"
  envs += "export ACFG_LOADGENERATOR_ARGS_ARGSSLEEPDURATION=\"0.5\"\n"
  envs += "export ACFG_TESTNAME=predefined_configs_" + workload_name + '_' + sub_name + "\n"
  
  
  command += envs + "\n"
  command += "$PATH_TO_ACFG autoconfig --config $PATH_TO_CONFIG_FILE bnv2 --initialdelta 2000 --initcpu 500 --initmem 256 --maxcpuperreplica 500 --mincpu 500 --mindelta 500\n\n"
  
command += "\nset +a"

print(command)
