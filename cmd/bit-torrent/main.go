package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func decodeBencode(bencodedString string, index *int) (interface{}, error) {

	switch bencodedString[*index] {
	case 'i':
		return decodeInt(bencodedString, index)
	case 'l':
		return decodeList(bencodedString, index)
	default:
		return decodeString(bencodedString, index)
	}

}

func decodeList(bencodedString string, index *int) (interface{}, error) {
	*index++

	list := make([]interface{}, 0)

	for bencodedString[*index] != 'e' {
		value, err := decodeBencode(bencodedString, index)
		if err != nil {
			return nil, err
		}
		list = append(list, value)
	}

	if bencodedString[*index] != 'e' {
		return nil, fmt.Errorf("current list in not properly closed")
	}

	*index++

	return list, nil
}

func decodeString(bencodedString string, index *int) (interface{}, error) {

	var current_index int = *index
	for bencodedString[current_index] != ':' {
		current_index++

	}

	lengthStr := bencodedString[*index:current_index]

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", err
	}

	start_index := current_index + 1
	end_index := start_index + length

	value := bencodedString[start_index:end_index]
	*index = end_index
	return value, nil
}

func decodeInt(bencodedString string, index *int) (interface{}, error) {
	*index++

	value := ""

	for bencodedString[*index] != 'e' {
		value += string(bencodedString[*index])
		*index++
	}

	*index++

	if value == "-0" {
		return nil, fmt.Errorf("negative zero is not allowed")
	}

	if len(value) > 1 && value[0] == '0' {
		return nil, fmt.Errorf("leading value cannot be zero")
	}

	return strconv.Atoi(value)
}

func main() {

	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]
		index := 0
		decoded, err := decodeBencode(bencodedValue, &index)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
