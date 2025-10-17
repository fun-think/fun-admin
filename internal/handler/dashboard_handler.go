package handler

import (
	"fun-admin/internal/service"
	"fun-admin/pkg/admin/i18n"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DashboardHandler 仪表板处理器
type DashboardHandler struct {
	*Handler
	dashboardService service.DashboardServiceInterface
}

// NewDashboardHandler 创建仪表板处理器
func NewDashboardHandler(
	handler *Handler,
	dashboardService service.DashboardServiceInterface,
) *DashboardHandler {
	return &DashboardHandler{
		Handler:          handler,
		dashboardService: dashboardService,
	}
}

// GetDashboard 获取仪表板数据
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	// 获取语言参数，默认为中文
	language := c.Query("language")
	if language == "" {
		language = "zh-CN"
	}

	// 获取各种统计数据
	stats, err := h.dashboardService.GetDashboardStats(c, language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": i18n.Translate(language, "error.load_failed"),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": stats,
	})
}
