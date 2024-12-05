package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	pflag "github.com/spf13/pflag"
)

// Package main implements a simple producer to send message.
func main() {
	pflag.String("role", "sender", "角色，可选值为：sender, receiver")
	pflag.String("namesrv", "192.168.255.160:9876", "rocketmq namesrv地址")
	pflag.String("topic", "test1", "rocketmq topic名称")
	pflag.Parse()

	role := pflag.Lookup("role").Value.String()
	fmt.Printf("role=%s", role)
	namesrv := pflag.Lookup("namesrv").Value.String()
	topic := pflag.Lookup("topic").Value.String()

	switch role {
	case "receiver":
		//os.Setenv("mq.consoleAppender.enabled", "true")
		c, err := rocketmq.NewPushConsumer(
			// 指定 Group 可以实现消费者负载均衡进行消费，并且保证他们的Topic+Tag要一样。
			// 如果同一个 GroupID 下的不同消费者实例，订阅了不同的 Topic+Tag 将导致在对Topic 的消费队列进行负载均衡的时候产生不正确的结果，最终导致消息丢失。(官方源码设计)
			consumer.WithGroupName("testGroup"),
			consumer.WithNameServer([]string{namesrv}))
		if err != nil {
			panic(err)
		}
		err = c.Subscribe(topic, consumer.MessageSelector{}, func(
			ctx context.Context,
			msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, msg := range msgs {
				fmt.Printf("<- %s \n", msg.Body)
			}
			// 消费成功，进行ack确认
			return consumer.ConsumeSuccess, nil
		})
		if err != nil {
			panic(err)
		}
		err = c.Start()
		if err != nil {
			panic(err)
		}
		defer func() {
			err = c.Shutdown()
			if err != nil {

				fmt.Printf("shutdown Consumer error: %s", err.Error())
			}
		}()
		<-(chan interface{})(nil)
	case "sender":
		p, _ := rocketmq.NewProducer(
			producer.WithNsResolver(primitive.NewPassthroughResolver([]string{namesrv})),
			producer.WithRetry(2),
		)
		err := p.Start()
		if err != nil {
			fmt.Printf("start producer error: %s", err.Error())
			os.Exit(1)
		}
		defer func() {
			err = p.Shutdown()
			if err != nil {
				fmt.Printf("shutdown producer error: %s", err.Error())
			}
		}()

		for i := 0; i < 10000; i++ {
			msg := &primitive.Message{
				Topic: topic,
				Body:  []byte("Hello RocketMQ Go Client! " + strconv.Itoa(i)),
			}
			res, err := p.SendSync(context.Background(), msg)
			if err != nil {
				fmt.Printf("--> faile: %s\n", err)
			} else {
				fmt.Printf("--> ok: result=%s\n", res.String())
			}
			time.Sleep(time.Millisecond * 50)
		}
	default:
		fmt.Printf("role=%s is not support.", role)
	}
}
