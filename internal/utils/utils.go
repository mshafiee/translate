package utils

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"sort"
	"strconv"
	"strings"
)

func CountLines(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var lines int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lines, nil
}

func GetCSVFieldCount(csvString string) (int, error) {
	reader := csv.NewReader(strings.NewReader(csvString))
	records, err := reader.ReadAll()
	if err != nil {
		return 0, err
	}
	return len(records[0]), nil
}

func GetMaxCommas(filePath string) (int, error) {
	maxCommas := 0
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		commaCount, err := GetCSVFieldCount(line)
		if err != nil {
			return 0, err
		}
		if commaCount > maxCommas {
			maxCommas = commaCount
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return maxCommas, nil
}

func AddCommasToFile(inputFilePath, outputFilePath string) error {
	// Get max number of commas in a line
	maxCommas, err := GetMaxCommas(inputFilePath)
	if err != nil {
		return err
	}

	// Open input and output files
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Iterate over each line in input file
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()

		// Add commas if line has less than maxCommas
		commaCount, err := GetCSVFieldCount(line)
		if err != nil {
			return err
		}
		if commaCount < maxCommas {
			diff := maxCommas - commaCount
			line += strings.Repeat(",", diff)
		}

		// Write modified line to output file
		_, err = fmt.Fprintln(outputFile, line)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func NumericalSortCSV(inputFile string, outputFile string) error {
	// Open input file
	f, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// Read CSV data
	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	if err != nil {
		return err
	}

	// Sort lines based on numerical values of first column
	sort.Slice(lines, func(i, j int) bool {
		num1, _ := strconv.ParseFloat(lines[i][0], 64)
		num2, _ := strconv.ParseFloat(lines[j][0], 64)
		return num1 < num2
	})

	// Open output file
	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write sorted lines to output file
	w := csv.NewWriter(out)
	err = w.WriteAll(lines)
	if err != nil {
		return err
	}

	return nil
}

func SortCSVByFirstColumn(filename string) error {
	// Open the CSV file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the CSV data
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Sort the records by the first column
	sort.Slice(records, func(i, j int) bool {
		return records[i][0] < records[j][0]
	})

	// Write the sorted records back to the file
	file, err = os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll(records)
	if err != nil {
		return err
	}

	return nil
}

func ExtractColumn(inputFile string, outputFile string, columnNumber int) error {
	// Open input file
	f, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// Read CSV data
	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	if err != nil {
		return err
	}

	// Extract column values
	var column []string
	for _, line := range lines {
		if len(line) >= columnNumber {
			column = append(column, line[columnNumber-1])
		}
	}

	// Open output file
	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write column values to output file
	w := csv.NewWriter(out)
	for _, value := range column {
		err = w.Write([]string{value})
		if err != nil {
			return err
		}
	}
	w.Flush()

	return nil
}

func ExtractColumnWithEmptyRows(inputFile string, outputFile string, columnNumber int) error {
	// Open input file
	f, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// Read CSV data
	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	if err != nil {
		return err
	}

	// Extract column values
	var column []string
	prevNum := -1 // previous numerical value of first cell
	for _, line := range lines {
		if len(line) >= columnNumber {
			num, err := strconv.Atoi(line[0])
			if err != nil {
				return err
			}
			if num-prevNum > 1 {
				column = append(column, "") // add empty row
			}
			column = append(column, line[columnNumber-1])
			prevNum = num
		}
	}

	// Open output file
	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write column values to output file
	w := csv.NewWriter(out)
	for _, value := range column {
		err = w.Write([]string{value})
		if err != nil {
			return err
		}
	}
	w.Flush()

	return nil
}

// TranslationDataType Define the data structure
type TranslationDataType struct {
	Row     string `json:"row"`
	English string `json:"english"`
	Persian string `json:"persian"`
}

func GenerateHtmlFile(inputFileName string) {
	// Open the CSV file
	file, err := os.Open(inputFileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Read the CSV data
	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	var translationData []TranslationDataType
	for i, row := range rows {
		if i == 0 {
			// Skip the header row
			continue
		}
		translationData = append(translationData, TranslationDataType{
			row[0],
			row[1],
			row[2],
		})
		if i%50 == 0 {
			// Generate a new HTML file every 50 records
			err = GenerateHTMLFile(translationData, i/50)
			if err != nil {
				panic(err)
			}
			translationData = []TranslationDataType{}
		}
	}

	// Generate the final HTML file
	err = GenerateHTMLFile(translationData, (len(rows)-1)/50+1)
	if err != nil {
		panic(err)
	}
}

// GenerateHTMLFile Helper function to generate an HTML file
func GenerateHTMLFile(translationData []TranslationDataType, fileIndex int) error {
	// Create a template object
	tmpl, err := template.New(fmt.Sprintf("TranslationData-%d", fileIndex)).Parse(htmlTemplate)
	if err != nil {
		return err
	}

	// Execute the template and generate the HTML output
	outputFile, err := os.Create(fmt.Sprintf("output-%d.html", fileIndex))
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Define the data to be passed to the template
	data := struct {
		Data string
	}{
		Data: toJSON(translationData),
	}
	err = tmpl.Execute(outputFile, data)
	if err != nil {
		return err
	}

	return nil
}

// Helper function to convert data to JSON
func toJSON(data interface{}) string {
	json, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	s := string(json)
	s = strings.ReplaceAll(s, "\"row\"", "row")
	s = strings.ReplaceAll(s, "\"english\"", "english")
	s = strings.ReplaceAll(s, "\"persian\"", "persian")
	s = strings.ReplaceAll(s, "\"", "'")
	return s
}

// Define the HTML template
var htmlTemplate = `<!DOCTYPE html>
<html>
<head>
	<title>Table with Editing</title>
	<style>
		table {
			border-collapse: collapse;
			width: 100%;
		}
		th, td {
			padding: 8px;
			text-align: left;
			border-bottom: 1px solid #ddd;
		}

        th {
            background-color: #f2f2f2;
        }

        tr:nth-child(even) {
            background-color: #f2f2f2;
        }

		td:nth-child(3) {
			text-align: right;
		}


		td:last-child {
			text-align: center;
		}

		textarea {
			width: 100%;
			height: 100px;
			text-align: right;
			direction: rtl;
			font-size: large;
		}
	</style>
</head>
<body>
    <table id="myTable">
        <thead>
          <tr>
            <th colspan="4">
              <button onclick="exportData()">Export as JSON</button>
            </th>
          </tr>
          <tr>
            <th>Row</th>
            <th>English Text</th>
            <th>Persian Translation</th>
            <th>Action</th>
          </tr>
        </thead>
        <tbody>
          <!-- Table rows will be added dynamically using JavaScript -->
        </tbody>
      </table>

	<script>
		// Define the data structure
		var translationData = {{.Data}};

		// Check if data exists in local storage
		if (localStorage.getItem("translationData")) {
			translationData = JSON.parse(localStorage.getItem("translationData"));
		}

		// Function to add rows to the table
		function addRow(rowData) {
			var table = document.getElementById("myTable").getElementsByTagName("tbody")[0];
			var newRow = table.insertRow();
			var cell1 = newRow.insertCell(0);
			var cell2 = newRow.insertCell(1);
			var cell3 = newRow.insertCell(2);
			var cell4 = newRow.insertCell(3);
			cell1.innerHTML = rowData.row;
			cell2.innerHTML = rowData.english;
			cell3.innerHTML = "<div dir='rtl'>" + rowData.persian + "</div>";
			cell4.innerHTML = "<button onclick='editCell(this)'>Edit</button>";
			cell3.setAttribute("data-persian", rowData.persian);
		}

		// Add rows to the table using the data structure
		for (var i = 0; i < translationData.length; i++) {
			addRow(translationData[i]);
		}

		// Function to enable editing of a cell
        function editCell(button) {
			var cell = button.parentNode.previousSibling;
			var persianText = cell.getAttribute("data-persian");
			cell.innerHTML = "<textarea>" + persianText + "</textarea>";
			button.innerHTML = "Save";
			button.setAttribute("onclick", "saveCell(this)");
			if (button.nextSibling) {
				button.nextSibling.removeAttribute("onclick");
			}
        }


		// Function to save the edited cell to local storage
        function saveCell(button) {
			var cell = button.parentNode.previousSibling;
			var persianText = cell.getElementsByTagName("textarea")[0].value;
			cell.innerHTML = "<div dir='rtl'>" + persianText + "</div>";
			button.innerHTML = "Edit";
			button.setAttribute("onclick", "editCell(this)");
			if (button.nextSibling) {
				button.nextSibling.removeAttribute("onclick");
			}
			cell.setAttribute("data-persian", persianText);
			// Update the data structure in local storage
			var row = cell.parentNode.cells[0].innerHTML;
			translationData[row - 1].persian = persianText;
			localStorage.setItem("translationData", JSON.stringify(translationData));
        }

        function exportData() {
			var dataStr = "data:text/json;charset=utf-8," + encodeURIComponent(JSON.stringify(translationData));
			var downloadAnchorNode = document.createElement('a');
			downloadAnchorNode.setAttribute("href", dataStr);
			downloadAnchorNode.setAttribute("download", "translationData.json");
			document.body.appendChild(downloadAnchorNode);
			downloadAnchorNode.click();
			downloadAnchorNode.remove();
        }
    </script>
</body>
</html>`
