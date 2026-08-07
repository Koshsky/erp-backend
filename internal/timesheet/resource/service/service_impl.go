package service

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/timesheet/resource/domain"
	"github.com/Koshsky/erp-backend/internal/timesheet/resource/dto"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
)

type ResourceService struct {
	logger     *slog.Logger
	repository ResourceRepository
	mapper     *ResourceMapper
	validator  *ResourceValidator
}

func NewResourceService(logger *slog.Logger, repository ResourceRepository) *ResourceService {
	return &ResourceService{
		logger:     logger,
		repository: repository,
		mapper:     NewResourceMapper(),
		validator:  &ResourceValidator{},
	}
}

// ownsResource проверяет право пользователя на ресурс:
// admin — любой; vp — только свои; остальные роли — без ограничений (как раньше).
func ownsResource(role string, userID int64, resource *domain.Resource) bool {
	if role == userdomain.Admin || role != userdomain.ProcessOwner {
		return true
	}
	return resource.OwnerID != nil && *resource.OwnerID == userID
}

// ListResources возвращает все ресурсы (admin и прочие роли) или только
// принадлежащие текущему vp.
func (s *ResourceService) ListResources(
	ctx context.Context,
	userID int64,
	role string,
) ([]dto.ResourceResponse, error) {
	var (
		resources []domain.Resource
		err       error
	)
	if role == userdomain.ProcessOwner {
		resources, err = s.repository.ListResourcesByOwnerID(ctx, userID)
	} else {
		resources, err = s.repository.ListResources(ctx)
	}
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTOs(resources), nil
}

func (s *ResourceService) CreateResource(
	ctx context.Context,
	req dto.CreateResourceRequest,
	userID int64,
) (*dto.ResourceResponse, error) {
	// owner_id обязателен: если не указан — владельцем становится создающий.
	// admin может указать любого владельца.
	if req.OwnerID == nil {
		req.OwnerID = &userID
	}

	resource := s.mapper.ToDomainFromCreate(req)
	if err := s.validator.ValidateResource(&resource); err != nil {
		return nil, err
	}

	created, err := s.repository.CreateResource(ctx, resource)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTO(created), nil
}

func (s *ResourceService) FindResource(
	ctx context.Context,
	id int64,
	userID int64,
	role string,
) (*dto.ResourceResponse, error) {
	resource, err := s.repository.FindResource(ctx, id)
	if err != nil {
		return nil, err
	}
	if resource == nil {
		return nil, ErrNotFound
	}
	if !ownsResource(role, userID, resource) {
		return nil, ErrForbidden
	}
	return s.mapper.ToDTO(resource), nil
}

func (s *ResourceService) UpdateResource(
	ctx context.Context,
	id int64,
	req dto.UpdateResourceRequest,
	userID int64,
	role string,
) (*dto.ResourceResponse, error) {
	resource, err := s.repository.FindResource(ctx, id)
	if err != nil {
		return nil, err
	}
	if resource == nil {
		return nil, ErrNotFound
	}
	if !ownsResource(role, userID, resource) {
		return nil, ErrForbidden
	}

	s.mapper.ApplyUpdateToDomain(resource, req)
	// vp не может передать ресурс другому владельцу.
	if role == userdomain.ProcessOwner {
		resource.OwnerID = &userID
	}
	if err = s.validator.ValidateResource(resource); err != nil {
		return nil, err
	}

	updated, err := s.repository.UpdateResource(ctx, *resource)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTO(updated), nil
}

func (s *ResourceService) DeleteResource(
	ctx context.Context,
	id int64,
	userID int64,
	role string,
) error {
	resource, err := s.repository.FindResource(ctx, id)
	if err != nil {
		return err
	}
	if resource == nil {
		return ErrNotFound
	}
	if !ownsResource(role, userID, resource) {
		return ErrForbidden
	}

	return s.repository.DeleteResource(ctx, id)
}
