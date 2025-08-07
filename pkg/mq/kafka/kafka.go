package kafka

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/usual2970/retell/pkg/config"
)

type Client struct {
	brokers       []string
	saslMechanism sasl.Mechanism // 添加SASL认证机制
}

type Config struct {
	Brokers    []string
	Username   string // 添加用户名
	Password   string // 添加密码
	EnableSASL bool   // 是否启用SASL认证
}

type Message struct {
	Key       []byte
	Value     []byte
	Partition int // 指定消息发送到的分区
}

// TopicConfig 主题配置
type TopicConfig struct {
	NumPartitions     int           // 分区数量
	ReplicationFactor int           // 复制因子
	RetentionTime     time.Duration // 消息保留时间
}

// DefaultTopicConfig 返回主题的默认配置
func DefaultTopicConfig() TopicConfig {
	return TopicConfig{
		NumPartitions:     3,
		ReplicationFactor: 1,
		RetentionTime:     7 * 24 * time.Hour, // 默认保留7天
	}
}

func NewClient(cfg Config) *Client {
	client := &Client{
		brokers: cfg.Brokers,
	}

	// 如果启用了SASL认证，创建SASL Plain机制
	if cfg.EnableSASL && cfg.Username != "" && cfg.Password != "" {
		client.saslMechanism = plain.Mechanism{
			Username: cfg.Username,
			Password: cfg.Password,
		}
	}

	return client
}

// 获取配置了SASL的连接器
func (c *Client) getDialer() *kafka.Dialer {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	if c.saslMechanism != nil {
		dialer.SASLMechanism = c.saslMechanism
	}

	return dialer
}

// CreateTopic 创建一个主题及其分区
func (c *Client) CreateTopic(ctx context.Context, topic string, config TopicConfig) error {
	dialer := c.getDialer()
	conn, err := dialer.DialContext(ctx, "tcp", c.brokers[0])
	if err != nil {
		return fmt.Errorf("连接Kafka失败: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("获取控制器信息失败: %w", err)
	}

	controllerConn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s", controller.Host))
	if err != nil {
		return fmt.Errorf("连接控制器失败: %w", err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     config.NumPartitions,
			ReplicationFactor: config.ReplicationFactor,
		},
	}

	if err := controllerConn.CreateTopics(topicConfigs...); err != nil {
		return fmt.Errorf("创建主题失败: %w", err)
	}

	return nil
}

// GetPartitions 获取主题的分区信息
func (c *Client) GetPartitions(ctx context.Context, topic string) ([]kafka.Partition, error) {
	dialer := c.getDialer()
	conn, err := dialer.DialLeader(ctx, "tcp", c.brokers[0], topic, 0)
	if err != nil {
		return nil, fmt.Errorf("连接Kafka失败: %w", err)
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions(topic)
	if err != nil {
		return nil, fmt.Errorf("获取分区信息失败: %w", err)
	}

	return partitions, nil
}

// ProduceToPartition 发送消息到指定分区
func (c *Client) ProduceToPartition(ctx context.Context, topic string, partition int, messages []Message) error {
	if len(messages) == 0 {
		return nil
	}

	// 检查并创建 topic（如果不存在）
	if err := c.ensureTopicExists(ctx, topic, partition+1); err != nil {
		return fmt.Errorf("确保主题存在失败: %w", err)
	}

	w := &kafka.Writer{
		Addr:         kafka.TCP(c.brokers...),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		BatchTimeout: 10 * time.Millisecond,
		Transport: &kafka.Transport{
			SASL: c.saslMechanism,
		},
	}
	defer w.Close()

	kafkaMessages := make([]kafka.Message, len(messages))
	for i, msg := range messages {
		kafkaMessages[i] = kafka.Message{
			Key:       msg.Key,
			Value:     msg.Value,
			Partition: partition,
		}
	}

	if err := w.WriteMessages(ctx, kafkaMessages...); err != nil {
		return fmt.Errorf("写入消息到分区%d失败: %w", partition, err)
	}

	return nil
}

// ensureTopicExists 确保主题存在，如果不存在则创建
func (c *Client) ensureTopicExists(ctx context.Context, topic string, minPartitions int) error {
	// 先检查主题是否存在
	partitions, err := c.GetPartitions(ctx, topic)
	if err == nil {
		// 主题存在，检查分区数是否足够
		if len(partitions) >= minPartitions {
			return nil
		}
		// TODO: 如需要可以在这里扩展分区数
		return nil
	}

	// 主题不存在，创建它
	config := DefaultTopicConfig()
	if minPartitions > config.NumPartitions {
		config.NumPartitions = minPartitions
	}

	return c.CreateTopic(ctx, topic, config)
}

func (c *Client) Produce(ctx context.Context, topic string, messages []Message) error {
	if len(messages) == 0 {
		return nil
	}

	w := &kafka.Writer{
		Addr:         kafka.TCP(c.brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{}, // 默认使用最小字节负载均衡
		BatchTimeout: 10 * time.Millisecond,
		Transport: &kafka.Transport{
			SASL: c.saslMechanism,
		},
	}
	defer w.Close()

	kafkaMessages := make([]kafka.Message, len(messages))
	for i, msg := range messages {
		kafkaMessages[i] = kafka.Message{
			Key:       msg.Key,
			Value:     msg.Value,
			Partition: msg.Partition, // 尊重消息指定的分区
		}
	}

	if err := w.WriteMessages(ctx, kafkaMessages...); err != nil {
		return fmt.Errorf("写入消息失败: %w", err)
	}

	return nil
}

func (c *Client) ProduceOne(ctx context.Context, topic string, key, value []byte) error {
	return c.Produce(ctx, topic, []Message{{Key: key, Value: value}})
}

type ConsumeFunc func(msg Message) error

// ConsumePartition 从指定分区消费消息
func (c *Client) ConsumePartition(ctx context.Context, topic string, partition int, offset int64, handler ConsumeFunc) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   c.brokers,
		Topic:     topic,
		Partition: partition,
		MinBytes:  10e3,
		MaxBytes:  10e6,
		Dialer:    c.getDialer(), // 使用带SASL认证的连接器
	})

	// 设置起始偏移量
	if offset >= 0 {
		r.SetOffset(offset)
	}

	defer r.Close()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := r.ReadMessage(ctx)
			if err != nil {
				return fmt.Errorf("从分区%d读取消息失败: %w", partition, err)
			}

			if err := handler(Message{
				Key:       msg.Key,
				Value:     msg.Value,
				Partition: msg.Partition,
			}); err != nil {
				return fmt.Errorf("处理分区%d的消息失败: %w", partition, err)
			}
		}
	}
}

