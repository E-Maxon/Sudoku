package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Sudoku struct {
	n           int
	grid        [][]int
	difficult   int
	was_edited  bool
	empty_cells int
	prediction  [][][]int
	solution    [][]int
	double      [][]bool
	result      [][]int
}

func (s *Sudoku) make_base_grid(n int) { // создаем базовую сетку
	s.n = n

	grid := make([][]int, s.n*s.n)
	for it := range grid {
		grid[it] = make([]int, s.n*s.n)
	}

	s.grid = grid

	for i := 0; i < n*n; i++ {
		for j := 0; j < n*n; j++ {
			s.grid[i][j] = (i*n+i/n+j)%(n*n) + 1
		}
	}
}

func (s *Sudoku) transpose_table() { // транспонируем сетку
	for i := 0; i < s.n*s.n; i++ {
		for j := i; j < s.n*s.n; j++ {
			s.grid[i][j], s.grid[j][i] = s.grid[j][i], s.grid[i][j]
		}
	}
}

func (s *Sudoku) swap_rows_in_area() { // меняем строки
	area := rand.Intn(s.n)

	x := rand.Intn(s.n)
	row1 := area*s.n + x

	y := rand.Intn(s.n)
	for ; x == y; y = rand.Intn(s.n) {
	}
	row2 := area*s.n + y

	s.grid[row1], s.grid[row2] = s.grid[row2], s.grid[row1]
}

func (s *Sudoku) swap_columns_in_area() { // меняем столбцы
	s.transpose_table()
	s.swap_rows_in_area()
	s.transpose_table()
}

func (s *Sudoku) swap_big_rows() { // меняем большие строки
	area1 := rand.Intn(s.n)
	area2 := rand.Intn(s.n)
	for ; area1 == area2; area2 = rand.Intn(s.n) {
	}

	for i := 0; i < s.n; i++ {
		s.grid[area1*s.n+i], s.grid[area2*s.n+i] = s.grid[area2*s.n+i], s.grid[area1*s.n+i]
	}
}

func (s *Sudoku) swap_big_columns() { // меняем большие столбцы
	s.transpose_table()
	s.swap_big_rows()
	s.transpose_table()
}

func (s *Sudoku) mix(cnt int) {
	for i := 0; i < cnt; i++ {
		mode := rand.Intn(5)
		if mode == 0 {
			s.transpose_table()
		} else if mode == 1 {
			s.swap_rows_in_area()
		} else if mode == 2 {
			s.swap_columns_in_area()
		} else if mode == 3 {
			s.swap_big_rows()
		} else {
			s.swap_big_columns()
		}
	}
}

func (s *Sudoku) find_least_prediction_cell() (int, int) {
	x := 0
	y := 0
	for i := 0; i < s.n*s.n; i++ {
		for j := 0; j < s.n*s.n; j++ {
			if s.solution[i][j] == 0 && len(s.prediction[i][j]) < len(s.prediction[x][y]) {
				x = i
				y = j
			}
		}
	}
	return x, y
}

func (s *Sudoku) fill_prediction() {
	s.prediction = make([][][]int, s.n*s.n)
	for i := 0; i < s.n*s.n; i++ {
		s.prediction[i] = make([][]int, s.n*s.n)
		for j := 0; j < s.n*s.n; j++ {
			if s.grid[i][j] != 0 {
				s.prediction[i][j] = append(s.prediction[i][j], s.grid[i][j])
			} else {
				for p := 1; p < 10; p++ {
					s.prediction[i][j] = append(s.prediction[i][j], p)
				}
			}
		}
	}
}

func (s *Sudoku) delete_prediction_cell(i int, j int, p int) {
	for ind := range s.prediction[i][j] {
		if s.prediction[i][j][ind] == p {
			s.prediction[i][j] = append(s.prediction[i][j][0:ind], s.prediction[i][j][ind+1:]...)
			s.was_edited = true
			break
		}
	}
}

func (s *Sudoku) delete_prediction_row(i int, p int) {
	for ind := 0; ind < s.n*s.n; ind++ {
		if s.solution[i][ind] == 0 && !s.double[i][ind] {
			s.delete_prediction_cell(i, ind, p)
		}
	}
}

func (s *Sudoku) delete_prediction_column(j int, p int) {
	for ind := 0; ind < s.n*s.n; ind++ {
		if s.solution[ind][j] == 0 && !s.double[ind][j] {
			s.delete_prediction_cell(ind, j, p)
		}
	}
}

