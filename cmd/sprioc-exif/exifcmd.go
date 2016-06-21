package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/devinmcgloin/sprioc/src/api/metadata"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No File provided")
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	meta, err := metadata.GetExif(f)
	if err != nil {
		log.Fatal(err)
	}
	w := detailedWalker{}
	meta.Walk(&w)

	sort.Sort(w)

	for _, stat := range w.tags {
		fmt.Printf("%-30s = %2d = %s\n", stat.name, stat.t, stat.content)
	}
}

type detailedWalker struct {
	tags []metatags
}

type metatags struct {
	name    exif.FieldName
	t       tiff.DataType
	content string
}

func (dw *detailedWalker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	dw.tags = append(dw.tags, metatags{
		name:    name,
		t:       tag.Type,
		content: tag.String(),
	})
	return nil
}

func (dw detailedWalker) Len() int {
	return len(dw.tags)
}

func (dw detailedWalker) Swap(i, j int) {
	dw.tags[i], dw.tags[j] = dw.tags[j], dw.tags[i]
}

func (dw detailedWalker) Less(i, j int) bool {
	return strings.Compare(string(dw.tags[i].name), string(dw.tags[j].name)) < 0
}
