package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"

	"golang.design/x/clipboard"
)

type JSONRecord = map[string]interface{}
type Pipeline = []JSONRecord

func adjustValue(value interface{}) interface{} {
	adjustedValue := value
	if value != nil && reflect.TypeOf(value).Kind() == reflect.String {
		if regexp.MustCompile(`^[0-9a-fA-F]{24}$`).MatchString(value.(string)) {
			adjustedValue = "ObjectId(\"" + value.(string) + "\")"
		} else if regexp.MustCompile(`.*T[0-9][0-9]:[0-9][0-9]:[0-9][0-9].*Z$`).MatchString(value.(string)) {
			adjustedValue = "ISODate(\"" + value.(string) + "\")"
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
    result = regexp.MustCompile(`"`+ r + `\(\\`).ReplaceAllString(result, r + "(")
  }
  result = regexp.MustCompile(`\\"\)"`).ReplaceAllString(result, "\")")
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
    fmt.Println("Error: Invalid JSON in clipboard. Please copy a valid JSON representation of a MongoDB pipeline.")
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
