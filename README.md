Toy load balancer.

Send incoming HTTP requests using the Go built in http reverse proxy to the hosts defined in the `domains` slice.

A routine check runs in its own goroutine, establishing a tcp connection with the hosts to determine if the service is still up, providing info via stdout if a service is up or down.

```sh
2022/12/08 21:55:48 status: http://localhost:8000 - [up]
2022/12/08 21:55:48 status: http://localhost:8001 - [down]
2022/12/08 21:55:53 status: http://localhost:8000 - [up]
2022/12/08 21:55:53 status: http://localhost:8001 - [down]
```

For now, make theres a service host inside the `domains` slice inside `main.go`.

TODOs:
-dump stdout into log files
-Figure out a nice way of passing through an array of ports
-Make the balancer `service configurable`