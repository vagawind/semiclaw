package handler

import (
	"net/http"

	"github.com/vagawind/semiclaw/internal/types/interfaces"
	"github.com/gin-gonic/gin"
)

// SemiClawCloudHandler 处理 SemiClawCloud 凭证管理
type SemiClawCloudHandler struct {
	svc interfaces.SemiClawCloudService
}

// NewSemiClawCloudHandler 构造函数
func NewSemiClawCloudHandler(svc interfaces.SemiClawCloudService) *SemiClawCloudHandler {
	return &SemiClawCloudHandler{svc: svc}
}

type semiClawCloudCredentialsRequest struct {
	AppID     string `json:"app_id"     binding:"required"`
	AppSecret string `json:"app_secret" binding:"required"`
}

// SaveCredentials POST /api/v1/semiclawcloud/credentials
// 仅保存 APPID/APPSECRET 凭证到租户配置，不自动创建模型
//
// SaveCredentials godoc
// @Summary      保存 SemiClawCloud 凭证
// @Description  保存 APPID/APPSECRET 到当前租户配置（不自动创建模型）
// @Tags         SemiClawCloud
// @Accept       json
// @Produce      json
// @Param        request  body      map[string]interface{}  true  "{app_id, app_secret}"
// @Success      200      {object}  map[string]interface{}  "success: true"
// @Failure      400      {object}  map[string]interface{}  "请求参数错误"
// @Security     Bearer
// @Security     ApiKeyAuth
// @Router       /semiclawcloud/credentials [post]
func (h *SemiClawCloudHandler) SaveCredentials(c *gin.Context) {
	var req semiClawCloudCredentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.SaveCredentials(c.Request.Context(), req.AppID, req.AppSecret); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "凭证保存成功"})
}

// Status GET /api/v1/models/semiclawcloud/status
// 检查当前租户的 SemiClawCloud 凭证是否完好，如需重新初始化则返回 needs_reinit=true
//
// Status godoc
// @Summary      检查 SemiClawCloud 凭证状态
// @Description  检查当前租户的 SemiClawCloud 凭证是否完好；needs_reinit=true 表示需要重新保存
// @Tags         SemiClawCloud
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "凭证状态"
// @Failure      500  {object}  map[string]interface{}  "服务器错误"
// @Security     Bearer
// @Security     ApiKeyAuth
// @Router       /models/semiclawcloud/status [get]
func (h *SemiClawCloudHandler) Status(c *gin.Context) {
	result, err := h.svc.CheckStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
