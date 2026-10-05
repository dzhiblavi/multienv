package httpserver

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"multienv/server/internal/dispatcher"
)

type handlerFunc func(http.ResponseWriter, *http.Request)

func makeHandler(disp dispatcher.Dispatcher) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(disp, w, r)
	}
}

func handler(
	disp dispatcher.Dispatcher,
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("request body read failed")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dispatcher.Response{Error: err})
		return
	}

	resp, err := disp.Dispatch(r.Context(), string(bodyBytes))
	if err != nil {
		log.Printf("server dispatch error: %v", err)
		w.WriteHeader(httpCodeForError(err))
		json.NewEncoder(w).Encode(dispatcher.Response{Error: err})
		return
	}

	log.Printf("request OK")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func httpCodeForError(err error) int {
	if errors.Is(err, dispatcher.ErrInternal) {
		return http.StatusBadGateway
	}
	if errors.Is(err, dispatcher.ErrInvalidParams) {
		return http.StatusBadRequest
	}
	if errors.Is(err, dispatcher.ErrUnknownAction) {
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}
