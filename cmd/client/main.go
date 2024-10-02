package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/VikaPaz/task_tracker/proto/task"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	cmdDescription = `What command to use: list or done
list - returns list of all tasks
done - mark task with given id as done. need to include -i param with task id. Default: list
`
	cmdList = "list"
	cmdDone = "done"
)

func main() {

	cmd := flag.String("c", cmdList, cmdDescription)
	id := flag.String("i", "",
		"id of a task. Cannot be empty if command is 'done'")

	flag.Parse()

	if *cmd == "" {
		*cmd = cmdList
	}

	if *cmd == cmdDone && *id == "" {
		log.Fatal("'i' cannot be empty")
	}

	conn, err := grpc.NewClient(":9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)

	switch *cmd {
	case cmdDone:
		updateTaskStatus(client, *id, pb.TaskStatus_DONE)
	case cmdList:
		getTasks(client)
	}
}

func getTasks(client pb.TaskServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// status := "in_progress"

	// var pbstatus pb.TaskStatus

	// switch status {
	// case "in_progress":
	// 	pbstatus = pb.TaskStatus_IN_PROGRESS
	// case "done":
	// 	pbstatus = pb.TaskStatus_DONE
	// default:
	// 	log.Fatal("invalid task status")
	// }

	filter := &pb.TaskFilter{
		Id: []string{"3c3b8244-1dc2-4dc4-8e0f-19fb840aa12b"},
		// Status: pb.TaskStatus(pbstatus),
	}

	req := &pb.GetTasksRequest{
		Filter: filter,
	}

	resp, err := client.GetTasks(ctx, req)
	if err != nil {
		log.Fatalf("Error when calling GetTasks: %v", err)
	}

	log.Printf("Response from GetTasks: %v", resp.Tasks)
}

func updateTaskStatus(client pb.TaskServiceClient, taskID string, newStatus pb.TaskStatus) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	req := &pb.UpdateTaskStatusRequest{
		TaskId:    taskID,
		NewStatus: newStatus,
	}

	resp, err := client.UpdateTaskStatus(ctx, req)
	if err != nil {
		log.Fatalf("Error when calling UpdateTaskStatus: %v", err)
	}

	log.Printf("Response from UpdateTaskStatus: Success: %v, Message: %v", resp.Success, resp.Message)
}
