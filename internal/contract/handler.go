package contract

import (
	"docusign/internal/get_token"
	"docusign/internal/helper"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ContractInterfaceHandler interface {
	CreateContractHandler(c echo.Context) error
}

type ContractHandler struct {
	ContractInterface ContractInterfaceService
}

func NewContractHandler(ContractInterface ContractInterfaceService) *ContractHandler {
	return &ContractHandler{ContractInterface}
}

// CreateContractHandler processes the creation of contracts.
// @Summary Contract Creation
// @Description Creates a contract based on the provided data
// @Tags Contracts
// @Accept json
// @Produce json
// @Param request body ContractRequestCreate true "Contract data to be created"
// @Success 200 {object} ContractResponse "Returns the created contract"
// @Failure 400 {object} string "Validation error in the payload"
// @Failure 500 {object} string "Internal server error"
// @Router /docusign/send [post]
func (h *ContractHandler) CreateContractHandler(c echo.Context) error {
	var request ContractRequestCreate
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err := helper.Validate(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	payload := get_token.GetPayloadToken(c)
	data := ContractRequestCreate{
		ProviderName:   request.ProviderName,
		ProviderEmail:  request.ProviderEmail,
		ContractType:   request.ContractType,
		IdBatchControl: request.IdBatchControl,
		TenantId:       payload.TenantID,
		AccessId:       payload.AccessID,
	}

	result, err := h.ContractInterface.CreateContractService(c.Request().Context(), data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}
