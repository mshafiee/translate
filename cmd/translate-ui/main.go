package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mshafiee/progressbar"
	"github.com/mshafiee/translate/internal/gtranslate"
	"github.com/mshafiee/translate/internal/po"
	"github.com/mshafiee/translate/internal/utils"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const extendedTranslation = true

type multilineWriter struct {
	*widget.Entry
}

func (w *multilineWriter) Write(p []byte) (n int, err error) {
	text := string(p)
	w.SetText(w.Text + text)
	return len(p), nil
}

func main() {
	var languageNames []string

	for _, lang := range languages {
		languageNames = append(languageNames, lang[0])
	}

	a := app.New()
	w := a.NewWindow("Translate")
	w.Resize(fyne.NewSize(800, 300))

	fromCombo := widget.NewSelect(languageNames, func(s string) {})
	fromCombo.SetSelected("English")

	inputEntry := widget.NewEntry()
	inputEntry.SetPlaceHolder("Input file path")

	inputButton := widget.NewButton("Choose", func() {
		dialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				path := reader.URI().Path()
				inputEntry.SetText(path)
			}
		}, w)
		dialog.Show()
	})

	outputEntry := widget.NewEntry()
	outputEntry.SetPlaceHolder("Output folder")

	outputButton := widget.NewButton("Choose", func() {
		dialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				path := uri.Path()
				outputEntry.SetText(path)
			}
		}, w)
		dialog.Show()
	})

	toCombo := widget.NewSelect(languageNames, func(s string) {})
	toCombo.SetSelected("Persian")

	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	outputMultiLineEntry := widget.NewMultiLineEntry()
	outputMultiLineEntryWriter := &multilineWriter{outputMultiLineEntry}

	translateButton := widget.NewButton("Translate", func() {
		var from, to string
		for _, l := range languages {
			if l[0] == fromCombo.Selected {
				from = l[1]
			}
			if l[0] == toCombo.Selected {
				to = l[1]
			}
		}
		input := inputEntry.Text
		output := outputEntry.Text

		go func() {
			translate(outputMultiLineEntryWriter, progressBar, from, input, output, to, extendedTranslation)
			progressBar.Hide()
		}()
		progressBar.Show()
	})
	translateButton.Disable()

	inputEntry.OnChanged = validate(outputMultiLineEntryWriter, fromCombo, inputEntry, outputEntry, toCombo, translateButton)
	outputEntry.OnChanged = validate(outputMultiLineEntryWriter, fromCombo, inputEntry, outputEntry, toCombo, translateButton)

	creditLabel := widget.NewLabel("Translate, Developed by Mohammad Shafiee, muhammad.shafiee@gmail.com")
	creditLabel.Alignment = fyne.TextAlignCenter

	formContainer := container.New(
		layout.NewFormLayout(),
		widget.NewLabel("From:"),
		fromCombo,
		widget.NewLabel("Input file:"),
		container.New(layout.NewBorderLayout(nil, nil, nil, inputButton), inputEntry, inputButton),
		widget.NewLabel("Output folder:"),
		container.New(layout.NewBorderLayout(nil, nil, nil, outputButton), outputEntry, outputButton),
		widget.NewLabel("To:"),
		toCombo,
		layout.NewSpacer(),
		translateButton,
		layout.NewSpacer(),
		progressBar,
	)

	mainContainer := container.New(
		layout.NewGridLayout(1),
		formContainer,
		outputMultiLineEntry,
	)

	w.SetContent(mainContainer)
	w.ShowAndRun()
}

func validate(w io.Writer, fromCombo *widget.Select, inputEntry *widget.Entry, outputEntry *widget.Entry, toCombo *widget.Select, translateButton *widget.Button) func(_ string) {
	return func(_ string) {
		from := fromCombo.Selected
		input := inputEntry.Text
		output := outputEntry.Text
		to := toCombo.Selected

		// Check if any field is empty
		if from == "" || input == "" || output == "" || to == "" {
			translateButton.Disable()
			fmt.Fprintln(os.Stdout, "Please fill all fields")
			widget.NewLabel("Please fill all fields").Show()
			return
		}
		translateButton.Enable()
	}
}

const MAX_CONCURRENCY = 10

