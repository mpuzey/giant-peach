package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

const METACRITIC_SCRAPE_BATCH_SIZE = 1000

func parse(tokenizer *html.Tokenizer) []html.Token {
	var reviewsHTML = make([]html.Token, METACRITIC_SCRAPE_BATCH_SIZE)

	for {
		tt := tokenizer.Next()

		switch {
		case tt == html.StartTagToken:
			t := tokenizer.Token()
			if t.Data == "li" {
				for _, a := range t.Attr {
					if a.Key == "class" {
						isReviewHTML := (a.Val == "review critic_review" || a.Val == "review critic_review first_review")
						if isReviewHTML {
							currentToken := t
							reviewsHTML = append(reviewsHTML, currentToken)
						}
					}
				}
			}
		case tt == html.ErrorToken:
			// End of document
			return reviewsHTML
		}
	}
}

func main() {

	publications := []string{"pitchfork", "consequence-of-sound", "rolling-stone",
		"the-guardian", "drowned-in-sound", "the-quietus",
		"sputnikmusic", "spin", "beats-per-minute-formerly-one-thirty-bpm",
		"the-observer-uk", "tiny-mix-tapes", "mojo", "musicomhcom", "under-the-radar",
		"exclaim", "paste-magazine", "american-songwriter", "now-magazine", "clash-music",
		"the-wire", "no-ripcord", "delusions-of-adequacy", "new-musical-express-nme",
		"the-independent-uk", "record-collector", "uncut", "diy-magazine",
		"alternative-press", "the-new-york-times", "the-405", "dusted-magazine",
		"the-av-club", "the-skinny"}

	client := &http.Client{}

	for _, publication := range publications {
		metacriticPublicationURL := fmt.Sprintf("https://www.metacritic.com/publication/%s?filter=music&num_items=%d", publication, METACRITIC_SCRAPE_BATCH_SIZE)
		req, err := http.NewRequest("GET", metacriticPublicationURL, nil)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()
		tokenizer := html.NewTokenizer(resp.Body)
		reviewsHTML := parse(tokenizer)

		fmt.Printf("%s", reviewsHTML)
	}

}