// ConsumePartitions 消费多个分区的消息
func (c *Client) ConsumePartitions(ctx context.Context, topic string, partitions []int, offset int64, handler ConsumeFunc) error {
	var wg sync.WaitGroup
	errs := make(chan error, len(partitions))

	// 为每个分区启动一个消费者
	for _, partition := range partitions {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			if err := c.ConsumePartition(ctx, topic, p, offset, handler); err != nil {
				errs <- err
			}
		}(partition)
	}

	// 等待所有消费者完成或出错
	go func() {
		wg.Wait()
		close(errs)
	}()

	// 收集第一个错误
	for err := range errs {
		return err
	}

	return nil
}

func (c *Client) Consume(ctx context.Context, topic, groupID string, handler ConsumeFunc) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  c.brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
		Dialer:   c.getDialer(), // 使用带SASL认证的连接器
	})
	defer r.Close()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := r.ReadMessage(ctx)
			if err != nil {
				return fmt.Errorf("读取消息失败: %w", err)
			}

			if err := handler(Message{
				Key:       msg.Key,
				Value:     msg.Value,
				Partition: msg.Partition, // 包含分区信息
			}); err != nil {
				return fmt.Errorf("处理消息失败: %w", err)
			}
		}
	}
}

func (c *Client) ConsumeWithTimeout(ctx context.Context, topic, groupID string, timeout time.Duration, handler ConsumeFunc) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return c.Consume(ctx, topic, groupID, handler)
}

// 获取主题的分区偏移量
func (c *Client) GetPartitionOffsets(ctx context.Context, topic string, partition int) (first, last int64, err error) {
	dialer := c.getDialer()
	conn, err := dialer.DialLeader(ctx, "tcp", c.brokers[0], topic, partition)
	if err != nil {
		return 0, 0, fmt.Errorf("连接分区%d失败: %w", partition, err)
	}
	defer conn.Close()

	first, err = conn.ReadFirstOffset()
	if err != nil {
		return 0, 0, fmt.Errorf("获取分区%d首个偏移量失败: %w", partition, err)
	}

	last, err = conn.ReadLastOffset()
	if err != nil {
		return 0, 0, fmt.Errorf("获取分区%d最后偏移量失败: %w", partition, err)
	}

	return first, last, nil
}

var (
	defaultClient     *Client
	defaultClientOnce sync.Once
)

// SetupDefaultClient 设置默认客户端（不使用认证）
func SetupDefaultClient(confs ...config.Kafka) error {
	if len(confs) == 0 {
		confs = []config.Kafka{config.GetConfig().Kafka}
	}

	if len(confs) == 0 {
		return fmt.Errorf("Kafka配置不能为空")
	}

	conf := confs[0]

	defaultClientOnce.Do(func() {
		defaultClient = NewClient(Config{
			Brokers:    conf.Brokers,
			Username:   conf.Username,
			Password:   conf.Password,
			EnableSASL: true,
		})
	})
	return nil
}

