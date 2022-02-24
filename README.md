# Tiny Genetic Programming in Golang (based on [1])

A minimalistic program implementing Koza-style (tree-based) genetic programming to solve a symbolic regression problem. 

It is a basic (and fully functional) version, which produces textual output of the evolutionary progression and evolved trees.


| Symbolic | Regression using GP  |
|-------------:|:-------------| 
| Objective | Find an expression with one input (independent variable x), whose output equals the value of the quartic function  x<sup>4</sup> + x<sup>3</sup> + x<sup>2</sup> + x + 1 |
| Function set | add, sub, mul |   
| Terminal set | x, -2, -1, 0, 1, 2  |   
| Fitness | Inverse mean absolute error over a dataset of 101 target values, normalized to [-1,1]
| Paremeters | POP_SIZE (population size), MIN_DEPTH (minimal initial random tree depth), MAX_DEPTH (maximal initial random tree depth), GENERATIONS (maximal number of generations), TOURNAMENT_SIZE (size of tournament for tournament selection), XO_RATE (crossover rate), PROB_MUTATION (per-node mutation probability) |
| Termination | Maximal number of generations reached or an individual with fitness = 1.0 found |

## References
<a id="1">[1]</a> 
[Tiny Genetic Programming in Python](https://github.com/moshesipper/tiny_gp)
