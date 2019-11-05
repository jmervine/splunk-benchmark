## splunk-benchmark

Simple search load generation tool for Splunk.

#### Install
```
go get github.com/jmervine/splunk-benchmark
```

#### Usage
```
jmervine@laptop splunk-benchmark $ ./splunk-benchmark -h
NAME:
   splunk-benchmark - search load generation

USAGE:
   splunk-benchmark [args...]

DESCRIPTION:
   Simple search load generation tool for Splunk

AUTHOR:
   Joshua Mervine <joshua@mervine.net>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --output value, -o value       Output method; [text json jsonsummary] (default: "text")
   --splunk-host value, -H value  Splunk hostname; e.g. https://user:pass@splunk.example.com:8089 [$SPLUNK_HOST]
   --query value, -q value        Splunk search query (default: "search * | head 1")
   --runs value, -r value         Number of search runs to perform; -1 runs until ctrl-c and then collects results (default: 1)
   --threads value, -T value      Number of threads (default: 1)
   --delay value, -d value        Delay in seconds between runs (default: 0)
   --verbose                      Verbose output [$VERBOSE]
   --very-verbose                 Very verbose output [$VERY_VERBOSE]
   --help, -h                     show help
```

#### Example

```
jmervine@laptop splunk-benchmark $ ./splunk-benchmark \
    -s "https://user:pass@splunk.example.com:8089" \
    -r 10 -T 10 -d 0.5

 Runs  | Average    | Median     | Min        | Max
------------------------------------------------------
 10    | 0.249      | 0.269      | 0.141      | 0.304


jmervine@laptop splunk-benchmark $ ./splunk-benchmark \
    -s "https://user:pass@splunk.example.com:8089" \
    -r 10 -T 10 -d 0.5 -o jsonsummary
{
  "average": 0.2131,
  "median": 0.24,
  "min": 0.125,
  "max": 0.271,
  "errors": 0
}
 ```
