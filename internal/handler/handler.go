package handler

import (
	"log/slog"
)

type Service interface {
	UserService
	ProjectService
	ProcessService
	MilestoneService
	TaskService
	ResourceService
	AssignmentService
}

type Handler struct {
	logger *slog.Logger

	*UserHandler
	*ProjectHandler
	*ProcessHandler
	*MilestoneHandler
	*TaskHandler
	*ResourceHandler
	*AssignmentHandler
}

func New(logger *slog.Logger, service Service) *Handler {
	return &Handler{
		logger:            logger,
		UserHandler:       NewUserHandler(logger, service),
		ProjectHandler:    NewProjectHandler(logger, service),
		ProcessHandler:    NewProcessHandler(logger, service),
		MilestoneHandler:  NewMilestoneHandler(logger, service),
		TaskHandler:       NewTaskHandler(logger, service),
		ResourceHandler:   NewResourceHandler(logger, service),
		AssignmentHandler: NewAssignmentHandler(logger, service),
	}
}

func (h *Handler) Logger() *slog.Logger {
	return h.logger
}

// func (h *Handler) Service() *service.Service {
// 	return h.service
// }

// Response structures
type response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}
