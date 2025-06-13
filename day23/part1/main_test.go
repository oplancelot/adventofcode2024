package main

import (
	"strings"
	"testing"
)

func TestSolve(t *testing.T) {
	// 谜题中提供的示例输入
	exampleInput := `
kh-tc
qp-kh
de-cg
ka-co
yn-aq
qp-ub
cg-tb
vc-aq
tb-ka
wh-tc
yn-cg
kh-ub
ta-co
de-co
tc-td
tb-wq
wh-td
ta-ka
td-qp
aq-cg
wq-ub
ub-vc
de-ta
wq-aq
wq-vc
wh-yn
ka-de
kh-ta
co-tc
wh-qp
tb-vc
td-yn
`
	// 使用 Table-Driven Tests
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "example",
			input: strings.TrimSpace(exampleInput),
			want:  7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := solve(tt.input); got != tt.want {
				t.Errorf("solve() = %v, want %v", got, tt.want)
			}
		})
	}
}
