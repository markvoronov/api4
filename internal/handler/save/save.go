package save

import (
	"fmt"
	"io"
	"mime"
	"net/http"
)

func RootHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil || mediaType != "text/plain" {
		http.Error(w, "Content-Type must be text/plain", http.StatusUnsupportedMediaType)
		return
	}

	// Ограничиваем размер тела
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Body read error: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	bodyStr := string(body)
	fmt.Println(bodyStr)

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", bodyStr)
	w.WriteHeader(http.StatusCreated)
	// Пишем тело и проверяем ошибку
	if _, err := w.Write([]byte("http://localhost:8080/EwHXdJfB")); err != nil {
		//log.Printf("Failed to write response: %v", err)
		return
	}

}
