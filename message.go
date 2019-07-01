package rocketmq

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	// "encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
)

const (
	CompressedFlag          = (0x1 << 0)
	MultiTagsFlag           = (0x1 << 1)
	TransactionNotType      = (0x0 << 2)
	TransactionPreparedType = (0x1 << 2)
	TransactionCommitType   = (0x2 << 2)
	TransactionRollbackType = (0x3 << 2)
)

const (
	NameValueSeparator = 1 + iota
	PropertySeparator
)

const (
	CharacterMaxLength = 255
)

type Message struct {
	Topic      string
	Flag       int32
	Properties map[string]string
	Body       []byte
}

func NewMessage(topic string, body []byte) *Message {
	return &Message{
		Topic:      topic,
		Body:       body,
		Properties: make(map[string]string),
	}
}

type MessageExt struct {
	Message
	QueueId                   int32  `json:"queueId"`
	StoreSize                 int32  `json:"storeSize"`
	QueueOffset               int64  `json:"queueOffset"`
	SysFlag                   int32  `json:"sysFlag"`
	BornTimestamp             int64  `json:"bornTimestamp"`
	BornHost      			  string `json:"bornHost"`
	StoreTimestamp            int64  `json:"storeTimestamp"`
	StoreHost                 string `json:"storeHost"`
	MsgId                     string `json:"msgId"`
	CommitLogOffset           int64  `json:"commitLogOffset"`
	BodyCRC                   int32  `json:"bodyCRC"`
	ReconsumeTimes            int32  `json:"reconsumeTimes"`
	PreparedTransactionOffset int64  `json:"preparedTransactionOffset"`
}

func decodeMessage(data []byte) []*MessageExt {
	buf := bytes.NewBuffer(data)
	var storeSize, magicCode, bodyCRC, queueId, flag, sysFlag, reconsumeTimes, bodyLength, bornPort, storePort int32
	var queueOffset, physicOffset, preparedTransactionOffset, bornTimeStamp, storeTimestamp int64
	var topicLen byte
	var topic, body, properties, bornHost, storeHost []byte
	var propertiesLength int16

	// var propertiesMap map[string]string
	msgs := make([]*MessageExt, 0, 32)

	
	for buf.Len() > 0 {
		msg := new(MessageExt)
		//1
		binary.Read(buf, binary.BigEndian, &storeSize)
		//2
		binary.Read(buf, binary.BigEndian, &magicCode)
		//3
		binary.Read(buf, binary.BigEndian, &bodyCRC)
		//4
		binary.Read(buf, binary.BigEndian, &queueId)
		//5
		binary.Read(buf, binary.BigEndian, &flag)
		//6
		binary.Read(buf, binary.BigEndian, &queueOffset)
		//7
		binary.Read(buf, binary.BigEndian, &physicOffset)
		//8
		binary.Read(buf, binary.BigEndian, &sysFlag)
		//9
		binary.Read(buf, binary.BigEndian, &bornTimeStamp)
		//10
		bornHost = make([]byte, 4)
		binary.Read(buf, binary.BigEndian, &bornHost)
		binary.Read(buf, binary.BigEndian, &bornPort)
		//11
		binary.Read(buf, binary.BigEndian, &storeTimestamp)
		//12
		storeHost = make([]byte, 4)
		binary.Read(buf, binary.BigEndian, &storeHost)
		binary.Read(buf, binary.BigEndian, &storePort)
		//13
		binary.Read(buf, binary.BigEndian, &reconsumeTimes)
		binary.Read(buf, binary.BigEndian, &preparedTransactionOffset)
		//14
		binary.Read(buf, binary.BigEndian, &bodyLength)
		if bodyLength > 0 {
			body = make([]byte, bodyLength)
			//15
			binary.Read(buf, binary.BigEndian, body)

			if (sysFlag & CompressedFlag) == CompressedFlag {
				b := bytes.NewReader(body)
				z, err := zlib.NewReader(b)
				if err != nil {
				Error.Println(err)
					return nil
				}
				defer z.Close()
				body, err = ioutil.ReadAll(z)
				if err != nil {
				Error.Println(err)
					return nil
				}

			}
		}
		//16
		binary.Read(buf, binary.BigEndian, &topicLen)
		topic = make([]byte, topicLen)
		binary.Read(buf, binary.BigEndian, &topic)
		//17
		binary.Read(buf, binary.BigEndian, &propertiesLength)
		if propertiesLength > 0 {
			properties = make([]byte, propertiesLength)
			binary.Read(buf, binary.BigEndian, &properties)
			msg.Properties = convertProperties(properties)
		}

		if magicCode != -626843481 {
		Warning.Printf("magic code is error %d", magicCode)
			return nil
		}
		msg.Topic = string(topic)
		msg.QueueId = queueId
		msg.SysFlag = sysFlag
		msg.QueueOffset = queueOffset
		msg.BodyCRC = bodyCRC
		msg.StoreSize = storeSize
		msg.BornHost = convertHostString(bornHost, bornPort)
		msg.BornTimestamp = bornTimeStamp
		msg.ReconsumeTimes = reconsumeTimes
		msg.Flag = flag
		msg.CommitLogOffset=physicOffset
		msg.MsgId = convertMessageId(storeHost, storePort, physicOffset)
		msg.StoreHost = convertHostString(storeHost, storePort)
		msg.StoreTimestamp = storeTimestamp
		msg.PreparedTransactionOffset = preparedTransactionOffset
		msg.Body = body
		// msg.Properties = propertiesMap

		msgs = append(msgs, msg)
	}

	return msgs
}

