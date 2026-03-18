# noidle

*I love you, but do not stop my machine, could you?*

## Usage

Download the executable from Releases and choose the correct build for your platform, such as `amd64` or `arm64`.

Run it inside a `screen` session or another terminal multiplexer if you want it to keep running after you disconnect.

```asciidoc
2026/03/17 12:00:00 INFO  : noidle 0.2.3 - Getting worse from here
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
2026/03/17 12:00:00 INFO  : noidle is running, press Ctrl+C to stop
```

Command arguments:

```shell
./noidle -cp 0.15 -m 2 -n 4h
```

Flags:

`-c` enables periodic CPU load generation.  
For example, to burn CPU every 12 hours, 23 minutes, and 34 seconds, use `-c 12h23m34s`.

`-cp` enables adaptive CPU load control.  
The value is a ratio in the range `[0, 1]`. For example, use `-cp 0.2` to target roughly 20% additional CPU usage. Do not use it together with `-c`.

`-m` reserves memory in GiB.  
After startup, the specified amount of memory stays allocated until the process exits.

`-n` enables periodic network usage by running bandwidth tests.  
The interval format is the same as `-c`. noidle runs Ookla Speedtest on that schedule and prints the results.

`-t` sets the number of concurrent connections used for network testing.  
The default is `10`. Higher values consume more resources, and most setups do not need to change it.

`-p` sets the process priority. If omitted, noidle uses the lowest priority available on the current platform.  
On UNIX-like systems such as Linux, FreeBSD, and macOS, the valid range is `[-20,19]`, and larger numbers mean lower priority.  
For Windows, see [the official documentation](https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-setpriorityclass).  
In most cases, leaving it at the default is the safest option.
On Linux, the default mode lowers both CPU scheduling priority and disk I/O priority. On other UNIX-like systems, it lowers only CPU scheduling priority. On Windows, it changes the process priority class. This mainly affects CPU-heavy work and can indirectly affect network tests, but it does not reduce the memory reserved by `-m`.

All enabled features run once immediately at startup, so you can confirm that your settings are working.

## Docker

You can also run noidle with Docker Compose:

```yaml
services:
  noidle:
    image: ghcr.io/by275/neveridle
    container_name: noidle
    restart: always
    security_opt:
      - no-new-privileges:true
    logging:
      driver: json-file
      options:
        max-size: "1024k"
        max-file: "5"
    command: "-cp 0.15 -m 2 -n 4h"
```

The published image still uses the existing registry path `ghcr.io/by275/neveridle` for compatibility.

## Reference

### [OCI Idle Compute Instances](https://docs.oracle.com/en-us/iaas/Content/FreeTier/freetier_topic-Always_Free_Resources.htm)

### Reclamation of Idle Compute Instances

Idle Always Free compute instances may be reclaimed by Oracle. Oracle will deem virtual machine and bare metal compute instances as idle if, during a 7-day period, the following are true:

- CPU utilization for the 95th percentile is less than 20%
- Network utilization is less than 20%
- Memory utilization is less than 20% (applies to [A1 shapes](https://docs.oracle.com/en-us/iaas/Content/FreeTier/freetier_topic-Always_Free_Resources.htm#Details_of_the_Always_Free_Compute_instance__a1_flex) only)
