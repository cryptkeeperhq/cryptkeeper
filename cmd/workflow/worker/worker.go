package main

import (
	"log"

	"github.com/sirupsen/logrus"
	logrusadapter "logur.dev/adapter/logrus"
	"logur.dev/logur"

	"github.com/cryptkeeperhq/cryptkeeper/internal/workflow"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const TaskQueueName = "test"

func main() {
	lloger := logrus.New()
	lloger.Level = logrus.DebugLevel
	lloger.Formatter = &logrus.JSONFormatter{}
	logger := logur.LoggerToKV(logrusadapter.New(lloger))

	c, err := client.Dial(client.Options{
		Logger: logger,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer c.Close()

	w := worker.New(c, TaskQueueName, worker.Options{})

	// This worker hosts both Workflow and Activity functions.
	// Workflow
	t := workflow.Temporal{}
	w.RegisterWorkflow(t.RunWorkflow)

	// Activity functions.
	w.RegisterActivity(t.TestPlugin)
	w.RegisterActivity(t.ExecuteNode)

	// Start listening to the Task Queue.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
