package permission

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"
)

func RegisterHandlers(r *echo.Group, service Service) {
	res := permission{service}
	router := r.Group("/:org/permission")
	router.POST("/create", res.create)
	router.POST("/check", res.check)
	router.POST("/list/resource", res.getobjectlist)
	router.POST("/list/role", res.getsubjectlist)
	router.POST("/check_by_username", res.checkbyusername)
	router.POST("/check_multi_actions", res.checkpermissions)
	router.POST("/check_multi_resources", res.checkall)
	router2 := r.Group("/:org_id/permission")
	router2.POST("/check_actions", res.checkActions)
	router2.PATCH("/update", res.patchPermissions)
}

type permission struct {
	service Service
}

// @Description Create tuple.
// @Tags        Permission
// @Accept      json
// @Param org path string true "Organization"
// @Param request body Tuple true "body"
// @Produce     json
// @Success     201
// @failure     400,403,500
// @Router      /{org}/permission/create [post]
func (r permission) create(c echo.Context) error {
	var input entity.Tuple
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}

	err := r.service.CreateTuple(context.Background(), c.Param("org"), "permission", input)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, "")
}

// @Description Check tuple.
// @Tags        Permission
// @Accept      json
// @Param org path string true "Organization"
// @Param request body Tuple true "body"
// @Produce     json
// @Success     201
// @failure     400,403,500
// @Router      /{org}/permission/check [post]
func (r permission) check(c echo.Context) error {
	var input entity.Tuple
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	allow, err := r.service.CheckTuple(context.Background(), c.Param("org"), "permission", input, true)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, allow)
}

// @Description Delete tuple.
// @Tags        Permission
// @Accept      json
// @Param org path string true "Organization"
// @Param request body Tuple true "body"
// @Produce     json
// @Success     201
// @failure     400,403,500
// @Router      /{org}/permission/delete [post]
func (r permission) delete(c echo.Context) error {
	var input entity.Tuple
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	allow, err := r.service.CheckTuple(context.Background(), c.Param("org"), "permission", input, false)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, allow)
}

// @Description Get objects.
// @Tags        Permission
// @Accept      json
// @Param org path string true "Organization"
// @Param request body Tuple true "body"
// @Produce     json
// @Success     201
// @failure     400,403,500
// @Router      /{org}/permission/list/resource [post]
func (r permission) getobjectlist(c echo.Context) error {
	var input entity.Tuple
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	allow, err := r.service.GetObjectListBySubject(context.Background(), c.Param("org"), "permission", input)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, allow)
}

// @Description Get subjects.
// @Tags        Permission
// @Accept      json
// @Param org path string true "Organization"
// @Param request body Tuple true "body"
// @Produce     json
// @Success     201
// @failure     400,403,500
// @Router      /{org}/permission/list/role [post]
func (r permission) getsubjectlist(c echo.Context) error {
	var input entity.Tuple
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	allow, err := r.service.GetSubjectListByObject(context.Background(), c.Param("org"), "permission", input)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, allow)
}

// @Description Check by username.
// @Tags        Permission
// @Accept      json
// @Param org path string true "Organization"
// @Param request body entity.CheckRequestWithPermissions true "body"
// @Produce     json
// @Success     201
// @failure     400,403,500
// @Router      /{org}/permission/check_by_username [post]
func (r permission) checkbyusername(c echo.Context) error {
	var input entity.Tuple
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	allow, err := r.service.CheckByUsername(context.Background(), c.Param("org"), "permission", input)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, allow)
}

// @Description Check by username.
// @Tags        Permission
// @Accept      json
// @Param org path string true "Organization"
// @Param request body entity.CheckRequestWithPermissions true "body"
// @Produce     json
// @Success     201
// @failure     400,403,500
// @Router      /{org}/permission/check_multi_actions [post]
func (r permission) checkpermissions(c echo.Context) error {
	var input entity.CheckRequestWithPermissions
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	allow, err := r.service.CheckPermissions(context.Background(), c.Param("org"), "permission", input)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, allow)
}

// @Description Check by username.
// @Tags        Permission
// @Accept      json
// @Param org path string true "Organization"
// @Param request body entity.CheckRequestAll true "body"
// @Produce     json
// @Success     201
// @failure     400,403,500
// @Router      /{org}/permission/check_multi_resources [post]
func (r permission) checkall(c echo.Context) error {
	var input entity.CheckRequestAll
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	allow, err := r.service.CheckAll(context.Background(), c.Param("org"), "permission", input)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, allow)
}

// @Description Check by username.
// @Tags        Permission
// @Accept      json
// @Param org_id path string true "Organization"
// @Param request body CheckActionsRequest true "body"
// @Produce     json
// @Success     201
// @failure     400,403,500
// @Router      /{org_id}/permission/check_actions [post]
func (r permission) checkActions(c echo.Context) error {
	var input CheckActionsRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	allow := r.service.CheckActions(context.Background(), c.Param("org_id"), "permission", input)

	return c.JSON(http.StatusOK, allow)
}

// @Description Patch Permissions.
// @Tags        Permission
// @Accept      json
// @Param org_id path string true "Organization"
// @Param request body PermissionPatchRequest true "body"
// @Produce     json
// @Success     201
// @failure     400,403,500
// @Router      /{org_id}/permission/update [patch]
func (r permission) patchPermissions(c echo.Context) error {
	var input PermissionPatchRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	_ = r.service.PatchPermissions(context.Background(), c.Param("org_id"), "permission", input)

	return c.JSON(http.StatusOK, "")
}
