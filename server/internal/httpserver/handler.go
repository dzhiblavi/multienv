package httpserver

import (
	"errors"
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
	defer r.Body.Close()

	bodyBytes, err := readRequestBody(r)
	if err != nil {
		code := http.StatusBadRequest
		if errors.Is(err, errUnsupportedContentEncoding) {
			code = http.StatusUnsupportedMediaType
		}
		writeResponse(w, r, code, dispatcher.Response{Error: err})
		return
	}

	dispatch(disp, w, r, string(bodyBytes))
}

func dispatch(disp dispatcher.Dispatcher, w http.ResponseWriter, r *http.Request, payload string) {
	resp, err := disp.Dispatch(r.Context(), payload)
	if err != nil {
		writeResponse(w, r, httpCodeForError(err), dispatcher.Response{Error: err})
		return
	}

	writeResponse(w, r, http.StatusOK, resp)
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
