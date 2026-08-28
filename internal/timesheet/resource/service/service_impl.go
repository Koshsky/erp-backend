package service

import (
	"context"
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/timesheet/resource/repository"

	"github.com/Koshsky/erp-backend/internal/timesheet/resource/dto"
	tracingpkg "github.com/Koshsky/erp-backend/internal/tracing"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/pkg/date"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

type ResourceService struct {
	logger     *slog.Logger
	repository ResourceRepository
	mapper     *ResourceMapper
	validator  *ResourceValidator
	tracer     *tracingpkg.Tracer
}

// NewResourceService builds the ResourceService service.
func NewResourceService(logger *slog.Logger, tracer *tracingpkg.Tracer, r *repo.ResourceRepository) *ResourceService {
	return &ResourceService{
		logger:     logger,
		repository: r,
		mapper:     NewResourceMapper(),
		validator:  &ResourceValidator{},
		tracer:     tracer,
	}
}

func (s *ResourceService) ListResources(
	ctx context.Context,
	userID int64,
	role string,
	ownerID int64,
	limit, offset int,
) ([]dto.ResourceResponse, int64, error) {
	ctx, end := s.tracer.Start(ctx, "resource.ListResources")
	defer end(nil)

	rows, err := s.repository.ListResources(ctx, userID, role, ownerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repository.CountResources(ctx, userID, role, ownerID)
	if err != nil {
		return nil, 0, err
	}
	return s.mapper.ToDTOs(rows), total, nil
}

// CreateResource creates a resource. The middleware checked permissions; here
// only owner normalization: if not set, the creator becomes the owner
// (owner_id is required).
func (s *ResourceService) CreateResource(
	ctx context.Context,
	req dto.CreateResourceRequest,
	userID int64,
) (*dto.ResourceResponse, error) {
	ctx, end := s.tracer.Start(ctx, "resource.CreateResource")
	defer end(nil)

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

func (s *ResourceService) FindResource(ctx context.Context, id int64) (*dto.ResourceResponse, error) {
	ctx, end := s.tracer.Start(ctx, "resource.FindResource")
	defer end(nil)
	resource, err := s.repository.FindResource(ctx, id)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil, errors.ErrResourceNotFound
		}
		return nil, err
	}
	if resource == nil {
		return nil, errors.ErrResourceNotFound
	}
	return s.mapper.ToDTO(resource), nil
}

func (s *ResourceService) UpdateResource(
	ctx context.Context,
	id int64,
	req dto.UpdateResourceRequest,
) (*dto.ResourceResponse, error) {
	ctx, end := s.tracer.Start(ctx, "resource.UpdateResource")
	defer end(nil)

	resource, err := s.repository.FindResource(ctx, id)
	if err != nil || resource == nil {
		return nil, errors.ErrResourceNotFound
	}

	s.mapper.ApplyUpdateToDomain(resource, req)
	if err = s.validator.ValidateResource(resource); err != nil {
		return nil, err
	}

	updated, err := s.repository.UpdateResource(ctx, *resource)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTO(updated), nil
}

func (s *ResourceService) DeleteResource(ctx context.Context, id int64) error {
	ctx, end := s.tracer.Start(ctx, "resource.DeleteResource")
	defer end(nil)

	resource, err := s.repository.FindResource(ctx, id)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil // idempotent delete: already deleted — not an error
		}
		return err
	}
	if resource == nil {
		return nil // idempotent delete
	}

	return s.repository.DeleteResource(ctx, id)
}

// ListMembers returns the users attached to a resource.
func (s *ResourceService) ListMembers(ctx context.Context, resourceID int64) ([]dto.ResourceMemberResponse, error) {
	ctx, end := s.tracer.Start(ctx, "resource.ListMembers")
	defer end(nil)

	members, err := s.repository.ListMembersByResourceID(ctx, resourceID)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToMemberDTOs(members), nil
}

// AddMember attaches a user to a resource. The middleware checks resource
// management rights (admin or the resource owner); here we enforce the
// hierarchy rule: vp can attach only their direct subordinates ("a vassal of my vassal is not my vassal"),
// admin — anyone. Exception — self-subordination: the owner can add themselves
// to their own resource.
func (s *ResourceService) AddMember(
	ctx context.Context,
	resourceID int64,
	userID int64,
	actorID int64,
	role string,
) error {
	ctx, end := s.tracer.Start(ctx, "resource.AddMember")
	defer end(nil)

	if err := s.validator.ValidatePositiveID(resourceID, "resource_id"); err != nil {
		return err
	}
	if err := s.validator.ValidatePositiveID(userID, "user_id"); err != nil {
		return err
	}

	if role != userdomain.Admin && userID != actorID {
		managerID, err := s.repository.FindUserManager(ctx, userID)
		if err != nil {
			return err
		}
		if managerID == nil || *managerID != actorID {
			return errors.ErrForbidden
		}
	}

	return s.repository.AddMember(ctx, resourceID, userID)
}

// RemoveMember detaches a user from a resource.
func (s *ResourceService) RemoveMember(ctx context.Context, resourceID, userID int64) error {
	ctx, end := s.tracer.Start(ctx, "resource.RemoveMember")
	defer end(nil)

	if err := s.validator.ValidatePositiveID(resourceID, "resource_id"); err != nil {
		return err
	}
	if err := s.validator.ValidatePositiveID(userID, "user_id"); err != nil {
		return err
	}
	return s.repository.RemoveMember(ctx, resourceID, userID)
}

// ListAbsence returns absence ranges (is_available=false states) of the
// resource members overlapping the window.
func (s *ResourceService) ListAbsence(
	ctx context.Context,
	resourceID int64,
	start, end date.Date,
) ([]dto.ResourceAbsenceResponse, error) {
	ctx, finish := s.tracer.Start(ctx, "resource.ListAbsence")
	defer finish(nil)

	if err := s.validator.ValidatePositiveID(resourceID, "resource_id"); err != nil {
		return nil, err
	}
	if err := s.validator.ValidateDayRange(start, end); err != nil {
		return nil, err
	}
	absences, err := s.repository.ListAbsence(ctx, resourceID, start.Time(), end.Time())
	if err != nil {
		return nil, err
	}
	return s.mapper.ToAbsenceDTOs(absences), nil
}
