package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"github.com/mshafiee/translate/internal/gtranslate"
	"github.com/mshafiee/translate/internal/po"
	"github.com/mshafiee/translate/internal/utils"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const MAX_CONCURRENCY = 10

func main() {
	var (
		inputFilePath string
		translateFrom string
		translateTo   string
		outputFolder  string
	)
	// Define flags for command-line arguments
	flag.StringVar(&inputFilePath, "input", "", "Path to the input file for translation")
	flag.StringVar(&translateFrom, "from", "", "Language code to translate from (ISO 639-1) e.g: en")
	flag.StringVar(&translateTo, "to", "", "Language code to translate to (ISO 639-1) e.g: fa")
	flag.StringVar(&outputFolder, "output", "", "Folder to store translated files")
	flag.Parse()

	// Validate input parameters
	if inputFilePath == "" {
		exitWithError(errors.New("missing required input file path"))
	}
	if translateFrom == "" {
		exitWithError(errors.New("missing required 'from' language code"))
	}
	if translateTo == "" {
		exitWithError(errors.New("missing required 'to' language code"))
	}
	if outputFolder == "" {
		exitWithError(errors.New("missing required output folder path"))
	}

	inputFileNameWithoutExt := filepath.Base(inputFilePath[:len(inputFilePath)-len(filepath.Ext(inputFilePath))])

	// Open the input file.
	file, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	err = os.MkdirAll(outputFolder, os.ModePerm)
	if err != nil {
		log.Println(err)
	}

	intermediateFileName := fmt.Sprintf("%s/%s", outputFolder, fmt.Sprintf("%s-intermed.csv", inputFileNameWithoutExt))

	// Create the output file.
	intermediateFile, err := os.Create(intermediateFileName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer intermediateFile.Close()

	// Create a CSV writer.
	writer := csv.NewWriter(intermediateFile)

	// Create a scanner to read the file line by line.
	scanner := bufio.NewScanner(file)

	// Create a WaitGroup to wait for all goroutines to finish.
	var wg sync.WaitGroup

	// Create a channel to limit the number of concurrent goroutines.
	concurrency := make(chan struct{}, MAX_CONCURRENCY)

	lineNumber := 0

	// Process each line in a separate goroutine.
	for scanner.Scan() {
		lineNumber++

		// Acquire a slot in the concurrency channel.
		concurrency <- struct{}{}

		// Increment the WaitGroup lineNumber.
		wg.Add(1)

		go consumer(concurrency, &wg, lineNumber, scanner.Text(), translateFrom, translateTo, writer)

	}

	// Wait for all goroutines to finish.
	wg.Wait()

	// Flush any remaining data to the CSV file.
	writer.Flush()

	intermediateFile.Close()
	sortedFileName := fmt.Sprintf("%s/%s-sorted.csv", outputFolder, inputFileNameWithoutExt)
	translatedTextFileName := fmt.Sprintf("%s/%s-%s.txt", outputFolder, inputFileNameWithoutExt, translateTo)
	poFileName := fmt.Sprintf("%s/%s.po", outputFolder, inputFileNameWithoutExt)
	postProccess(intermediateFileName, sortedFileName, translatedTextFileName, poFileName, 3)
}

var mutex = &sync.Mutex{}

// Consumer function that consumes elements of the buffer and writes them to a CSV file.
func consumer(concurrency chan struct{}, wg *sync.WaitGroup, rowID int, originalText, translateFrom, translateTo string, writer *csv.Writer) {
	// Release the slot in the concurrency channel when done.
	defer func() { <-concurrency }()
	defer wg.Done()

	// Create a slice to store the paragraphs.
	var paragraphs [][]string

	if len(strings.TrimSpace(originalText)) > 0 {
		googleTranslated, err := gtranslate.TranslateWithParams(
			originalText,
			gtranslate.TranslationParams{
				From: "en",
				To:   "fa",
			},
		)
		if err != nil {
			panic(err)
		}
		log.Println(rowID)

		row := []string{
			strconv.Itoa(rowID),
			originalText,
		}
		//row = append(row, googleTranslated...)
		row = append(row, googleTranslated[0])

		sentence, err := gtranslate.SentenceWithParams(
			originalText,
			gtranslate.TranslationParams{
				From: translateFrom,
				To:   translateTo,
			},
		)
		row = append(row, sentence...)

		//vocabulary, err := gtranslate.VocabularyWithParams(
		//	originalText,
		//	gtranslate.TranslationParams{
		//		From: "en",
		//		To:   "fa",
		//	},
		//)
		//row = append(row, vocabulary...)

		// Add the pair to the slice.
		paragraphs = append(paragraphs, row)
	}

	// Write the paragraphs to the CSV file.
	if len(paragraphs) > 0 {
		mutex.Lock()
		defer mutex.Unlock()

		for _, pair := range paragraphs {
			// Insert the paragraph number into the first column.
			var record []string
			for _, s := range pair {
				record = append(record, s)
			}
			if err := writer.Write(record); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}
		writer.Flush()
	}
}

func postProccess(inputFileName, sortedFileName, translatedTextFileName, poFileName string, columnNumber int) {
	normalizedCommasFileName := "normalizedCommasFileName.csv"
	err := utils.AddCommasToFile(inputFileName, normalizedCommasFileName)
	if err != nil {
		exitWithError(fmt.Errorf("Error AddCommasToCSV: %v\n", err))
	}

	err = utils.NumericalSortCSV(normalizedCommasFileName, sortedFileName)
	if err != nil {
		exitWithError(fmt.Errorf("Error NumericalSortCSV: %v\n", err))
	}
	err = utils.ExtractColumnWithEmptyRows(sortedFileName, translatedTextFileName, columnNumber)
	if err != nil {
		exitWithError(fmt.Errorf("Error ExtractColumnWithEmptyRows: %v\n", err))
	}
	err = po.CSVtoPo(sortedFileName, poFileName)
	if err != nil {
		exitWithError(fmt.Errorf("Error converting CSV to PO: %v\n", err))
	}
}

func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
