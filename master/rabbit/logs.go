package rabbit

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/rabbit"
	"gitlab.com/systemz/aimpanel2/master/db"
	"gitlab.com/systemz/aimpanel2/master/model"
)

func ListenWrapperLogsQueue() {
	msgs, err := channel.Consume(
		"wrapper_logs",
		"",
		true,
		false,
		false,
		false,
		nil)
	lib.FailOnError(err, "Failed to register a consumer")

	go func() {
		for msg := range msgs {
			var msgBody rabbit.QueueMsg
			err = json.Unmarshal(msg.Body, &msgBody)
			if err != nil {
				logrus.Warn(err)
			}

			if msgBody.TaskId == rabbit.SERVER_LOG {
				var gsLog model.GameServerLog
				gsLog.GameServerID = msgBody.GameServerID

				if len(msgBody.Stdout) > 0 {
					gsLog.Log = msgBody.Stdout
					gsLog.Type = model.STDOUT
				}

				if len(msgBody.Stderr) > 0 {
					gsLog.Log = msgBody.Stderr
					gsLog.Type = model.STDERR
				}

				err = db.DB.Save(&gsLog).Error
				if err != nil {
					logrus.Warn(err)
				}

				logrus.Printf("Received a message: %s", gsLog)
			}
		}
	}()
}
