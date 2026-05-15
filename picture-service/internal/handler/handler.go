package handler

import (
	"context"
	"log"

	"picture-service/internal/config"
	"picture-service/internal/service"
	pb "picturepb"
)

type PictureHandler struct {
	pb.UnimplementedPictureServiceServer
	cfg     *config.Config
	service *service.Service
}

func NewPictureHandler(cfg *config.Config, srv *service.Service) *PictureHandler {
	return &PictureHandler{
		cfg:     cfg,
		service: srv,
	}
}

func (h *PictureHandler) GenerateIllustrations(ctx context.Context, req *pb.GenerateRequest) (*pb.GenerateResponse, error) {
	log.Printf("GenerateIllustrations called: text_len=%d, max_images=%d", len(req.Text), req.MaxImages)

	if req.Text == "" {
		return &pb.GenerateResponse{
			Success: false,
			Message: "text is required",
		}, nil
	}

	// Получаем user_id из контекста (устанавливается API Gateway)
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		userID = "guest"
	}

	taskID, err := h.service.Generation.StartGeneration(ctx, userID, req.Text)
	if err != nil {
		log.Printf("Failed to start generation: %v", err)
		return &pb.GenerateResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.GenerateResponse{
		Success: true,
		TaskId:  taskID,
		Message: "generation started",
	}, nil
}

func (h *PictureHandler) GetTaskStatus(ctx context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {
	log.Printf("GetTaskStatus called: task_id=%s", req.TaskId)

	status, imageURLs, err := h.service.Generation.GetTaskStatus(ctx, req.TaskId)
	if err != nil {
		return &pb.StatusResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	var results []*pb.ImageResult
	for i, url := range imageURLs {
		results = append(results, &pb.ImageResult{
			SceneIndex: int32(i + 1),
			Text:       "", // текст сцены не возвращается, только URL
			ImageUrl:   url,
		})
	}

	progress := 0
	if status == "completed" {
		progress = 100
	} else if status == "processing" {
		progress = 50
	}

	return &pb.StatusResponse{
		Success:  true,
		TaskId:   req.TaskId,
		Status:   status,
		Progress: int32(progress),
		Images:   results,
	}, nil
}

func (h *PictureHandler) RegenerateImage(ctx context.Context, req *pb.RegenerateRequest) (*pb.RegenerateResponse, error) {
	log.Printf("RegenerateImage called: task_id=%s, scene_index=%d", req.TaskId, req.SceneIndex)

	// TODO: реализовать перегенерацию конкретной сцены
	return &pb.RegenerateResponse{
		Success: false,
		Message: "not implemented yet",
	}, nil
}
