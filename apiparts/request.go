package epicrm_apiparts

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var (
	ErrReadBody  = errReadBody()
	ErrUnmarshal = errUnmarshal()
)

func errReadBody() error {
	return errors.New("error reading body")
}

func errUnmarshal() error {
	return errors.New("error unmarshalling")
}

func ReadJson(r *http.Request, obj interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO debug log
		return err
	}

	err = json.Unmarshal(body, obj)
	if err != nil {
		// TODO debug log
		return err
	}
	
	return nil
}
