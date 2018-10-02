# emon

Eventstore monitoring over HTTP

**NOTE: NOT ready for production!!!**

## Configuration

| Environment variable       | Description                          | Default               |
|----------------------------|--------------------------------------|-----------------------|
| EMON_HTTP_BIND_ADDRESS     | The address emon will listen on      | :8113                 |
| EMON_SLOW_CHECK_THRESHOLD  | The threshold for slow_check warning | 20ms                  |
| EMON_CLUSTER_HTTP_ENDPOINT | The eventstore HTTP address          | http://localhost:2113 |
| EMON_CLUSTER_SIZE          | The eventstore cluster size          | 3                     |


## Checks

| Type   | Name              | Description                                                                                     | Implemented        |
|--------|-------------------|-------------------------------------------------------------------------------------------------|--------------------|
| gossip | collect_gossip    | Collecting gossip from gossip endpoint                                                          | :white_check_mark: |
| gossip | alive_master      | Expects exacly 1 master                                                                         | :white_check_mark: |
| gossip | alive_slaves      | Expects `n` number of slaves where `n = ceil(cluster size / 2) - 1`.                            | :white_check_mark: |
| gossip | alive_nodes       | Expects `n` number of alive nodes where `n = ceil(cluster size / 2)`.                           | :white_check_mark: |
| stats  | collect_stats     | Collecting stats from stats endpoint                                                            | :white_check_mark: |
| stats  | sys_cpu           | Expects system cpu usage to be less then `n` percent. (`default: 90%`)                          | :white_check_mark: |
| stats  | sys_mem           | Expects `n` MB free system memory. (`default: 200MB`)                                           | :white_check_mark: |
| stats  | proc_cpu          | Expects process cpu usage to be less then `n` percent. (`default: 90%`)                         | :white_check_mark: |
| stats  | proc_mem          | Expects memory usage to be less then `n` MB. (`default: 1000MB`)                                | :white_check_mark: |
| ?      | storage_disk_full | Expects there to be at least `n` percent free disk space for the storage drive (`default: 25%`) | -                  |
| ?      | os_disk_full      | Expects there to be at least `n` percent free disk space for the os drive (`default: 25%`)      | -                  |
| gossip | clock_drift       | Expects diff between nodes to be `diff <= n`. (`default: 60s`)                                  | -                  |
| timer  | slow_checks       | Expects each check to complete within `n` milliseconds (`default: 100ms`)                       | :white_check_mark: |


## Development

```bash
docker-compose up -d
./run
```

## Releasing

Pre-requisites:
- [goreleaser](https://goreleaser.com/)
- A github token

```bash
goreleaser --rm-dist
```
