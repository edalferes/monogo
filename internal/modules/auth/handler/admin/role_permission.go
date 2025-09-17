package admin

import (
	"net/http"

	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/repository"
	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	RoleRepo       repository.RoleRepository
	PermissionRepo repository.PermissionRepository
}

// ListRoles godoc
// @Summary Lista todas as roles
// @Tags admin
// @Security BearerAuth
// @Success 200 {array} map[string]interface{}
// @Router /v1/admin/roles [get]
func (h *AdminHandler) ListRoles(c echo.Context) error {
	roles, err := h.RoleRepo.ListAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, roles)
}

// CreateRole godoc
// @Summary Cria uma nova role
// @Tags admin
// @Security BearerAuth
// @Param role body map[string]string true "Role name"
// @Success 201 {object} map[string]string
// @Router /v1/admin/roles [post]
func (h *AdminHandler) CreateRole(c echo.Context) error {
	var req map[string]string
	if err := c.Bind(&req); err != nil || req["name"] == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid role name"})
	}
	if err := h.RoleRepo.Create(&domain.Role{Name: req["name"]}); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]string{"message": "role created"})
}

// DeleteRole godoc
// @Summary Remove uma role
// @Tags admin
// @Security BearerAuth
// @Param name path string true "Role name"
// @Success 204 {object} map[string]string
// @Router /v1/admin/roles/{name} [delete]
func (h *AdminHandler) DeleteRole(c echo.Context) error {
	name := c.Param("name")
	if err := h.RoleRepo.DeleteByName(name); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

// ListPermissions godoc
// @Summary Lista todas as permissões
// @Tags admin
// @Security BearerAuth
// @Success 200 {array} map[string]interface{}
// @Router /v1/admin/permissions [get]
func (h *AdminHandler) ListPermissions(c echo.Context) error {
	perms, err := h.PermissionRepo.ListAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, perms)
}

// CreatePermission godoc
// @Summary Cria uma nova permissão
// @Tags admin
// @Security BearerAuth
// @Param permission body map[string]string true "Permission name"
// @Success 201 {object} map[string]string
// @Router /v1/admin/permissions [post]
func (h *AdminHandler) CreatePermission(c echo.Context) error {
	var req map[string]string
	if err := c.Bind(&req); err != nil || req["name"] == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid permission name"})
	}
	if err := h.PermissionRepo.Create(&domain.Permission{Name: req["name"]}); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]string{"message": "permission created"})
}

// DeletePermission godoc
// @Summary Remove uma permissão
// @Tags admin
// @Security BearerAuth
// @Param name path string true "Permission name"
// @Success 204 {object} map[string]string
// @Router /v1/admin/permissions/{name} [delete]
func (h *AdminHandler) DeletePermission(c echo.Context) error {
	name := c.Param("name")
	if err := h.PermissionRepo.DeleteByName(name); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
