package finder

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func FindLinks(startURL string, skipRegexes []*regexp.Regexp, allowRegexes []*regexp.Regexp) ([]string, error) {
	resp, err := http.Get(startURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)

	var links []string

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return links, nil
		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "a"
			if isAnchor {
				for _, a := range t.Attr {
					if a.Key == "href" {
						// проверка относительной ссылки или нет
						if !strings.HasPrefix(a.Val, "http") {
							u, err := url.Parse(startURL)
							if err != nil {
								return nil, err
							}
							baseURL := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
							a.Val = fmt.Sprintf("%s%s", baseURL, a.Val)
						}
						// Если ссылка равна startURL, пропускаем её
						if a.Val == startURL {
							continue
						}

						// Проверяем, удовлетворяет ли ссылка хотя бы одному из шаблонов allowRegexes
						allow := false
						for _, regex := range allowRegexes {
							if regex.MatchString(a.Val) {
								allow = true
								break
							}
						}

						// Если ссылка удовлетворяет хотя бы одному из шаблонов skipRegexes, пропускаем её
						for _, regex := range skipRegexes {
							if regex.MatchString(a.Val) {
								allow = false
								break
							}
						}

						if allow {
							links = append(links, a.Val)
						}
					}
				}
			}
		}
	}
}
