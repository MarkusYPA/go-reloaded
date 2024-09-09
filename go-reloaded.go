package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"unicode"

	"go-reloaded/piscinefuncs"
)

// isVowel checks if a given rune is vowel-like at the start of a word
func isVowel(r rune) bool {
	return r == 'a' || r == 'e' || r == 'i' || r == 'o' || r == 'u' || r == 'h'
}

// getNumber returns the numbers in a string as an integer
func getNumber(s string) int {
	numString := ""
	for _, r := range s {
		if unicode.IsNumber(r) {
			numString += string(r)
		}
	}
	if len(numString) != 0 {
		return piscinefuncs.ToDec(numString, "0123456789")
	} else {
		return 1
	}
}

// findStartAndEnd finds the beginning and end of a defined
// number of words preceding a given index i
func findStartAndEnd(i int, in string, howMany int) (start, end int) {
	start = 0
	end = 0

	spaces := 0
	isSpace := false
	prevIsSpace := false
	endFound := false

	// scan the input in reverse starting from the beginning of the match
	for {
		if unicode.IsSpace(rune(in[i])) {
			isSpace = true
			if !prevIsSpace {
				spaces++
			}
		} else {
			isSpace = false
		}
		if !unicode.IsSpace(rune(in[i])) && prevIsSpace && !endFound {
			end = i + 1 // word to convert ends before first found space
			endFound = true
		}
		prevIsSpace = isSpace
		if spaces == 1+howMany {
			break
		}
		i--
		if i < 0 {
			break
		}
	}
	start = i + 1 // word(s) to convert start after all spaces are found

	return
}

// getAllStartsAndEnds detects starts and ends of words to be modified according to
// instructions in the input string
func getAllStartsAndEnds(matches [][]int, input string) (starts, ends []int) {
	starts = []int{}
	ends = []int{}

	for _, pair := range matches {
		howMany := getNumber(input[pair[0]:pair[1]])
		wStart, wEnd := findStartAndEnd(pair[0], input, howMany)
		starts = append(starts, wStart)
		ends = append(ends, wEnd)
	}

	return
}

// fixPunctuation corrects erraneous spaces around punctuation
func fixPunctuation(in string) string {
	output := ""
	for i, r := range in {
		output += string(r)
		if unicode.IsPunct(r) && r != '\'' {
			spaces := 0
			for i > spaces && unicode.IsSpace(rune(in[i-1-spaces])) { // count spaces before punctuation
				spaces++
			}
			output = output[:len(output)-1-spaces] + string(r) // remove found spaces
		}
		if unicode.IsPunct(r) && r != '\'' {
			if len(in) > i+2 && !unicode.IsSpace(rune(in[i+1])) && !unicode.IsPunct(rune(in[i+1])) { // No space or other punctuation after punctuation
				output += " " // add the missing space
			}
		}
	}
	return output
}

// fixAsToAns adds an 'n' to the article 'a' when needed
func fixAsToAns(in string) string {
	output := ""
	for i, r := range in {
		output += string(r)
		if i > 2 && r == ' ' {
			// find articles a or A
			if (in[i-1] == 'a' && in[i-2] == ' ') || (in[i-1] == 'A' && unicode.IsPunct(rune(in[i-2]))) {
				if i < len(in)-2 && isVowel(rune(in[i+1])) { // next word starts with vowel
					output = output[:len(output)-1] + "n" + output[len(output)-1:] // add n to a or A
				}
			}
		}
	}
	return output
}

