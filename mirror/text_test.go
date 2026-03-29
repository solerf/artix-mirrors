package mirror

import (
	"bytes"
	"os"
	"slices"
	"strings"
	"testing"
)

func Test_Parse_Text_No_Http(t *testing.T) {
	file, _ := os.ReadFile("testdata/artixmirrors.txt")
	group, _ := FromText(false, bytes.NewBuffer(file))

	idx := slices.IndexFunc(group, func(t Server) bool {
		return strings.HasPrefix("http://", t.Url)
	})

	if idx != -1 {
		t.Fatalf("unexpected http server: %v", group)
	}
}

func Test_Parse_Text_Missing_Mirrors(t *testing.T) {
	file, _ := os.ReadFile("testdata/artix/artix.txt")
	_, err := FromText(false, bytes.NewBuffer(file))

	if err.Error() != "no mirrors found" {
		t.Fatal("unexpected error:", err)
	}
}
