package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/davecgh/go-spew/spew"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	eventtypes "github.com/docker/engine-api/types/events"
	"github.com/docker/engine-api/types/filters"
	"golang.org/x/net/context"
)

type eventProcessor func(event eventtypes.Message) error

// DecodeEvents : to decode the events
func decodeEvents(input io.Reader, ep eventProcessor) error {
	dec := json.NewDecoder(input)
	for {
		var event eventtypes.Message
		err := dec.Decode(&event)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		if procErr := ep(event); procErr != nil {
			return procErr
		}
	}
	return nil
}

func inspect(cli *client.Client, containerID string) (types.ContainerJSON, error) {
	r, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return types.ContainerJSON{}, err
	}
	return r, nil
}

var healthstatusRegex = regexp.MustCompile("health_status: .+")

func main() {
	fmt.Printf("Starting events monitor for %s \n", os.Getenv("COMPOSE_PROJECT_NAME"))

	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
	if err != nil {
		panic(err)
	}

	filters := filters.NewArgs()
	filters.Add("label", fmt.Sprintf("com.docker.compose.project=%s", os.Getenv("COMPOSE_PROJECT_NAME")))

	body, err := cli.Events(context.Background(), types.EventsOptions{
		Filters: filters,
	})
	if err != nil {
		panic(err)
	}
	defer body.Close()

	err = decodeEvents(body, func(event eventtypes.Message) error {
		fmt.Printf("Docker event : \nType : %s Action : %s Id : %s\n", event.Type, event.Action, event.Actor.ID)
		if event.Action == "start" || healthstatusRegex.MatchString(event.Action) {
			info, err := inspect(cli, event.Actor.ID)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Health state information : \n%s health status : \n %s\n", info.Name, spew.Sdump(info.State.Health))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
