package qgomysql

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func parsePath(path string) []string {
	out := []string{}
	arr := strings.Split(path, "/")
	for _, value := range arr {
		if value != "" {
			out = append(out, value)
		}
	}
	return out
}

func buildQueryPart(params []string, delimiter string) string {
	vars := []string{}

	for _, param := range params {
		if param != "" {
			vars = append(vars, strings.Trim(param, " "))
		}
	}

	if len(vars)%2 == 1 {
		vars = vars[0 : len(vars)-1]
	}

	query_parts := []string{}
	query_part := ""

	for index, query_var := range vars {

		remain := index % 2

		if remain == 0 {
			query_part += query_var + " = "
		} else {
			query_part += escapeValue(query_var)
			query_parts = append(query_parts, query_part)
			query_part = ""
		}
	}

	return strings.Join(query_parts, delimiter)
}

func buildLimit(params []string) string {
	var out string
	if len(params) == 0 {
		return out
	}
	return " limit " + params[0]
}

func (config *Config) buildQuery(path string) {
	request_parts := strings.Split(path, config.SERVER.Delimiter)
	query_array := parsePath(request_parts[0])

	if len(query_array) < 4 {
		config.raw = append(config.raw, map_si{
			"Build error": "not enought arguments",
		})
		return
	}

	method := query_array[0]
	table := query_array[1]

	var second_params []string
	if len(request_parts) == 2 {
		second_params = parsePath(request_parts[1])
	}

	var query string

	if method == "select" {

		query = method + " * from " + table + " where " + buildQueryPart(query_array[2:], " and ") + buildLimit(second_params)

	} else if method == "delete" {

		query = method + " from " + table + " where " + buildQueryPart(query_array[2:], " and ") + buildLimit(second_params)

	} else if method == "insert" {

		query = method + " into " + table + " set " + buildQueryPart(query_array[2:], ", ")

	} else if method == "update" {

		if len(second_params) < 2 {
			config.raw = append(config.raw, map_si{
				"Build error": "not enought arguments for method update",
			})
			return
		}

		query = method + " " + table + " set " + buildQueryPart(second_params, ", ") + " where " + buildQueryPart(query_array[2:], " and ")

	} else {
		config.raw = append(config.raw, map_si{
			"Build error": "method not supported",
		})
		return
	}

	config.query = query
}

func (config *Config) RunServer() {

	if config.SERVER.Delimiter == "" {
		config.SERVER.Delimiter = "@@@"
	}

	if config.SERVER.Port == 0 {
		config.SERVER.Port = 5555
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {

		config.buildQuery(request.URL.Path)
		config.getData()

		writer.Header().Set("Content-Type", "application/json")
		io.WriteString(writer, config.result)
	})

	err := http.ListenAndServe(":"+fmt.Sprintf("%v", config.SERVER.Port), nil)
	if err != nil {
		panic(err)
	}
}
