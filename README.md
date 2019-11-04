## splunk-benchmark

Simple search load generation tool for Splunk.

#### Usage
```
jmervine@jmervine-ltm1 splunk-benchmark $ ./splunk-benchmark -h
Usage of ./splunk-benchmark:
  -S string
        Splunk query
  -T int
        Number of threads, e.g. 10 runs * 2 threads will run 20 total searches (default 1)
  -d float
        Delay in seconds between runs
  -n int
        Number of search runs to perform (default 10)
  -s string
        Splunk hostname (https://uname:pword@host:port)
  -v    Verbose output
  -version
        Print version and exit
  -vv
        Very verbose output


jmervine@jmervine-ltm1 splunk-benchmark $ ./splunk-benchmark \
    -s "https://user:pass@splunk.example.com:8089" \
    -S "search index=main | head 100000" -n 5 -d 0 -T 2
 Thread     | Runs       | Average    | Median     | Min        | Max
--------------------------------------------------------------------------------
 0          | 5          | 2.2734     | 2.3060     | 2.1490     | 2.2390
 1          | 5          | 2.2170     | 2.2625     | 2.0620     | 2.2230
--------------------------------------------------------------------------------
 Query: search index=main | head 100000...
 ```