func translate(w io.Writer, progressBarUI *widget.ProgressBar, translateFrom string, inputFilePath string, outputFolder string, translateTo string, extendedTranslation bool) {
	fmt.Fprintf(w, "Source: %s\n", inputFilePath)
	fmt.Fprintf(w, "Translation files path: %s\n", outputFolder)
	fmt.Fprintf(w, "Start Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	// Validate input parameters
	if inputFilePath == "" {
		exitWithError(w, errors.New("missing required input file path"))
	}
	if translateFrom == "" {
		exitWithError(w, errors.New("missing required 'from' language code"))
	}
	if translateTo == "" {
		exitWithError(w, errors.New("missing required 'to' language code"))
	}
	if outputFolder == "" {
		exitWithError(w, errors.New("missing required output folder path"))
	}
	inputFileNameWithoutExt := filepath.Base(inputFilePath[:len(inputFilePath)-len(filepath.Ext(inputFilePath))])

	outputFolder = fmt.Sprintf("%s/%s", outputFolder, inputFileNameWithoutExt)
	totalLineNumber, err := utils.CountLines(inputFilePath)
	if err != nil {
		exitWithError(w, err)
	}

	// Open the input file.
	file, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Fprintln(w, "Error:", err)
		return
	}
	defer file.Close()

	err = os.MkdirAll(outputFolder, os.ModePerm)
	if err != nil {
		fmt.Fprintln(w, "Error:", err)
		return
	}

	intermediateFileName := fmt.Sprintf("%s/%s", outputFolder, fmt.Sprintf("%s-intermed.csv", inputFileNameWithoutExt))

	// Create the output file.
	intermediateFile, err := os.Create(intermediateFileName)
	if err != nil {
		fmt.Fprintln(w, "Error:", err)
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

		go consumer(w, concurrency, &wg, progressBarUI, totalLineNumber, lineNumber, scanner.Text(), translateFrom, translateTo, extendedTranslation, writer)

	}

	// Wait for all goroutines to finish.
	wg.Wait()

	// Flush any remaining data to the CSV file.
	writer.Flush()

	progressbar.ColorArrowProgressBar(100, 100)
	intermediateFile.Close()
	normalizedCommasFileName := fmt.Sprintf("%s/%s-normalized.csv", outputFolder, inputFileNameWithoutExt)
	sortedFileName := fmt.Sprintf("%s/%s-sorted.csv", outputFolder, inputFileNameWithoutExt)
	translatedTextFileName := fmt.Sprintf("%s/%s-%s.txt", outputFolder, inputFileNameWithoutExt, translateTo)
	poFileName := fmt.Sprintf("%s/%s.po", outputFolder, inputFileNameWithoutExt)
	postProccess(w, intermediateFileName, normalizedCommasFileName, sortedFileName, translatedTextFileName, poFileName, 3)
}

var mutex = &sync.Mutex{}

