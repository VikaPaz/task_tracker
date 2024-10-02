package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/VikaPaz/task_tracker/internal/models"
	pb "github.com/VikaPaz/task_tracker/proto/task"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskServise interface {
	Update(ctx context.Context, req models.Task) (models.Task, error)
	List(ctx context.Context, filter models.TaskFilter) ([]models.Task, error)
}

type TaskHandler struct {
	pb.UnimplementedTaskServiceServer
	router   *grpc.Server
	service  TaskServise
	validate *validator.Validate
	log      *zerolog.Logger
}

func NewTaskHandler(svc TaskServise, log *zerolog.Logger) *TaskHandler {
	router := grpc.NewServer()
	validate := validator.New()
	return &TaskHandler{
		router:   router,
		service:  svc,
		validate: validate,
		log:      log,
	}
}

func Run(server *TaskHandler, port string) {
	server.router.RegisterService(&pb.TaskService_ServiceDesc, server)
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	err = server.router.Serve(l)
	if err != nil {
		log.Fatal(err)
	}
}

func (h *TaskHandler) GetTasks(ctx context.Context, req *pb.GetTasksRequest) (*pb.GetTasksResponse, error) {
	filter := modelsTaskFilter(req.Filter)

	if err := h.validate.Struct(filter); err != nil {
		return &pb.GetTasksResponse{}, fmt.Errorf("failed to bind request: %w", err)
	}
	h.log.Debug().Msgf("validated a filter: %v", filter)

	tasks, err := h.service.List(ctx, filter)
	if err != nil {
		return &pb.GetTasksResponse{}, fmt.Errorf("failed to list tasks: %w", err)

	}

	pbTasks := pbTasks(tasks)

	return &pb.GetTasksResponse{
		Tasks: pbTasks,
	}, nil
}

// TODO: add auth
func genOwner() string {
	return uuid.New().String()
}

func (h *TaskHandler) UpdateTaskStatus(ctx context.Context, req *pb.UpdateTaskStatusRequest) (*pb.UpdateTaskStatusResponse, error) {
	task := models.Task{
		ID: req.TaskId,
	}

	task.ID = req.TaskId
	switch req.NewStatus {
	case pb.TaskStatus_IN_PROGRESS:
		task.Status = models.InProgress
	case pb.TaskStatus_DONE:
		task.Status = models.Done
	default:
		task.Status = ""
	}

	task.OwnerID = genOwner()
	_, err := uuid.Parse(task.ID)
	if err != nil {
		return &pb.UpdateTaskStatusResponse{}, fmt.Errorf("invalid UUID: %w", err)
	}
	h.log.Debug().Msgf("validated update task: %v", task)

	updatedTask, err := h.service.Update(ctx, task)
	if err != nil {
		return &pb.UpdateTaskStatusResponse{}, fmt.Errorf("failed to update task: %w", err)
	}

	return &pb.UpdateTaskStatusResponse{
		Success: true,
		Message: fmt.Sprint(updatedTask),
	}, nil
}

func modelsTaskFilter(filter *pb.TaskFilter) models.TaskFilter {
	res := models.TaskFilter{
		ID:          filter.Id,
		Title:       filter.Title,
		Description: filter.Description,
		OwnerID:     filter.OwnerId,
	}

	if filter.Sort != nil {
		res.TaskSort = models.TaskSort{
			Field: filter.Sort.Field,
			Order: filter.Sort.Order,
		}
	}

	if filter.Range != nil {
		res.Range = models.Range{
			Limit:  filter.Range.Limit,
			Offset: filter.Range.Offset,
		}
	}

	switch filter.Status {
	case pb.TaskStatus_IN_PROGRESS:
		res.Status = models.InProgress.String()
	case pb.TaskStatus_DONE:
		res.Status = models.Done.String()
	default:
		res.Status = ""
	}

	return res
}

func pbTasks(tasks []models.Task) []*pb.Task {
	req := make([]*pb.Task, len(tasks))
	for i, task := range tasks {
		req[i] = &pb.Task{
			Id:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			OwnerId:     task.OwnerID,
			Created:     timestamppb.New(task.Created),
			Updated:     timestamppb.New(task.Updated),
		}
		switch task.Status {
		case models.InProgress:
			req[i].Status = pb.TaskStatus_IN_PROGRESS
		case models.Done:
			req[i].Status = pb.TaskStatus_DONE
		default:
			req[i].Status = pb.TaskStatus_TASK_STATUS_UNSPECIFIED
		}
	}
	return req
}
