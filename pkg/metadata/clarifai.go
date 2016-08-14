package metadata

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/sprioc/clarifai-go"
	"github.com/sprioc/conductor/pkg/model"
	"github.com/sprioc/conductor/pkg/refs"
	"github.com/sprioc/conductor/pkg/store"
)

// TODO would like to check for NSFW
func init() {
	clarifaiClient = clarifai.NewClient(os.Getenv("CLARIFAI_ID"), os.Getenv("CLARIFAI_SECRET"))
	_, err := clarifaiClient.Info()
	if err != nil {
		log.Fatal(err)
	}

	info, err := clarifaiClient.Info()
	if err != nil {
		log.Println(info)
	}
}

var clarifaiClient *clarifai.Client
var workChan chan string
var quit chan bool
var ticker *time.Ticker

func SetupClarifai(imageChan chan string) {
	quit = make(chan bool)
	workChan = imageChan
	ticker = time.NewTicker(time.Minute)
}

func Start() {
	go work()
}

func Stop() {
	quit <- true
}

// Tags the image with the given shortcode
func work() {
	for {
		select {
		case <-ticker.C:

			var dbRefs []model.DBRef
			var urls []string

			for i := 0; i < 45; i++ {
				select {
				case code := <-workChan:
					dbRefs = append(dbRefs, refs.GetImageRef(code))
					urls = append(urls, fmt.Sprintf("http://images.sprioc.xyz/content/%s?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=1080&fit=max", code))
				default:
					break
				}
			}

			if len(urls) == 0 {
				continue
			}

			data, err := clarifaiClient.Tag(clarifai.TagRequest{URLs: urls})
			if err != nil {
				log.Println(err)
				return
			}

			var finalTags = filterTags(data.Results, .90)

			for i := 0; i < len(urls); i++ {
				err = store.Modify(dbRefs[i], bson.M{"$set": bson.M{"machine_tags": finalTags[i]}})
				if err != nil {
					log.Println(err)
				}
			}
		case <-quit:
			return
		}
	}
}

func filterTags(results []clarifai.TagResult, threshold float32) [][]string {
	var totalResults = make([][]string, len(results))

	for i := 0; i < len(results); i++ {
		tagProbs := results[i].Result.Tag

		var finalTags []string

		for j, prob := range tagProbs.Probs {
			if prob > threshold {
				finalTags = append(finalTags, tagProbs.Classes[j])
			}
		}
		totalResults[i] = finalTags
	}
	return totalResults
}
