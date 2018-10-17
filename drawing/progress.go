package drawing

import (
	"math"
	"strings"
)

// ProgressBar draws an ASCII progress bar.
func ProgressBar(complete float64, size int) string {
	progress := make([]string, size+2)

	for i := range progress {
		progress[i] = " "
	}

	progress[0] = "["
	progress[size+1] = "]"

	segments := math.Ceil(float64(size) * (complete / 100))

	for i := 1; i < int(segments)+1; i++ {
		progress[i] = "="
	}

	return strings.Join(progress, "")
}
