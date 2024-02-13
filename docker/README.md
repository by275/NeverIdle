# Dockerized [NeverIdle](https://github.com/layou233/NeverIdle)

## [OCI Idle Compute Instances](https://docs.oracle.com/en-us/iaas/Content/FreeTier/freetier_topic-Always_Free_Resources.htm)

### Reclamation of Idle Compute Instances

Idle Always Free compute instances may be reclaimed by Oracle. Oracle will deem virtual machine and bare metal compute instances as idle if, during a 7-day period, the following are true:

- CPU utilization for the 95th percentile is less than 20%
- Network utilization is less than 20%
- Memory utilization is less than 20% (applies to [A1 shapes](https://docs.oracle.com/en-us/iaas/Content/FreeTier/freetier_topic-Always_Free_Resources.htm#Details_of_the_Always_Free_Compute_instance__a1_flex) only)

## Usage

```asciidoc
NeverIdle 0.2.3 - Getting worse from here.
Platform: linux , amd64 , go1.21.1
GitHub: https://github.com/layou233/NeverIdle
[PRIORITY] Use the worst priority by default.
  -c duration
        Interval for CPU waste
  -cp float
        Percent of CPU waste
  -m int
        GiB of memory waste
  -n duration
        Interval for network speed test
  -p int
        Set process priority value (default 666)
  -t int
        Set concurrent connections for network speed test (default 10)
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
