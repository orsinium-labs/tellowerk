package plugins

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type WebAPI struct {
	state  *State
	logger *zap.Logger
}

func (api *WebAPI) Start() error {
	http.HandleFunc("/api/", api.handleAPI)
	go func() {
		err := http.ListenAndServe("127.0.0.1:8080", nil)
		if err != nil {
			api.logger.Error("cannot serve api", zap.Error(err))
		}
	}()
	return nil
}

func (api *WebAPI) handleAPI(w http.ResponseWriter, r *http.Request) {
	resp := APIResponse{
		Battery:  api.state.battery,
		Flying:   api.state.Flying(),
		Exposure: api.state.Exposure(),
	}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		api.logger.Error("cannot send response", zap.Error(err))
	}
}

type APIResponse struct {
	Battery  int8
	Flying   bool
	Exposure int8
}
