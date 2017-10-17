package main

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/johnsudaar/ink-server/config"
	"github.com/johnsudaar/ink-server/fetcher"
	"github.com/johnsudaar/ink-server/sink"
)

func main() {
	config := config.InitConfig()
	fetcher := fetcher.NewPrinterFetcher(config)
	var s sink.Sink
	var err error
	switch config.Sink {
	case "http":
		s = sink.NewHTTPSink(config.SinkURL)
	case "influx":
		s, err = sink.NewInfluxSink(config.SinkURL)
		if err != nil {
			panic(err)
		}

	default:
		panic("invalid sink: " + config.Sink)
	}

	ticker := time.NewTicker(1 * time.Hour)
	for {
		logrus.Info("Starting ink fetcher")
		status, err := fetcher.GetInkStatus()
		if err != nil {
			logrus.Error(err, "fail to get status")
			continue
		}
		logrus.Info("Sending status")
		err = s.Send(status)
		if err != nil {
			logrus.Error(err, "fail to send status")
			continue
		}
		logrus.Info("Done !")
		<-ticker.C
	}
}
