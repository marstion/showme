package img

import "testing"

func TestTrimming(t *testing.T) {
	var position [][]int = [][]int{{1000, 800}, {0, 0}, {2222, 1065}, {0, 0}}
	var offset []int = []int{-10, -10, 10, 10}

	TrimmingForFile("image.jpg", "test.jpg", position, offset)
}
