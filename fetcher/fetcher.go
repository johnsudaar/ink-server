package fetcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"github.com/johnsudaar/ink-server/config"
	"github.com/johnsudaar/ink-server/models"
	"github.com/pkg/errors"
)

type Fetcher interface {
	GetInkStatus() (models.InkStatus, error)
}

type PrinterFetcher struct {
	Config config.Config
}

func NewPrinterFetcher(c config.Config) PrinterFetcher {
	return PrinterFetcher{
		Config: c,
	}
}

func (f PrinterFetcher) GetInkStatus() (models.InkStatus, error) {
	res := models.InkStatus{}
	client := http.Client{}
	url := fmt.Sprintf("http://%s/PRESENTATION/HTML/TOP/PRTINFO.HTML", f.Config.PrinterIP)

	resp, err := client.Get(url)
	if err != nil {
		return res, errors.Wrap(err, "fail to make http request")
	}
	defer resp.Body.Close()

	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, errors.Wrap(err, "fail to read body")
	}
	// <img class='color' src='../../IMAGE/Ink_C.PNG' height='50'>

	regex := regexp.MustCompile(`<img class='color' src='../../IMAGE/Ink_(?P<color>[KCMY])\.PNG' height='(?P<value>[0-9]+)'>`)
	matches := regex.FindAllSubmatch(buffer, -1)
	for _, match := range matches {
		inkLevel, err := strconv.Atoi(string(match[2]))
		if err != nil {
			return res, errors.Wrap(err, "fail to get ink level")
		}
		switch string(match[1]) {
		case "K":
			res.Black = inkLevel
		case "C":
			res.Cyan = inkLevel
		case "M":
			res.Magenta = inkLevel
		case "Y":
			res.Yellow = inkLevel
		}
	}
	return res, nil
}