func (s *Sudoku) delete_prediction_area(i int, j int, p int) {
	i_start := s.n * (i / s.n)
	j_start := s.n * (j / s.n)
	for ind1 := i_start; ind1 < i_start+s.n; ind1++ {
		for ind2 := j_start; ind2 < j_start+s.n; ind2++ {
			if s.solution[ind1][ind2] == 0 && !s.double[ind1][ind2] {
				s.delete_prediction_cell(ind1, ind2, p)
			}
		}
	}
}

func (s *Sudoku) delete_prediction(i int, j int) {
	s.delete_prediction_row(i, s.solution[i][j])
	s.delete_prediction_column(j, s.solution[i][j])
	s.delete_prediction_area(i, j, s.solution[i][j])
}

func (s *Sudoku) calc_prediction() {
	for i := 0; i < s.n*s.n; i++ {
		for j := 0; j < s.n*s.n; j++ {
			if s.solution[i][j] != 0 {
				s.delete_prediction(i, j)
			}
		}
	}
}

func (s *Sudoku) fill_cells() {
	for i := 0; i < s.n*s.n; i++ {
		for j := 0; j < s.n*s.n; j++ {
			if s.solution[i][j] == 0 && len(s.prediction[i][j]) == 1 {
				s.solution[i][j] = s.prediction[i][j][0]
				s.delete_prediction(i, j)
				s.was_edited = true
				s.empty_cells--
			}
		}
	}
}

func (s *Sudoku) fill_double() {
	s.double = make([][]bool, s.n*s.n)
	for i := 0; i < s.n*s.n; i++ {
		s.double[i] = make([]bool, s.n*s.n)
	}
}

func (s *Sudoku) fill_solution() {
	s.solution = make([][]int, s.n*s.n)
	for i := 0; i < s.n*s.n; i++ {
		s.solution[i] = make([]int, s.n*s.n)
		for j := 0; j < s.n*s.n; j++ {
			s.solution[i][j] = s.grid[i][j]
		}
	}
}

func (s *Sudoku) is_similar_prediction(i1 int, j1 int, i2 int, j2 int) bool {
	if len(s.prediction[i1][j1]) != len(s.prediction[i2][j2]) {
		return false
	}
	for ind := 0; ind < len(s.prediction[i1][j1]); ind++ {
		if s.prediction[i1][j1][ind] != s.prediction[i2][j2][ind] {
			return false
		}
	}
	return true
}

func (s *Sudoku) find_doubles_row(i int, j int) {
	cnt := 0

	for ind := 0; ind < s.n*s.n; ind++ {
		if s.is_similar_prediction(i, ind, i, j) {
			s.double[i][ind] = true
			cnt++
		}
	}

	if cnt == len(s.prediction[i][j]) {
		for _, p := range s.prediction[i][j] {
			s.delete_prediction_row(i, p)
		}
	}
}

func (s *Sudoku) find_doubles_column(i int, j int) {
	cnt := 0

	for ind := 0; ind < s.n*s.n; ind++ {
		if s.is_similar_prediction(ind, j, i, j) {
			s.double[ind][j] = true
			cnt++
		}
	}

	if cnt == len(s.prediction[i][j]) {
		for _, p := range s.prediction[i][j] {
			s.delete_prediction_column(j, p)
		}
	}
}

func (s *Sudoku) find_doubles_area(i int, j int) {
	cnt := 0

	i_start := s.n * (i / s.n)
	j_start := s.n * (j / s.n)
	for ind1 := i_start; ind1 < i_start+s.n; ind1++ {
		for ind2 := j_start; ind2 < j_start+s.n; ind2++ {
			if s.is_similar_prediction(ind1, ind2, i, j) {
				s.double[ind1][ind2] = true
				cnt++
			}
		}
	}

	if cnt == len(s.prediction[i][j]) {
		for _, p := range s.prediction[i][j] {
			s.delete_prediction_area(i, j, p)
		}
	}
}

func (s *Sudoku) zero() {
	for i := 0; i < s.n*s.n; i++ {
		for j := 0; j < s.n*s.n; j++ {
			s.double[i][j] = false
		}
	}
}

func (s *Sudoku) find_doubles() {
	for i := range s.prediction {
		for j := range s.prediction[i] {
			if s.solution[i][j] == 0 {
				s.find_doubles_row(i, j)
				s.zero()
				s.find_doubles_column(i, j)
				s.zero()
				s.find_doubles_area(i, j)
				s.zero()
			}
		}
	}
}

func (s *Sudoku) find_same_prediction_cell(i int, j int, p int) bool {
	for _, cell_prediction := range s.prediction[i][j] {
		if cell_prediction == p {
			return true
		}
	}
	return false
}