func convertMessageId(storeHost []byte, storePort int32, offset int64) string {
	buffMsgId := make([]byte, MSG_ID_LENGTH)
	input := bytes.NewBuffer(buffMsgId)
	input.Reset()
	input.Grow(MSG_ID_LENGTH)
	input.Write(storeHost)

	storePortBytes := int32ToBytes(storePort)
	input.Write(storePortBytes)

	offsetBytes := int64ToBytes(offset)
	input.Write(offsetBytes)

	return bytesToHexString(input.Bytes())
}

// 根据ip和port，解析host
func convertHostString(ipBytes []byte, port int32) string {
	ip := bytesToIPv4String(ipBytes)

	return  fmt.Sprintf("%s:%s", ip, strconv.FormatInt(int64(port), 10))
}
func convertProperties(buf []byte) map[string]string {

	tbuf := buf
	properties := make(map[string]string)
	for len(tbuf) > 0 {
		pi := bytes.IndexByte(tbuf, PROPERTY_SEPARATOR)
		if pi == -1 {
			break
		}

		propertie := tbuf[0:pi]

		ni := bytes.IndexByte(propertie, NAME_VALUE_SEPARATOR)
		if ni == -1 || ni > pi {
			break
		}

		key := string(propertie[0:ni])
		properties[key] = string(propertie[ni+1:])

		tbuf = tbuf[pi+1:]
	}

	return properties
}
func messageProperties2String(properties map[string]string) string {
	StringBuilder := bytes.NewBuffer([]byte{})
	if properties != nil && len(properties) != 0 {
		for k, v := range properties {
			binary.Write(StringBuilder, binary.BigEndian, k)                  // 4
			binary.Write(StringBuilder, binary.BigEndian, NameValueSeparator) // 4
			binary.Write(StringBuilder, binary.BigEndian, v)                  // 4
			binary.Write(StringBuilder, binary.BigEndian, PropertySeparator)  // 4
		}
	}
	return StringBuilder.String()
}

func (msg Message) checkMessage(producer *DefaultProducer) (err error) {
	if err = checkTopic(msg.Topic); err != nil {
		if len(msg.Body) == 0 {
			err = errors.New("ResponseCode:" + strconv.Itoa(MsgIllegal) + ", the message body is null")
		} else if len(msg.Body) > producer.maxMessageSize {
			err = errors.New("ResponseCode:" + strconv.Itoa(MsgIllegal) + ", the message body size over max value, MAX:" + strconv.Itoa(producer.maxMessageSize))
		}
	}
	return
}

func checkTopic(topic string) (err error) {
	if topic == "" {
		err = errors.New("the specified topic is blank")
	}
	if len(topic) > CharacterMaxLength {
		err = errors.New("the specified topic is longer than topic max length 255")
	}
	if topic == DefaultTopic {
		err = errors.New("the topic[" + topic + "] is conflict with default topic")
	}
	return
}
