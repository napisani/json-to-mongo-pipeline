package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"

	"golang.design/x/clipboard"
)

const LEFT_PARENTHESIS_TOKEN = "__LEFT_PARENTHESIS__"
const RIGHT_PARENTHESIS_TOKEN = "__RIGHT_PARENTHESIS__"

func adjustValue(value interface{}) interface{} {
	adjustedValue := value
	if value != nil && reflect.TypeOf(value).Kind() == reflect.String {
		if regexp.MustCompile(`^[0-9a-fA-F]{24}$`).MatchString(value.(string)) {
			adjustedValue = "ObjectId" + LEFT_PARENTHESIS_TOKEN + value.(string) + RIGHT_PARENTHESIS_TOKEN
		} else if regexp.MustCompile(`.*T[0-9][0-9]:[0-9][0-9]:[0-9][0-9].*Z$`).MatchString(value.(string)) {
			adjustedValue = "ISODate" + LEFT_PARENTHESIS_TOKEN + value.(string) + RIGHT_PARENTHESIS_TOKEN
		}
	}
	return adjustedValue
}

func traverseMap(aMap map[string]interface{}) {
	for key, val := range aMap {
		switch val.(type) {
		case map[string]interface{}:
			traverseMap(val.(map[string]interface{}))
		case []interface{}:
			traverseArray(val.([]interface{}))
		default:
			aMap[key] = adjustValue(val)
		}
	}
}

func traverseArray(anArray []interface{}) {
	for i, val := range anArray {
		switch val.(type) {
		case map[string]interface{}:
			traverseMap(val.(map[string]interface{}))
		case []interface{}:
			traverseArray(val.([]interface{}))
		default:
			anArray[i] = adjustValue(val)
		}
	}
}

func doRegexReplacements(jsonText string) string {
	result := jsonText
	for _, r := range []string{"ObjectId", "ISODate"} {
    result = regexp.MustCompile(`"`+r+LEFT_PARENTHESIS_TOKEN).ReplaceAllString(result, r+"(\"")
	}
	result = regexp.MustCompile(RIGHT_PARENTHESIS_TOKEN+"\"").ReplaceAllString(result, "\")")
	return result
}

func main() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	fmt.Println("Reading clipboard... ")
	fmt.Println("")
	text := clipboard.Read(clipboard.FmtText)
	var pipeline interface{}
	err = json.Unmarshal(text, &pipeline)
	if err != nil {

		panic(err)
	}
	fmt.Println("Converting...")
	fmt.Println("")
	traverseArray(pipeline.([]interface{}))
	jsonBytes, err := json.MarshalIndent(pipeline, "", "  ")

	result := string(jsonBytes)
	result = doRegexReplacements(result)
	fmt.Println("Writing clipboard...")
	fmt.Println("")
	clipboard.Write(clipboard.FmtText, []byte(result))
	fmt.Println("----------------------------")
	fmt.Println("")
	fmt.Println(result)
	fmt.Println("")
	fmt.Println("----------------------------")
	fmt.Println("Done!")

}
