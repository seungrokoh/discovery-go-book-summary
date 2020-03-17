package problem02

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrWrongExpr = errors.New("it is wrong expr")
	ErrWrongOperator = errors.New("contain wrong operator")
)

type BinOp func(int, int) int

type StrSet map[string]struct{}
type PrecMap map[string]StrSet

func NewStrSet(strs ...string) StrSet {
	m := StrSet{}
	for _, str := range strs {
		m[str] = struct{}{}
	}
	return m
}

func NewEvaluator(opMap map[string]BinOp, prec PrecMap) func(expr string) (int, error) {
	return func(expr string) (int, error) {
		return Eval(opMap, prec, expr)
	}
}

// The expression can have +, -, *, /, (, ) operators and
// decimal integers. Operators and operands should be
func Eval(opMap map[string]BinOp, prec PrecMap, expr string) (int, error) {
	ops := []string{"("}
	var nums []int

	pop := func() int {
		last := nums[len(nums)-1]
		nums = nums[:len(nums)-1]
		return last
	}

	reduce := func(nextOp string) {
		for len(ops) > 0 {
			// 스택의 가장 맨 위에 있는 operator를 꺼냄
			op := ops[len(ops)-1]
			if _, higher := prec[nextOp][op]; nextOp != ")" && !higher {
				// 더 낮은 순위 연산자이므로 여기서 계산 종료
				return
			}
			ops = ops[:len(ops) - 1]
			if op == "(" {
				return
			}
			b, a := pop(), pop()
			if f := opMap[op]; f != nil {
				nums = append(nums, f(a, b))
			}
		}
	}

	// expr의 마지막이 ) 또는 숫자로 끝나는지 확인하는 코드
	lastVal := string(expr[len(expr) - 1])
	if _, err := strconv.Atoi(lastVal); err != nil {
		if lastVal != ")" {
			return 0, ErrWrongExpr
		}
	}

	// 사용자에게 입력받은 expr을 처리하는 for 문
	for _, token := range strings.Fields(expr) {
		if token == "(" {
			ops = append(ops, token)
		} else if _, ok := prec[token]; ok {
			reduce(token)
			ops = append(ops, token)
		} else if token == ")" {
			// 닫는 괄호는 여는 괄호까지 계산하고 제거
			reduce(token)
		} else {
			num, err := strconv.Atoi(token)
			// 숫자가 아닌 잘못된 oper가 들어왔을 때 에러처리
			if err != nil {
				return 0, ErrWrongOperator
			}
			nums = append(nums, num)
		}
	}
	reduce(")")
	if len(nums) > 1 || len(ops) > 0 {
		// 잘못된 수식
		return 0, ErrWrongExpr
	}
	return nums[0], nil
}
