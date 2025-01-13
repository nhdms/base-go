package handlers

import (
	"github.com/gorilla/mux"
	transhttp "github.com/nhdms/base-go/pkg/transport"
	"github.com/nhdms/base-go/proto/exmsg/services"
	"github.com/spf13/cast"
	"math/rand"
	"net/http"
	"time"
)

type GetUserByIdHandler struct {
	UserClient services.UserService
}

func (h *GetUserByIdHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	resp, err := h.UserClient.GetUserByID(r.Context(), &services.UserRequest{
		UserId: cast.ToInt64(id), // Replace with actual user ID
	})

	if err != nil {
		transhttp.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	transhttp.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"now":  time.Now().String(),
		"resp": resp,
	})
}

func randomInt(i int) int64 {
	return int64(rand.Intn(i) + 1)
}
