package worker

import (
	logger "consumer/pkg/log"

	"golang.org/x/sync/errgroup"
)

func closeWorkers(group *errgroup.Group) {
	if err := group.Wait(); err != nil {
		logger.Error("One or more workers failed: %v", err)
	} else{
		logger.Info("All workers completed successfully")
	}
}