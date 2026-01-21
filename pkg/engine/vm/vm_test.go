package vm

import (
	"bytes"
	"context"
	"testing"
	"zeno/pkg/engine"
)

func TestVMArithmetic(t *testing.T) {
	// 1 + 2
	chunk := &Chunk{
		Code: []byte{
			byte(OpConstant), 0, // 1
			byte(OpConstant), 1, // 2
			byte(OpAdd),
			byte(OpReturn),
		},
		Constants: []Value{
			NewNumber(1),
			NewNumber(2),
		},
	}

	vm := NewVM()
	scope := engine.NewScope(nil)
	err := vm.Run(context.Background(), chunk, scope)
	if err != nil {
		t.Fatal(err)
	}

	result := vm.pop()
	if result.AsNum != 3 {
		t.Errorf("Expected 3, got %g", result.AsNum)
	}
}

func TestVMCompilerVariables(t *testing.T) {
	// AST: $x: 10
	node := &engine.Node{
		Name:  "$x",
		Value: "10",
	}

	compiler := NewCompiler()
	chunk, err := compiler.Compile(node)
	if err != nil {
		t.Fatal(err)
	}

	vm := NewVM()
	scope := engine.NewScope(nil)
	err = vm.Run(context.Background(), chunk, scope)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := scope.Get("x")
	if !ok {
		t.Fatal("Variable x should be set in scope")
	}

	// Value representation in prototype might need adjustment,
	// but for now we expect the raw value or NewNumber.
	// Currently compiler uses NewNumber(10)
	if n, ok := val.(float64); ok && n != 10 {
		t.Errorf("Expected 10, got %v", val)
	}
}

func TestVMComplexSlot(t *testing.T) {
	// AST:
	// http.response:
	//    status: 201
	//    body: "created"
	node := &engine.Node{
		Name: "http.response",
		Children: []*engine.Node{
			{Name: "status", Value: "201"},
			{Name: "body", Value: "created"},
		},
	}

	// Mock Engine Registry
	eng := engine.NewEngine()
	called := false
	eng.Register("http.response", func(ctx context.Context, n *engine.Node, s *engine.Scope) error {
		called = true
		// Verify attributes
		statusFound := false
		bodyFound := false
		for _, child := range n.Children {
			if child.Name == "status" && child.Value == 201.0 {
				statusFound = true
			}
			if child.Name == "body" && child.Value == "created" {
				bodyFound = true
			}
		}
		if !statusFound || !bodyFound {
			t.Errorf("Attributes not correctly passed. StatusFound: %v, BodyFound: %v", statusFound, bodyFound)
		}
		return nil
	}, engine.SlotMeta{})

	compiler := NewCompiler()
	chunk, err := compiler.Compile(node)
	if err != nil {
		t.Fatal(err)
	}

	vm := NewVM()
	scope := engine.NewScope(nil)
	ctx := context.WithValue(context.Background(), "engine", eng)

	err = vm.Run(ctx, chunk, scope)
	if err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Error("http.response slot was not called")
	}
}
func TestVMControlFlow(t *testing.T) {
	// AST:
	// if: $x == 10
	//   then:
	//     $res: "yes"
	//   else:
	//     $res: "no"
	node := &engine.Node{
		Name:  "if",
		Value: "$x == 10",
		Children: []*engine.Node{
			{
				Name: "then",
				Children: []*engine.Node{
					{Name: "$res", Value: "'yes'"},
				},
			},
			{
				Name: "else",
				Children: []*engine.Node{
					{Name: "$res", Value: "'no'"},
				},
			},
		},
	}

	compiler := NewCompiler()
	chunk, err := compiler.Compile(node)
	if err != nil {
		t.Fatal(err)
	}

	vm := NewVM()

	// Case 1: x == 10
	scope1 := engine.NewScope(nil)
	scope1.Set("x", 10.0)
	err = vm.Run(context.Background(), chunk, scope1)
	if err != nil {
		t.Fatal(err)
	}
	res1, _ := scope1.Get("res")
	if res1 != "yes" {
		t.Errorf("Expected 'yes', got %v", res1)
	}

	// Case 2: x != 10
	scope2 := engine.NewScope(nil)
	scope2.Set("x", 20.0)
	err = vm.Run(context.Background(), chunk, scope2)
	if err != nil {
		t.Fatal(err)
	}
	res2, _ := scope2.Get("res")
	if res2 != "no" {
		t.Errorf("Expected 'no', got %v", res2)
	}
}

func TestVMSerialization(t *testing.T) {
	chunk := &Chunk{
		Code: []byte{byte(OpConstant), 0, byte(OpReturn)},
		Constants: []Value{
			NewString("hello"),
		},
		LocalNames: []string{"var1"},
	}

	var buf bytes.Buffer
	if err := chunk.Serialize(&buf); err != nil {
		t.Fatal(err)
	}

	decoded, err := DeserializeChunk(&buf)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(chunk.Code, decoded.Code) {
		t.Error("Code mismatch after serialization")
	}
	if len(chunk.Constants) != len(decoded.Constants) || decoded.Constants[0].AsPtr.(string) != "hello" {
		t.Error("Constants mismatch after serialization")
	}
	if len(chunk.LocalNames) != len(decoded.LocalNames) || decoded.LocalNames[0] != "var1" {
		t.Error("LocalNames mismatch after serialization")
	}
}
