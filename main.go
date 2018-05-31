package main

import "log"

func main() {
	app, err := newApp(onAction)
	if err != nil {
		log.Fatal(err)
	}
	app.Run()
}

func onAction(action Action, data interface{}) {
	switch action {
	case ActionContainerStarted:
		log.Printf("Container started: %v\n", data)
	case ActionContainerStopped:
		log.Printf("Container stopped: %v\n", data)
	}
}
