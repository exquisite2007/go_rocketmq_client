# Introduction

A RocketMQ client for golang supportting producer and consumer.
Rocketmq作为国内少有的apache开源项目，官方对go用户的支持却很弱。
[https://github.com/apache/rocketmq-client-go]官方的做法是通过go调用c++版的客户端，依赖编译对于不同的linux版本不友好，依赖的包有24M的大小，研究半天太繁琐放弃。
然后就是各种搜索，比较好的实现有didapinchegit/go_rocket_mq,但是只有消费者。
后面找到sevennt/rocketmq，遂基于这个版本开发。

本质上，rocketmq的客户端都是基于java客户端开发的，所有声明定义都可以在java实现里找到原型。

需要注意的是， 
- 客户端要维持与nameserver，broker的连接心跳
- 客户端需要从nameserver获取topic的路由信息
- 消费者需要做负载均衡

# Import package
import "exquisite2007/go_rocketmq_client"

# Getting started
### Getting message with consumer
```
group := "dev-VodHotClacSrcData"
topic := "canal_vod_collect__video_collected_count_live"
var timeSleep = 30 * time.Second
conf := &rocketmq.Config{
    Nameserver:   "192.168.7.101:9876;192.168.7.102:9876;192.168.7.103:9876",
    InstanceName: "DEFAULT",
}

consumer, err := rocketmq.NewDefaultConsumer(consumerGroup, consumerConf)
if err != nil {
    return err
}
consumer.Subscribe(consumerTopic, "*")
consumer.RegisterMessageListener(
    func(msgs []*MessageExt) error {
        for i, msg := range msgs {
            fmt.Println("msg", i, msg.Topic, msg.Flag, msg.Properties, string(msg.Body))
        }
        fmt.Println("Consume success!")
        return nil
    })
consumer.Start()

time.Sleep(timeSleep)
```

### Sending message with producer
- Synchronous sending
```
group := "dev-VodHotClacSrcData"
topic := "canal_vod_collect__video_collected_count_live"
conf := &rocketmq.Config{
    Nameserver:   "192.168.7.101:9876;192.168.7.102:9876;192.168.7.103:9876",
    ClientIp:     "192.168.1.23",
    InstanceName: "DEFAULT",
}

producer, err := rocketmq.NewDefaultProducer(group, conf)
producer.Start()
if err != nil {
    return errors.New("NewDefaultProducer err")
}
msg := NewMessage(topic, []byte("Hello RocketMQ!")
if sendResult, err := producer.Send(msg); err != nil {
    return errors.New("Sync sending fail!")
} else {
    fmt.Println("Sync sending success!, ", sendResult)
}
```

- Asynchronous sending
```
group := "dev-VodHotClacSrcData"
topic := "canal_vod_collect__video_collected_count_live"
conf := &rocketmq.Config{
    Nameserver:   "192.168.7.101:9876;192.168.7.102:9876;192.168.7.103:9876",
    InstanceName: "DEFAULT",
}
producer, err := rocketmq.NewDefaultProducer(group, conf)
producer.Start()
if err != nil {
    return err
}
msg := NewMessage(topic, []byte("Hello RocketMQ!")
sendCallback := func() error {
    fmt.Println("I am callback")
    return nil
}
if err := producer.SendAsync(msg, sendCallback); err != nil {
    return err
} else {
    fmt.Println("Async sending success!")
}
```

- Oneway sending
```
group := "dev-VodHotClacSrcData"
topic := "canal_vod_collect__video_collected_count_live"
conf := &rocketmq.Config{
    Nameserver:   "192.168.7.101:9876;192.168.7.102:9876;192.168.7.103:9876",
    InstanceName: "DEFAULT",
}
producer, err := rocketmq.NewDefaultProducer(group, conf)
producer.Start()
if err != nil {
    return err
}
msg := NewMessage(topic, []byte("Hello RocketMQ!")
if err := producer.SendOneway(msg); err != nil {
    return err
} else {
    fmt.Println("Oneway sending success!")
}
```