package main

import (
	"fmt"
	"io"
	"os"
	"testing"
)

type testCase struct {
	name     string
	input    string
	expected string
}

var testCases = []testCase{
	{
		name:     "Low, cap, up and punctuation",
		input:    `If I make you BREAKFAST IN BED (low, 3) just say thank you instead of: how (cap) did you get in my house (up, 2) ?`,
		expected: `If I make you breakfast in bed just say thank you instead of: How did you get in MY HOUSE?`,
	},
	{
		name:     "Bin and hex conversions",
		input:    `I have to pack 101 (bin) outfits. Packed 1a (hex) just to be sure`,
		expected: `I have to pack 5 outfits. Packed 26 just to be sure`,
	},
	{
		name:     "Punctuation",
		input:    "Don not be sad ,because sad backwards is das . And das not good",
		expected: "Don not be sad, because sad backwards is das. And das not good",
	},
	{
		name:     "Cap, punctuation and single quotes",
		input:    `harold wilson (cap, 2) : ' I am a optimist ,but a optimist who carries a raincoat . '`,
		expected: `Harold Wilson: 'I am an optimist, but an optimist who carries a raincoat.'`,
	},
	{
		name:     "Cap, up, low and puntuation again",
		input:    `it (cap) was the best of times, it was the worst of times (up) , it was the age of wisdom, it was the age of foolishness (cap, 6) , it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, IT WAS THE (low, 3) winter of despair.`,
		expected: `It was the best of times, it was the worst of TIMES, it was the age of wisdom, It Was The Age Of Foolishness, it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, it was the winter of despair.`,
	},
	{
		name:     "Hex and bin again",
		input:    `Simply add 42 (hex) and 10 (bin) and you will see the result is 68.`,
		expected: `Simply add 66 and 2 and you will see the result is 68.`,
	},
	{
		name:     "A to an",
		input:    `There is no greater agony than bearing a untold story inside you.`,
		expected: `There is no greater agony than bearing an untold story inside you.`,
	},
}

// TestCreateOutput1 calls createOutput with input from test cases and compares
// the results with known right answers
func TestCreateOutput1(t *testing.T) {
	for _, tc := range testCases {
		// using t.Run(string, func) produces more details in go test -v
		t.Run(tc.name, func(t *testing.T) {
			result := createOutput(tc.input)
			if tc.expected != result {
				t.Errorf("\ncreateOutput(%s) \ngot:  %#v,\nwant: %#v", tc.input, result, tc.expected)
			}
		})
	}
}

// TestCreateOutput2 calls createOutput with the input file sample.txt and compares
// the result with the correct result from resultTest.txt
func TestCreateOutput2(t *testing.T) {
	fileIn, err := os.Open("sample.txt")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer fileIn.Close()

	arrIn, err := io.ReadAll(fileIn)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fileWant, err := os.Open("resultTest.txt")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer fileWant.Close()

	arrWant, err := io.ReadAll(fileWant)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	res := createOutput(string(arrIn))
	want := string(arrWant)

	if res != want {
		t.Fatalf("The function createOutput returns:\n%s\nand not\n%s", res, want)
	}
}

// TestGetNumber calls getNumber with a sting containing numbers,
// checking if the return value is correct
func TestGetNumber(t *testing.T) {
	str := "Le77ers and (0ther!) stuff w1th numbers 23"
	want := 770123
	res := getNumber(str)

	if res != want {
		t.Fatalf(`"%s" returns %v and not %v`, str, res, want)
	}
}

// TestIsVowel calls isVowel with a string of different characters
// checking if all return what they should
func TestIsVowel(t *testing.T) {
	str := "abc?eh"
	results := []bool{}
	want := []bool{true, false, false, false, true, true}

	for _, r := range str {
		results = append(results, isVowel(r))
	}

	errMsg := fmt.Sprintf(`isVowel returns "%v" and not "%v" with the input %s`, results, want, str)

	for i := 0; i < len(results); i++ {
		if len(results) != len(want) || results[i] != want[i] {
			t.Fatal(errMsg)
		}
	}
}
