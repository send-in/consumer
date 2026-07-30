package worker

import (
	mq "consumer/internal/queue"
	logger "consumer/pkg/log"
	
	"encoding/json"
	"fmt"
	"bytes"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func closeWorkers(group *errgroup.Group) {
	err := group.Wait()

	if err != nil {
		logger.Error("One or more workers failed: %v", err)
	} else{
		logger.Info("All workers completed successfully")
	}
}

func sendAck(message mq.Message) (bool, error) {
	body, err := json.Marshal(message)
	if err != nil {
		return false, err
	}
	
	request, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8000/api/v1/messages/%s/ack",
		bytes.NewBuffer(body),
	)

	if err != nil {
		return false, err
	}

	request.Header = http.Header{
		"Authorization": []string{
			fmt.Sprintf("Bearer: %s", ""),
		},
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return false, err
	}

	defer response.Body.Close()
	isOk := response.StatusCode == http.StatusOK
	return isOk, nil
}