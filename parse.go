package qgomysql

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func escapeValue(value interface{}) string {
	_, is_string := value.(string)

	if is_string {
		value = html.EscapeString(fmt.Sprintf("%v", value))
		value = fmt.Sprintf("%s%s%s", "\"", value, "\"")
	}

	return fmt.Sprintf("%v", value)
}

func (config *Config) tableToMap() {

	titles := []string{}
	values := []interface{}{}

	tokenizer := html.NewTokenizer(strings.NewReader(config.result))

	for {
		tag := tokenizer.Next()
		if tag == html.ErrorToken {
			break
		}

		if tag == html.StartTagToken {

			token := tokenizer.Token()
			tag_name := token.Data

			if tag_name == "th" || tag_name == "td" {

				inner := tokenizer.Next()

				if inner == html.TextToken || inner == html.EndTagToken {

					value := strings.TrimSpace((string)(tokenizer.Text()))

					if tag_name == "th" {
						titles = append(titles, value)
					} else if tag_name == "td" {
						values = append(values, value)
					}

				}
			}
		}
	}

	titles_count := len(titles) - 1
	count := 0

	part := map_si{}

	for _, value := range values {

		part[titles[count]] = value

		if count == titles_count {

			config.raw = append(config.raw, part)
			part = map_si{}
			count = 0

		} else {
			count++
		}

	}

	config.result = ""
}
