#!/bin/sh

echo "$STACKENGINE_API_TOKEN"
echo "$STACKENGINE_LEADER_IP"
echo "Searching for <$STACKENGINE_LOGSTASH_KEY> in the kv store"

i=0

while [ $i -le 1 ]
do
	echo "$STACKENGINE_LOGSTASH_IP"

	if [ "$STACKENGINE_LOGSTASH_IP" != "" ]; then
		echo "Setting logstash ip to $STACKENGINE_LOGSTASH_IP"
		i=$((i+2))
	else
		echo "No logstash IP yet. Checking in 5 seconds"
		sleep 5
		/bin/kvnator -env=STACKENGINE_LOGSTASH_IP
		source /tmp/kvnator.txt
	fi
done

/bin/logspout syslog://$STACKENGINE_LOGSTASH_IP