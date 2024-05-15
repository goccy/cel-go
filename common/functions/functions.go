// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package functions defines the standard builtin functions supported by the interpreter
package functions

import (
	"context"

	"github.com/google/cel-go/common/types/ref"
)

// Overload defines a named overload of a function, indicating an operand trait
// which must be present on the first argument to the overload as well as one
// of either a unary, binary, or function implementation.
//
// The majority of  operators within the expression language are unary or binary
// and the specializations simplify the call contract for implementers of
// types with operator overloads. Any added complexity is assumed to be handled
// by the generic FunctionOp.
type Overload struct {
	// Operator name as written in an expression or defined within
	// operators.go.
	Operator string

	// Operand trait used to dispatch the call. The zero-value indicates a
	// global function overload or that one of the Unary / Binary / Function
	// definitions should be used to execute the call.
	OperandTrait int

	// Unary defines the overload with a UnaryOp implementation. May be nil.
	Unary UnaryOp

	// Binary defines the overload with a BinaryOp implementation. May be nil.
	Binary BinaryOp

	// Function defines the overload with a FunctionOp implementation. May be
	// nil.
	Function FunctionOp

	// NonStrict specifies whether the Overload will tolerate arguments that
	// are types.Err or types.Unknown.
	NonStrict bool
}

// OverloadContext defines a named overload of a function, indicating an operand trait
// which must be present on the first argument to the overload as well as one
// of either a unary, binary, or function implementation.
//
// The majority of  operators within the expression language are unary or binary
// and the specializations simplify the call contract for implementers of
// types with operator overloads. Any added complexity is assumed to be handled
// by the generic FunctionOp.
type OverloadContext struct {
	// Operator name as written in an expression or defined within
	// operators.go.
	Operator string

	// Operand trait used to dispatch the call. The zero-value indicates a
	// global function overload or that one of the Unary / Binary / Function
	// definitions should be used to execute the call.
	OperandTrait int

	// Unary defines the overload with a UnaryOp implementation. May be nil.
	Unary UnaryOpContext

	// Binary defines the overload with a BinaryOp implementation. May be nil.
	Binary BinaryOpContext

	// Function defines the overload with a FunctionOp implementation. May be
	// nil.
	Function FunctionOpContext

	// NonStrict specifies whether the Overload will tolerate arguments that
	// are types.Err or types.Unknown.
	NonStrict bool
}

// UnaryOp is a function that takes a single value and produces an output.
type UnaryOp func(value ref.Val) ref.Val

// UnaryOpContext is a function that takes a single value and produces an output.
type UnaryOpContext func(ctx context.Context, value ref.Val) ref.Val

// BinaryOp is a function that takes two values and produces an output.
type BinaryOp func(lhs ref.Val, rhs ref.Val) ref.Val

// BinaryOpContext is a function that takes two values and produces an output.
type BinaryOpContext func(ctx context.Context, lhs ref.Val, rhs ref.Val) ref.Val

// FunctionOp is a function with accepts zero or more arguments and produces
// a value or error as a result.
type FunctionOp func(values ...ref.Val) ref.Val

// FunctionOpContext is a function with accepts zero or more arguments and produces
// a value or error as a result.
type FunctionOpContext func(ctx context.Context, values ...ref.Val) ref.Val

// ToOverloadContext convert to OverloadContext.
func (o *Overload) ToOverloadContext() *OverloadContext {
	ret := &OverloadContext{
		Operator:     o.Operator,
		OperandTrait: o.OperandTrait,
		NonStrict:    o.NonStrict,
	}
	if o.Unary != nil {
		ret.Unary = func(_ context.Context, value ref.Val) ref.Val {
			return o.Unary(value)
		}
	}
	if o.Binary != nil {
		ret.Binary = func(_ context.Context, lhs ref.Val, rhs ref.Val) ref.Val {
			return o.Binary(lhs, rhs)
		}
	}
	if o.Function != nil {
		ret.Function = func(_ context.Context, values ...ref.Val) ref.Val {
			return o.Function(values...)
		}
	}
	return ret
}

// ToOverload convert to Overload.
func (o *OverloadContext) ToOverload() *Overload {
	ret := &Overload{
		Operator:     o.Operator,
		OperandTrait: o.OperandTrait,
		NonStrict:    o.NonStrict,
	}
	if o.Unary != nil {
		ret.Unary = func(value ref.Val) ref.Val {
			return o.Unary(context.Background(), value)
		}
	}
	if o.Binary != nil {
		ret.Binary = func(lhs ref.Val, rhs ref.Val) ref.Val {
			return o.Binary(context.Background(), lhs, rhs)
		}
	}
	if o.Function != nil {
		ret.Function = func(values ...ref.Val) ref.Val {
			return o.Function(context.Background(), values...)
		}
	}
	return ret
}
