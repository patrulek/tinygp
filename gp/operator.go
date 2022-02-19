package gp

type Operator interface {
	Eval(...interface{}) interface{}
	String() string
}

type Mul struct{}

func (this *Mul) Eval(args ...interface{}) interface{} {
	var sum float64 = 1
	for _, v := range args {
		switch v.(type) {
		case float64:
			sum *= v.(float64)
		}
	}
	return sum
}

func (this *Mul) String() string {
	return "Mul"
}

type Add struct{}

func (this *Add) Eval(args ...interface{}) interface{} {
	var sum float64 = 0
	for _, v := range args {
		switch v.(type) {
		case float64:
			sum += v.(float64)
		}
	}
	return sum
}

func (this *Add) String() string {
	return "Add"
}

type Sub struct{}

func (this *Sub) Eval(args ...interface{}) interface{} {
	var sum float64 = args[0].(float64)
	for _, v := range args {
		switch v.(type) {
		case float64:
			sum -= v.(float64)
		}
	}

	return sum + args[0].(float64)
}

func (this *Sub) String() string {
	return "Sub"
}
