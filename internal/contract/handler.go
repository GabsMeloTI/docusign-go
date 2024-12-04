package contract

import (
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

func (h *ContractHandler) CreateContractHandler(c echo.Context) error {
	var request ContractRequestCreate
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	data := ContractRequestCreate{
		ProviderName: request.ProviderName,
	}

	result, err := h.ContractInterface.CreateContractService(c.Request().Context(), data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}
