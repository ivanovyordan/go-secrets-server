package secret

import (
	"net/http"

	"github.com/ivanovyordan/go-secrets-server/model/secret"
	"github.com/ivanovyordan/go-secrets-server/tools/respond"

	"github.com/gorilla/mux"
)

func Post(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	message, err := secret.New(
		request.FormValue("secret"),
		request.FormValue("expireAfterViews"),
		request.FormValue("expireAfter"),
	)

	if err != nil {
		respond.Fail(response, http.StatusMethodNotAllowed, err.Error())
		return
	}

	respond.Success(response, request, message)
}

func Get(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	message, err := secret.Find(params["hash"])

	if err != nil {
		respond.Fail(response, http.StatusNotFound, err.Error())
		return
	}

	respond.Success(response, request, message)
}
