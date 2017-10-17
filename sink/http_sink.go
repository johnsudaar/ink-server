package sink

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/johnsudaar/ink-server/models"
	"github.com/pkg/errors"
)

type HTTPSink struct {
	URL string
}

func NewHTTPSink(url string) HTTPSink {
	return HTTPSink{
		URL: url,
	}
}

func (s HTTPSink) Send(status models.InkStatus) error {
	res, err := json.Marshal(status)
	if err != nil {
		return errors.Wrap(err, "fail to marshal ink status")
	}

	client := http.Client{}
	resp, err := client.Post(s.URL, "application/json", bytes.NewBuffer(res))
	if err != nil {
		return errors.Wrap(err, "fail to make request")
	}
	resp.Body.Close()
	return nil
}
