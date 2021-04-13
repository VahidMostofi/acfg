```
python extract_replica_infos.py train_23_test_24_1x workload_ranges_top_10_bin_25_no_replica_train_sorted.json /home/vahid/acfg-results/bookstore/OnlineShoppingStore/BNV2/ workload_ranges_top_10_bin_25_with_replica_train_sorted.json
```

```
python generate-autoconfig-run-file.py workload_ranges_top_10_bin_25_no_replica_train_sorted.json train_23_test_24_1x > autoconfig-script.sh
```