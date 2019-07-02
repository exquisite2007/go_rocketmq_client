package main
import (
    "github.com/exquisite2007/go_rocketmq_client"
    // "fmt"
    "strconv"
    )
func main() {
	group := "test"
	topic := "BenchmarkTest1"
	conf := &rocketmq.Config{
    	Namesrv:   "localhost:9876",
    	InstanceName: "DEFAULT",
	}

	producer, err := rocketmq.NewDefaultProducer(group, conf)
	producer.Start()
	if err != nil {
    	return 
	}
	for i:=0;i<1000;i++ {
		msg := rocketmq.NewMessage(topic, []byte("wangli"+strconv.Itoa(i)))
		if _, err := producer.Send(msg); err != nil {
    		return 
		}	
	}
	
}