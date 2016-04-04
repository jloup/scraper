package html

import (
	"os"
	"reflect" // to compare maps
	"strings"
	"testing"
)

func TestHtmlScraper(t *testing.T) {
	s := `<div>
	      <p>Links:</p>
        <ul>
          <li class='myList'>
            <a id="1" href="link1" data="data">Foo</a>
            <a id="2" href="link2">Foo</a>
          </li>
          <li class='myList'>NOTHING</li>
        </ul>
        </div>`

	var want = []map[string]interface{}{
		{"class": "myList", "data": "data", "type": "ul", "href": "link1"},
		{"class": "myList", "type": "ul", "href": "link1"},
		{"class": "myList", "type": "ul", "href": "link2"},
	}

	f, err := os.Open("testdata/TestScraper/config.json")

	if err != nil {
		t.Fatalf("cannot read file : %v", err)
	}

	scrapers, err := JSONToHtmlScraper(f)

	if err != nil {
		t.Fatalf("cannot parse JSON : %v", err)
	}

	arr, err := ScrapHTML(scrapers, strings.NewReader(s), "")

	if err != nil {
		t.Fatalf("error while scrapping %v", err)
	} else {
		if len(arr) != len(want) {
			t.Fatalf("wrong number of item returned %v in %v", len(arr), arr)
		}

		for i, _ := range arr {
			if !reflect.DeepEqual(want[i], arr[i]) {
				t.Fatalf("DATA SCRAPPED NOT VALID => at index %v returned %v WANT %v\n", i, want[i], arr[i])
			}
		}
		t.Log(arr)
	}

}
