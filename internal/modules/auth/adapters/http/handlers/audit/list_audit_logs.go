package audit

import (
	"net/http"

	"github.com/edalferes/monetics/internal/modules/auth/adapters/http/responses"
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
	pkgresponses "github.com/edalferes/monetics/pkg/responses"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	auditLogRepo interfaces.AuditLogRepository
}

func NewHandler(auditLogRepo interfaces.AuditLogRepository) *Handler {
	return &Handler{
		auditLogRepo: auditLogRepo,
	}
}

// ListAuditLogs godoc
// @Summary List all audit logs
// @Description Get a list of all audit logs
// @Tags Audit
// @Security BearerAuth
// @Produce json
// @Success 200 {object} responses.SuccessResponse{data=[]dto.AuditLogResponse}
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/audit-logs [get]
func (h *Handler) ListAuditLogs(c echo.Context) error {
	logs, err := h.auditLogRepo.ListAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, pkgresponses.ErrorResponse{
			Error: "Failed to retrieve audit logs",
		})
	}

	// Convert domain to response
	response := make([]responses.AuditLogResponse, len(logs))
	for i, log := range logs {
		response[i] = responses.AuditLogResponse{
			ID:        log.ID,
			UserID:    log.UserID,
			Username:  log.Username,
			Action:    log.Action,
			Status:    log.Status,
			IP:        log.IP,
			Details:   log.Details,
			CreatedAt: log.CreatedAt,
		}
	}

	return pkgresponses.OK(c, response)
}
