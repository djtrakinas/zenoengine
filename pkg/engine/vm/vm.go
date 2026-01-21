package vm

import (
	"context"
	"fmt"
	"zeno/pkg/engine"
)

const StackMax = 256

// VM is the bytecode execution engine.
type VM struct {
	chunk *Chunk
	ip    int // Instruction Pointer
	stack [StackMax]Value
	sp    int // Stack Pointer
}

// Chunk stores a sequence of bytecode and constants.
type Chunk struct {
	Code      []byte
	Constants []Value
}

func NewVM() *VM {
	return &VM{}
}

func (vm *VM) push(val Value) {
	vm.stack[vm.sp] = val
	vm.sp++
}

func (vm *VM) pop() Value {
	vm.sp--
	return vm.stack[vm.sp]
}

func (vm *VM) Run(ctx context.Context, chunk *Chunk, scope *engine.Scope) error {
	vm.chunk = chunk
	vm.ip = 0
	vm.sp = 0

	for {
		instruction := OpCode(vm.readByte())
		switch instruction {
		case OpReturn:
			return nil

		case OpConstant:
			constant := vm.readConstant()
			vm.push(constant)

		case OpNil:
			vm.push(NewNil())
		case OpTrue:
			vm.push(NewBool(true))
		case OpFalse:
			vm.push(NewBool(false))

		case OpGetGlobal:
			name := vm.readConstant().AsPtr.(string)
			val, ok := scope.Get(name)
			if ok {
				vm.push(NewObject(val))
			} else {
				vm.push(NewNil())
			}

		case OpSetGlobal:
			name := vm.readConstant().AsPtr.(string)
			val := vm.pop()
			scope.Set(name, val.ToNative())

		case OpAdd:
			b := vm.pop()
			a := vm.pop()
			// Basic numeric add
			vm.push(NewNumber(a.AsNum + b.AsNum))

		case OpSubtract:
			b := vm.pop()
			a := vm.pop()
			vm.push(NewNumber(a.AsNum - b.AsNum))

		case OpCallSlot:
			slotName := vm.readConstant().AsPtr.(string)
			argCount := int(vm.readByte())

			// 1. Get Engine from context
			eng, ok := ctx.Value("engine").(*engine.Engine)
			if !ok {
				return fmt.Errorf("engine not found in context")
			}

			// 2. Lookup Handler
			handler, exists := eng.Registry[slotName]
			if !exists {
				return fmt.Errorf("slot not found: %s", slotName)
			}

			// 3. Prepare Mock Node for arguments (Stack to Node bridge)
			mockNode := &engine.Node{Name: slotName}
			if argCount > 0 {
				mockNode.Children = make([]*engine.Node, argCount)
				// Pop in reverse order because they were pushed in original order
				for i := argCount - 1; i >= 0; i-- {
					val := vm.pop()
					nameVal := vm.pop()
					name := nameVal.AsPtr.(string)

					mockNode.Children[i] = &engine.Node{
						Name:   name,
						Value:  val.ToNative(),
						Parent: mockNode,
					}
				}
			}

			// 4. Execution
			err := handler(ctx, mockNode, scope)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("unsupported opcode: %d", instruction)
		}
	}
}

func (vm *VM) readByte() byte {
	b := vm.chunk.Code[vm.ip]
	vm.ip++
	return b
}

func (vm *VM) readConstant() Value {
	index := vm.readByte()
	return vm.chunk.Constants[index]
}
