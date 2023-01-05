package epicrm_apiparts

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type SuccessJson struct {
	Data interface{} `json:"data"`
}

// Do not forget to manually return from your handler after calling this
func HandleBodyUnmarshalError(w http.ResponseWriter, service string, activity string, err error) {
	if errors.Is(err, ErrUnmarshal) {
		// TODO why don't use LogAndSendError() for this too?
		Send400Unmarshal(w)
	} else if err != nil {
		LogAndSendError(w, service, activity, err)
	}
}

// Do not forget to manually return from your handler after calling this
func LogAndSend500(w http.ResponseWriter, service string, activity string, fmt string, args ...interface{}) {
	Log(Error, service, activity, fmt, args)
	Send500(w)
}


// Do not forget to manually return from your handler after calling this
func LogAndSendError(w http.ResponseWriter, service string, activity string, err error) {
	LogError(service, activity, err)

	if errors.Is(err, ErrNotFound) {
		Send404(w)
	} else if errors.Is(err, ErrServerError) {
		Send500(w)
	} else {
		LogError("epicrm_apiparts", "LogErrorAndSend", fmt.Errorf("unhandled error type: %w", err))
		Send500(w)
	}
}

// Do not forget to manually return from your handler after calling this
func LogErrorAndSend500(w http.ResponseWriter, service string, activity string, err error) {
	LogError(service, activity, err)
	Send500(w)
}

func Send400InvalidOrgHierarchy(w http.ResponseWriter) {
	SendJsonError(w, 400, "INVALID_ORG_HIERARCHY")
}

func Send400InvalidOrgtype(w http.ResponseWriter) {
	SendJsonError(w, 400, "INVALID_ORGTYPE")
}

func Send400ParentOrgNoExist(w http.ResponseWriter) {
	SendJsonError(w, 400, "PARENT_ORG_NOEXIST")
}

func Send400Unmarshal(w http.ResponseWriter) {
	SendJsonError(w, 400, "BAD_REQUEST_BODY")
}

func Send403RestrictedOrgtype(w http.ResponseWriter) {
	SendJsonError(w, 403, "RESTRICTED_ORGTYPE")
}

// User not allowed to perform what was being attempted
func Send403UserNotAllowed(w http.ResponseWriter) {
	SendJsonError(w, 403, "USER_NOT_ALLOWED")
}

func Send500ReadOnlyMode(w http.ResponseWriter) {
	SendJsonError(w, 500, "READ_ONLY_MODE")
}

// Based on Go's http.Error() (but it sets the Content-Type to be plaintext)
func SendJsonError(w http.ResponseWriter, httpStatus int, publicError string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(httpStatus)
	fmt.Fprintln(w, "{\"error\":\"" + publicError + "\"}")
	closeResponseWriter(w) // XXX Not guaranteed to work
}

func Send401(w http.ResponseWriter, message string, scheme string) {
	w.Header().Set("WWW-Authenticate", scheme)
	SendJsonError(w, http.StatusUnauthorized, message)
}

func Send401Bearer(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	SendJsonError(w, http.StatusUnauthorized,
		"Unauthorized: missing or ill-formed bearer token")
}

func Send403(w http.ResponseWriter) {
	SendJsonError(w, http.StatusForbidden, "Forbidden")
}

func Send404(w http.ResponseWriter) {
	SendJsonError(w, http.StatusForbidden, "NOT_FOUND")
}

func Send500(w http.ResponseWriter) {
	SendJsonError(w, http.StatusInternalServerError, "Internal server error")
}

func SendJsonResponse(w http.ResponseWriter, dataJson interface{}) {
	SendJsonResponseNon200(w, dataJson, http.StatusOK)
}

func SendJsonResponseNon200(w http.ResponseWriter, dataJson interface{}, statusCode int) {
	json, err := json.Marshal(&SuccessJson{Data: dataJson})
	if err != nil {
		log.Print("Error on json.Marshal(): " + err.Error())
		Send500(w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(statusCode)
	w.Write(json)
	closeResponseWriter(w) // XXX Not guaranteed to work
}

func SendJsonStringResponse(w http.ResponseWriter, dataJson string) {
	json := "{\"data\":" + dataJson + "}\n"

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(json))
	closeResponseWriter(w) // XXX Not guaranteed to work
}

func SendOK(w http.ResponseWriter) {
	SendJsonStringResponse(w, "\"OK\"")
}

// XXX No reliable way to do this as of know, AFAIK. Just make sure all handlers
// issue `return` after error handling (even after calls to helpers for error
// handling).
func closeResponseWriter(w http.ResponseWriter) {}
