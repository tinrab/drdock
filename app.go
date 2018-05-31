package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
)

type Container struct {
	ID   string
	Name string
}

type Action int

const (
	ActionContainerStarted = iota
	ActionContainerStopped = iota
)

type ActionListener func(a Action, data interface{})

type App struct {
	containers     map[string]Container
	client         *client.Client
	mutex          sync.Mutex
	actionListener ActionListener
}

func newApp(actionListener ActionListener) (*App, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	return &App{
		containers:     make(map[string]Container),
		client:         cli,
		mutex:          sync.Mutex{},
		actionListener: actionListener,
	}, nil
}

func (a *App) addContainer(id string) error {
	c, err := a.client.ContainerInspect(context.Background(), id)
	if err != nil {
		return err
	}
	a.mutex.Lock()
	defer a.mutex.Unlock()
	container := Container{
		ID:   c.ID,
		Name: c.Name,
	}
	a.containers[id] = container
	a.actionListener(ActionContainerStarted, container)
	return nil
}

func (a *App) removeContainer(id string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if _, ok := a.containers[id]; ok {
		delete(a.containers, id)
		a.actionListener(ActionContainerStopped, id)
	}
}

func (a *App) Run() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	go func() {
		a.fetchContainers()
		eventsCh, errorsCh := a.client.Events(context.Background(), types.EventsOptions{})

		for {
			select {
			case err := <-errorsCh:
				log.Println(err)
			case msg := <-eventsCh:
				if msg.Type == events.ContainerEventType {
					// Are there action constants somewhere?
					switch msg.Action {
					case "start":
						if err := a.addContainer(msg.ID); err != nil {
							log.Println(err)
						}
					case "die":
						fallthrough
					case "stop":
						a.removeContainer(msg.ID)
					}
				}
			}
		}
	}()
	<-shutdown
}

func (a *App) fetchContainers() {
	containers, err := a.client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Println(err)
		return
	}
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.containers = make(map[string]Container)
	for _, c := range containers {
		a.containers[c.ID] = Container{
			ID:   c.ID,
			Name: c.Names[0],
		}
	}
}
