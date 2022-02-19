package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"tinygp/gp"
	"tinygp/mathext"
)

// tiny genetic programming by © patrulek based on moshesipper's code: https://github.com/moshesipper/tiny_gp

func target_func(x float64) float64 { // evolution's target
	return x*x*x*x + x*x*x + x*x + x + 1
}

func generate_dataset() map[float64]float64 { // generate 101 data points from target_func
	dataset := make(map[float64]float64)
	var x float64 = -50
	for ; x <= 50; x++ {
		y := x / 50
		dataset[y] = target_func(y)
	}
	return dataset
}

func init_population() []*gp.GPTree { // ramped half-and-half
	pop := make([]*gp.GPTree, 0, gp.POP_SIZE)
	for md := 3; md < gp.MAX_DEPTH+1; md++ {
		for i := 0; i < int(gp.POP_SIZE/6); i++ {
			t := gp.NewGPTree(nil, nil, nil)
			t.RandomTree(true, md, 0) // grow
			pop = append(pop, t)
		}
		for i := 0; i < int(gp.POP_SIZE/6); i++ {
			t := gp.NewGPTree(nil, nil, nil)
			t.RandomTree(false, md, 0) // full
			pop = append(pop, t)
		}
	}
	return pop
}

func fitness(individual *gp.GPTree, dataset map[float64]float64) float64 { // inverse mean absolute error over dataset normalized to [0,1]
	abs_errors := make([]float64, 0, len(dataset))
	sum := float64(0.0)
	for k, v := range dataset {
		eval := individual.ComputeTree(k)
		err := float64(math.Abs(float64(eval - v)))
		abs_errors = append(abs_errors, err)
		sum += err
	}
	return 1.0 / (1.0 + sum/float64(len(abs_errors)))
}

func selection(population []*gp.GPTree, fitnesses []float64) *gp.GPTree { // select one individual using tournament selection
	tournament := make([]int, 0, gp.TOURNAMENT_SIZE)
	for i := 0; i < gp.TOURNAMENT_SIZE; i++ {
		tournament = append(tournament, rand.Intn(len(population)))
	}

	tournament_fitnesses := make([]float64, 0, gp.TOURNAMENT_SIZE)
	for i := 0; i < gp.TOURNAMENT_SIZE; i++ {
		tournament_fitnesses = append(tournament_fitnesses, fitnesses[tournament[i]])
	}

	i, _ := mathext.Argmax(tournament_fitnesses)
	return population[tournament[i]].Clone()
}

func main() {
	// init stuff
	rand.Seed(time.Now().UnixNano())
	dataset := generate_dataset()
	population := init_population()
	var (
		best_of_run     *gp.GPTree
		best_of_run_f   float64
		best_of_run_gen int
	)
	fitnesses := make([]float64, 0, gp.POP_SIZE)
	for i := 0; i < gp.POP_SIZE; i++ {
		fitnesses = append(fitnesses, fitness(population[i], dataset))
	}

	// go evolution!
	for gen := 0; gen < gp.GENERATIONS; gen++ {
		nextgen_population := make([]*gp.GPTree, 0, gp.POP_SIZE)

		for i := 0; i < gp.POP_SIZE; i++ {
			parent1 := selection(population, fitnesses)
			parent2 := selection(population, fitnesses)
			parent1.Crossover(parent2)
			parent1.Mutation()
			nextgen_population = append(nextgen_population, parent1)
		}
		population = nextgen_population

		for i := 0; i < gp.POP_SIZE; i++ {
			fitnesses[i] = fitness(population[i], dataset)
		}

		i, v := mathext.Argmax(fitnesses)
		if v > best_of_run_f {
			best_of_run_f = v
			best_of_run_gen = gen
			best_of_run = population[i].Clone()
			fmt.Println("________________________")
			fmt.Printf("gen: %v, best_of_run_f: %v, best_of_run:\n", gen, mathext.RoundTo(float64(v), 3))
			best_of_run.PrintTree("")
		}
		if best_of_run_f == 1 {
			break
		}
	}

	fmt.Printf("\n\n_________________________________________________\nEND OF RUN\nbest_of_run attained at gen %v"+
		" and has f=%v\n\n", best_of_run_gen, mathext.RoundTo(float64(best_of_run_f), 3))
	best_of_run.PrintTree("")
}
