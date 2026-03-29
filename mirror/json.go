package mirror

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

type mirrorsJson struct {
	Cutoff         int       `json:"cutoff"`
	LastCheck      time.Time `json:"last_check"`
	NumChecks      int       `json:"num_checks"`
	CheckFrequency int       `json:"check_frequency"`
	Urls           []struct {
		Url            string     `json:"url"`
		Protocol       string     `json:"protocol"`
		LastSync       *time.Time `json:"last_sync"`
		CompletionPct  float64    `json:"completion_pct"`
		Delay          *int       `json:"delay"`
		DurationAvg    *float64   `json:"duration_avg"`
		DurationStddev *float64   `json:"duration_stddev"`
		Score          *float64   `json:"score"`
		Active         bool       `json:"active"`
		Country        string     `json:"country"`
		CountryCode    string     `json:"country_code"`
		Isos           bool       `json:"isos"`
		Ipv4           bool       `json:"ipv4"`
		Ipv6           bool       `json:"ipv6"`
		Details        string     `json:"details"`
	} `json:"urls"`
	Version int `json:"version"`
}

func FromJson(addHttp bool, content io.Reader) (Group, error) {
	var mj mirrorsJson
	if err := json.NewDecoder(content).Decode(&mj); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, errors.New("no mirrors found")
		}
		return nil, fmt.Errorf("json parsing: %w", err)
	}

	group := make(Group, 0, len(mj.Urls))
	for i := 0; i < len(mj.Urls); i++ {
		server := mj.Urls[i]

		if !server.Active {
			continue
		}

		if strings.HasPrefix(server.Url, "http://") && !addHttp {
			continue
		}

		if !strings.HasPrefix(server.Url, "https://") {
			continue
		}

		group = append(group, Server{Country: server.Country, Url: server.Url})
	}
	return group, nil
}
