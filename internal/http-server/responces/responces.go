package responces

import (
	"encoding/json"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, r *http.Request, statuscode int, resp interface{}) {
	w.Header().Set("content-Type", "application/json")
	w.WriteHeader(statuscode)
	jresp, _ := json.Marshal(resp)
	w.Write(jresp)
}
