package main

import (
	"git.kuschku.de/justjanne/imghost-frontend/configuration"
	"git.kuschku.de/justjanne/imghost-frontend/environment"
	"git.kuschku.de/justjanne/imghost-frontend/task"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

func main() {
	var config configuration.BackendConfiguration
	configFile, err := os.Open("config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.NewDecoder(configFile).Decode(&config)
	if err != nil {
		panic(err)
	}

	env, err := environment.NewBackendEnvironment(config)
	if err != nil {
		panic(err)
	}
	defer env.Destroy()

	mux := asynq.NewServeMux()
	mux.Handle(config.Conversion.TaskId, task.NewImageProcessor(env))
	if err := env.QueueServer.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
