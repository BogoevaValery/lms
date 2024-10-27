package main

import (
	"reflect"
	"testing"
)

// Тесты для функции Calc
func TestCalc(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
		shouldFail bool
	}{
		// Тесты на корректные выражения
		{"3 + 5", 8, false},
		{"10 - 2 * 3", 4, false},
		{"(1 + 2) * 4", 12, false},
		{"10 / 2 + 3", 8, false},
		{"3 + 5 * (2 - 8)", -27, false},
		{"3.5 + 2.5", 6, false}, // Тест для чисел с плавающей точкой

		// Тесты на некорректные выражения
		{"10 / 0", 0, true},     // Деление на ноль
		{"3 + 5 *", 0, true},    // Недостаточно операндов
		{"3 + * 5", 0, true},    // Неправильный порядок операторов
		{"5 + (3 - 2", 0, true}, // Несовпадение скобок
		{"3 + )5 * 2(", 0, true}, // Неверные символы скобок
		{"a + b", 0, true},      // Неподдерживаемые символы
		{"3 + * 5", 0, true},	//два оператора
	}

	for _, tt := range tests {
		result, err := Calc(tt.expression)
		if tt.shouldFail {
			if err == nil {
				t.Errorf("Calc(%s) expected an error, got result %f", tt.expression, result)
			}
		} else {
			if err != nil {
				t.Errorf("Calc(%s) unexpected error: %v", tt.expression, err)
			}
			if result != tt.expected {
				t.Errorf("Calc(%s) expected %f, got %f", tt.expression, tt.expected, result)
			}
		}
	}
}

// Тесты для функции tokenize
func TestTokenize(t *testing.T) {
	tests := []struct {
		input       string
		expected    []string
		shouldError bool
	}{
		// Корректные выражения
		{"3 + 5", []string{"3", "+", "5"}, false},
		{"10 - 2 * 3", []string{"10", "-", "2", "*", "3"}, false},
		{"(1 + 2) * 4", []string{"(", "1", "+", "2", ")", "*", "4"}, false},
		{"10 / 2 + 3", []string{"10", "/", "2", "+", "3"}, false},
		{"3.14 + 2.71", []string{"3.14", "+", "2.71"}, false},

		// Некорректные выражения
		{"3 + a", nil, true}, // Неподдерживаемый символ
		{"3..14 + 2", nil, true}, // Неправильное число с двумя точками
	}

	for _, tt := range tests {
		result, err := tokenize(tt.input)
		if tt.shouldError {
			if err == nil {
				t.Errorf("tokenize(%s) expected an error, got none", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("tokenize(%s) unexpected error: %v", tt.input, err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("tokenize(%s) expected %v, got %v", tt.input, tt.expected, result)
			}
		}
	}
}

// Тесты для функции toRPN
func TestToRPN(t *testing.T) {
	tests := []struct {
		tokens      []string
		expected    []string
		shouldError bool
	}{
		// Корректные выражения
		{[]string{"3", "+", "5"}, []string{"3", "5", "+"}, false},
		{[]string{"10", "-", "2", "*", "3"}, []string{"10", "2", "3", "*", "-"}, false},
		{[]string{"(", "1", "+", "2", ")", "*", "4"}, []string{"1", "2", "+", "4", "*"}, false},
		{[]string{"10", "/", "2", "+", "3"}, []string{"10", "2", "/", "3", "+"}, false},

		// Некорректные выражения
		//{[]string{"3", "+", "*"}, nil, true},           // Лишний оператор
		{[]string{"(", "3", "+", "5"}, nil, true},      // Несовпадающие скобки
		{[]string{"3", "+", "5", ")"}, nil, true},      // Несовпадающие скобки
	}

	for _, tt := range tests {
		result, err := toRPN(tt.tokens)
		if tt.shouldError {
			if err == nil {
				t.Errorf("toRPN(%v) expected an error, got none", tt.tokens)
			}
		} else {
			if err != nil {
				t.Errorf("toRPN(%v) unexpected error: %v", tt.tokens, err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("toRPN(%v) expected %v, got %v", tt.tokens, tt.expected, result)
			}
		}
	}
}

// Тесты для функции evaluateRPN
func TestEvaluateRPN(t *testing.T) {
	tests := []struct {
		rpn         []string
		expected    float64
		shouldError bool
	}{
		// Корректные выражения
		{[]string{"3", "5", "+"}, 8, false},
		{[]string{"10", "2", "3", "*", "-"}, 4, false},
		{[]string{"1", "2", "+", "4", "*"}, 12, false},
		{[]string{"10", "2", "/", "3", "+"}, 8, false},
		{[]string{"3", "5", "2", "8", "-", "*", "+"}, -27, false},

		// Некорректные выражения
		{[]string{"3", "+"}, 0, true},           // Недостаточно операндов
		{[]string{"10", "0", "/"}, 0, true},     // Деление на ноль
		{[]string{"3", "5", "*", "+"}, 0, true}, // Лишний оператор
	}

	for _, tt := range tests {
		result, err := evaluateRPN(tt.rpn)
		if tt.shouldError {
			if err == nil {
				t.Errorf("evaluateRPN(%v) expected an error, got none", tt.rpn)
			}
		} else {
			if err != nil {
				t.Errorf("evaluateRPN(%v) unexpected error: %v", tt.rpn, err)
			}
			if result != tt.expected {
				t.Errorf("evaluateRPN(%v) expected %f, got %f", tt.rpn, tt.expected, result)
			}
		}
	}
}
