package main

import (
	"context"
	"fmt"

	"github.com/open-policy-agent/opa/rego"
)

func main() {
	ctx := context.TODO()

	query, err := rego.New(
		rego.Query("aa = data.example.authz.allow"),
		rego.Load([]string{"./example.rego"}, nil),
		// rego.Module("example.rego", module),
	).PrepareForEval(ctx)

	if err != nil {
		panic("damn")
	}

	input := map[string]interface{}{
		"method": "GET",
		"path":   []interface{}{"salary", "bob"},
		"subject": map[string]interface{}{
			"user":   "bob",
			"groups": []interface{}{"sales", "marketing"},
		},
		"expiry_year": 2050,
	}

	type hehe struct {
		User   string   `json:"user"`
		Groups []string `json:"groups"`
	}

	type hoho struct {
		Subject hehe `json:"subject"`
	}

	zz := hoho{
		Subject: hehe{
			User:   "admin",
			Groups: []string{"testing"},
		},
	}

	results, err := query.Eval(ctx, rego.EvalInput(input))
	fmt.Printf("%+v\n", results)

	results, err = query.Eval(ctx, rego.EvalInput(zz))
	fmt.Printf("%+v\n", results)
}
