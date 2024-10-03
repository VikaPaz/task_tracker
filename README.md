# About

Task Tracker is a microservice for task management (to-do list? sic!) developed using the following technologies, libraries and frameworks:

1. Gin framework for processing HTTP requests.
2. gRPC for interaction with other microservices.
3. PostgreSQL as a database.
4. Bun ORM for interacting with the database.
5. Zerolog for event logging.

# Usage

## make commands
| Command         | Description                                                              |
|-----------------|--------------------------------------------------------------------------|
| `proto`         | Generate Go code from proto files                                        |
| `run`           | Run server on local machine                                              |
| `docker_build`  | Build and run app + DB in Docker Compose                                 |
| `run_postgres`  | Run PostgreSQL in Docker                                                 |
| `docker_stop`   | Stop Docker Compose with app + DB                                        |
| `build_client`  | Build client                                                             |
| `client_done`   | Set status to DONE for a task with a given ID (use `id=*needed id*`)     |
| `client_list`   | Get list of tasks using existing client                                  |
| `help`          | Display this help screen                                                 |
| `swag`          | Generate swagger                                                         |


## Swagger
http://{host}:{port}/swagger/index.html

Default URL: http://localhost:8900/swagger/index.html

All files are stored in [docs](docs)

## Proto
All files are stored in [proto/task](proto/task)

```protobuf
service TaskService {
    rpc GetTasks (GetTasksRequest) returns (GetTasksResponse);

    rpc UpdateTaskStatus (UpdateTaskStatusRequest) returns (UpdateTaskStatusResponse);
}
```

# TODO 

- [ ] swagger for http router 
- [ ] unit tests
- [ ] add auth
- [ ] add middleware to parse token with owner_id
- [ ] add interceptor to validate that request were made from authorizad client
- [ ] integration tests