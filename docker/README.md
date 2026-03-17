# Dockerized [NeverIdle](https://github.com/by275/NeverIdle)

## [OCI Idle Compute Instances](https://docs.oracle.com/en-us/iaas/Content/FreeTier/freetier_topic-Always_Free_Resources.htm)

### Reclamation of Idle Compute Instances

Idle Always Free compute instances may be reclaimed by Oracle. Oracle will deem virtual machine and bare metal compute instances as idle if, during a 7-day period, the following are true:

- CPU utilization for the 95th percentile is less than 20%
- Network utilization is less than 20%
- Memory utilization is less than 20% (applies to [A1 shapes](https://docs.oracle.com/en-us/iaas/Content/FreeTier/freetier_topic-Always_Free_Resources.htm#Details_of_the_Always_Free_Compute_instance__a1_flex) only)

## Usage

```asciidoc
2026/03/17 12:00:00 INFO  : NeverIdle 0.2.3 - Getting worse from here
2026/03/17 12:00:00 INFO  : Platform: linux, amd64, go1.25.0
2026/03/17 12:00:00 INFO  : GitHub: https://github.com/by275/NeverIdle
2026/03/17 12:00:00 PRIOR : Use the worst priority by default.
  -c duration
        Interval for CPU waste
  -cp float
        Target CPU waste ratio between 0 and 1
  -m int
        GiB of memory waste
  -n duration
        Interval for network speed test
  -p int
        Set process priority value (default 666)
  -t int
        Set concurrent connections for network speed test (default 10)
2026/03/17 12:00:00 MEM   : Reserving 2 GiB in the background until shutdown
2026/03/17 12:00:00 CPU   : Maintaining background CPU occupancy with target ratio 0.15
2026/03/17 12:00:00 NET   : Starting network speed testing with interval 4h0m0s
2026/03/17 12:00:00 INFO  : NeverIdle is running, press Ctrl+C to stop
```

```yaml

  neveridle:
    restart: always
    security_opt:
      - no-new-privileges:true
    logging:
      driver: json-file
      options:
        max-size: "1024k"
        max-file: "5"
    container_name: neveridle
    image: ghcr.io/by275/neveridle
    command: "-cp 0.15 -m 2 -n 4h"
```
