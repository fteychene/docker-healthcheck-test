# Docker 1.12 healthcheck tests

This project's goal is not to show a use case scenario of the Healthcheck new feature in docker.
It's more a project to do tests of varios cases of healthcheck and have a easy setup to be able to look on how it's done, available and could use it on a project.

## How-to use
Run `docker-compose up --build` to launch the tests.
Look at the monitor logs to see what append on the containers with healthcheck.

Modify the monitor code or the healthcheck if you want to check something specific

## Webapp :
A nodejs webapp with a curl based `HEALTHCHECK`.
You can change the response to the healthcheck by doing a `GET /switch` to siwtch de status code of the reponse from 200 to 500.


## Monitor :
A go program that listen to docker events for the current docker-compose project and launch a goroutine that pritn the health status of the running containers periodically.

## Note
A event type `health_status: {new_status}` is triggered when the health state of a container change.
The logs of the lasts healthcheck execution is available in the info (TODO check the purge or the timeout)
A check is a exec of the command done by docker himself (you can see the events `exec_create` and `exec_start` with the healthcheck configured command)
