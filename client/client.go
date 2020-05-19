package client

import (
	"context"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
)

// ServiceBusClient represents an Azure
// service bus client instance.
type ServiceBusClient struct {
	connectionString string
	timeout          time.Duration
}

// Stats represents metrics of topics and queues.
type Stats struct {
	Queues *[]QueueStats
	Topics *[]TopicStats
}

// MessageCounts represents various message counts.
type MessageCounts struct {
	ActiveMessages             int32
	DeadLetterMessages         int32
	ScheduledMessages          int32
	TransferDeadLetterMessages int32
	TransferMessages           int32
}

// Sizes represents metrics unit of measure.
type Sizes struct {
	SizeInBytes  int64
	MaxSizeBytes int64
}

// QueueStats represents metrics of queues.
type QueueStats struct {
	Name string
	MessageCounts
	Sizes
}

// SubscriptionStats represents metrics of subscriptions.
type SubscriptionStats struct {
	Name string
	MessageCounts
}

// TopicStats represents metrics of topics.
type TopicStats struct {
	Name string
	MessageCounts
	Sizes

	Subscriptions *[]SubscriptionStats
}

// New accespts a connection string and a time duration
// to create and return a new service bus client instance.
func New(connectionString string, timeout time.Duration) *ServiceBusClient {
	return &ServiceBusClient{
		connectionString: connectionString,
		timeout:          timeout,
	}
}

// GetServiceBusStats returns servicebus metrics returned by azure.
func (c *ServiceBusClient) GetServiceBusStats() (*Stats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(c.connectionString))
	if err != nil {
		return nil, err
	}

	queues, err := getQueueStats(ctx, ns)
	if err != nil {
		return nil, err
	}

	topics, err := getTopicStats(ctx, ns)
	if err != nil {
		return nil, err
	}

	return &Stats{
		Queues: queues,
		Topics: topics,
	}, nil
}

func getQueueStats(ctx context.Context, ns *servicebus.Namespace) (*[]QueueStats, error) {
	var result []QueueStats

	queueManager := ns.NewQueueManager()

	queues, err := queueManager.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, queue := range queues {
		result = append(result, QueueStats{
			Name:          queue.Name,
			MessageCounts: countDetailsToMessageCounts(queue.CountDetails),
			Sizes: Sizes{
				SizeInBytes:  *queue.QueueDescription.SizeInBytes,
				MaxSizeBytes: int64(*queue.QueueDescription.MaxSizeInMegabytes) * 1024 * 1024,
			},
		})
	}

	return &result, nil
}

func getTopicStats(ctx context.Context, ns *servicebus.Namespace) (*[]TopicStats, error) {
	var result []TopicStats

	topicManager := ns.NewTopicManager()

	topics, err := topicManager.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, topic := range topics {
		subs, err := getSubscriptionStats(ctx, ns, topic.Name)
		if err != nil {
			return nil, err
		}

		result = append(result, TopicStats{
			Name:          topic.Name,
			MessageCounts: countDetailsToMessageCounts(topic.CountDetails),
			Sizes: Sizes{
				SizeInBytes:  *topic.TopicDescription.SizeInBytes,
				MaxSizeBytes: int64(*topic.TopicDescription.MaxSizeInMegabytes) * 1024 * 1024,
			},
			Subscriptions: subs,
		})
	}

	return &result, nil
}

func getSubscriptionStats(ctx context.Context, ns *servicebus.Namespace, topicName string) (*[]SubscriptionStats, error) {
	var result []SubscriptionStats

	subsManager, err := ns.NewSubscriptionManager(topicName)
	if err != nil {
		return nil, err
	}

	subs, err := subsManager.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, sub := range subs {
		result = append(result, SubscriptionStats{
			Name:          sub.Name,
			MessageCounts: countDetailsToMessageCounts(sub.CountDetails),
		})
	}
	return &result, nil
}

func countDetailsToMessageCounts(countDetails *servicebus.CountDetails) MessageCounts {
	return MessageCounts{
		ActiveMessages:             *countDetails.ActiveMessageCount,
		DeadLetterMessages:         *countDetails.DeadLetterMessageCount,
		ScheduledMessages:          *countDetails.ScheduledMessageCount,
		TransferDeadLetterMessages: *countDetails.TransferDeadLetterMessageCount,
		TransferMessages:           *countDetails.TransferMessageCount,
	}
}
