package main

import (
	"fmt"
	"time"

	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
)

func lis(msg *message.PublishMessage) error {
	fmt.Println("X")
	fmt.Println(msg)
	return nil
}

func stop(msg, ack message.Message, err error) error {

	fmt.Println(ack)
	//time.Sleep(10000000)
	return nil
}

func main() {

	/// Instantiates a new Client
	c := &service.Client{}

	// Creates a new MQTT CONNECT message and sets the proper parameters
	msg := message.NewConnectMessage()
	msg.SetWillQos(1)
	msg.SetVersion(4)
	msg.SetCleanSession(true)
	msg.SetClientId([]byte("surgemq"))
	msg.SetKeepAlive(10)
	msg.SetWillTopic([]byte("will"))
	msg.SetWillMessage([]byte("send me home"))

	// Connects to the remote server at 127.0.0.1 port 1883
	c.Connect("tcp://127.0.0.1:1883", msg)

	// Creates a new SUBSCRIBE message to subscribe to topic "abc"
	submsg := message.NewSubscribeMessage()

	submsg.AddTopic([]byte("test"), 1)

	// Subscribes to the topic by sending the message. The first nil in the function
	// call is a OnCompleteFunc that should handle the SUBACK message from the server.
	// Nil means we are ignoring the SUBACK messages. The second nil should be a
	// OnPublishFunc that handles any messages send to the client because of this
	// subscription. Nil means we are ignoring any PUBLISH messages for this topic.
	//  c.Subscribe(submsg, stop, lis)

	// Creates a new PUBLISH message with the appropriate contents for publishing
	pubmsg := message.NewPublishMessage()
	pubmsg.SetTopic([]byte("test"))

	pubmsg.SetPayload([]byte("abc_msg"))
	pubmsg.SetQoS(1)

	//var wg sync.WaitGroup
	//wg.Add(1)

	// // Publishes to the server by sending the message

	c.Publish(pubmsg, nil)
	//pubmsg.SetPayload([]byte("abc_msg2"))
	////defer wg.Done()
	// c.Publish(pubmsg, nil)
	//pubmsg.SetPayload([]byte("abc_msg3"))
	////defer wg.Done()
	//	c.Publish(pubmsg, nil)
	//

	time.Sleep(time.Second * 20)
	// Disconnects from the server
	// wg.Wait()
	c.Disconnect()
}
