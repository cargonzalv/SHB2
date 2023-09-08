# 2. Approval after ARB meeting

Date: 2022-09-21

## Status

Accepted

## Context
Initial discussion on header-bidder design https://docs.google.com/document/d/1LEhCs5Dud61DjvdOUMQH0zI3y45zUYGS6vyozHf0Tik/edit 

## Decision

1. Move the creation of inventory event from KSQL DB to sps-header-bidder
2. Add strategy for handling Kafka downtime

## Consequences
1. System will become more resilient to Kafka failure/downtime if we adopt the secondary Kafka cluster approach and plan for saving events to S3 during downtime
2. Kafka cost will reduce if we move inventory topic creation in sps-header-bidder
