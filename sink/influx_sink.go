package sink

import (
	"fmt"
	"net/url"

	client "github.com/influxdata/influxdb/client/v2"
	"github.com/johnsudaar/ink-server/models"
	"github.com/pkg/errors"
)

type InfluxSink struct {
	Addr     string
	Username string
	Password string
	Database string
}

func NewInfluxSink(databaseUrl string) (InfluxSink, error) {
	// url: https://user:pass@host:port/db_name

	res := InfluxSink{}
	url, err := url.Parse(databaseUrl)
	if err != nil {
		return res, errors.Wrap(err, "fail to get influx string")
	}

	res.Addr = fmt.Sprintf("%s://%s", url.Scheme, url.Host)
	if url.User != nil {
		res.Username = url.User.Username()
		res.Password, _ = url.User.Password()
	}

	res.Database = url.Path[1:]

	return res, nil
}

func (s InfluxSink) Send(st models.InkStatus) error {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:               s.Addr,
		Username:           s.Username,
		Password:           s.Password,
		InsecureSkipVerify: true,
	})

	if err != nil {
		return errors.Wrap(err, "fail to open connection to database")
	}
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  s.Database,
		Precision: "s",
	})

	if err != nil {
		return errors.Wrap(err, "fail to create batchpoint")
	}

	pt, err := client.NewPoint("printer", map[string]string{}, map[string]interface{}{
		"black_ink_level":   st.Black,
		"yellow_ink_level":  st.Yellow,
		"magenta_ink_level": st.Magenta,
		"cyan_ink_level":    st.Cyan,
	})

	if err != nil {
		return errors.Wrap(err, "fail to create point")
	}
	bp.AddPoint(pt)

	err = c.Write(bp)
	if err != nil {
		return errors.Wrap(err, "fail to write point")
	}

	return nil
}
