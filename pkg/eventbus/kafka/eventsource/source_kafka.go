package eventsource

import (
	"github.com/IBM/sarama"
	"go.uber.org/zap"

	eventbusv1alpha1 "github.com/nholuongut/argo-events/pkg/apis/events/v1alpha1"
	eventbuscommon "github.com/nholuongut/argo-events/pkg/eventbus/common"
	"github.com/nholuongut/argo-events/pkg/eventbus/kafka/base"
)

type KafkaSource struct {
	*base.Kafka
	topic string
}

func NewKafkaSource(config *eventbusv1alpha1.KafkaBus, logger *zap.SugaredLogger) *KafkaSource {
	return &KafkaSource{
		Kafka: base.NewKafka(config, logger),
		topic: config.Topic,
	}
}

func (s *KafkaSource) Initialize() error {
	return nil
}

func (s *KafkaSource) Connect(string) (eventbuscommon.EventSourceConnection, error) {
	config, err := s.Config()
	if err != nil {
		return nil, err
	}

	// eventsource specific config
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true

	client, err := sarama.NewClient(s.Brokers(), config)
	if err != nil {
		return nil, err
	}

	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}

	conn := &KafkaSourceConnection{
		KafkaConnection: base.NewKafkaConnection(s.Logger),
		Topic:           s.topic,
		Client:          client,
		Producer:        producer,
	}

	return conn, nil
}
