import os
import sys
import json
from yaml import load, dump
try:
    from yaml import CLoader as Loader, CDumper as Dumper
except ImportError:
    from yaml import Loader, Dumper

workload_name = sys.argv[1] #'train_23_test_24_1x'
workload_range_conditions_file = sys.argv[2] # 'workload_ranges_top_10_bin_25_no_replica_train_sorted.json'
results_path = sys.argv[3] # "/home/vahid/acfg-results/bookstore/OnlineShoppingStore/BNV2/"
output_name = sys.argv[4] # 'workload_ranges_top_10_bin_25_with_replica_train_sorted.json' 
with open(workload_range_conditions_file) as f:
  workload_ranges=json.load(f)
durations = []
for idx, item in enumerate(workload_ranges):
  sub_name = str(idx)
  if len(sub_name) == 1:
    sub_name = "0" + sub_name
  file_name = "predefined_configs_" + workload_name + '_' + sub_name + ".yaml"

  if not os.path.exists(results_path + file_name):
    break
  print('reading', file_name)
  with open(results_path + file_name) as f:
    data = load(f, Loader=Loader)
  best_config = {}
  total_resouce = 1000000
  last_start = 0
  for iteration in data["iterations"]:
    t = 0
    for key, value in iteration['configurations'].items():
      t += value['replicacount']
  
    if t < total_resouce and iteration['strategyInfo']['doMeet']:
      total_resouce = t
      best_config = iteration['configurations']
    print(iteration['aggregatedData']['startTime'])
    if iteration['aggregatedData']['startTime'] > 0 and last_start > 0:
      print(iteration['aggregatedData']['startTime'] - last_start)
      durations.append(iteration['aggregatedData']['startTime'] - last_start)
      last_start = iteration['aggregatedData']['startTime']
    elif iteration['aggregatedData']['startTime'] > 0:
      last_start = iteration['aggregatedData']['startTime']
      print('updated last start to ', last_start)

  best_replicas = {}
  for key, value in best_config.items():
    best_replicas[key] = value['replicacount']
  
  workload_ranges[idx]["replicas"] = best_replicas

print(durations, len(durations))
print(sum(durations) / len(durations))

with open(output_name,'w') as f:
  json.dump(workload_ranges, f,indent=4, sort_keys=True)