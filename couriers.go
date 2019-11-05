package my_tracker

import (
	"github.com/PuerkitoBio/goquery"
	"gopkg.in/headzoo/surf.v1"
	"regexp"
	"strings"
)

func PosLaju(serial string) ([]Log, error) {
	bow := surf.NewBrowser()
	err := bow.Open("http://track.pos.com.my/postal-services/quick-access?track-trace")
	if err != nil {
		return nil, err
	}
	fm, errForm := bow.Form("form#tracking03-form")
	if err == nil {
		fm.Input("trackingNo03", serial)
		if fm.Submit() != nil {
			return nil, errForm
		}
		str, _ := goquery.OuterHtml(bow.Dom())
		re := regexp.MustCompile(`<table id='tbDetails' class='table table-striped table-hover table-bordered'>([^"]+)</table>`)
		match := re.FindStringSubmatch(str)
		var logs []Log
		data := `<html><body>` + match[0] + `</body></html>`
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
		if err != nil {
			return nil, err
		}
		doc.Find("table").Each(func(index int, tablehtml *goquery.Selection) {
			tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
				l := Log{}
				rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
					if indexth == 0 {
						l.Date = tablecell.Text()
					} else if indexth == 1 {
						l.Process = tablecell.Text()
					} else if indexth == 2 {
						l.Event = tablecell.Text()
					}
				})
				if l.Date != "" {
					logs = append(logs, l)
				}
			})
		})
		return logs, nil
	} else {
		return nil, errForm
	}
}

func JnT(serial string) ([]Log, error) {
	var logs []Log
	bow := surf.NewBrowser()
	err := bow.Open("https://www.jtexpress.my/track")
	if err != nil {
		return nil, err
	}
	fm, errForm := bow.Form("form#track-package-form")
	if errForm == nil {
		fm.Input("billcode", serial)
		if fm.Submit() != nil {
			return nil, errForm
		}
		bow.Dom().Find("div.entry").Each(func(_ int, s *goquery.Selection) {
			l := Log{}
			for idx, dom := range s.Children().Nodes {
				if idx == 0 {
					l.Date = s.Find("div."+dom.Attr[0].Val).Find("h3").Text() + ", " + s.Find("div."+dom.Attr[0].Val).Find("p").Text()
				} else if idx == 1 {
					var events []string
					var result string
					s.Find("div."+dom.Attr[0].Val).Children().Each(func(i int, s2 *goquery.Selection) {
						if i <= 1 {
							events = append(events, s2.Text())
						} else {
							result = s2.Text()
						}
					})
					l.Event = strings.Join(events, ", ")
					l.Process = result
				}
			}
			logs = append(logs, l)
		})
		return logs, nil
	} else {
		return nil, errForm
	}
}
