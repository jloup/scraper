package js_test

import (
	"fmt"
	"strings"

	"github.com/jloup/scraper/js"
)

func ExampleScrapJS() {
	const jscode = `
	(function() {

		var element = MyObject.play([{"videoid": "Yfhuyr-T", "views": 10},
		                                   {"videoid": "dhy-Ggdz", "views": 356980}
													 ]);
	})();
	`

	const scrapconfig = `
{
"config": [
	       {
			   "type": "var",
				"identifier": "element",
				"childs": [
				          {
							 "type": "*",
							 "identifier": "play",
							 "childs": [
								        {
								        "type": "[]",
										  "childs": [
										            {
										            "type": "{}",
														"agg": "1",
														"childs": [
														          {
																	   "type": "*",
																		"identifier": "videoid",
																		"childs": [
															                   {
																	             "type": "literal",
																		          "extractors": [{"name": "literal", "key": "id"}]
																	             }
																					 ]
																	 },
														          {
																	   "type": "*",
																		"identifier": "views",
																		"childs": [
															                   {
																	             "type": "literal",
																		          "extractors": [{"name": "literal", "key": "views"}]
																	             }
																					 ]
																	 }
																	 ]
														}
													   ]
										  }
								        ]
							 }
				          ]
			 } 
          ]
}
`
	const scrapconfig2 = `
{
"config": [
	       {
			   "type": "var",
				"identifier": "element",
				"childs": [
				          {
							 "type": "*",
							 "identifier": "play",
							 "childs": [
								        {
								        "type": "[]",
										  "childs": [
										            {
										            "type": "{}->map",
														"config": {"videoid": "1", "views": "1"}
														}
													   ]
										  }
								        ]
							 }
				          ]
			 } 
          ]
}
`
	if scrapconfig != "" {
	}

	scrapers, err := js.JSONToJsScraper(strings.NewReader(scrapconfig2))
	if err != nil {
		fmt.Println(err, scrapers)
	}

	m, err := js.ScrapJS(scrapers, strings.NewReader(jscode))
	if err != nil {
		fmt.Println(err)
	}

	for i, res := range m {
		fmt.Printf("#%v videoid '%s' views %v\n", i, res["videoid"], res["views"])
	}
	//Output:
	//#0 videoid 'Yfhuyr-T' views 10
	//#1 videoid 'dhy-Ggdz' views 356980

}
