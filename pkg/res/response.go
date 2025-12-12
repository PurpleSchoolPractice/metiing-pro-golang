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
		w.Write([]byte("ошибка: не удалось преобразовать JSON"))
		return
	}

	w.WriteHeader(statusCode)
	w.Write(buf.Bytes())
}
