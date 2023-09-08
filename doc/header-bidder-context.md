## Context

Header Bidder Service - This service will expose public HTTP Post API for creating a data record of ad event based on open-RTB (or similar) JSON payload coming in each ad-request. The JSON data will ensure privacy compliance by redacting PII data based on privacy flags coming along with data.


Each ad and inventory event will be converted to Kafka record in Avro format and will be sent to Kafka cluster hosted by Hermes team. The service will produce telemetry logs and events to be stored in Prometheus and to be visualized in Grafana dashboard