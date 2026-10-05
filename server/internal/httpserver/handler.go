package httpserver

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"multienv/server/internal/dispatcher"
)

type handlerFunc func(http.ResponseWriter, *http.Request)

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
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dispatcher.Response{Error: err})
		return
	}

	resp, err := disp.Dispatch(ctx, string(bodyBytes))
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
