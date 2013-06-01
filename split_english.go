package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Trie struct {
	children []*Trie
	letter   string
	end      bool
}

func readDictionary(fileName string) []string {
	/*
	   Read list of common words from a file into a string array
	*/
	file, error := ioutil.ReadFile(fileName)
	if error != nil {
		fmt.Println("Cannot open dictionary file.")
	}
	return strings.Split(string(file), "\n")
}

func buildTrie(dict []string, reverse bool) *Trie {
	/*
	   Build a trie for efficient retrieval of all our dictionary words
	*/
	var root *Trie = &Trie{[]*Trie{}, "", false}
	for _, w := range dict {
		if reverse {
			w = Reverse(w)
		}
		var t *Trie = root
		for i, c := range w {
			if i == len(w)-1 {
				t = insert(t, string(c), true)
			} else {
				t = insert(t, string(c), false)
			}
		}
	}
	return root
}

func insert(t *Trie, c string, end bool) *Trie {
	// if clean tree
	if t == nil {
		// return new tree
		return &Trie{[]*Trie{}, c, end}
	}

	// if node has children
	if len(t.children) > 0 {
		// check children
		for i := range t.children {
			// if letter in child
			if t.children[i].letter == c {
				t.children[i].end = t.children[i].end || end
				// go on from there
				return t.children[i]
			}
		}
	}
	// letter was not in child
	var newNode *Trie = &Trie{[]*Trie{}, c, end}
	t.children = append(t.children, newNode)
	return newNode
}

func findLongestWord(t *Trie, phrase string) string {
	/*
	   Find the longest word from the start of a string.
	   Returns empty string if no word is found.
	*/
	var word string = ""
	var progress string = ""
	var r *Trie = t
	for _, c := range phrase {
		var foundLetter bool = false
		for _, child := range r.children {
			if child.letter == string(c) {
				foundLetter = true
				r = child
				progress += r.letter
				if r.end == true {
					word = progress[:]
				}
			}
		}
		if !foundLetter {
			break
		}
	}
	return word
}

func splitPhrase(t *Trie, phrase string, sep string) string {
	/*
			   Split a spaceless phrase into words separated by separator string.

			   Just chooses the longest word at each point.

		       Attempts to leave the minimum amount of characters in the phrase 
		       unidentified by the dictionary.
	*/
	var newPhrase []string = []string{}
	unmatched := 0
	for i := 0; i < len(phrase); i++ {
		word := findLongestWord(t, phrase[i:])
		if len(word) > 0 {
			if unmatched > 0 {
				newPhrase = append(newPhrase, phrase[i-unmatched:i])
			}
			i += len(word) - 1
			newPhrase = append(newPhrase, word)
			unmatched = 0
		} else {
			unmatched++
		}
	}
	if unmatched > 0 {
		newPhrase = append(newPhrase, phrase[len(phrase)-unmatched:])
	}
	return strings.Join(newPhrase, sep)
}

func findPossibleWords(t *Trie, phrase string, start int) []string {
	/*
		Find all possible dictionary words that could form the start
		of the phrase.

		start: integer index of starting letter to look at 
		(cheaper than copying stringsin memory)
	*/
	var words []string = []string{}
	var progress string = ""
	var r *Trie = t
	for i := start; i < len(phrase); i++ {
		c := phrase[i]
		var foundLetter bool = false
		for _, child := range r.children {
			if child.letter == string(c) {
				foundLetter = true
				r = child
				progress += r.letter
				if r.end == true {
					words = append(words, progress)
				}
			}
		}
		if !foundLetter {
			break
		}
	}
	return words
}

func _minEntropyRecurse(t *Trie, phrase string, start int, maxLevel int) int {
	if start >= len(phrase) {
		return 0
	}
	if maxLevel <= 0 {
		return 1
	}
	words := findPossibleWords(t, phrase, start)
	if len(words) == 0 {
		return -3 + _minEntropyRecurse(t, phrase, start+1, maxLevel-1)
	}
	var min int = -1
	for _, word := range words {
		num := len(word) + _minEntropyRecurse(t, phrase, start+len(word), maxLevel-1)
		fmt.Println("---", Reverse(word), num)
		if min < 0 || num >= min {
			min = num
		}
	}
	return min
}

func splitPhraseMinEntropy(t *Trie, phrase string, sep string, reverse bool) string {
	/*
	   Split a spaceless phrase into words separated by separator string.

	   Attempts to leave the minimum amount of characters in the phrase 
	   unidentified by the dictionary.
	*/
	if reverse {
		phrase = Reverse(phrase)
	}
	bestWords := []string{}
	unmatched := 0
	for u := 0; u < len(phrase); u++ {
		words := findPossibleWords(t, phrase, u)
		min := -1
		bestWord := ""
		for _, w := range words {
			fmt.Println("Testing", Reverse(w))
			n := len(w) + _minEntropyRecurse(t, phrase, u+len(w), 5)
			fmt.Println(Reverse(w), n)
			if min == -1 || n >= min {
				bestWord = w
				min = n
			}
		}
		fmt.Println("USING", Reverse(bestWord))
		if len(bestWord) == 0 {
			unmatched++
		} else {
			if unmatched > 0 {
				bestWords = append(bestWords, phrase[u-unmatched:u])
				unmatched = 0
			}
			bestWords = append(bestWords, bestWord)
			u += len(bestWord) - 1
		}
	}
	if unmatched > 0 {
		bestWords = append(bestWords, phrase[len(phrase)-unmatched:])
	}
	if reverse {
		return Reverse(strings.Join(bestWords, sep))
	}
	return strings.Join(bestWords, sep)
}

func StripPunctuation(word string) string {
	return strings.Trim(word, " .-+,—")
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func main() {
	file, error := ioutil.ReadFile("doc.in")
	if error != nil {
		fmt.Println("Cannot open doc.in, make sure you have it in the current directory")
	}
	text := string(file)
	paragraphs := strings.Split(text, "\n")
	var words []string
	for _, paragraph := range paragraphs {
		for _, word := range strings.Split(paragraph, " ") {
			words = append(words, StripPunctuation(strings.ToLower(word)))
		}
	}

	dict := readDictionary("english_dict.txt")
	for i := range dict {
		dict[i] = strings.Trim(strings.ToLower(dict[i]), " \n,.")
	}
	var trie *Trie = buildTrie(dict, true)
	phrase := strings.Join(words[:100], "")
	fmt.Println(splitPhraseMinEntropy(trie, phrase, " ", true))
}
