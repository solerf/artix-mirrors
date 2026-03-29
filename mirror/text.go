package mirror

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

func FromText(addHttp bool, content io.Reader) (Group, error) {
	scanner := bufio.NewScanner(content)

	sanitize := func(l string) string {
		return strings.ReplaceAll(strings.ReplaceAll(l, "#", ""), " ", "")
	}

	getServerUrl := func(l string) string {
		if split := strings.Split(l, "Server="); len(split) == 2 {
			return split[1]
		}
		return ""
	}

	group := make(Group, 0, 15)
	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("text parsing: %w", err)
		}

		line := strings.TrimSpace(sanitize(scanner.Text()))
		if len(line) == 0 {
			continue
		}

		url := getServerUrl(line)
		if len(url) == 0 {
			continue
		}

		if strings.HasPrefix(url, "http://") && !addHttp {
			continue
		}

		group = append(group, Server{Url: url})
	}

	if len(group) == 0 {
		return nil, errors.New("no mirrors found")
	}
	return group, nil
}
