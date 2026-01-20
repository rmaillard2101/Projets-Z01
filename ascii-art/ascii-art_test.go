package main

import (
	art "ascii-art/functions"
	"reflect"
	"testing"
)

func TestAscii(t *testing.T) {
	letter := " "
	got := art.ConvertitEnMatrice(letter[0], "shadowtagged.txt")
	expected := []string{"      ", "      ", "      ", "      ", "      ", "      ", "      ", "      "}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %q, got %q", expected, got)
	}

	letter1 := []string{"A ", "B ", "c ", "d ", "e ", "f ", "g ", "h "}
	letter2 := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	got = art.AddLetter(letter1, letter2)
	expected = []string{"A 1", "B 2", "c 3", "d 4", "e 5", "f 6", "g 7", "h 8"}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}
