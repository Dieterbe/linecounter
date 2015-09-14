Linecounter: every `x` ms, print out how many lines seen on stdin, possibly seggregated by a field in your lines.


## simple line counting

the easiest way to get real-time numbers on the command line.
for example, to count 404's by tail following a log file, grepping for what you want to count, and feeding the lines into linecounter,
which tells you the number of lines seen every second.

```
$ tail -f /var/log/apache/access.log | grep 404 | linecounter
102
156
234
98
```

## seggregating by an input field

often I want to know what the rates are for several different forms of input lines.
we can use the field flag to select a field.  That field will be used as a bucket.

For example,
I work on a website monitoring platform, and in my development stack I often create some fake accounts with fake monitors which send measurement data through NSQ.
I can tap into the stream with a tool called nsq_metrics_to_stdout and print out the data which looks like so:
```
litmus.fake_org_1_endpoint_4.dev1.ping.min [endpoint_id:45 monitor_id:135 collector:dev1]
litmus.fake_org_2_endpoint_1.dev1.http.dns [endpoint_id:46 monitor_id:138 collector:dev1]
```

I often need to get insights on this multi-dimensional stream and for example, check which monitor (http vs ping vs ...) creates more volume, etc.
I can simply use a tool like sed or tr to split the lines in fields and then select the field I need (here the fourth)
```
$ nsq_metrics_to_stdout | tr '.' ' ' | linecounter -f 4 --freq 5000
2015/09/14 15:24:30 INF    1 [metrics/stdout] (nsqd:4150) connecting to nsqd
2015/09/14 15:24:30 connected to nsqd
2015/09/14 15:24:30 INFO starting listener for http/debug on :6060
===============================
dns 7
http 132
ping 129
===============================
dns 56
http 86
ping 81
===============================
dns 14
http 85
ping 80
===============================
dns 7
http 118
ping 100
===============================
dns 7
http 88
ping 84
```

## usage syntax
```
./linecounter -h
Usage of ./linecounter:
  -f int
      seggregate by given whitespace separated field value. fields are numbered from 1 like cut and awk. (default -1 to disable) (default -1)
  -freq int
      report frequency in ms (default:1000) (default 1000 i.e. 1 second)
```

