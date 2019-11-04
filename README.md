## splunk-benchmark

Simple search load generation tool for Splunk.

#### Install
```
go get github.com/jmervine/splunk-benchmark
```

#### Usage
```
jmervine@laptop splunk-benchmark $ ./splunk-benchmark -h
Usage of ./splunk-benchmark:
  -S string
        Splunk query (default "search * | head 1")
  -T int
        Number of threads, e.g. 10 runs * 2 threads will run 20 total searches (default 1)
  -V
        Print version and exit
  -d float
        Delay in seconds between runs (default 0.0)
  -n int
        Number of search runs to perform; 0 runs until SIGINT (default 1)
  -o string
        Output method: text or json (default "text")
  -s string
        Splunk hostname (https://uname:pword@host:port)
  -v
        Verbose output
  -vv
        Very verbose output
```

#### Example

```
jmervine@laptop splunk-benchmark $ ./splunk-benchmark \
    -s "https://user:pass@splunk.example.com:8089" \
    -n 10 -t 10 -d 0.5

 Thread     | Runs       | Average    | Median     | Min        | Max
--------------------------------------------------------------------------------
 0          | 10         | 0.1815     | 0.2490     | 0.1140     | 0.2490
 1          | 10         | 0.1275     | 0.1460     | 0.1090     | 0.1460
 2          | 10         | 0.1800     | 0.2450     | 0.1150     | 0.2450
 3          | 10         | 0.1855     | 0.2640     | 0.1070     | 0.2640
 4          | 10         | 0.1860     | 0.2610     | 0.1110     | 0.2610
 5          | 10         | 0.1810     | 0.2540     | 0.1080     | 0.2540
 6          | 10         | 0.1230     | 0.1340     | 0.1120     | 0.1340
 7          | 10         | 0.1920     | 0.2670     | 0.1170     | 0.2670
 8          | 10         | 0.1865     | 0.2610     | 0.1120     | 0.2610
 9          | 10         | 0.1235     | 0.1340     | 0.1130     | 0.1340
--------------------------------------------------------------------------------
     -  aggregate  -     | 0.1667     | 0.1340     | 0.1070     | 0.2670
--------------------------------------------------------------------------------
 Query: search * | head 1...


jmervine@laptop splunk-benchmark $ ./splunk-benchmark \
    -s "https://user:pass@splunk.example.com:8089" \
    -n 2 -t 2 -d 0.5 -o json
{
  "query": "search * | head 1",
  "average": 0.13225,
  "median": 0.14300000000000002,
  "min": 0.12,
  "max": 0.146,
  "thread": [
    {
      "average": 0.1315,
      "median": 0.1315,
      "min": 0.123,
      "max": 0.14,
      "run": [
        0.123,
        0.14
      ]
    },
    {
      "average": 0.133,
      "median": 0.133,
      "min": 0.12,
      "max": 0.146,
      "run": [
        0.12,
        0.146
      ]
    }
  ]
}
 ```
