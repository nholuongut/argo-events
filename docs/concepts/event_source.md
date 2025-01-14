# Event Source

An `EventSource` defines the configurations required to consume events from external sources like AWS SNS, SQS, GCP PubSub, Webhooks, etc. It further
transforms the events into the [cloudevents](https://github.com/cloudevents/spec) and dispatches them over to the eventbus.

Available event-sources:

1. AMQP
1. AWS SNS
1. AWS SQS
1. Azure Events Hub
1. Azure Queue Storage
1. Bitbucket
1. Bitbucket Server
1. Calendar
1. Emitter
1. File Based Events
1. GCP PubSub
1. Generic EventSource
1. GitHub
1. GitLab
1. HDFS
1. K8s Resources
1. Kafka
1. Minio
1. NATS
1. NetApp StorageGrid
1. MQTT
1. NSQ
1. Pulsar
1. Redis
1. Slack
1. Stripe
1. Webhooks

## Specification

The complete specification is available [here](../APIs.md#argoproj.io/v1alpha1.EventSourceSpec).

## Examples

Examples are located under [examples/event-sources](https://github.com/nholuongut/argo-events/tree/main/examples/event-sources).
