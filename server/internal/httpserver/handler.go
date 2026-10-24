package httpserver

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"multienv/server/internal/dispatcher"
)

func makeHandler(disp dispatcher.Dispatcher) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		handler(disp, w, r)
	}

	return AccessLogMiddleware()(handler)
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
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dispatcher.Response{Error: err})
		return
	}

	dispatch(disp, w, r, string(bodyBytes))
}

func dispatch(disp dispatcher.Dispatcher, w http.ResponseWriter, r *http.Request, payload string) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := disp.Dispatch(r.Context(), payload)
	if err != nil {
		w.WriteHeader(httpCodeForError(err))
		json.NewEncoder(w).Encode(dispatcher.Response{Error: err})
		return
	}

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
