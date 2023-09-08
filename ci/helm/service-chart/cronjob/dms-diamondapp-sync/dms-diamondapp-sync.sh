#!/usr/bin/bash

export PATH=$PATH:/usr/bin:/bin:/usr/local/bin

HOUR=`date -d "now -2 hours" +"%Y%m%d%H"`
DATE=`date -d "now -2 hours" +"%Y%m%d"`

LOG=$DATE.log

echo "Syncing Diamond app tables for $DATE"  | tee -a $LOG
echo "Start time: `date +%r`" | tee -a $LOG | tee -a $LOG

#DIAMOND_APP_PATH=vd-prd-smarttvdata-cdw-virginia/data/meta_mms_snapshot
DIAMOND_APP_PATH=arn:aws:s3:us-east-1:833376745199:accesspoint/na-service-biz-smarttvdata-cdw-virginia/data/meta_mms_snapshot
UDW_PATH=udw-core-prod-us-east-1-data-sps/landing
DMS_TABLES="tsap_prod_m_parquet tsap_prod_ver_l_parquet"

sleep 120

for tables in $DMS_TABLES ; do
   aws s3 cp --recursive --request-payer requester s3://$DIAMOND_APP_PATH/$tables/partition_datehour=$HOUR/ s3://$UDW_PATH/udw_prod.sps_prod_$tables/partition_date=$DATE/ --acl bucket-owner-full-control
   if [ $? -eq 0 ]; then
      touch _SUCCESS
      aws s3 cp _SUCCESS s3://$UDW_PATH/udw_prod.sps_prod_$tables/partition_date=$DATE/_SUCCESS --acl bucket-owner-full-control
      echo "Success: Table $tables copied to UDW bucket" | tee -a $LOG
   else
      echo "Failed: Couldn't copy table $tables to UDW bucket" | tee -a $LOG
      exit 1
   fi ;
done

echo "Completed: Uploaded all the DMS tables to UDW bucket" | tee -a $LOG

echo "End time: `date +%r`" | tee -a $LOG

exit 0
