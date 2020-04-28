package watch

import (
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/common/version"
	"k8s.io/apimachinery/pkg/labels"
)

const acceptHeader = `application/openmetrics-text; version=0.0.1,text/plain;version=0.0.4;q=0.5,*/*;q=0.1`

var userAgentHeader = fmt.Sprintf("Prometheus/%s", version.Version)

// Target refers to a singular HTTP or HTTPS endpoint.
type Target struct {
	// Labels before any processing.
	discoveredLabels labels.Labels
	// Any labels that are added to this target and its metrics.
	labels labels.Labels
	// Additional URL parameters that are part of the target URL.
	params url.Values

	mtx                sync.RWMutex
	lastError          error
	lastScrape         time.Time
	lastScrapeDuration time.Duration
	health             TargetHealth
	metadata           MetricMetadataStore
}

//TargetScraper implements the scraper interface for a target.
type TargetScraper struct {
	*Target

	client  *http.Client
	req     *http.Request
	timeout time.Duration

	gzipr *gzip.Reader
	buf   *bufio.Reader
}

//Createtargetscraper takes in the target specified and returns a scraper which can then be used
//to scrape the metrics off the target endpoint
func Createtargetscraper(target Target) *TargetScraper {
	//TODO
	ts := TargetScraper{
		//Fill up details
		client:,
	}
	return &ts
}

func (s *targetScraper) scrape(ctx context.Context, w io.Writer) (string, error) {
	if s.req == nil {
		req, err := http.NewRequest("GET", s.URL().String(), nil)
		if err != nil {
			return "", err
		}
		req.Header.Add("Accept", acceptHeader)
		req.Header.Add("Accept-Encoding", "gzip")
		req.Header.Set("User-Agent", userAgentHeader)
		req.Header.Set("X-Prometheus-Scrape-Timeout-Seconds", fmt.Sprintf("%f", s.timeout.Seconds()))

		s.req = req
	}

	resp, err := s.client.Do(s.req.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("server returned HTTP status %s", resp.Status)
	}

	if resp.Header.Get("Content-Encoding") != "gzip" {
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			return "", err
		}
		return resp.Header.Get("Content-Type"), nil
	}

	if s.gzipr == nil {
		s.buf = bufio.NewReader(resp.Body)
		s.gzipr, err = gzip.NewReader(s.buf)
		if err != nil {
			return "", err
		}
	} else {
		s.buf.Reset(resp.Body)
		if err = s.gzipr.Reset(s.buf); err != nil {
			return "", err
		}
	}

	_, err = io.Copy(w, s.gzipr)
	s.gzipr.Close()
	if err != nil {
		return "", err
	}
	return resp.Header.Get("Content-Type"), nil
}
