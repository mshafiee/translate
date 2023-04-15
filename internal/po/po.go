package po

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"text/template"
)

type PoEntry struct {
	MsgId      string
	MsgStr     string
	MsgCtxt    string
	MsgPlurals []string
	Comment    string
}

func CSVtoPo(inputFile string, outputFile string) error {
	// Open the CSV file
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the CSV file into a slice of PoEntry structs
	reader := csv.NewReader(file)
	entries, err := parseCsv(reader)
	if err != nil {
		return err
	}

	// Write the PoEntry structs to the .po file
	return writePo(entries, outputFile)
}

func parseCsv(reader *csv.Reader) ([]PoEntry, error) {
	// Read the CSV file line by line and convert each line to a PoEntry struct
	var entries []PoEntry
	lineNum := 0
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		lineNum++
		if len(record) < 2 {
			continue
		}

		var entry PoEntry

		for i, r := range record {
			switch i {
			case 0:
				lineNo, _ := strconv.Atoi(r)
				entry.MsgCtxt = fmt.Sprintf("%08d", lineNo)
				break
			case 1:
				entry.MsgId = escape(r)
				break
			case 2:
				entry.MsgStr = escape(r)
				break
			default:
				if len(r) > 0 {
					entry.Comment += fmt.Sprintf("# %s\n", r)
					if len(r) > 120 {
						entry.Comment += fmt.Sprintf("#‌\n")
					}
				}
			}
		}

		entries = append(entries, entry)
	}
	return entries, nil
}

func writePo(entries []PoEntry, outputFile string) error {
	// Create the output file
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Define the template for the .po file
	templateString := `msgid ""
msgstr ""
"Content-Type: text/plain; charset=UTF-8\n"
"Content-Transfer-Encoding: 8bit\n"

{{range .}}
{{.Comment}}
#, fuzzy
{{if .MsgCtxt}}msgctxt "{{.MsgCtxt}}"{{end}}
msgid "{{.MsgId}}"
msgstr "{{.MsgStr}}"
{{range $index, $plural := .MsgPlurals}}msgid_plural "{{$plural}}"
msgstr[{{$index}}] ""
{{end}}
{{end}}`
	poTemplate, err := template.New("poTemplate").Parse(templateString)
	if err != nil {
		return err
	}

	// Write the PoEntry structs to the output file using the .po template
	err = poTemplate.Execute(file, entries)
	if err != nil {
		return err
	}

	return nil
}

func escape(s string) string {
	var buf bytes.Buffer
	for _, c := range s {
		switch c {
		case '"':
			buf.WriteString("'")
		default:
			buf.WriteRune(c)
		}
	}
	return buf.String()
}
