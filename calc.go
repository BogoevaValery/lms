package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Calc принимает строку выражения и возвращает результат вычисления или ошибку.
func Calc(expression string) (float64, error) {
	tokens, err := tokenize(expression)
	if err != nil {
		return 0, err
	}
	
	rpn, err := toRPN(tokens)
	if err != nil {
		return 0, err
	}

	result, err := evaluateRPN(rpn)
	if err != nil {
		return 0, err
	}

	return result, nil
}

// tokenize разбивает строку выражения на токены (числа, операторы, скобки)
func tokenize(expression string) ([]string, error) {
	var tokens []string
	var number strings.Builder
	dotCount := 0 // Счетчик точек в числе

	for _, ch := range expression {
		switch {
		case unicode.IsDigit(ch):
			number.WriteRune(ch)

		case ch == '.':
			if dotCount > 0 {
				return nil, errors.New("invalid number format: multiple decimal points")
			}
			dotCount++
			number.WriteRune(ch)

		case ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '(' || ch == ')':
			if number.Len() > 0 {
				tokens = append(tokens, number.String())
				number.Reset()
				dotCount = 0 // Сбросим счетчик точек
			}
			tokens = append(tokens, string(ch))

		case unicode.IsSpace(ch):
			if number.Len() > 0 {
				tokens = append(tokens, number.String())
				number.Reset()
				dotCount = 0
			}

		default:
			return nil, errors.New("invalid character in expression")
		}
	}

	if number.Len() > 0 {
		tokens = append(tokens, number.String())
	}

	return tokens, nil
}

// toRPN преобразует токены в обратную польскую нотацию (RPN) с использованием алгоритма сортировочной станции.
func toRPN(tokens []string) ([]string, error) {
	var output []string
	var stack []string

	precedence := map[string]int{
		"+": 1, "-": 1,
		"*": 2, "/": 2,
	}

	for _, token := range tokens {
		if isNumber(token) {
			output = append(output, token)
		} else if token == "(" {
			stack = append(stack, token)
		} else if token == ")" {
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, errors.New("mismatched parentheses")
			}
			stack = stack[:len(stack)-1]
		} else if prec, isOperator := precedence[token]; isOperator {
			for len(stack) > 0 && precedence[stack[len(stack)-1]] >= prec {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		} else {
			return nil, errors.New("invalid token")
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == "(" {
			return nil, errors.New("mismatched parentheses")
		}
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output, nil
}


// evaluateRPN вычисляет значение выражения, представленное в виде обратной польской нотации (RPN).
func evaluateRPN(rpn []string) (float64, error) {
	var stack []float64

	for _, token := range rpn {
		if isNumber(token) {
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, err
			}
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, errors.New("invalid expression")
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch token {
			case "+":
				stack = append(stack, a+b)
			case "-":
				stack = append(stack, a-b)
			case "*":
				stack = append(stack, a*b)
			case "/":
				if b == 0 {
					return 0, errors.New("division by zero")
				}
				stack = append(stack, a/b)
			default:
				return 0, errors.New("invalid operator")
			}
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("invalid expression")
	}

	return stack[0], nil
}

// isNumber проверяет, является ли строка числом.
func isNumber(token string) bool {
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}


