package base

import (
	"strings"

	"github.com/IBM/sarama"
	"github.com/nholuongut/argo-events/pkg/apis/events/v1alpha1"
	sharedutil "github.com/nholuongut/argo-events/pkg/shared/util"
	"go.uber.org/zap"
)

type Kafka struct {
	Logger *zap.SugaredLogger
	config *v1alpha1.KafkaBus
}

func NewKafka(config *v1alpha1.KafkaBus, logger *zap.SugaredLogger) *Kafka {
	// set defaults
	if config.ConsumerGroup == nil {
		config.ConsumerGroup = &v1alpha1.KafkaConsumerGroup{}
	}

	return &Kafka{
		Logger: logger,
		config: config,
	}
}

func (k *Kafka) Brokers() []string {
	return strings.Split(k.config.URL, ",")
}

func (k *Kafka) Config() (*sarama.Config, error) {
	config := sarama.NewConfig()

	// consumer config
	config.Consumer.IsolationLevel = sarama.ReadCommitted
	config.Consumer.Offsets.AutoCommit.Enable = false

	switch k.config.ConsumerGroup.Oldest {
	case true:
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	case false:
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	}

	switch k.config.ConsumerGroup.RebalanceStrategy {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	default:
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	}

	// producer config
	config.Producer.Idempotent = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Net.MaxOpenRequests = 1

	// common config
	if k.config.Version != "" {
		version, err := sarama.ParseKafkaVersion(k.config.Version)
		if err != nil {
			return nil, err
		}

		config.Version = version
	}

	// sasl
	if k.config.SASL != nil {
		config.Net.SASL.Enable = true
		config.Net.SASL.Mechanism = sarama.SASLMechanism(k.config.SASL.GetMechanism())

		switch config.Net.SASL.Mechanism {
		case "SCRAM-SHA-512":
			config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
				return &sharedutil.XDGSCRAMClient{HashGeneratorFcn: sharedutil.SHA512New}
			}
		case "SCRAM-SHA-256":
			config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
				return &sharedutil.XDGSCRAMClient{HashGeneratorFcn: sharedutil.SHA256New}
			}
		}

		user, err := sharedutil.GetSecretFromVolume(k.config.SASL.UserSecret)
		if err != nil {
			return nil, err
		}
		config.Net.SASL.User = user

		password, err := sharedutil.GetSecretFromVolume(k.config.SASL.PasswordSecret)
		if err != nil {
			return nil, err
		}
		config.Net.SASL.Password = password
	}

	// tls
	if k.config.TLS != nil {
		tls, err := sharedutil.GetTLSConfig(k.config.TLS)
		if err != nil {
			return nil, err
		}

		config.Net.TLS.Config = tls
		config.Net.TLS.Enable = true
	}

	return config, nil
}
