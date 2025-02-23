package main

import (
	"bytes"
)

/*type TestServerHttpParsingResult struct {
	is_ok bool
}*/

func is_data_ok(data []byte) bool {

	method := data[0:3]
	if string(method) != "GET" {
		return false
	}

	//Skip space, so it's 4, not 3

	path_end_index := bytes.IndexByte(data[4:], byte(' '))

	path := data[4 : 4+path_end_index]
	if string(path) != "/test_plain" {
		return false
	}

	return true
}
