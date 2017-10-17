package sink

import "github.com/johnsudaar/ink-server/models"

type Sink interface {
	Send(models.InkStatus) error
}
