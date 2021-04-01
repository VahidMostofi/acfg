import os
import sys
import json
from yaml import load, dump
try:
    from yaml import CLoader as Loader, CDumper as Dumper
except ImportError:
    from yaml import Loader, Dumper
workload_range_conditions_file = "workload_ranges_no_replicas.json"
results_path = "/home/vahid/acfg-results/bookstore/OnlineShoppingStore/BNV2/"
with open(workload_range_conditions_file) as f:
  workload_ranges=json.load(f)

for idx, item in enumerate(workload_ranges):
  sub_name = str(idx)
  if len(sub_name) == 1:
    sub_name = "0" + sub_name
  file_name = "predefined_configs_" + sub_name + ".yaml"

  with open(results_path + file_name) as f:
    data = load(f, Loader=Loader)
  best_config = {}
  total_resouce = 1000000
  for iteration in data["iterations"]:
    t = 0
    for key, value in iteration['configurations'].items():
      t += value['replicacount']
  
    if t < total_resouce and iteration['strategyInfo']['doMeet']:
      total_resouce = t
      best_config = iteration['configurations']
  best_replicas = {}
  for key, value in best_config.items():
    best_replicas[key] = value['replicacount']
  
  workload_ranges[idx]["replicas"] = best_replicas
  break
with open('./workload_ranges_with_replicas.json','w') as f:
  json.dump(workload_ranges, f,indent=4, sort_keys=True)