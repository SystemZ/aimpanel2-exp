package handler

import (
	"github.com/gorilla/mux"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/master/service/host"
	"net/http"
)

// FIXME consider moving it to REST API with docs
func DeploymentScript(w http.ResponseWriter, r *http.Request) {
	httpParams := mux.Vars(r)
	hostToken := httpParams["hostToken"]
	err, result := host.GenerateDeploymentScript(hostToken)

	if err == nil {
		w.Header().Set("Content-Type", "text/x-sh;charset=utf-8")
		w.Write([]byte(result))
		return
	} else if err.(*lib.Error).ErrorCode == ecode.HostNotFound {
		lib.ReturnError(w, http.StatusNotFound, ecode.HostNotFound, err)
		return
	} else if err.(*lib.Error).ErrorCode == ecode.TemplateGenerationFailed {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.TemplateGenerationFailed, err)
		return
	}

}
