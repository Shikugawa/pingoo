package consumer

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Shikugawa/pingoo/pkg/provider"
	"github.com/slack-go/slack"
)

type SlackHandler struct {
	callbacks provider.ProviderCallbacks
}

func NewSlackHandler(callbacks provider.ProviderCallbacks) *SlackHandler {
	return &SlackHandler{callbacks: callbacks}
}

func (s *SlackHandler) Start(port int) {
	http.HandleFunc("/actions", s.actionHandler)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func (s *SlackHandler) actionHandler(w http.ResponseWriter, r *http.Request) {
	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := s.callbacks.OnDelete(payload.CallbackID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
