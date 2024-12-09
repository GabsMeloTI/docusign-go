package get_token

import (
	"docusign/internal/model"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"time"
)

func GetPayloadToken(c echo.Context) model.PayloadDTO {
	strID, _ := c.Get("token_id").(uuid.UUID)
	strUserID, _ := c.Get("token_user_id").(string)
	strUserName, _ := c.Get("token_user_name").(string)
	strExpiryAt, _ := c.Get("token_expiry_at").(time.Time)
	strAccessKey, _ := c.Get("token_access_key").(int64)
	strAccessID, _ := c.Get("token_access_ID").(int64)
	strTenantID, _ := c.Get("token_tenant_id").(string)

	tenantID, _ := uuid.Parse(strTenantID)

	return model.PayloadDTO{
		ID:           strID,
		UserID:       strUserID,
		UserNickname: strUserName,
		ExpiryAt:     strExpiryAt,
		AccessKey:    strAccessKey,
		AccessID:     strAccessID,
		TenantID:     tenantID,
	}
}
