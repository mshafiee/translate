
Translation Tool
================

This is a command-line tool for translating text in a text file from one language to another language using Google Translate API. The tool translates the text in each row of the CSV file and writes the original and translated text to an intermediate CSV file. It then sorts the intermediate CSV file based on the first column (line number), extracts the translated text column into a text file, and converts the sorted CSV file to a PO file.

Requirements
------------

*   Go 1.16 or later

Installation
------------

1.  Clone the repository: `git clone https://github.com/mshafiee/translate.git`
2.  `cd translate`
3.  Build the binary: `go build -o translate cmd/main.go`

Usage
-----

The tool can be used as follows:

bash

```bash
./translate -input <input-file> -from <from-language-code> -to <to-language-code> -output <output-folder>
```

where:

*   `<input-file>` is the path to the input file for translation
*   `<from-language-code>` is the language code to translate from (ISO 639-1)
*   `<to-language-code>` is the language code to translate to (ISO 639-1)
*   `<output-folder>` is the folder to store the translated files

Example:

bash

```bash
./translate -input input.csv -from en -to fa -output output
```

The tool will generate the following files in the `<output-folder>`:

*   `<input-file>-intermed.csv`: the intermediate CSV file containing the original and translated text
*   `<input-file>-sorted.csv`: the sorted intermediate CSV file
*   `<input-file>-<to-language-code>.txt`: the translated text file
*   `<input-file>.po`: the PO file containing the translated text
