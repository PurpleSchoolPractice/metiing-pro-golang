package res

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		errorBody := map[string]string{
			"error": "ошибка: не удалось преобразовать JSON",
		}
		b, _ := json.Marshal(errorBody)
		w.Write(b)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(buf.Bytes())
}
