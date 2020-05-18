package client

import (
	"time"
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
