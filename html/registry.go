package html

import (
	"fmt"

	"github.com/jloup/scraper/html/extractor"
	"github.com/jloup/scraper/html/validator"
	"github.com/jloup/scraper/node"
)

// validator constructor prototype
type ValidatorGenerator func(config map[string]string) (validator.Validator, error)

// extractor constructor prototype
type ExtractorGenerator func(config map[string]string) (extractor.Extractor, error)

var validatorRegistry = map[string]ValidatorGenerator{
	"exists":                  validator.NewExists,
	"attrEquals":              validator.NewAttrEquals,
	"attrContains":            validator.NewAttrContains,
	"regexp":                  validator.NewRegexp,
	"ytI":                     validator.NewYTIframe,
	"ytO":                     validator.NewYTObject,
	"scI":                     validator.NewSCIframe,
	"scO":                     validator.NewSCObject,
	"ytIFast":                 validator.NewYTIframeFast,
	"ytOFast":                 validator.NewYTObjectFast,
	"scIFast":                 validator.NewSCIframeFast,
	"scOFast":                 validator.NewSCObjectFast,
	"facebookVideoIframeFast": validator.NewFacebookVideoIframeFast,
	"facebookVideoIframe":     validator.NewFacebookVideoIframe,
	"dailymotionIframeFast":   validator.NewDailymotionIframeFast,
	"dailymotionIframe":       validator.NewDailymotionIframe,
	"vimeoIframeFast":         validator.NewVimeoIframeFast,
	"vimeoIframe":             validator.NewVimeoIframe,
}

var extractorRegistry = map[string]ExtractorGenerator{
	"extractAttr":          extractor.NewAttribute,
	"extractText":          extractor.NewTextContent,
	"extractTestNoNewLine": extractor.NewTextContentStripNewLine,
	"regexp":               extractor.NewRegexp,
	"setType":              extractor.NewSetType,
	"setKV":                extractor.NewSetKV,
	"sc":                   extractor.NewSoundcloudR,
	"scStreamUrl":          extractor.NewSoundcloudStreamUrl,
	"yt":                   extractor.NewYoutubeId,
	"js":                   extractor.NewJs,
	"facebookVideo":        extractor.NewFacebookVideo,
	"dailymotionVideo":     extractor.NewDailymotionVideo,
	"vimeoVideo":           extractor.NewVimeoVideo,
}

var scraperRegistry = map[string][]*node.ScraperNode{}

type generatorError struct {
	What string
}

func (e generatorError) Error() string {
	return fmt.Sprintf("Registry : %v", e.What)
}

// register a validator generator into registry. The associated validator can be then used in JSON config files
func AddValidatorGenerator(name string, v ValidatorGenerator) error {
	if _, ok := validatorRegistry[name]; ok {
		return generatorError{What: fmt.Sprintf("a validator already has this name '%v'", name)}
	}

	validatorRegistry[name] = v

	return nil

}

// register an extractor generator into registry. The associated extractor can be then used in JSON config files
func AddExtractorGenerator(name string, e ExtractorGenerator) error {
	if _, ok := extractorRegistry[name]; ok {
		return generatorError{What: fmt.Sprintf("an extractor already has this name '%v'", name)}
	}

	extractorRegistry[name] = e

	return nil

}

// register a scraper into registry. Useful for common or complex scraper that needs to be duplicated
func AddScraperConfig(name string, s []*node.ScraperNode) error {
	if _, ok := scraperRegistry[name]; ok {
		return generatorError{What: fmt.Sprintf("a scraper config already has this name '%v'", name)}
	}

	scraperRegistry[name] = s

	return nil
}

func getScraperConfig(name string) ([]*node.ScraperNode, error) {
	if _, ok := scraperRegistry[name]; !ok {
		return nil, generatorError{What: fmt.Sprintf("no scraper config found for '%v'", name)}
	}

	length := len(scraperRegistry[name])
	scrapers := make([]*node.ScraperNode, length, length)

	for i, scraper := range scraperRegistry[name] {
		s := scraper.Copy()
		scrapers[i] = s
	}

	return scrapers, nil
}
