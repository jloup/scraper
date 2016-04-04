package html_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/jloup/scraper/html"
	"github.com/jloup/scraper/js"
	"github.com/jloup/scraper/node"
)

func ExampleScrapHTML() {
	// we avoid checking error for reading ease
	soundcloud, _ := html.JSONFileToHtmlScraper("testdata/exampledata/soundcloud.json")
	html.AddScraperConfig("soundcloud", soundcloud)

	youtube, _ := html.JSONFileToHtmlScraper("testdata/exampledata/youtube.json")
	html.AddScraperConfig("youtube", youtube)

	s, _ := html.JSONFileToHtmlScraper("testdata/exampledata/config.json")

	f, _ := os.Open("testdata/exampledata/website.html")

	items, _ := html.ScrapHTML(s, f, "text/html")

	for _, item := range items {
		if _, ok := item["type"]; ok {
			fmt.Printf("Soundcloud type '%s' id %s\n", item["type"], item["id"])
		} else {
			fmt.Printf("Youtube id %s\n", item["id"])
		}
	}

	//Output:
	//Soundcloud type 'tracks' id 196618442
	//Soundcloud type 'tracks' id 195052762
	//Soundcloud type 'tracks' id 195529989
	//Soundcloud type 'tracks' id 196295430
	//Youtube id LUP89HZBWhI
	//Soundcloud type 'tracks' id 187956230
	//Soundcloud type 'tracks' id 196308250
	//Youtube id 3VvQlPjY878
	//Soundcloud type 'tracks' id 195962464
	//Soundcloud type 'tracks' id 196130447
}

func ExampleJS() {

	const jsScrapConfig = `
{
"config": [
{
	"type": "*",
	"identifier": "createElement",
	"childs": [
					{
						"type": "*",
						"identifier": "videos",
						"childs": [
										{
											"type": "{}->map",
											"config": {"id": "1", "title": "1", "views": "1"}
										}
									 ]
					}
	          ]
}
]
}
`

	const htmlScrapConfig = `
{
"config": [
{
	"tag": "script",
	"validators": [
						  {"name": "exists", "attr": "type"},
                    {"name": "attrEquals", "attr": "type", "value": "text/javascript"}
					  ],
	"childs": [
					{
						"tag": "text",
						"extractors": [{"name": "js", "config": "arte-videos"}]
					}
	          ]
}
]
}
`

	f, _ := os.Open("testdata/arte.html")

	js.AddScraperConfig("arte-videos", func(config map[string]string) ([]*node.ScraperNode, error) {
		return js.JSONToJsScraper(strings.NewReader(jsScrapConfig))
	})

	htmlscrapers, _ := html.JSONToHtmlScraper(strings.NewReader(htmlScrapConfig))
	videos, _ := html.ScrapHTML(htmlscrapers, f, "text/html")

	for i, video := range videos {
		fmt.Printf("#%v %v (%v views) '%v'\n", i, video["id"], video["views"], video["title"])
	}

	//Output:
	//#0 064537-000-A (139197 views) 'Daech, paroles de d\u00e9serteurs'
	//#1 057849-000-A (109419 views) 'Ces nouvelles drogues qui submergent l'Europe'
	//#2 053958-000-A (93879 views) '\u00c9thologie : ce que ressentent les animaux '
	//#3 053986-001-A (87453 views) 'La fin des Ottomans (1\/2)'
	//#4 057482-000-A (73036 views) 'Branchez les guitares\u00a0!'
	//#5 064867-000-A (60366 views) 'Les armes des djihadistes'
	//#6 061390-000-A (49000 views) 'La plan\u00e8te Fifa'
	//#7 061725-000-A (46587 views) 'Elisabeth I - Au service secret de sa Majest\u00e9'
	//#8 060139-010-A (45730 views) 'Le dessous des cartes'
	//#9 063685-000-A (45302 views) 'Kurdistan, la guerre des filles'
	//#10 063711-005-A (44793 views) 'Le dessous des cartes'
	//#11 060139-011-A (43954 views) 'Le dessous des cartes'
	//#12 053331-000-A (43346 views) 'Hannah Arendt - Du devoir de la d\u00e9sob\u00e9issance civile'
	//#13 062928-001-A (42990 views) 'Personne ne bouge !'
	//#14 057884-004-A (42376 views) 'Peaky Blinders - Saison 2 (4\/6)'
	//#15 062943-000-A (40638 views) 'Les femmes de pouvoir'
	//#16 062928-002-A (39768 views) 'Personne ne bouge !'
	//#17 053986-002-A (38192 views) 'La fin des Ottomans (2\/2)'
	//#18 062928-004-A (37797 views) 'Personne ne bouge !'
	//#19 060139-002-A (37089 views) 'Le dessous des cartes'
	//#20 060139-013-A (36870 views) 'Le dessous des cartes'
	//#21 060139-012-A (36357 views) 'Le Dessous des cartes'
	//#22 050143-000-A (33972 views) 'Aviation \u00e9lectrique '
	//#23 050294-000-A (32456 views) 'Le Baron Rouge'

}
