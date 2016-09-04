package refs

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sprioc/composer/pkg/model"
)

func init() {
	if os.Getenv("TYPE") == "dev" {
		baseurl = "http://localhost:8080/v0/"
	} else {
		baseurl = "https://api.sprioc.xyz/v0/"
	}
}

var baseurl string

func GetImageRef(shortcode string) model.Ref {
	return model.Ref{Collection: model.Images, ShortCode: shortcode}
}
func GetCollectionRef(shortcode string) model.Ref {
	return model.Ref{Collection: model.Collections, ShortCode: shortcode}
}

func GetUserRef(shortcode string) model.Ref {
	return model.Ref{Collection: model.Users, ShortCode: shortcode}
}

func GetRedisRef(redisString string) model.Ref {
	splitTag := strings.Split(redisString, ":")
	return model.Ref{
		ShortCode:  splitTag[1],
		Collection: model.RString(splitTag[0]),
	}
}

func GetRedisRefs(redisStrings []string) []model.Ref {
	var arr = make([]model.Ref, len(redisStrings))
	for i, ref := range redisStrings {
		arr[i] = GetRedisRef(ref)
	}
	return arr
}

func GetURL(ref model.Ref) string {
	if ref.Collection != "" && ref.ShortCode != "" {
		return fmt.Sprintf("%s%s/%s", baseurl, ref.Collection, ref.ShortCode)
	}
	return ""
}

func GetURLs(refs []model.Ref) []string {
	var arr = make([]string, len(refs))
	for i, ref := range refs {
		arr[i] = GetURL(ref)
	}
	return arr
}

func GetRef(url string) (model.Ref, error) {
	if !strings.HasPrefix(url, baseurl) {
		return model.Ref{}, errors.New("URL is of incorrect schema")
	}

	parts := strings.Split(strings.Replace(url, baseurl, "", 1), "/")
	return model.Ref{Collection: model.RString(parts[0]), ShortCode: parts[1]}, nil
}

func GetRefs(urls []string) []model.Ref {
	var arr = make([]model.Ref, len(urls))
	for i, url := range urls {
		ref, err := GetRef(url)
		if err != nil {
			continue
		}
		arr[i] = ref
	}
	return arr
}
