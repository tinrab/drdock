package main

import (
	"context"
	"log"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type Container struct {
	ID         string            `json:"id"`
	Names      []string          `json:"names"`
	Image      string            `json:"image"`
	ImageID    string            `json:"image_id"`
	Command    string            `json:"command"`
	SizeRw     int64             `json:"size_rw,omitempty"`
	SizeRootFs int64             `json:"size_root_fs,omitempty"`
	Labels     map[string]string `json:"labels"`
}

func newContainer(c types.Container) Container {
	container := Container{
		ID:         c.ID,
		Names:      c.Names,
		Image:      c.Image,
		Command:    c.Command,
		SizeRw:     c.SizeRw,
		SizeRootFs: c.SizeRootFs,
		Labels:     c.Labels,
	}
	return container
}

type Action int

const (
	ActionContainerStarted = iota + 1
	ActionContainerStopped
)

type ActionListener func(a Action, data interface{})

type App struct {
	Containers     map[string]Container
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
		Containers:     make(map[string]Container),
		client:         cli,
		mutex:          sync.Mutex{},
		actionListener: actionListener,
	}, nil
}

func (a *App) addContainer(ctx context.Context, id string) error {
	fc := filters.NewArgs()
	fc.Add("id", id)
	containers, err := a.client.ContainerList(ctx, types.ContainerListOptions{
		Size:    true,
		All:     true,
		Filters: fc,
	})
	if err != nil {
		log.Println(err)
		return err
	}

	a.mutex.Lock()
	defer a.mutex.Unlock()
	c := newContainer(containers[0])
	a.Containers[id] = c
	a.actionListener(ActionContainerStarted, c)

	return nil
}

func (a *App) removeContainer(id string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if _, ok := a.Containers[id]; ok {
		delete(a.Containers, id)
		a.actionListener(ActionContainerStopped, id)
	}
}

func (a *App) Run(ctx context.Context) {
	a.fetchContainers(ctx)

	cf := filters.NewArgs()
	cf.Add("type", events.ContainerEventType)
	eventsCh, _ := a.client.Events(ctx, types.EventsOptions{
		Filters: cf,
	})

	for msg := range eventsCh {
		switch msg.Action {
		case "start":
			if err := a.addContainer(ctx, msg.ID); err != nil {
				log.Println(err)
			}
		case "die":
			fallthrough
		case "stop":
			a.removeContainer(msg.ID)
		}
	}
}

func (a *App) fetchContainers(ctx context.Context) {
	containers, err := a.client.ContainerList(ctx, types.ContainerListOptions{
		Size: true,
		All:  true,
	})
	if err != nil {
		log.Println(err)
		return
	}
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.Containers = make(map[string]Container)
	for _, c := range containers {
		a.Containers[c.ID] = newContainer(c)
	}
}
