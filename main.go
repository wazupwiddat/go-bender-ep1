package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type BenderTracks []Bender

var tracks = BenderTracks{}

func (t BenderTracks) debug() {
	for _, bender := range t {
		fmt.Fprintln(os.Stderr, Directions[bender.direction])
	}
}

func (t BenderTracks) print() {
	if len(t) == 0 {
		fmt.Println("LOOP")
	}
	for _, bender := range t {
		fmt.Println(Directions[bender.direction])
	}
}

type Direction int

const (
	SOUTH Direction = 0
	EAST  Direction = 1
	NORTH Direction = 2
	WEST  Direction = 3
)

var Directions = [...]string{
	"SOUTH",
	"EAST",
	"NORTH",
	"WEST",
}

type Bender struct {
	x, y      int
	breaker   bool
	inverted  bool
	direction Direction
}

func (b Bender) southStep() (int, int) {
	return b.x, b.y + 1
}

func (b Bender) northStep() (int, int) {
	return b.x, b.y - 1
}

func (b Bender) eastStep() (int, int) {
	return b.x + 1, b.y
}

func (b Bender) westStep() (int, int) {
	return b.x - 1, b.y
}

func (b *Bender) stepSouth() {
	b.y++
}

func (b *Bender) stepNorth() {
	b.y--
}

func (b *Bender) stepEast() {
	b.x++
}

func (b *Bender) stepWest() {
	b.x--
}

type stateFn func(gameMap *GameMap, position *Bender) stateFn

type GameMap [][]string

func NewGameMap(l, c int) GameMap {
	gm := make(GameMap, l)
	for i := range gm {
		gm[i] = make([]string, c)
	}
	return gm
}

func printGameMap(gameMap [][]string) {
	for _, x := range gameMap {
		for _, v := range x {
			fmt.Fprint(os.Stderr, v)
		}
		fmt.Fprintln(os.Stderr)
	}
}

func locateBender(gameMap [][]string) Bender {
	for i, x := range gameMap {
		for j, v := range x {
			if v == "@" {
				return Bender{j, i, false, false, -1}
			}
		}
	}
	return Bender{0, 0, false, false, -1}
}

func locateOtherTransport(gameMap [][]string, bender Bender) (int, int) {
	for i, x := range gameMap {
		for j, v := range x {
			if bender.y == i && bender.x == j {
				continue
			}
			if v == "T" {
				return j, i
			}
		}
	}
	return bender.x, bender.y
}

func benderModifiers(gameMap GameMap, bender *Bender) {
	breaker := gameMap[bender.y][bender.x] == "B"
	if breaker {
		bender.breaker = !bender.breaker
	}
	inverted := gameMap[bender.y][bender.x] == "I"
	if inverted {
		bender.inverted = !bender.inverted
	}
	transport := gameMap[bender.y][bender.x] == "T"
	if transport {
		x, y := locateOtherTransport(gameMap, *bender)
		bender.x, bender.y = x, y
	}
}

func headWest(gameMap *GameMap, bender *Bender) stateFn {
	x, y := bender.westStep()
	if booth(*gameMap, x, y) {
		tracks = append(tracks, Bender{x: bender.x, y: bender.y, direction: WEST})
		return nil
	}
	if !bender.breaker && obstacle(*gameMap, x, y) {
		return handleMove
	}
	if obstacle(*gameMap, x, y) {
		(*gameMap)[y][x] = " "
	}
	if wall(*gameMap, x, y) {
		return handleMove
	}
	bender.stepWest()
	tracks = append(tracks, Bender{x: bender.x, y: bender.y, direction: WEST})
	if benderLooping(bender) {
		return loop
	}

	benderModifiers(*gameMap, bender)
	mod := directionModifier(*gameMap, bender.x, bender.y)
	if mod != nil {
		return mod
	}
	return headWest(gameMap, bender)
}

func headNorth(gameMap *GameMap, bender *Bender) stateFn {
	x, y := bender.northStep()
	if booth(*gameMap, x, y) {
		tracks = append(tracks, Bender{x: bender.x, y: bender.y, direction: NORTH})
		return nil
	}
	if !bender.breaker && obstacle(*gameMap, x, y) {
		return handleMove
	}
	if obstacle(*gameMap, x, y) {
		(*gameMap)[y][x] = " "
	}
	if wall(*gameMap, x, y) {
		return handleMove
	}
	bender.stepNorth()
	tracks = append(tracks, Bender{x: bender.x, y: bender.y, direction: NORTH})
	if benderLooping(bender) {
		return loop
	}

	benderModifiers(*gameMap, bender)
	mod := directionModifier(*gameMap, bender.x, bender.y)
	if mod != nil {
		return mod
	}
	return headNorth(gameMap, bender)
}

func headEast(gameMap *GameMap, bender *Bender) stateFn {
	x, y := bender.eastStep()
	if booth(*gameMap, x, y) {
		tracks = append(tracks, Bender{x: bender.x, y: bender.y, direction: EAST})
		return nil
	}
	if !bender.breaker && obstacle(*gameMap, x, y) {
		return handleMove
	}
	if obstacle(*gameMap, x, y) {
		(*gameMap)[y][x] = " "
	}
	if wall(*gameMap, x, y) {
		return handleMove
	}
	bender.stepEast()
	tracks = append(tracks, Bender{x: bender.x, y: bender.y, direction: EAST})
	if benderLooping(bender) {
		return loop
	}

	benderModifiers(*gameMap, bender)
	mod := directionModifier(*gameMap, bender.x, bender.y)
	if mod != nil {
		return mod
	}
	return headEast(gameMap, bender)
}