// CreateTopicWithPartitions 包级别方法：创建主题和分区
func CreateTopicWithPartitions(ctx context.Context, topic string, numPartitions, replicationFactor int) error {
	if defaultClient == nil {
		return fmt.Errorf("默认客户端未初始化，请先调用 SetupDefaultClient")
	}
	config := DefaultTopicConfig()
	config.NumPartitions = numPartitions
	config.ReplicationFactor = replicationFactor
	return defaultClient.CreateTopic(ctx, topic, config)
}

// GetTopicPartitions 包级别方法：获取主题的分区信息
func GetTopicPartitions(ctx context.Context, topic string) ([]kafka.Partition, error) {
	if defaultClient == nil {
		return nil, fmt.Errorf("默认客户端未初始化，请先调用 SetupDefaultClient")
	}
	return defaultClient.GetPartitions(ctx, topic)
}

// PushToPartition 包级别方法：发送消息到指定分区
func PushToPartition(ctx context.Context, topic string, partition int, key, value []byte) error {
	if defaultClient == nil {
		return fmt.Errorf("默认客户端未初始化，请先调用 SetupDefaultClient")
	}
	return defaultClient.ProduceToPartition(ctx, topic, partition, []Message{{Key: key, Value: value}})
}

func Push(ctx context.Context, topic string, key, value []byte) error {
	if defaultClient == nil {
		return fmt.Errorf("默认客户端未初始化，请先调用 SetupDefaultClient")
	}
	return defaultClient.ProduceOne(ctx, topic, key, value)
}

func PushBatch(ctx context.Context, topic string, messages []Message) error {
	if defaultClient == nil {
		return fmt.Errorf("默认客户端未初始化，请先调用 SetupDefaultClient")
	}
	return defaultClient.Produce(ctx, topic, messages)
}

func PushString(ctx context.Context, topic string, key, value string) error {
	return Push(ctx, topic, []byte(key), []byte(value))
}

// PushStringToPartition 包级别方法：发送字符串消息到指定分区
func PushStringToPartition(ctx context.Context, topic string, partition int, key, value string) error {
	return PushToPartition(ctx, topic, partition, []byte(key), []byte(value))
}

// ListenPartition 包级别方法：监听指定分区的消息
func ListenPartition(ctx context.Context, topic string, partition int, offset int64, handler ConsumeFunc) error {
	if defaultClient == nil {
		return fmt.Errorf("默认客户端未初始化，请先调用 SetupDefaultClient")
	}
	return defaultClient.ConsumePartition(ctx, topic, partition, offset, handler)
}

// ListenPartitions 包级别方法：监听多个分区的消息
func ListenPartitions(ctx context.Context, topic string, partitions []int, offset int64, handler ConsumeFunc) error {
	if defaultClient == nil {
		return fmt.Errorf("默认客户端未初始化，请先调用 SetupDefaultClient")
	}
	return defaultClient.ConsumePartitions(ctx, topic, partitions, offset, handler)
}

func Listen(ctx context.Context, topic, groupID string, handler ConsumeFunc) error {
	if defaultClient == nil {
		return fmt.Errorf("默认客户端未初始化，请先调用 SetupDefaultClient")
	}
	return defaultClient.Consume(ctx, topic, groupID, handler)
}

func ListenWithTimeout(ctx context.Context, topic, groupID string, timeout time.Duration, handler ConsumeFunc) error {
	if defaultClient == nil {
		return fmt.Errorf("默认客户端未初始化，请先调用 SetupDefaultClient")
	}
	return defaultClient.ConsumeWithTimeout(ctx, topic, groupID, timeout, handler)
}

// GetOffsets 包级别方法：获取分区的偏移量范围
func GetOffsets(ctx context.Context, topic string, partition int) (first, last int64, err error) {
	if defaultClient == nil {
		return 0, 0, fmt.Errorf("默认客户端未初始化，请先调用 SetupDefaultClient")
	}
	return defaultClient.GetPartitionOffsets(ctx, topic, partition)
}

func HandleStringMessage(handler func(key, value string) error) ConsumeFunc {
	return func(msg Message) error {
		return handler(string(msg.Key), string(msg.Value))
	}
}

// 扩展的处理函数，包含分区信息
func HandleDetailedMessage(handler func(key, value string, partition int) error) ConsumeFunc {
	return func(msg Message) error {
		return handler(string(msg.Key), string(msg.Value), msg.Partition)
	}
}
