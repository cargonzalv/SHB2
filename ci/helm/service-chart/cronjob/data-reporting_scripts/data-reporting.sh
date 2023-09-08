#!/usr/bin/bash

export PATH=$PATH:/usr/bin:/bin:/usr/local/bin

log=`date +%F`.txt
if [ -f "$log" ]; then
    rm $log
fi

# The function expects 4 arguments
# $1 - number of retries
# $2 - time in seconds between retires
# $3 - header string reference to event
# $4 - S3 bucket location up to partitions
#     e.x: s3://udw-core-prod-us-east-1-data-sps/landing/udw_prod.sps_prod_adevent/use1-rprod/sps.adevent
gatherLocationData() {
    echo "$3 :: " >> $log
    local result=""
    max_retries=$1
    retrying=false
    tries=3
    day=1
    while [[ $day -le 5 ]]; do
        if [[ $retrying == false ]]; then
            tries=1
        fi
        d=`date -d "${day} days ago" +%d`
        month=`date -d "${day} days ago" +%m`
        year=`date -d "${day} days ago" +%Y`
        aws_report="$(aws s3 ls "${4}/year=$year/month=$month/day=$d/" --recursive || echo false)"
        if [[ $aws_report != "false" ]]; then
            retrying=false
            tmp="$(echo "$aws_report" |  awk 'BEGIN {total=0}{total+=$3}END{printf "%s %.2f %s\n", "'"$month/$d/$year"' ",  total/1024/1024 , "MiB"}')" # >> $log
            echo "$tmp" >> $log
            result="${result}\n${tmp}"
            let "day+=1"
        elif [[ $aws_report == "false" && $tries -lt $max_retries ]]; then
            echo "Command failes for $3 on $month/$d/$year. Attempt $tries/$max_retries" >> $log
            retrying=true
            tries=$((tries + 1))
            sleep $2
        else
            echo "The command failed for $3 for date $month/$d/$year. after $max_retries attempts."
            return 1
        fi
    done
    echo "--------------------------------------------------------------" >> $log
    echo $result
    return 0
}

# sps_prod_adevent location
sps_prod_adevent="$(gatherLocationData 3 2 sps_prod_adevent s3://udw-core-prod-us-east-1-data-sps/landing/udw_prod.sps_prod_adevent/use1-rprod/sps.adevent)"

# sps_prod_inventoryevent location
sps_prod_inventoryevent="$(gatherLocationData 3 2 sps_prod_inventoryevent s3://udw-core-prod-us-east-1-data-sps/landing/udw_prod.sps_prod_inventoryevent/use1-rprod/sps.inventoryevent)"

# sps_prod_sessionevent location
sps_prod_sessionevent="$(gatherLocationData 3 2 sps_prod_sessionevent s3://udw-core-prod-us-east-1-data-sps/landing/udw_prod.sps_prod_sessionevent/use1-rprod/sps.sessionevent)"

data="$(<$log)"
creationTimestamp: null
IFS='' read -r -d '' data << EOF
{
	"blocks": [
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "*nSync S3 Data Check*"
			}
		},
		{
			"type": "divider"
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "*sps_prod_adevent*\n$sps_prod_adevent"
			}
		},
		{
			"type": "divider"
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "*sps_prod_inventoryevent*\n$sps_prod_inventoryevent"
			}
		},
		{
			"type": "divider"
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "*sps_prod_sessionevent*\n$sps_prod_sessionevent"
			}
		},
		{
			"type": "divider"
		}
	]
}
EOF
curl -X POST -H 'Content-type: text/plain' --data "$data" "$slack_webhookURL"

exit 0