// Consumer function that consumes elements of the buffer and writes them to a CSV file.
func consumer(w io.Writer, concurrency chan struct{}, wg *sync.WaitGroup, progressBarUI *widget.ProgressBar, totalRows, rowID int, originalText, translateFrom, translateTo string, extendedTranslation bool, writer *csv.Writer) {
	// Release the slot in the concurrency channel when done.
	defer func() { <-concurrency }()
	defer wg.Done()

	// Create a slice to store the paragraphs.
	var paragraphs [][]string

	if len(strings.TrimSpace(originalText)) > 0 {
		googleTranslated, err := gtranslate.TranslateWithParams(
			originalText,
			gtranslate.TranslationParams{
				From: translateFrom,
				To:   translateTo,
			},
		)
		if err != nil {
			fmt.Fprintln(w, err)
			//panic(err)
		}
		progressBarUI.SetValue(float64(rowID) / float64(totalRows))
		progressbar.ColorArrowProgressBar(rowID, totalRows)

		row := []string{
			strconv.Itoa(rowID),
			originalText,
		}
		//row = append(row, googleTranslated...)
		row = append(row, googleTranslated[0])

		if extendedTranslation {
			sentence, err := gtranslate.SentenceWithParams(
				originalText,
				gtranslate.TranslationParams{
					From: translateFrom,
					To:   translateTo,
				},
			)
			if err != nil {

				fmt.Fprintln(w, err)
			}
			row = append(row, sentence...)
		}
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

func postProccess(w io.Writer, inputFileName, normalizedCommasFileName, sortedFileName, translatedTextFileName, poFileName string, columnNumber int) {
	err := utils.AddCommasToFile(inputFileName, normalizedCommasFileName)
	if err != nil {
		exitWithError(w, fmt.Errorf("Error AddCommasToCSV: %v\n", err))
	}

	err = utils.NumericalSortCSV(normalizedCommasFileName, sortedFileName)
	if err != nil {
		exitWithError(w, fmt.Errorf("Error NumericalSortCSV: %v\n", err))
	}
	err = utils.ExtractColumnWithEmptyRows(sortedFileName, translatedTextFileName, columnNumber)
	if err != nil {
		exitWithError(w, fmt.Errorf("Error ExtractColumnWithEmptyRows: %v\n", err))
	}
	err = po.CSVtoPo(sortedFileName, poFileName)
	if err != nil {
		exitWithError(w, fmt.Errorf("Error converting CSV to PO: %v\n", err))
	}
}

func exitWithError(w io.Writer, err error) {
	fmt.Fprintln(w, err)
}

var languages = [][]string{
	{"Automatic Detection", "auto"},
	{"Afrikaans", "af"},
	{"Albanian", "sq"},
	{"Amharic", "am"},
	{"Arabic", "ar"},
	{"Armenian", "hy"},
	{"Assamese", "as"},
	{"Aymara", "ay"},
	{"Azerbaijani", "az"},
	{"Bambara", "bm"},
	{"Basque", "eu"},
	{"Belarusian", "be"},
	{"Bengali", "bn"},
	{"Bhojpuri", "bho"},
	{"Bosnian", "bs"},
	{"Bulgarian", "bg"},
	{"Catalan", "ca"},
	{"Cebuano", "ceb"},
	{"Chichewa", "ny"},
	{"Chinese (Simplified)", "zh-CN"},
	{"Chinese (Traditional)", "zh-TW"},
	{"Corsican", "co"},
	{"Croatian", "hr"},
	{"Czech", "cs"},
	{"Danish", "da"},
	{"Dhivehi", "dv"},
	{"Dogri", "doi"},
	{"Dutch", "nl"},
	{"English", "en"},
	{"Esperanto", "eo"},
	{"Estonian", "et"},
	{"Ewe", "ee"},
	{"Filipino", "fil"},
	{"Finnish", "fi"},
	{"French", "fr"},
	{"Frisian", "fy"},
	{"Galician", "gl"},
	{"Georgian", "ka"},
	{"German", "de"},
	{"Greek", "el"},
	{"Guarani", "gn"},
	{"Gujarati", "gu"},
	{"Haitian Creole", "ht"},
	{"Hausa", "ha"},
	{"Hawaiian", "haw"},
	{"Hebrew", "he"},
	{"Hindi", "hi"},
	{"Hmong", "hmn"},
	{"Hungarian", "hu"},
	{"Icelandic", "is"},
	{"Igbo", "ig"},
	{"Ilocano", "ilo"},
	{"Indonesian", "id"},
	{"Irish", "ga"},
	{"Italian", "it"},
	{"Japanese", "ja"},
	{"Javanese", "jv"},
	{"Kannada", "kn"},
	{"Kazakh", "kk"},
	{"Khmer", "km"},
	{"Kinyarwanda", "rw"},
	{"Konkani", "kok"},
	{"Korean", "ko"},
	{"Krio", "kri"},
	{"Kurdish (Kurmanji)", "kmr"},
	{"Kurdish (Sorani)", "ckb"},
	{"Kyrgyz", "ky"},
	{"Lao", "lo"},
	{"Latin", "la"},
	{"Latvian", "lv"},
	{"Lingala", "ln"},
	{"Lithuanian", "lt"},
	{"Luganda", "lg"},
	{"Luxembourgish", "lb"},
	{"Macedonian", "mk"},
	{"Maithili", "mai"},
	{"Malagasy", "mg"},
	{"Malay", "ms"},
	{"Malayalam", "ml"},
	{"Maltese", "mt"},
	{"Maori", "mi"},
	{"Marathi", "mr"},
	{"Meiteilon (Manipuri)", "mni"},
	{"Mizo", "lus"},
	{"Mongolian", "mn"},
	{"Myanmar (Burmese)", "my"},
	{"Nepali", "ne"},
	{"Norwegian", "no"},
	{"Odia (Oriya)", "or"},
	{"Oromo", "om"},
	{"Pashto", "ps"},
	{"Persian", "fa"},
	{"Polish", "pl"},
	{"Portuguese", "pt"},
	{"Punjabi", "pa"},
	{"Quechua", "qu"},
	{"Romanian", "ro"},
	{"Russian", "ru"},
	{"Samoan", "sm"},
	{"Sanskrit", "sa"},
	{"Scots Gaelic", "gd"},
	{"Sepedi", "nso"},
	{"Serbian", "sr"},
	{"Sesotho", "st"},
	{"Shona", "sn"},
	{"Sindhi", "sd"},
	{"Sinhala", "si"},
	{"Slovak", "sk"},
	{"Slovenian", "sl"},
	{"Somali", "so"},
	{"Spanish", "es"},
	{"Sundanese", "su"},
	{"Swahili", "sw"},
	{"Swedish", "sv"},
	{"Tajik", "tg"},
	{"Tamil", "ta"},
	{"Tatar", "tt"},
	{"Telugu", "te"},
	{"Thai", "th"},
	{"Tigrinya", "ti"},
	{"Tsonga", "ts"},
	{"Turkish", "tr"},
	{"Turkmen", "tk"},
	{"Twi", "tw"},
	{"Ukrainian", "uk"},
	{"Urdu", "ur"},
	{"Uyghur", "ug"},
	{"Uzbek", "uz"},
	{"Vietnamese", "vi"},
	{"Welsh", "cy"},
	{"Xhosa", "xh"},
	{"Yiddish", "yi"},
	{"Yoruba", "yo"},
	{"Zulu", "zu"},
}
