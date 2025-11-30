package admin

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/edalferes/monetics/internal/modules/auth/domain"
)

type AdminHandler struct {
	ListRolesUC  interface{ Execute() ([]domain.Role, error) }
	CreateRoleUC interface {
		Execute(name string, permissionIDs []uint) error
	}
	DeleteRoleUC      interface{ Execute(name string) error }
	ListPermissionsUC interface {
		Execute() ([]domain.Permission, error)
	}
	CreatePermissionUC interface{ Execute(name string) error }
	DeletePermissionUC interface{ Execute(name string) error }
}

// ListRoles godoc
// @Summary List all roles
// @Tags Auth - Admin
// @Security BearerAuth
// @Success 200 {array} map[string]interface{}
// @Router /v1/admin/roles [get]
func (h *AdminHandler) ListRoles(c echo.Context) error {
	roles, err := h.ListRolesUC.Execute()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, roles)
}

// CreateRole godoc
// @Summary Create a new role
// @Tags Auth - Admin
// @Security BearerAuth
// @Param role body map[string]string true "Role name"
// @Success 201 {object} map[string]string
// @Router /v1/admin/roles [post]
func (h *AdminHandler) CreateRole(c echo.Context) error {
	var req map[string]string
	if err := c.Bind(&req); err != nil || req["name"] == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid role name"})
	}
	// permissionIDs can be extracted from body if needed
	if err := h.CreateRoleUC.Execute(req["name"], nil); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]string{"message": "role created"})
}

// DeleteRole godoc
// @Summary Delete a role
// @Tags Auth - Admin
// @Security BearerAuth
// @Param name path string true "Role name"
// @Success 204 {object} map[string]string
// @Router /v1/admin/roles/{name} [delete]
func (h *AdminHandler) DeleteRole(c echo.Context) error {
	name := c.Param("name")
	if err := h.DeleteRoleUC.Execute(name); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

// ListPermissions godoc
// @Summary List all permissions
// @Tags Auth - Admin
// @Security BearerAuth
// @Success 200 {array} map[string]interface{}
// @Router /v1/admin/permissions [get]
func (h *AdminHandler) ListPermissions(c echo.Context) error {
	perms, err := h.ListPermissionsUC.Execute()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, perms)
}

// CreatePermission godoc
// @Summary Create a new permission
// @Tags Auth - Admin
// @Security BearerAuth
// @Param permission body map[string]string true "Permission name"
// @Success 201 {object} map[string]string
// @Router /v1/admin/permissions [post]
func (h *AdminHandler) CreatePermission(c echo.Context) error {
	var req map[string]string
	if err := c.Bind(&req); err != nil || req["name"] == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid permission name"})
	}
	if err := h.CreatePermissionUC.Execute(req["name"]); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]string{"message": "permission created"})
}

// DeletePermission godoc
// @Summary Delete a permission
// @Tags Auth - Admin
// @Security BearerAuth
// @Param name path string true "Permission name"
// @Success 204 {object} map[string]string
// @Router /v1/admin/permissions/{name} [delete]
func (h *AdminHandler) DeletePermission(c echo.Context) error {
	name := c.Param("name")
	if err := h.DeletePermissionUC.Execute(name); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
