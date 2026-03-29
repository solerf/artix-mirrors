package mirror

import (
	"bytes"
	"os"
	"slices"
	"strings"
	"testing"
)

func Test_Parse_Json_No_Http(t *testing.T) {
	file, _ := os.ReadFile("testdata/archmirrors.txt")
	group, _ := FromJson(false, bytes.NewBuffer(file))

	idx := slices.IndexFunc(group, func(t Server) bool {
		return strings.HasPrefix("http://", t.Url)
	})

	if idx != -1 {
		t.Fatalf("unexpected http server: %v", group[idx])
	}
}

func Test_Parse_Json_Missing_Mirrors(t *testing.T) {
	file, _ := os.ReadFile("testdata/arcch/one.txt")
	_, err := FromJson(false, bytes.NewBuffer(file))

	if err.Error() != "no mirrors found" {
		t.Fatal("unexpected error:", err)
	}
}
