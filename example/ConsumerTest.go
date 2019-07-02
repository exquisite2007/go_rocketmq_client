package main
import (
    "github.com/exquisite2007/go_rocketmq_client"
    "time"
    "log"
    )




func main() {
group := "testGroup4"
topic := "BenchmarkTest"
log.SetFlags(log.Ldate|log.Ltime|log.Lshortfile)
var timeSleep =1800 * time.Second
conf := &rocketmq.Config{
    Namesrv:   "localhost:9876",

    InstanceName: "DEFAULT",
    MaxMsgNums: 100,
    LogFileName:"rmq1.log",
}

consumer, err := rocketmq.NewDefaultConsumer(group, conf)
if err != nil {
    return 
}
consumer.Subscribe(topic, "*")
consumer.RegisterMessageListener(
    func(msgs []*rocketmq.MessageExt) error {
        // for i, msg := range msgs {
        //     fmt.Printf("msg %d %+v\n", i,msg)
        // }
        log.Print("Consume success!",len(msgs))
        return nil
    })
consumer.Start()
log.Print("  END")
time.Sleep(timeSleep)
}