// fixSingleQuotes makes sure there's a space only on the correct side of a single quote
func fixSingleQuotes(in string) string {
	output := ""
	singlesCount := 0
	skip := false
	for i, r := range in {
		output += string(r)
		if r == '\'' {
			if i > 3 && in[i-3:i] == "don" && in[i+1] == 't' { // exception for "don't"
				skip = true
				continue
			}

			if singlesCount%2 == 0 { // left single quotes
				if i > 0 && !unicode.IsSpace(rune(in[i-1])) { // no space before left '
					output = output[:len(output)-2] + " " + string(r) // add the missing space
				}
			} else { // right single quotes
				spaces := 0
				for i > spaces && unicode.IsSpace(rune(in[i-1-spaces])) { // space(s) before right '
					spaces++
				}
				output = output[:len(output)-1-spaces] + string(r) // remove the space(s)
			}
			singlesCount++
		}

		if singlesCount%2 == 1 && unicode.IsSpace(r) && i > 0 && in[i-1] == '\'' { // space after left '
			output = output[:len(output)-1] // remove the space
		}
		if singlesCount%2 == 0 && !unicode.IsSpace(r) && i > 0 && in[i-1] == '\'' && !skip { // no space after right '
			output = output[:len(output)-1] + " " + string(r) // add the space
		}
		skip = false
	}
	return output
}

// createOutput makes the desired changes to the input text and corrects errors
func createOutput(input string) string {

	// regular expressions to find different instructions in the input
	reHex := regexp.MustCompile(`\(hex\)`)
	reBin := regexp.MustCompile(`\(bin\)`)
	reLow := regexp.MustCompile(`\((low)(,\s*-?\d+)?\)`) // catches both (low) and (low, <integer>)
	reUp := regexp.MustCompile(`\((up)(,\s*-?\d+)?\)`)
	reCap := regexp.MustCompile(`\((cap)(,\s*-?\d+)?\)`)

	var matches [][]int // indexes in input of instructions in parentheses

	baseREs := []*regexp.Regexp{reHex, reBin}
	baseInputs := []string{"0123456789abcdef", "01"}
	stringREs := []*regexp.Regexp{reLow, reUp, reCap}
	stringFuncs := []func(string) string{piscinefuncs.ToLower, piscinefuncs.ToUpper, piscinefuncs.Capitalize}

	// find and convert hex and bin numbers
	for i := 0; i < 2; i++ {
		matches = baseREs[i].FindAllIndex([]byte(input), -1)
		prevWordStarts, prevWordEnds := getAllStartsAndEnds(matches, input)

		diff := 0 // compensates for change in index when input string is modified
		for j, start := range prevWordStarts {
			lenBefore := len(input)
			insert := fmt.Sprintf("%v", piscinefuncs.ToDec(input[start-diff:prevWordEnds[j]-diff], baseInputs[i]))
			input = input[:start-diff] + insert + input[matches[j][1]-diff:]
			diff += lenBefore - len(input)
		}
	}

	// find and put words to lowercase or uppercase or capitalize them
	for i := 0; i < 3; i++ {
		matches = stringREs[i].FindAllIndex([]byte(input), -1)
		prevWordStarts, prevWordEnds := getAllStartsAndEnds(matches, input)

		diff := 0
		for j, start := range prevWordStarts {
			lenBefore := len(input)
			input = input[:start-diff] + stringFuncs[i](input[start-diff:prevWordEnds[j]-diff]) + input[matches[j][1]-diff:]
			diff += lenBefore - len(input)
		}
	}

	// fix punctuation
	input = fixPunctuation(input)

	// fix single quotes
	input = fixSingleQuotes(input)

	// fix wrong "a":s to "an":s
	input = fixAsToAns(input)

	return input
}

// main loads the user specified file, sends its content as a string
// to createOutput() and saves the result in another user specified
// file
func main() {
	if len(os.Args) < 3 {
		fmt.Println("Provide input and output file names")
		fmt.Println("Example: ./goReloaded sample.txt result.txt")
		os.Exit(1)
	}

	if len(os.Args) > 3 {
		fmt.Println("Too many arguments")
		os.Exit(1)
	}

	// open and read the file from the first argument
	fileIn, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer fileIn.Close()

	bytes, err := io.ReadAll(fileIn)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// manipulate the input as a string to the desired form
	output := createOutput(string(bytes))

	// create and write the file from the second argument to store the result
	fileOut, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer fileOut.Close()

	_, err = fileOut.WriteString(output)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}
