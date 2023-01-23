package resource

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"
	"go.uber.org/zap"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Service interface {
	Get(ctx context.Context, org_id string, id string) (Resource, error)
	Query(ctx context.Context, org_id string, filter Filter) ([]Resource, error)
	Create(ctx context.Context, org_id string, input CreateResourceRequest) (Resource, error)
	Update(ctx context.Context, org_id string, id string, input UpdateResourceRequest) (Resource, error)
	Delete(ctx context.Context, org_id string, id string) (Resource, error)
}

type Resource struct {
	entity.Resource
}

type CreateResourceRequest struct {
	Key   string `json:"resource_key" db:"resource_key"`
	Name  string `json:"name" db:"name"`
	OrgID string `json:"-" db:"org_id"`
}

func (m CreateResourceRequest) Validate() error {

	return validation.ValidateStruct(&m,
		validation.Field(&m.Key, validation.Required),
	)
}

type UpdateResourceRequest struct {
	Name string `json:"name" db:"name"`
}

func (m UpdateResourceRequest) Validate() error {

	return validation.ValidateStruct(&m)
}

type service struct {
	repo   Repository
	logger *zap.Logger
}

func NewService(repo Repository, logger *zap.Logger) Service {

	return service{repo: repo, logger: logger}
}

// Get resource by id.
func (s service) Get(ctx context.Context, org_id string, id string) (Resource, error) {

	resource, err := s.repo.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Error("Error while getting the resource.",
			zap.String("organization_id", org_id),
			zap.String("resource_id", id))
		return Resource{}, &util.NotFoundError{Path: "Resource"}
	}
	return Resource{resource}, nil
}

// Create new resource.
func (s service) Create(ctx context.Context, org_id string, req CreateResourceRequest) (Resource, error) {

	// Validate resource request.
	if err := req.Validate(); err != nil {
		s.logger.Error("Error while validating resource request.")
		return Resource{}, &util.InvalidInputError{Path: "Invalid input for resource."}
	}

	// Check resource already exists.
	exists, _ := s.repo.ExistByKey(ctx, req.Key)
	if exists {
		s.logger.Debug("Resource already exists.")
		return Resource{}, &util.AlreadyExistsError{Path: "Resource : " + req.Key + " already exists."}
	}

	// Generate resource id.
	id := entity.GenerateID()

	err := s.repo.Create(ctx, org_id, entity.Resource{
		ID:   id,
		Key:  req.Key,
		Name: req.Name,
	})
	if err != nil {
		s.logger.Error("Error while creating resource.",
			zap.String("organization_id", org_id),
			zap.String("resource key", req.Key))
		return Resource{}, err
	}
	return s.Get(ctx, org_id, id)
}

// Update resource.
func (s service) Update(ctx context.Context, org_id string, id string, req UpdateResourceRequest) (Resource, error) {

	// Validate resource request.
	if err := req.Validate(); err != nil {
		s.logger.Error("Error while validating resource request.")
		return Resource{}, &util.InvalidInputError{Path: "Invalid input for resource."}
	}

	resource, err := s.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Debug("Resource not exists.")
		return Resource{}, &util.NotFoundError{Path: "Resource " + id + " not exists."}
	}
	resource.Name = req.Name
	if err := s.repo.Update(ctx, org_id, resource.Resource); err != nil {
		s.logger.Error("Error while creating resource.",
			zap.String("organization_id", org_id),
			zap.String("resource_id", id))
		return Resource{}, err
	}
	return resource, err
}

// Delete resource.
func (s service) Delete(ctx context.Context, org_id string, id string) (Resource, error) {

	resource, err := s.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Error("Resource not exists.")
		return resource, &util.NotFoundError{Path: "Resource " + id + " not exists."}
	}
	if err = s.repo.Delete(ctx, org_id, id); err != nil {
		s.logger.Error("Error while deleting resource.",
			zap.String("organization_id", org_id),
			zap.String("resource_id", id))
		return Resource{}, err
	}
	return resource, nil
}

// Pagination filter.
type Filter struct {
	Cursor int    `json:"cursor" query:"cursor"`
	Limit  int    `json:"limit" query:"limit"`
	Name   string `json:"name" query:"name"`
}

// Get all resources.
func (s service) Query(ctx context.Context, org_id string, filter Filter) ([]Resource, error) {

	result := []Resource{}
	items, err := s.repo.Query(ctx, org_id, filter)
	if err != nil {
		s.logger.Error("Error while retrieving all resources.",
			zap.String("organization_id", org_id))
		return []Resource{}, err
	}

	for _, item := range items {
		result = append(result, Resource{item})
	}
	return result, err
}