func headSouth(gameMap *GameMap, bender *Bender) stateFn {
	x, y := bender.southStep()
	if booth(*gameMap, x, y) {
		tracks = append(tracks, Bender{x: bender.x, y: bender.y, direction: SOUTH})
		return nil
	}
	if !bender.breaker && obstacle(*gameMap, x, y) {
		return handleMove
	}
	if obstacle(*gameMap, x, y) {
		(*gameMap)[y][x] = " "
	}
	if wall(*gameMap, x, y) {
		return handleMove
	}
	bender.stepSouth()
	tracks = append(tracks, Bender{x: bender.x, y: bender.y, direction: SOUTH})
	if benderLooping(bender) {
		return loop
	}

	benderModifiers(*gameMap, bender)
	mod := directionModifier(*gameMap, bender.x, bender.y)
	if mod != nil {
		return mod
	}

	return headSouth(gameMap, bender)
}

func benderLooping(bender *Bender) bool {
	currentPath := [2]BenderTracks{}
	pathIndex := 0
	for i := len(tracks) - 1; i > 0; i-- {
		track := tracks[i]
		if i == len(tracks)-1 {
			currentPath[pathIndex] = append(currentPath[pathIndex], track)
			continue
		}
		if track.x == bender.x && track.y == bender.y {
			if pathIndex == 1 {
				break
			}
			pathIndex = 1
		}
		currentPath[pathIndex] = append(currentPath[pathIndex], track)
	}
	if len(currentPath[0]) > 1 && len(currentPath[0]) == len(currentPath[1]) {
		for i, bender := range currentPath[0] {
			if bender.x != currentPath[1][i].x && bender.y != currentPath[1][i].y {
				return false
			}
		}
		return true
	}
	return false
}

func directionModifier(gameMap GameMap, x, y int) stateFn {
	if y >= len(gameMap) || x >= len(gameMap[y]) {
		return nil
	}
	mod := gameMap[y][x]
	if mod != "" {
		switch mod {
		case "E":
			return headEast
		case "S":
			return headSouth
		case "N":
			return headNorth
		case "W":
			return headWest
		}
	}
	return nil
}

func handleMove(gameMap *GameMap, bender *Bender) stateFn {
	if bender.inverted {
		if x, y := bender.westStep(); !obstacle(*gameMap, x, y) && !wall(*gameMap, x, y) {
			return headWest
		}
		if x, y := bender.northStep(); !obstacle(*gameMap, x, y) && !wall(*gameMap, x, y) {
			return headNorth
		}
		if x, y := bender.eastStep(); !obstacle(*gameMap, x, y) && !wall(*gameMap, x, y) {
			return headEast
		}
		if x, y := bender.southStep(); !obstacle(*gameMap, x, y) && !wall(*gameMap, x, y) {
			return headSouth
		}
	}
	if x, y := bender.southStep(); !obstacle(*gameMap, x, y) && !wall(*gameMap, x, y) {
		return headSouth
	}
	if x, y := bender.eastStep(); !obstacle(*gameMap, x, y) && !wall(*gameMap, x, y) {
		return headEast
	}
	if x, y := bender.northStep(); !obstacle(*gameMap, x, y) && !wall(*gameMap, x, y) {
		return headNorth
	}
	if x, y := bender.westStep(); !obstacle(*gameMap, x, y) && !wall(*gameMap, x, y) {
		return headWest
	}
	return loop
}

func booth(gameMap GameMap, x, y int) bool {
	if y >= len(gameMap) || x >= len(gameMap[y]) {
		return false
	}
	if gameMap[y][x] == "$" {
		return true
	}
	return false
}

func loop(gameMap *GameMap, position *Bender) stateFn {
	tracks.debug()
	tracks = tracks[:0]
	return nil
}

func obstacle(gameMap GameMap, x, y int) bool {
	if y >= len(gameMap) || x >= len(gameMap[y]) {
		return false
	}

	if gameMap[y][x] == "X" {
		return true
	}
	return false
}

func wall(gameMap GameMap, x, y int) bool {
	if y >= len(gameMap) || x >= len(gameMap[y]) {
		return false
	}

	if gameMap[y][x] == "#" {
		return true
	}
	return false
}

func getMapWithBender(gameMap GameMap, bender Bender) GameMap {
	benderMap := NewGameMap(len(gameMap), len(gameMap[0]))
	for i, x := range gameMap {
		for j := range x {
			benderMap[i][j] = gameMap[i][j]
		}
	}
	benderMap[bender.y][bender.x] = "+"
	return benderMap
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	var L, C int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &L, &C)

	gameMap := NewGameMap(L, C)

	for i := 0; i < L; i++ {
		scanner.Scan()
		row := scanner.Text()
		cols := strings.Split(row, "")
		gameMap[i] = cols
	}

	printGameMap(gameMap)
	bender := locateBender(gameMap)
	startState := handleMove

	for state := startState; state != nil; {
		state = state(&gameMap, &bender)
		// outputMap := getMapWithBender(gameMap, bender)
		// printGameMap(outputMap)
	}
	tracks.print()
}
