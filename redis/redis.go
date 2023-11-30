package redis

import (
	"context"

	"github.com/redis/go-redis/v9" // Import the Redis client package
	log "github.com/troneras/gorews/logger"
	"github.com/troneras/gorews/redis/config"
)

var conf *config.Config
var client *redis.Client

type RedisBase struct {
	redisChannel string
	MessageChan  chan interface{}
	cancel       context.CancelFunc
}

type RedisReader struct {
	RedisBase
}

type RedisWriter struct {
	RedisBase
}

func Configure() {
	conf = config.Configure()

	client = redis.NewClient(&redis.Options{
		Addr:     conf.RedisAddr, // Redis server address
		Password: "",             // No password set
		DB:       0,              // Use default DB
	})
	// Example of a Redis Ping to test connectivity
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("[REDIS] Error connecting to redis", log.Fields{"addr": conf.RedisAddr, "error": err})
	}
	log.Info("[REDIS] Client connected", log.Fields{"addr": conf.RedisAddr})
}

func NewRedisReader(redisChannel string) *RedisReader {
	ctx, cancel := context.WithCancel(context.Background())
	rr := &RedisReader{
		RedisBase: RedisBase{
			redisChannel: redisChannel,
			MessageChan:  make(chan interface{}),
			cancel:       cancel,
		},
	}
	log.Debug("[REDIS] Subscribing to redis channel ", log.Fields{"channel": redisChannel})

	go func(ctx context.Context) {
		// subscribe to redis channel
		pubsub := client.Subscribe(ctx, redisChannel)
		defer pubsub.Close()
		// Wait for confirmation that subscription is created before publishing anything.
		_, err := pubsub.Receive(ctx)
		if err != nil {
			log.Error("[REDIS] Error subscribing to redis channel", log.Fields{"channel": redisChannel, "error": err})
		}
		log.Debug("[REDIS] Subscribed to redis channel", log.Fields{"channel": redisChannel})
		// Go channel which receives messages.
		ch := pubsub.Channel()

		// Consume messages.
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					return
				}
				log.Debug("[REDIS] Received message from redis channel ", log.Fields{"channel": redisChannel, "message": msg.Payload})
				rr.RedisBase.MessageChan <- msg.Payload
			case <-ctx.Done():
				log.Debug("[REDIS] Closing redis reader for channel ", log.Fields{"channel": redisChannel})
				return
			}
		}
	}(ctx)

	return rr
}

func NewRedisWriter(redisChannel string) *RedisWriter {
	ctx, cancel := context.WithCancel(context.Background())
	rw := &RedisWriter{
		RedisBase: RedisBase{
			redisChannel: redisChannel,
			MessageChan:  make(chan interface{}),
			cancel:       cancel,
		},
	}

	go func(ctx context.Context) {
		for msg := range rw.MessageChan {
			log.Debug("[REDIS] Writing message to redis channel", log.Fields{"channel": redisChannel, "message": msg})
			client.Publish(ctx, redisChannel, msg)
		}
	}(ctx)

	return rw
}

func (rb *RedisBase) Close() {
	rb.cancel()

	close(rb.MessageChan)
}
