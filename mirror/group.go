package mirror

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"slices"
	"sync"

	"github.com/solerf/artix-mirrors/conf"
	"golang.org/x/sync/semaphore"
)

type Group []Server

func (g Group) Rank(c *conf.Settings) Group {
	rankedSorted := g.rank(c)
	slices.SortFunc(rankedSorted, func(a, b Server) int {
		return int(a.measure.Nanoseconds() - b.measure.Nanoseconds())
	})
	return rankedSorted[:c.TopMirrors]
}

func (g Group) Content() io.Reader {
	var buff bytes.Buffer
	writer := bufio.NewWriter(&buff)

	suffix := "$repo/os/$arch"
	var country string
	for _, gi := range g {
		if country != gi.Country {
			country = gi.Country
			_, _ = writer.Write([]byte(fmt.Sprintf("\n## %s\n", country)))
		}
		_, _ = writer.Write([]byte(fmt.Sprintf("Server = %s\n", gi.Url+suffix)))
	}

	_ = writer.Flush()
	return &buff
}

func (g Group) rank(c *conf.Settings) Group {
	run := func(s *semaphore.Weighted, w *sync.WaitGroup, t *Server /*, add func(*Server)*/) {
		defer w.Done()
		defer s.Release(1)

		err := t.measureConnection(c.MirrorTimeout)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "[ERROR] %s - %s => %v\n", t.Country, t.Url, err)
			return
		}
	}

	sem := semaphore.NewWeighted(int64(c.MaxConcurrency))
	var wg sync.WaitGroup
	wg.Add(len(g))

	for i := 0; i < len(g); i++ {
		_ = sem.Acquire(context.Background(), 1)
		go run(sem, &wg, &g[i])
	}
	wg.Wait()

	result := make(Group, 0, len(g))
	for i := 0; i < len(g); i++ {
		if g[i].measure > 0 {
			result = append(result, g[i])
		}
	}
	return slices.Clip(result)
}