func (s *Sudoku) find_same_prediction_row(i int, p int) bool {
	was_prediction := false
	for ind := 0; ind < s.n*s.n; ind++ {
		if s.find_same_prediction_cell(i, ind, p) {
			if was_prediction {
				return true
			}
			was_prediction = true
		}
	}
	return false
}

func (s *Sudoku) find_same_prediction_column(j int, p int) bool {
	was_prediction := false
	for ind := 0; ind < s.n*s.n; ind++ {
		if s.find_same_prediction_cell(ind, j, p) {
			if was_prediction {
				return true
			}
			was_prediction = true
		}
	}
	return false
}

func (s *Sudoku) find_same_prediction_area(i int, j int, p int) bool {
	i_start := s.n * (i / s.n)
	j_start := s.n * (j / s.n)
	was_prediction := false
	for ind1 := i_start; ind1 < i_start+s.n; ind1++ {
		for ind2 := j_start; ind2 < j_start+s.n; ind2++ {
			if s.find_same_prediction_cell(ind1, ind2, p) {
				if was_prediction {
					return true
				}
				was_prediction = true
			}
		}
	}
	return false
}

func (s *Sudoku) find_same_prediction() {
	for i := 0; i < s.n*s.n; i++ {
		for j := 0; j < s.n*s.n; j++ {
			for _, p := range s.prediction[i][j] {
				if s.solution[i][j] == 0 && (!s.find_same_prediction_row(i, p) || !s.find_same_prediction_column(j, p) || !s.find_same_prediction_area(i, j, p)) {
					s.solution[i][j] = p
					s.was_edited = true
					s.empty_cells--
					s.delete_prediction(i, j)
					return
				}
			}
		}
	}
}

func (s *Sudoku) solve() bool {
	s.fill_solution()
	s.was_edited = true
	has_solution := false
	for s.empty_cells > 0 && s.was_edited {
		s.was_edited = false
		s.fill_cells()
		s.find_doubles()
		s.find_same_prediction()
		if !s.was_edited {
			i, j := s.find_least_prediction_cell()
			s.empty_cells--
			for num := range s.prediction[i][j] {
				s.solution[i][j] = num
				if s.solve() {
					if has_solution {
						return false
					}
					has_solution = true
				} else {
					s.was_edited = false
					s.calc_prediction()
				}
			}
			if !s.was_edited {
				s.solution[i][j] = 0
				s.empty_cells++
				return false
			}
		}
	}
	return s.empty_cells == 0
}

func (s *Sudoku) delete_cells() {
	used := make([][]bool, s.n*s.n)
	for it := 0; it < s.n*s.n; it++ {
		used[it] = make([]bool, s.n*s.n)
	}

	s.difficult = 0

	for it := 0; it < s.n*s.n*s.n*s.n; {
		i := rand.Intn(s.n * s.n)
		j := rand.Intn(s.n * s.n)
		if !used[i][j] {
			it++
			used[i][j] = true
			s.difficult++

			num := s.grid[i][j]
			s.grid[i][j] = 0
			s.empty_cells = s.difficult
			s.fill_prediction()
			s.fill_double()
			s.fill_solution()
			s.calc_prediction()
			if !s.solve() {
				s.grid[i][j] = num
				s.difficult--
			}
		}

	}
}

func (s *Sudoku) fill_result() {
	s.result = make([][]int, s.n*s.n)
	for i := 0; i < s.n*s.n; i++ {
		s.result[i] = make([]int, s.n*s.n)
		for j := 0; j < s.n*s.n; j++ {
			s.result[i][j] = s.grid[i][j]
		}
	}
}

func (s *Sudoku) print() {
	for i := 0; i < s.n*s.n; i++ {
		if i%s.n == 0 && i != 0 {
			fmt.Println("------+-------+------")
		}

		for j := 0; j < s.n*s.n; j++ {
			if j%s.n == 0 && j != 0 {
				fmt.Print("| ")
			}

			if s.grid[i][j] == 0 {
				fmt.Print("  ")
			} else {
				fmt.Printf("%d ", s.grid[i][j])
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func (s *Sudoku) print_result() {
	for i := 0; i < s.n*s.n; i++ {
		if i%s.n == 0 && i != 0 {
			fmt.Println("------+-------+------")
		}

		for j := 0; j < s.n*s.n; j++ {
			if j%s.n == 0 && j != 0 {
				fmt.Print("| ")
			}

			fmt.Printf("%d ", s.result[i][j])
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var s Sudoku
	s.make_base_grid(3)
	s.mix(10)
	s.fill_result()
	s.delete_cells()
	s.print()
	s.print_result()
}
