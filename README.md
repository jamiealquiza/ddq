# ddq

ddq is an overly simple Datadog metric query CLI tool. It's intended to view metrics as single-value rollups.

### Installation
```
 - go get github.com/jamiealquiza/ddq
 - go install github.com/jamiealquiza/ddq
```

Binary will be found at `$GOPATH/bin/ddq`

### Usage
```
Usage of ddq:
  -api-key string
    	Datadog API key [DDQ_API_KEY]
  -app-key string
    	Datadog app key [DDQ_APP_KEY]
  -by-tags string
    	Metric tags to reference data by (comma delimited) [DDQ_BY_TAGS] (default "host")
  -query string
    	Datadog metric query [DDQ_QUERY] (default "avg:system.load.1{*}")
  -span int
    	Query range in seconds (now - span) [DDQ_SPAN] (default 300)
```

The `-span` parameter specifies both the query range and the rollup in order to summarize metrics as single values. The `-by-tags` parameter will print "null" for metrics where the provided tag wasn't present.

### Example

```
% ddq -query "max:kafka.log.partition.size{role:some-kafka-cluster,topic:test_topic} by {topic,partition}" -span 7200  --by-tags topic,partition
submitting max:kafka.log.partition.size{role:some-kafka-cluster,topic:test_topic} by {topic,partition}.rollup(avg, 7200)

      test_topic,0: 1144499495262.09
      test_topic,1: 105342171367.44
     test_topic,10: 964324554703.97
     test_topic,11: 1384769522690.18
     test_topic,12: 763105537662.91
     test_topic,13: 1271944903583.73
     test_topic,14: 825473348879.32
     test_topic,15: 1862294121410.74
     test_topic,16: 1079007028250.26
     test_topic,17: 1001678621284.65
     test_topic,18: 557069021779.15
     test_topic,19: 827626127709.34
      test_topic,2: 1405778257636.16
     test_topic,20: 864819133483.76
     test_topic,21: 1219782711873.64
     test_topic,22: 1189880143784.67
     test_topic,23: 695106051185.78
     test_topic,24: 1362661517342.57
     test_topic,25: 1068553277675.80
     test_topic,26: 961372789928.12
     test_topic,27: 1181415376524.83
     test_topic,28: 1326462561078.70
     test_topic,29: 1167357320865.91
      test_topic,3: 401411176605.54
     test_topic,30: 1270905714635.49
     test_topic,31: 457939298893.51
      test_topic,4: 1227198899786.39
      test_topic,5: 296621871909.20
      test_topic,6: 1215867335811.28
      test_topic,7: 1305938202789.94
      test_topic,8: 1246741111274.12
      test_topic,9: 583089020666.00
```
