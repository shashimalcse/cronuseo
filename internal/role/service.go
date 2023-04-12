package role

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/cache"
	"github.com/shashimalcse/cronuseo/internal/mongo_entity"
	"github.com/shashimalcse/cronuseo/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Service interface {
	Get(ctx context.Context, org_id string, id string) (Role, error)
	Query(ctx context.Context, org_id string, filter Filter) ([]Role, error)
	// QueryByUserID(ctx context.Context, org_id string, user_id string, filter Filter) ([]Role, error)
	Create(ctx context.Context, org_id string, input CreateRoleRequest) (Role, error)
	Update(ctx context.Context, org_id string, id string, input UpdateRoleRequest) (Role, error)
	Delete(ctx context.Context, org_id string, id string) error
	PatchPermission(ctx context.Context, org_id string, role_id string, req PatchRolePermissionRequest) error
}

type Role struct {
	mongo_entity.Role
}

type CreateRoleRequest struct {
	Identifier  string               `json:"identifier"`
	DisplayName string               `json:"display_name"`
	Users       []primitive.ObjectID `json:"users"`
}

func (m CreateRoleRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Identifier, validation.Required),
	)
}

type UpdateRoleRequest struct {
	DisplayName *string              `json:"display_name"`
	AddedUsers  []primitive.ObjectID `json:"added_users"`
	RemovedUser []primitive.ObjectID `json:"removed_users"`
}

type UpdateRole struct {
	DisplayName *string              `json:"display_name"`
	AddedUsers  []primitive.ObjectID `json:"added_users"`
	RemovedUser []primitive.ObjectID `json:"removed_users"`
}

type PatchRolePermissionRequest struct {
	AddedPermission   []mongo_entity.Permission `json:"added_permissions"`
	RemovedPermission []mongo_entity.Permission `json:"removed_permissions"`
}

type PatchRolePermission struct {
	AddedPermission   []mongo_entity.Permission `json:"added_permissions"`
	RemovedPermission []mongo_entity.Permission `json:"removed_permissions"`
}

func (m UpdateRoleRequest) Validate() error {

	return validation.ValidateStruct(&m)
}

type service struct {
	repo            Repository
	permissionCache cache.PermissionCache
	logger          *zap.Logger
}

func NewService(repo Repository, permissionCache cache.PermissionCache, logger *zap.Logger) Service {

	return service{repo: repo, permissionCache: permissionCache, logger: logger}
}

// Get role by id.
func (s service) Get(ctx context.Context, org_id string, id string) (Role, error) {

	role, err := s.repo.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Error("Error while getting the role.",
			zap.String("organization_id", org_id),
			zap.String("role_id", id))
		return Role{}, &util.NotFoundError{Path: "Role"}
	}
	return Role{*role}, nil
}

// Create role.
func (s service) Create(ctx context.Context, org_id string, req CreateRoleRequest) (Role, error) {

	// Validate role request.
	if err := req.Validate(); err != nil {
		s.logger.Error("Error while validating role creation request.")
		return Role{}, &util.InvalidInputError{Path: "Invalid input for role."}
	}

	exists, _ := s.repo.CheckRoleExistsByIdentifier(ctx, org_id, req.Identifier)
	if exists {
		s.logger.Debug("Role already exists.")
		return Role{}, &util.AlreadyExistsError{Path: "Role : " + req.Identifier + " already exists."}

	}

	// Generate role id.
	roleId := primitive.NewObjectID()

	for _, userId := range req.Users {
		exists, _ := s.repo.CheckUserExistById(ctx, org_id, userId.Hex())
		if !exists {
			return Role{}, &util.InvalidInputError{Path: "Invalid user id " + userId.String()}
		}
	}

	err := s.repo.Create(ctx, org_id, mongo_entity.Role{
		ID:          roleId,
		Identifier:  req.Identifier,
		DisplayName: req.DisplayName,
		Users:       req.Users,
	})

	if err != nil {
		s.logger.Error("Error while creating role.",
			zap.String("organization_id", org_id),
			zap.String("role identifier", req.Identifier))
		return Role{}, err
	}
	return s.Get(ctx, org_id, roleId.Hex())
}

