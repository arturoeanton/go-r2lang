package r2core

import (
	"fmt"
	"strings"
)

// ============================================================
// 3) AST - Node interface
// ============================================================

// PositionInfo holds location information for error reporting
type PositionInfo struct {
	Line     int
	Col      int
	Pos      int
	Filename string
}

// StackFrame represents a single frame in the R2Lang call stack
type StackFrame struct {
	FunctionName string
	Position     *PositionInfo
	Arguments    []interface{} // Optional: for debugging
}

// CallStack represents the complete R2Lang call stack
type CallStack struct {
	Frames []StackFrame
}

// Node interface for all AST nodes
type Node interface {
	Eval(env *Environment) interface{}
}

// PositionedNode interface for nodes that can report their position
// This is optional to maintain backward compatibility
type PositionedNode interface {
	Node
	GetPosition() *PositionInfo
}

// BaseNode provides default position tracking for AST nodes
type BaseNode struct {
	Position *PositionInfo
}

// GetPosition returns the position information
func (b *BaseNode) GetPosition() *PositionInfo {
	return b.Position
}

// CreatePositionInfo creates position info from a token
func CreatePositionInfo(token Token, filename string) *PositionInfo {
	return &PositionInfo{
		Line:     token.Line,
		Col:      token.Col,
		Pos:      token.Pos,
		Filename: filename,
	}
}

// CreatePositionError creates a formatted error message with position
func CreatePositionError(pos *PositionInfo, message string) string {
	if pos == nil {
		return message
	}
	if pos.Filename != "" {
		return fmt.Sprintf("%s:%d:%d: %s", pos.Filename, pos.Line, pos.Col, message)
	}
	return fmt.Sprintf("line %d:%d: %s", pos.Line, pos.Col, message)
}

// PanicWithPosition creates a panic with position information for VSCode linking
func PanicWithPosition(pos *PositionInfo, message string) {
	panic(CreatePositionError(pos, message))
}

// PanicWithStack creates a panic with position and call stack information
func PanicWithStack(pos *PositionInfo, message string, callStack *CallStack) {
	errorMsg := CreatePositionError(pos, message)
	if callStack != nil {
		stackTrace := callStack.FormatStackTrace()
		if stackTrace != "" {
			errorMsg += stackTrace
		}
	}
	panic(errorMsg)
}

// GetNodePosition extracts position from a node if it implements PositionedNode
func GetNodePosition(node Node) *PositionInfo {
	if posNode, ok := node.(PositionedNode); ok {
		return posNode.GetPosition()
	}
	return nil
}

// FormatStackTrace creates a formatted string representation of the call stack
func (cs *CallStack) FormatStackTrace() string {
	if len(cs.Frames) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("\nR2Lang call stack (most recent call first):\n")

	for i := len(cs.Frames) - 1; i >= 0; i-- {
		frame := cs.Frames[i]
		if frame.Position != nil && frame.Position.Filename != "" {
			sb.WriteString(fmt.Sprintf("    %s:%d:%d in %s()\n",
				frame.Position.Filename, frame.Position.Line, frame.Position.Col, frame.FunctionName))
		} else {
			sb.WriteString(fmt.Sprintf("    in %s()\n", frame.FunctionName))
		}
	}

	return sb.String()
}

// PushFrame adds a new frame to the call stack
func (cs *CallStack) PushFrame(functionName string, pos *PositionInfo, args []interface{}) {
	frame := StackFrame{
		FunctionName: functionName,
		Position:     pos,
		Arguments:    args,
	}
	cs.Frames = append(cs.Frames, frame)
}

// PopFrame removes the top frame from the call stack
func (cs *CallStack) PopFrame() {
	if len(cs.Frames) > 0 {
		cs.Frames = cs.Frames[:len(cs.Frames)-1]
	}
}

// Clone creates a copy of the call stack
func (cs *CallStack) Clone() *CallStack {
	clone := &CallStack{
		Frames: make([]StackFrame, len(cs.Frames)),
	}
	copy(clone.Frames, cs.Frames)
	return clone
}
