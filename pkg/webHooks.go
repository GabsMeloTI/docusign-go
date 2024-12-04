package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Falha ao ler o corpo da requisição", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var notification map[string]interface{}
	if err := json.Unmarshal(body, &notification); err != nil {
		http.Error(w, "Erro ao parsear JSON", http.StatusBadRequest)
		return
	}

	fmt.Println("Notificação recebida do DocuSign:")
	fmt.Printf("%+v\n", notification)

	if event, ok := notification["event"].(string); ok && event == "envelope-completed" {
		data, ok := notification["data"].(map[string]interface{})
		if ok {
			envelopeID := data["envelopeId"].(string)
			fmt.Printf("Envelope concluído! Envelope ID: %s\n", envelopeID)

			//err := updateContractStatus(envelopeID, "completed")
			//if err != nil {
			//	fmt.Printf("Erro ao atualizar o status do contrato: %v\n", err)
			//}
		}
	} else {
		fmt.Println("Evento recebido não é de conclusão do envelope.")
	}

	w.WriteHeader(http.StatusOK)
}
