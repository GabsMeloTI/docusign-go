package contract

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type WebhookHandler struct {
	Service ContractService
}

func NewWebhookHandler(service ContractService) *WebhookHandler {
	return &WebhookHandler{Service: service}
}

func (h *WebhookHandler) HandleWebhook(c echo.Context) error {
	body := c.Request().Body
	defer body.Close()

	var notification map[string]interface{}
	if err := json.NewDecoder(body).Decode(&notification); err != nil {
		fmt.Println("Erro ao parsear JSON:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Erro ao parsear JSON",
		})
	}

	event, ok := notification["event"].(string)
	if !ok || event != "envelope-completed" {
		fmt.Println("Evento recebido não é de conclusão do envelope.")
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Evento ignorado",
		})
	}

	data, ok := notification["data"].(map[string]interface{})
	if !ok {
		fmt.Println("Dados da notificação inválidos.")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Dados inválidos",
		})
	}

	envelopeID, ok := data["envelopeId"].(string)
	if !ok {
		fmt.Println("Envelope ID não encontrado na notificação.")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Envelope ID não encontrado",
		})
	}

	_, err := h.Service.AssignContractService(c.Request().Context(), envelopeID)
	if err != nil {
		fmt.Println("Erro ao atualizar contrato:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Erro ao atualizar contrato",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Webhook processado com sucesso!",
	})
}
