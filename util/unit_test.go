package util

import (
	"fmt"
	"testing"
)

func TestUnit(t *testing.T) {
	unit := "30M"
	fmt.Println(Unit2Int(unit))
}
