
filebeat fields if logs go directly to elasticsearch
```
    @timestamp	            Apr 21, 2020 @ 17:04:06.662
	_id	                    CZwNnXEB2bWS8Caznhp9
	_index	                filebeat-7.6.2-2020.04.21-000001
	_score	                - 
	_type	                _doc
	input.type	            log
	log.file.path	        /logs/my.log
	log.offset	            34
	message	                hello
	suricata.eve.timestamp	Apr 21, 2020 @ 17:04:06.662
```


```
    @timestamp	Apr 21, 2020 @ 17:57:18.074
	_id	frE-nXEBoYoqSaN3UPTq
	_index	filebeat-7.6.2-2020.04.21-000001
	_score	 - 
	_type	_doc
	agent.ephemeral_id	30d48b19-9bb1-4b5b-b8b0-275eeb7821f1
	agent.hostname	9dfcf6f5dec1
	agent.id	aaa6b414-b356-4a6b-8d18-c467c7aa7359
	agent.type	filebeat
	agent.version	7.6.2
	ecs.version	1.4.0
	host.name	9dfcf6f5dec1
	input.type	log
	log.file.path	/logs/my.log
	log.offset	40
	message	привет
	suricata.eve.timestamp	Apr 21, 2020 @ 17:57:18.074

```
logstash
```
@timestamp	Apr 21, 2020 @ 18:06:51.298
	@version	1
	_id	SbhHnXEBVhgm1dpfF_03
	_index	filebeat-2020.04.21
	_score	 - 
	_type	_doc
	agent.ephemeral_id	c026ceaa-ee0b-4381-8179-31c0cf0dfd46
	agent.hostname	9b9cd372403e
	agent.id	35e942a3-08db-4259-8227-09f9e3d18fcc
	agent.type	filebeat
	agent.version	7.6.2
	ecs.version	1.4.0
	host.name	9b9cd372403e
	input.type	log
	log.file.path	/logs/my.log
	log.offset	53
	message	привет снова
	tags	beats_input_codec_plain_applied


```
