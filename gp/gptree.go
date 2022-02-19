package gp

import (
	"fmt"
	"math/rand"
)

var (
	POP_SIZE        = 60  // population Size
	MIN_DEPTH       = 2   // minimal initial random tree depth
	MAX_DEPTH       = 5   // maximal initial random tree depth
	GENERATIONS     = 500 // maximal number of generations to run evolution
	TOURNAMENT_SIZE = 5   // Size of tournament for tournament selection
	XO_RATE         = 0.8 // Crossover rate
	PROB_MUTATION   = 0.2 // per-node Mutation probability

	FUNCTIONS = []Operator{&Add{}, &Sub{}, &Mul{}}
	TERMINALS = []interface{}{"x", -2.0, -1.0, 0.0, 1.0, 2.0}
)

type GPTree struct {
	data  interface{}
	left  *GPTree
	right *GPTree
}

func NewGPTree(data interface{}, left *GPTree, right *GPTree) *GPTree {
	return &GPTree{data, left, right}
}

func (this *GPTree) Clone() *GPTree {
	if this == nil {
		return nil
	}

	cloned := NewGPTree(this.data, nil, nil)
	if this.left != nil {
		cloned.left = this.left.Clone()
	}
	if this.right != nil {
		cloned.right = this.right.Clone()
	}

	return cloned
}

func (this *GPTree) NodeLabel() string {
	switch this.data.(type) {
	case Operator:
		return this.data.(Operator).String()
	case float64:
		return fmt.Sprintf("%.0f", this.data)
	default:
		return this.data.(string)
	}
}

func (this *GPTree) PrintTree(prefix string) { // textual printout
	fmt.Printf("%s%s\n", prefix, this.NodeLabel())
	if this.left != nil {
		this.left.PrintTree(prefix + "  ")
	}
	if this.right != nil {
		this.right.PrintTree(prefix + "  ")
	}
}

func (this *GPTree) ComputeTree(x float64) float64 {
	switch this.data.(type) {
	case Operator:
		op := this.data.(Operator)
		result := op.Eval(this.left.ComputeTree(x), this.right.ComputeTree(x))
		return result.(float64)
	case string:
		return x
	default:
		return this.data.(float64)
	}
}

func (this *GPTree) RandomTree(grow bool, max_depth int, depth int) { // create random tree using either grow or full method
	if depth < MIN_DEPTH || (depth < max_depth && !grow) {
		this.data = FUNCTIONS[rand.Intn(len(FUNCTIONS))]
	} else if depth >= max_depth {
		this.data = TERMINALS[rand.Intn(len(TERMINALS))]
	} else { // # intermediate depth, grow
		if rand.Float64() > 0.5 {
			this.data = TERMINALS[rand.Intn(len(TERMINALS))]
		} else {
			this.data = FUNCTIONS[rand.Intn(len(FUNCTIONS))]
		}
	}

	_, ok := this.data.(Operator)
	if ok {
		this.left = NewGPTree(nil, nil, nil)
		this.left.RandomTree(grow, max_depth, depth+1)
		this.right = NewGPTree(nil, nil, nil)
		this.right.RandomTree(grow, max_depth, depth+1)
	}
}

func (this *GPTree) Mutation() {
	if rand.Float64() < PROB_MUTATION { // mutate at this node
		this.RandomTree(true, 2, 0)
	} else if this.left != nil {
		this.left.Mutation()
	} else if this.right != nil {
		this.right.Mutation()
	}
}

func (this *GPTree) Size() int { // tree Size in nodes
	switch this.data.(type) {
	case string, int:
		return 1
	}

	l, r := 0, 0
	if this.left != nil {
		l = this.left.Size()
	}
	if this.right != nil {
		r = this.right.Size()
	}
	return 1 + l + r
}

func (this *GPTree) BuildSubtree() *GPTree { // count is list in order to pass "by reference"
	t := NewGPTree(nil, nil, nil)
	t.data = this.data
	if this.left != nil {
		t.left = this.left.BuildSubtree()
	}
	if this.right != nil {
		t.right = this.right.BuildSubtree()
	}
	return t
}

func (this *GPTree) ScanTree(count *int, second *GPTree) *GPTree { // note: count is list, so it's passed "by reference"
	*count -= 1
	if *count <= 1 {
		if second == nil { // return subtree rooted here
			return this.BuildSubtree()
		} else { // glue subtree here
			this.data = second.data
			this.left = second.left
			this.right = second.right
			return nil
		}
	}

	var ret *GPTree
	if this.left != nil && *count > 1 {
		ret = this.left.ScanTree(count, second)
	}
	if this.right != nil && *count > 1 {
		ret = this.right.ScanTree(count, second)
	}
	return ret
}

func (this *GPTree) Crossover(other *GPTree) { // xo 2 trees at random nodes
	if rand.Float64() < XO_RATE {
		count := 1 + rand.Intn(other.Size())
		// 2nd random subtree
		second := other.ScanTree(&count, nil)

		// 2nd subtree "glued" inside 1st tree
		count = 1 + rand.Intn(this.Size())
		_ = this.ScanTree(&count, second)
	}
}

// end of class
