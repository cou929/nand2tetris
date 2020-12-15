package main

import "fmt"

type VmCode struct{}

func NewVmCode() *VmCode {
	return &VmCode{}
}

func (v *VmCode) function(class string, fn string, localC int) string {
	return fmt.Sprintf("function %s.%s %d", class, fn, localC)
}

func (v *VmCode) call(class string, fn string, argC int) string {
	return fmt.Sprintf("call %s.%s %d", class, fn, argC)
}

func (v *VmCode) pushConstant(in int) string {
	return fmt.Sprintf("push constant %d", in)
}

func (v *VmCode) add() string {
	return "add"
}

func (v *VmCode) sub() string {
	return "sub"
}

func (v *VmCode) mul() string {
	return v.call("Math", "multiply", 2)
}

func (v *VmCode) div() string {
	return v.call("Math", "divide", 2)
}

func (v *VmCode) neg() string {
	return "neg"
}

func (v *VmCode) eq() string {
	return "eq"
}

func (v *VmCode) gt() string {
	return "gt"
}

func (v *VmCode) lt() string {
	return "lt"
}

func (v *VmCode) and() string {
	return "and"
}

func (v *VmCode) or() string {
	return "or"
}

func (v *VmCode) not() string {
	return "not"
}

func (v *VmCode) ret() string {
	return "return"
}
