package utils

import (
	"fmt"
	"strings"
)

func PrintProgressBar(percentage int) {
	barWidth := 50
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 100 {
		percentage = 100
	}
	filledWidth := (percentage * barWidth) / 100
	emptyWidth := barWidth - filledWidth
	bar := "[" + strings.Repeat("=", filledWidth) + strings.Repeat(" ", emptyWidth) + "]"
	fmt.Printf("\r%s %d%%", bar, percentage)
	if percentage == 100 {
		fmt.Println()
	}
}
