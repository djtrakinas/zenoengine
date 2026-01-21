package vm

import (
	"fmt"
	"strconv"
	"strings"
	"zeno/pkg/engine"
)

type Compiler struct {
	chunk *Chunk
}

func NewCompiler() *Compiler {
	return &Compiler{
		chunk: &Chunk{
			Code:      []byte{},
			Constants: []Value{},
		},
	}
}

func (c *Compiler) Compile(node *engine.Node) (*Chunk, error) {
	err := c.compileNode(node)
	if err != nil {
		return nil, err
	}
	c.emitByte(byte(OpReturn))
	return c.chunk, nil
}

func (c *Compiler) compileNode(node *engine.Node) error {
	// Simple Arithmetic Example: "1 + 2" atau "$x: 10"

	// If it's a variable assignment: $varName: value
	if strings.HasPrefix(node.Name, "$") {
		varName := node.Name[1:]
		// Evaluate value (simplified for now: only constants)
		if err := c.compileValue(node.Value); err != nil {
			return err
		}
		c.emitByte(byte(OpSetGlobal))
		c.emitByte(c.addConstant(NewString(varName)))
		return nil
	}

	// If it's an expression like "1 + 2" (Current Zeno stores this in Value)
	if node.Value != nil {
		valStr := fmt.Sprintf("%v", node.Value)
		parts := strings.Fields(valStr)
		if len(parts) == 3 && parts[1] == "+" {
			// Very basic arithmetic parser
			v1, _ := strconv.ParseFloat(parts[0], 64)
			v2, _ := strconv.ParseFloat(parts[2], 64)

			c.emitByte(byte(OpConstant))
			c.emitByte(c.addConstant(NewNumber(v1)))

			c.emitByte(byte(OpConstant))
			c.emitByte(c.addConstant(NewNumber(v2)))

			c.emitByte(byte(OpAdd))
			return nil
		}
	}

	// If it's a regular slot call (e.g., http.response)
	if node.Name != "" && !strings.HasPrefix(node.Name, "$") && node.Name != "root" {
		// Compile children as named arguments
		for _, child := range node.Children {
			// Push Name
			c.emitByte(byte(OpConstant))
			c.emitByte(c.addConstant(NewString(child.Name)))
			// Push Value
			if err := c.compileValue(child.Value); err != nil {
				return err
			}
		}

		c.emitByte(byte(OpCallSlot))
		c.emitByte(c.addConstant(NewString(node.Name)))
		c.emitByte(byte(len(node.Children))) // Argument count
		return nil
	}

	return nil
}

func (c *Compiler) compileValue(v interface{}) error {
	if s, ok := v.(string); ok {
		// Variable reference?
		if strings.HasPrefix(s, "$") {
			c.emitByte(byte(OpGetGlobal))
			c.emitByte(c.addConstant(NewString(s[1:])))
			return nil
		}
		// Number?
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			c.emitByte(byte(OpConstant))
			c.emitByte(c.addConstant(NewNumber(f)))
			return nil
		}
		// String literal
		c.emitByte(byte(OpConstant))
		c.emitByte(c.addConstant(NewString(s)))
		return nil
	}
	// Fallback raw values
	c.emitByte(byte(OpConstant))
	c.emitByte(c.addConstant(NewObject(v)))
	return nil
}

func (c *Compiler) emitByte(b byte) {
	c.chunk.Code = append(c.chunk.Code, b)
}

func (c *Compiler) addConstant(v Value) byte {
	c.chunk.Constants = append(c.chunk.Constants, v)
	return byte(len(c.chunk.Constants) - 1)
}
