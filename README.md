# About

Task Tracker is a microservice for task management (to-do list? sic!) developed using the following technologies, libraries and frameworks:

1. Gin framework for processing HTTP requests.
2. gRPC for interaction with other microservices.
3. PostgreSQL as a database.
4. Bun ORM for interacting with the database.
5. Zerolog for event logging.

# Usage

## make commands
| command | description |
|----------|----------|
|run|                            run server in local machine|
|build        |                  build and run app+DB in docker compose|
|run_postgres  |                 run postgres in docker|
|build_client  |                 build client|
|client_done  |                  set status DONE of a task with given id. Need to add id=*needed id*|
|client_list |                   get list of tasks using existion client|
|help      |                     Display this help screen|


## Swagger

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