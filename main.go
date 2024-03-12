package main

import (
	"fmt"
	"math/rand/v2"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Word struct {
	letters   []*Letter
	word      string
	len       int32
	x, y      float32
	completed bool
}

type Letter struct {
	letter  string
	color   rl.Color
	written bool
}

func NewWord(word string) *Word {
	var w Word = Word{}
	w.len = 0

	for _, l := range word {
		var letter Letter = Letter{letter: string(l), color: rl.Black, written: false}
		w.word += letter.letter
		w.letters = append(w.letters, &letter)
		w.len += 1
	}
	w.completed = false

	var random float32 = float32(rand.Int32N(850-50) + 50)
	w.x = 0
	w.y = random

	return &w
}

func (w *Word) HandleText() {
	var x float32 = 0
	for i, letter := range w.letters {
		if i == 0 {
			rl.DrawText(letter.letter, int32(w.x+x), int32(w.y), 25, letter.color)
		} else {
			if w.letters[i-1].letter == "i" || w.letters[i-1].letter == "l" {
				x -= 9
				rl.DrawText(letter.letter, int32(w.x+x), int32(w.y), 25, letter.color)
			} else if w.letters[i-1].letter == "j" {
				x -= 7
				rl.DrawText(letter.letter, int32(w.x+x), int32(w.y), 25, letter.color)
			} else if w.letters[i].letter == "j" {
				x -= 2
				rl.DrawText(letter.letter, int32(w.x+x), int32(w.y), 25, letter.color)
			} else {
				rl.DrawText(letter.letter, int32(w.x+x), int32(w.y), 25, letter.color)
			}
		}
		x += 15
	}
}

func (w *Word) HandleKeyPressed(points *int) {
	key := rl.GetCharPressed()
	if string(key) == w.letters[0].letter {
		w.letters[0].written = true
		w.letters[0].color = rl.Blue
	}

	if key != 0 && w.letters[0].written {
		for i, letter := range w.letters {
			if i > 0 && string(key) == letter.letter && w.letters[i-1].written {
				letter.color = rl.Blue
				letter.written = true
			}
		}
	}

	if w.letters[w.len-1].written {
		if w.word != "start" {
			*points += 1
		}
	}
}

func (w *Word) Draw(points *int, word_l *[]*Word) {
	if !w.letters[w.len-1].written {
		w.HandleKeyPressed(points)
		w.HandleText()

		if w.word != "start" {
			w.x += 150 * rl.GetFrameTime()
		} else {
			w.x = (1600/2)-(float32(rl.MeasureText("start", 25))/2)
			w.y = 900/2
		}
	} else {
		copy_w := *word_l
		for i := range len(*word_l) {
			if w.word == copy_w[i].word {
				fmt.Println("BEFORE: ", len(*word_l))
				*word_l = append(copy_w[:i], copy_w[i+1:]...)
				fmt.Println("AFTER: ", len(*word_l))
				break
			}
		}
	}
}

func contains(used_random *[]int, i int) bool {
	for _, x := range *used_random {
		if i == x {
			return true
		}
	}

	return false
}

func find_index(word_list *[]*Word, word *Word) int {
	for i, w := range *word_list {
		if word.word == w.word {
			return i
		}
	}
	return -1
}

func randomize(max_range int, word_list *[]*Word, used_random *[]int, random_list *[]*Word) {
	copy_word_list := *word_list
	
	for i := 0; i < max_range; i++ {
		var random int = rand.IntN(len(*word_list))

		for contains(used_random, random) {
			random = rand.IntN(len(*word_list))
		}

		*random_list = append(*random_list, NewWord(copy_word_list[random].word))
		*used_random = append(*used_random, random)
	}
}

// func add_random(used_random *[]int, word_list *[]*Word, word_list_random []*Word) {
// 	for i := range word_list_random {
// 		r := find_index(word_list, word_list_random[i])
// 		*used_random = append(*used_random, r)
// 	}
// }

func main() {
	var width int32 = 1600
	var height int32 = 900

	rl.InitWindow(width, height, "Words")
	defer rl.CloseWindow()
	rl.SetTargetFPS(240)

	var word_list []*Word = []*Word{NewWord("stock"), NewWord("remember"), NewWord("goroutines"),
		NewWord("hope"), NewWord("fought"), NewWord("jesus"), NewWord("throught"), NewWord("throughout"),
		NewWord("overwrought"), NewWord("smartwatch"), NewWord("thinking"), NewWord("playground"),
		NewWord("modify"), NewWord("random"), NewWord("window"), NewWord("released"), NewWord("package"),
	}

	var word_list_random []*Word
	var used_random []int

	randomize(len(word_list), &word_list, &used_random, &word_list_random)

	var points int = 0
	var started bool = false
	var start *Word = NewWord("start")
	var start_written bool = false
	var first_time bool = true

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		var fps string = strconv.Itoa(int(rl.GetFPS()))
		var mouseX string = strconv.Itoa(int(rl.GetMouseX()))
		var mouseY string = strconv.Itoa(int(rl.GetMouseY()))

		rl.DrawText("FPS: "+fps, 1400, 10, 30, rl.Black)
		rl.DrawText("MouseX: "+mouseX, 1400, 40, 30, rl.Black)
		rl.DrawText("MouseY: "+mouseY, 1400, 70, 30, rl.Black)
		rl.DrawText("Points: "+strconv.Itoa(points), 1400, 100, 30, rl.Black)

		if !start_written {
			start.Draw(&points, &word_list_random)
			if start.letters[start.len-1].written { start_written = true }
		} else {
			// create more words
			if len(word_list_random) < 4 {
				used_random = used_random[:0]

				for i := range word_list_random {
					r := find_index(&word_list, word_list_random[i])
					used_random = append(used_random, r)
				}
				
				randomize(len(word_list)-len(used_random), &word_list, &used_random, &word_list_random)
				started = false
			}
			
			// start at random posX
			if !started {
				var spacing float32 = 200
				for i, w := range word_list_random {
					if first_time {
						if i > 0 {
							w.x -= 300 + spacing
							spacing += 200
						}
					} else {
						if i > 2 {
							w.x -= 300 + spacing
							spacing += 200
						}
					}
				}
				first_time = false
				started = true
			}
	
			// checks if outside screen
			if word_list_random[0].x < float32(width) {
				for _, w := range word_list_random {
					w.Draw(&points, &word_list_random)
				}
			} else {
				rl.DrawText("you lost", (width/2)-(rl.MeasureText("you lost", 35)/2), height/2, 35, rl.Red)
			}
		}

		rl.EndDrawing()
	}
}
