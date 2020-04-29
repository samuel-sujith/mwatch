package watch

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/version"
	"github.com/samuel-sujith/mwatch/pkg/types"
)

const acceptHeader = `application/openmetrics-text; version=0.0.1,text/plain;version=0.0.4;q=0.5,*/*;q=0.1`

var userAgentHeader = fmt.Sprintf("Prometheus/%s", version.Version)

//Targetwatching watches the given listenaddress and returns the http get response
func Targetwatching(cfg types.Cfg, w io.Writer, logger log.Logger) (string, error) {

	var gzipr *gzip.Reader
	var buf *bufio.Reader
	//TODO
	loggertarget := log.With(logger, "component", "targetwatcher")
	client := &http.Client{}

	level.Info(loggertarget).Log("msg", "Sending GET request to scraper", "listenaddress", cfg.Listenaddress)
	//fmt.Println("The address received in watcher is ", cfg.Listenaddress)

	req, err := http.NewRequest("GET", cfg.Listenaddress, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", acceptHeader)
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Set("User-Agent", userAgentHeader)
	req.Header.Set("X-Prometheus-Scrape-Timeout-Seconds", "5.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	level.Info(loggertarget).Log("msg", "Sent GET request to scraper")
	//fmt.Println("Response from server is ", resp)

	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	/*if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("server returned HTTP status %s", resp.Status)
	}*/

	if resp.Header.Get("Content-Encoding") != "gzip" {
		level.Info(loggertarget).Log("msg", "Encoding is not Zip")
		//fmt.Println("response is ", resp.Body)
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			return "", err
		}
		return resp.Header.Get("Content-Type"), nil
	}

	if gzipr == nil {
		//fmt.Println("Entered gzipr area")
		buf = bufio.NewReader(resp.Body)
		gzipr, err = gzip.NewReader(buf)
		if err != nil {
			return "", err
		}
	} else {
		buf.Reset(resp.Body)
		if err = gzipr.Reset(buf); err != nil {
			return "", err
		}
	}

	_, err = io.Copy(w, gzipr)
	gzipr.Close()
	if err != nil {
		return "", err
	}

	return resp.Header.Get("Content-Type"), nil

}
