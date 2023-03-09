package user

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/vnnyx/golang-dot-api/model/message"
)

func (controller *UserControllerImpl) HandleMessage() {
	err := controller.Consumer.SubscribeTopics([]string{message.USER_OTP_TOPIC}, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	msg_count := 0
	run := true
	MIN_COMMIT_COUNT := 10
	for run {
		ev := controller.Consumer.Poll(100)
		switch e := ev.(type) {
		case kafka.AssignedPartitions:
			fmt.Fprintf(os.Stderr, "%% %v\n", e)
			controller.Consumer.Assign(e.Partitions)
		case kafka.RevokedPartitions:
			fmt.Fprintf(os.Stderr, "%% %v\n", e)
			controller.Consumer.Unassign()
		case *kafka.Message:
			msg_count += 1
			if msg_count%MIN_COMMIT_COUNT == 0 {
				controller.Consumer.Commit()
			}
			payload := new(message.Message)
			err := json.Unmarshal(e.Value, &payload)
			fmt.Printf("%% Message on %s:\n{UserID: %s}\n",
				e.TopicPartition, payload.User.UserID)
			if err != nil {
				fmt.Println(err.Error())
			}
			err = controller.UserService.SendOTP(context.Background(), payload.User, payload.OTP)
			if err != nil {
				fmt.Println(err.Error())
			}
		case kafka.PartitionEOF:
			fmt.Printf("%% Reached %v\n", e)
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			run = false
		default:

		}
	}
}
