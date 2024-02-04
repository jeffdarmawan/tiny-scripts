package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"reflect"
)

// Concats
// - open all csv files located in './join' directory
// - create a map of first column

func main() {
	// 1. discover all csv files
	entries, err := os.ReadDir("./csv/join/")
	if err != nil {
		log.Fatal(err)
	}

	var (
		allFileContents = make(map[string][]string)
		originalOrder   []string
	)

	// 2. iterate all files
	//    - open file
	//    - convert into map string of string
	//    - append the value to the main map string of strings
	for i, e := range entries {
		// // ignore output file if it is contained in the f
		// if e.Name() == "output.csv" {

		// }

		// convert file into map of values
		fileContent, fileOrder := openFile(e.Name())

		// check if the original order of the file is the same
		if i > 0 && !reflect.DeepEqual(originalOrder, fileOrder) {
			for i := range originalOrder {
				if originalOrder[i] != fileOrder[i] {
					log.Println("not match: ", originalOrder[i], fileOrder[i]) // log debug
				}
			}
			log.Fatalf("contents between files does not match. filename: %s", e.Name())
		}

		originalOrder = fileOrder

		// append value according to the key
		allFileContentsNew, err := appendValues(allFileContents, fileContent)
		if err != nil {
			log.Fatalf("error when appending values. err: %s. filename: %s", err.Error(), e.Name())
		}

		// TODO: not sure if allFileContents is already updated inside the function or not
		allFileContents = allFileContentsNew
	}

	// 3. output

	// create a file
	file, err := os.Create("result.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// initialize csv writer
	writer := csv.NewWriter(file)

	defer writer.Flush()

	// write all rows at once
	writer.WriteAll(makeItWritable(allFileContents, originalOrder))
}

func openFile(filename string) (map[string]string, []string) {
	// open file
	f, err := os.Open(fmt.Sprintf("./csv/join/%s", filename))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var (
		allRowData    = make(map[string]string)
		originalOrder = make([]string, 0, len(data))
	)
	for i, d := range data {
		if i == 0 {
			allRowData[d[0]] = filename[:len(filename)-4]
			originalOrder = append(originalOrder, d[0])
			continue
		}
		allRowData[d[0]] = d[5]
		originalOrder = append(originalOrder, d[0])
	}

	return allRowData, originalOrder
}

func appendValues(original map[string][]string, addition map[string]string) (map[string][]string, error) {
	for k, val := range addition {
		original[k] = append(original[k], val)
	}

	return original, nil
}

func makeItWritable(allFileContents map[string][]string, originalOrder []string) [][]string {
	output := make([][]string, 0, len(originalOrder))
	for _, s := range originalOrder {
		output = append(output, append([]string{s}, allFileContents[s]...))
	}

	return output
}