// Update role.
func (s service) Update(ctx context.Context, org_id string, id string, req UpdateRoleRequest) (Role, error) {

	_, err := s.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Debug("Role not exists.", zap.String("role_id", id))
		return Role{}, &util.NotFoundError{Path: "Role " + id + " not exists."}
	}

	for _, userId := range req.AddedUsers {
		exists, _ := s.repo.CheckUserExistById(ctx, org_id, userId.Hex())
		if !exists {
			return Role{}, &util.InvalidInputError{Path: "Invalid role id " + userId.String()}
		}
	}
	for _, userId := range req.RemovedUser {
		exists, _ := s.repo.CheckUserExistById(ctx, org_id, userId.Hex())
		if !exists {
			return Role{}, &util.InvalidInputError{Path: "Invalid role id " + userId.String()}
		}
	}
	added_users := []primitive.ObjectID{}
	for _, userId := range req.AddedUsers {
		already_added, _ := s.repo.CheckUserAlreadyAssignToRoleById(ctx, org_id, id, userId.Hex())
		if !already_added {
			added_users = append(added_users, userId)
		}
	}

	removed_roles := []primitive.ObjectID{}
	for _, userId := range req.RemovedUser {
		already_added, _ := s.repo.CheckUserAlreadyAssignToRoleById(ctx, org_id, id, userId.Hex())
		if already_added {
			removed_roles = append(removed_roles, userId)
		}
	}

	if err := s.repo.Update(ctx, org_id, id, UpdateRole{
		DisplayName: req.DisplayName,
		AddedUsers:  added_users,
		RemovedUser: removed_roles,
	}); err != nil {
		s.logger.Error("Error while updating role.", zap.String("organization_id", org_id), zap.String("role_id", id))
		return Role{}, err
	}
	updatedRole, err := s.repo.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Debug("Role not exists.", zap.String("role_id", id))
		return Role{}, &util.NotFoundError{Path: "Role " + id + " not exists."}
	}
	return Role{*updatedRole}, nil
}

// Delete role.
func (s service) Delete(ctx context.Context, org_id string, id string) error {

	_, err := s.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Error("Resource not exists.", zap.String("resource_id", id))
		return &util.NotFoundError{Path: "Resource " + id + " not exists."}
	}
	if err = s.repo.Delete(ctx, org_id, id); err != nil {
		s.logger.Error("Error while deleting resource.",
			zap.String("organization_id", org_id),
			zap.String("resource_id", id))
		return err
	}
	return nil
}

type Filter struct {
	Cursor int    `json:"cursor" query:"cursor"`
	Limit  int    `json:"limit" query:"limit"`
	Name   string `json:"name" query:"name"`
}

// Get all roles.
func (s service) Query(ctx context.Context, org_id string, filter Filter) ([]Role, error) {

	result := []Role{}
	items, err := s.repo.Query(ctx, org_id)
	if err != nil {
		s.logger.Error("Error while retrieving all resources.",
			zap.String("organization_id", org_id))
		return []Role{}, err
	}

	for _, item := range *items {
		result = append(result, Role{item})
	}
	return result, err
}

// Get all roles by user id.
// func (s service) QueryByUserID(ctx context.Context, org_id string, user_id string, filter Filter) ([]Role, error) {

// 	items, err := s.repo.QueryByUserID(ctx, org_id, user_id, filter)
// 	if err != nil {
// 		s.logger.Error("Error while retrieving all roles.",
// 			zap.String("organization_id", org_id),
// 			zap.String("user_id", user_id))
// 		return nil, err
// 	}
// 	result := []Role{}
// 	for _, item := range items {
// 		result = append(result, Role{item})
// 	}
// 	return result, err
// }

// Patch permissions to role.
func (s service) PatchPermission(ctx context.Context, org_id string, role_id string, req PatchRolePermissionRequest) error {

	for _, permission := range req.AddedPermission {

		exists, _ := s.repo.CheckResourceActionExists(ctx, org_id, permission.Resource, permission.Action)
		if !exists {
			return &util.InvalidInputError{Path: "Invalid permission, Resource : " + permission.Resource + " Action : " + permission.Action}
		}
	}

	if err := s.repo.PatchPermission(ctx, org_id, role_id, PatchRolePermission{
		AddedPermission: req.AddedPermission,
	}); err != nil {
		s.logger.Error("Error while patching role permissions.", zap.String("organization_id", org_id), zap.String("role_id", role_id))
		return err
	}

	return nil
}
