package utils

import (
	"fmt"
	"strings"
	"time"
)

func ColorArrowProgressBar(completedUnits int, totalUnits int) {
	// Calculate the percentage of completion.
	percentage := float64(completedUnits) / float64(totalUnits) * 100

	// Determine the color of the progress bar based on the percentage of completion.
	var color string
	if percentage >= 90 {
		color = "\033[1;32m" // Green
	} else if percentage >= 50 {
		color = "\033[1;33m" // Yellow
	} else {
		color = "\033[1;31m" // Red
	}

	// Calculate the number of completed and remaining blocks.
	completedBlocks := int(percentage / 2)
	remainingBlocks := 50 - completedBlocks

	// Create the progress bar string with the color code.
	progressBar := fmt.Sprintf("[%s%s>\033[0m%s%s] %.2f%%%s", color, strings.Repeat("=", completedBlocks), strings.Repeat(".", remainingBlocks), "\033[0m", percentage, "\033[0m")

	if remainingBlocks == 0 {
		progressBar = fmt.Sprintf("[%s%s\033[0m%s%s] %.2f%%%s", color, strings.Repeat("=", completedBlocks), strings.Repeat(".", remainingBlocks), "\033[0m", percentage, "\033[0m")
	}
	// Print the progress bar string.
	fmt.Printf("\r%s", progressBar)
}

func RotatingProgressBarWithDots(completedUnits int, totalUnits int) {
	// Define the list of characters to use in the rotating progress bar.
	progressChars := []string{"\\", "|", "/", "-"}

	// Calculate the percentage of completion.
	percentage := float64(completedUnits) / float64(totalUnits) * 100

	// Calculate the index of the current progress character based on the completed units.
	progressIndex := (completedUnits / 2) % len(progressChars)

	// Calculate the number of completed and remaining blocks.
	completedBlocks := int(percentage / 2)
	remainingBlocks := 50 - completedBlocks

	// Create the progress bar string with the rotating progress character.
	progressBar := fmt.Sprintf("[%s%s%s] %.2f%%", strings.Repeat("=", completedBlocks), progressChars[progressIndex], strings.Repeat(".", remainingBlocks), percentage)

	// Print the progress bar string and add a small delay to make the progress bar appear to rotate.
	fmt.Printf("\r%s", progressBar)
	time.Sleep(100 * time.Millisecond)
}

func AnimateProgressBar(completedUnits int, totalUnits int) {
	// Define the list of characters to use in the rotating progress bar.
	progressChars := []string{"|", "/", "-", "\\"}

	// Calculate the percentage of completion.
	percentage := float64(completedUnits) / float64(totalUnits) * 100

	// Calculate the number of completed and remaining blocks.
	completedBlocks := int(percentage / 2)
	remainingBlocks := 50 - completedBlocks

	// Calculate the index of the current progress character based on the current time.
	progressIndex := int(time.Now().UnixNano()/100000000) % len(progressChars)

	// Create the progress bar string with the rotating progress character.
	progressBar := fmt.Sprintf("[%s%s%s] %.2f%%", strings.Repeat("=", completedBlocks), progressChars[progressIndex], strings.Repeat(".", remainingBlocks), percentage)

	// Print the progress bar string.
	fmt.Printf("\r%s", progressBar)
}

func ColorBlockProgressBar(completedUnits int, totalUnits int) {
	// Calculate the percentage of completion.
	percentage := float64(completedUnits) / float64(totalUnits) * 100

	// Calculate the number of completed and remaining blocks.
	completedBlocks := int(percentage / 2)
	remainingBlocks := 50 - completedBlocks

	// Create the progress bar string with ANSI escape codes for colors.
	progressBar := fmt.Sprintf("\r[%s%s]%s %.2f%%",
		strings.Repeat("\033[32m█\033[0m", completedBlocks),
		strings.Repeat("\033[31m░\033[0m", remainingBlocks),
		strings.Repeat(" ", 2),
		percentage)

	// Print the progress bar string.
	fmt.Printf("%s", progressBar)
}

func BlockProgressBar(completedUnits int, totalUnits int) {
	// Calculate the percentage of completion.
	percentage := float64(completedUnits) / float64(totalUnits) * 100

	// Calculate the number of completed and remaining blocks.
	completedBlocks := int(percentage / 2)
	remainingBlocks := 50 - completedBlocks

	// Create the progress bar string.
	progressBar := fmt.Sprintf("[%s%s]", strings.Repeat("█", completedBlocks), strings.Repeat("░", remainingBlocks))

	// Print the progress bar string.
	fmt.Printf("\r%s %v/%v", progressBar, completedUnits, totalUnits)
}

func ArrowProgressBar(completedUnits int, totalUnits int) {
	// Calculate the percentage of completion.
	percentage := float64(completedUnits) / float64(totalUnits) * 100

	// Calculate the number of completed and remaining blocks.
	completedBlocks := int(percentage / 2)
	remainingBlocks := 50 - completedBlocks

	// Create the progress bar string.
	progressBar := fmt.Sprintf("[%s>%s] %.2f%%", strings.Repeat("=", completedBlocks), strings.Repeat(".", remainingBlocks), percentage)

	// Print the progress bar string.
	fmt.Printf("\r%s", progressBar)
}

func RotatingProgressBar(currentRow int, totalRows int) {
	// Define the list of characters to use in the rotating progress bar.
	progressChars := []string{"\\", "|", "/", "-"}

	// Calculate the index of the current progress character based on the current row.
	progressIndex := (currentRow / 2) % len(progressChars)

	// Calculate the percentage of completion.
	percentage := float64(currentRow) / float64(totalRows) * 100

	// Calculate the number of completed and remaining blocks.
	completedBlocks := int(percentage / 2)
	remainingBlocks := 50 - completedBlocks

	// Create the progress bar string with the rotating progress character.
	progressBar := fmt.Sprintf("[%s%s%s] %.2f%%", strings.Repeat("=", completedBlocks), progressChars[progressIndex], strings.Repeat(" ", remainingBlocks), percentage)

	// Print the progress bar string and percentage of completion.
	fmt.Printf("\r%s", progressBar)
	time.Sleep(100 * time.Millisecond) // Add a small delay to make the progress bar appear to rotate.
}

func ProgressBar(currentRow int, totalRows int) {
	// Calculate the percentage of completion.
	percentage := float64(currentRow) / float64(totalRows) * 100

	// Calculate the number of completed and remaining blocks.
	completedBlocks := int(percentage / 2)
	remainingBlocks := 50 - completedBlocks

	// Create the progress bar string.
	progressBar := "[" + strings.Repeat("=", completedBlocks) + strings.Repeat(" ", remainingBlocks) + "]"

	// Print the progress bar string and percentage of completion.
	fmt.Printf("\r%s %.2f%%\r", progressBar, percentage)
}
