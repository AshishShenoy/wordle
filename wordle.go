package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

const WORDS_URL = "https://raw.githubusercontent.com/dwyl/english-words/master/words_alpha.txt"
const WORD_LENGTH = 5
const MAX_GUESSES = 6

func get_filled_color_vector(color string) [WORD_LENGTH]string {
	color_vector := [WORD_LENGTH]string{}
	for i := range color_vector {
		color_vector[i] = color
	}
	return color_vector
}

func display_word(word string, color_vector [WORD_LENGTH]string) {
	for i, c := range word {
		switch color_vector[i] {
		case "Green":
			fmt.Print("\033[42m\033[1;30m")
		case "Yellow":
			fmt.Print("\033[43m\033[1;30m")
		case "Grey":
			fmt.Print("\033[40m\033[1;37m")
		}
		fmt.Printf(" %c ", c)
		fmt.Print("\033[m\033[m")
	}
	fmt.Println()
}

func main() {
	rand.Seed(time.Now().Unix())

	res, err := http.Get(WORDS_URL)
	if err != nil {
		log.Fatalln(err)
	}

	body, _ := ioutil.ReadAll(res.Body)
	words := strings.Split(string(body), "\r\n")

	wordle_words := []string{}
	for _, word := range words {
		if len(word) == WORD_LENGTH {
			wordle_words = append(wordle_words, strings.ToUpper(word))
		}
	}
	sort.Strings(wordle_words)

	selected_word := wordle_words[rand.Intn(len(wordle_words))]

	reader := bufio.NewReader(os.Stdin)
	guesses := []map[string][WORD_LENGTH]string{}
	var guess_count int
	for guess_count = 0; guess_count < MAX_GUESSES; guess_count++ {
		fmt.Printf("Enter your guess (%v/%v): ", guess_count+1, MAX_GUESSES)

		guess_word, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		guess_word = strings.ToUpper(guess_word[:len(guess_word)-1])

		if guess_word == selected_word {
			fmt.Println("You guessed right!")
			color_vector := get_filled_color_vector("Green")

			guesses = append(guesses, map[string][WORD_LENGTH]string{guess_word: color_vector})

			fmt.Println("Your wordle matrix is: ")
			for _, guess := range guesses {
				for guess_word, color_vector := range guess {
					display_word(guess_word, color_vector)
				}
			}
			break
		} else {
			i := sort.SearchStrings(wordle_words, guess_word)
			if i < len(wordle_words) && wordle_words[i] == guess_word {
				color_vector := get_filled_color_vector("Grey")

				// stores whether an index is allowed to cause another index to be yellow
				yellow_lock := [WORD_LENGTH]bool{}

				for j, guess_letter := range guess_word {
					for k, letter := range selected_word {
						if guess_letter == letter && j == k {
							color_vector[j] = "Green"
							// now the kth index can no longer cause another index to be yellow
							yellow_lock[k] = true
							break

						}
					}
				}
				for j, guess_letter := range guess_word {
					for k, letter := range selected_word {
						if guess_letter == letter && color_vector[j] != "Green" && yellow_lock[k] == false {
							color_vector[j] = "Yellow"
							yellow_lock[k] = true
						}
					}
				}
				guesses = append(guesses, map[string][WORD_LENGTH]string{guess_word: color_vector})
				display_word(guess_word, color_vector)
			} else {
				guess_count--
				fmt.Printf("Please guess a valid %v letter word from the wordlist", WORD_LENGTH)
				fmt.Println()
			}
		}
	}

	if guess_count == MAX_GUESSES {
		fmt.Println("Better luck next time!")
		color_vector := get_filled_color_vector("Green")
		fmt.Print("The correct word is : ")
		display_word(selected_word, color_vector)
	}
}
