package miniIM

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
	"strconv"
)

/*
User: Dpro
Date: 2025/9/11  周四
Time: 19:20
*/

type Service struct {
	// mq
	producer sarama.SyncProducer // 到了producer  就要考虑event
}

func (s *Service) Receive(ctx context.Context, sender int64, msg Message) error {
	for _, member := range s.findMembers() {
		if member == sender {
			continue
		}
		event := Event{
			Msg:      msg,
			Receiver: member,
		}
		val, _ := json.Marshal(event)
		// what need to do next
		// producer.sendMessage & producer.sendMessages
		_, _, err := s.producer.SendMessage(&sarama.ProducerMessage{
			Topic: eventName,
			Key:   sarama.StringEncoder(strconv.FormatInt(member, 10)), // 就是int 转string,只不过用stringEncoder再转一下，适配key
			Value: sarama.ByteEncoder(val),
		})
		if err != nil {
			log.Println("send message fail: ", err.Error())
			return err
		}

	}

	return nil
}

func (s *Service) findMembers() []int64 {
	// simple mock

	return []int64{1, 2, 3, 4, 5}

}
