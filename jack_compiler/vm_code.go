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

func (v *VmCode) push(seg string, idx int) string {
	return fmt.Sprintf("push %s %d", seg, idx)
}

func (v *VmCode) pushConstant(in int) string {
	return fmt.Sprintf("push constant %d", in)
}

func (v *VmCode) pop(seg string, idx int) string {
	return fmt.Sprintf("pop %s %d", seg, idx)
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

func (v *VmCode) true() []string {
	return []string{v.pushConstant(1), v.neg()}
}

func (v *VmCode) false() string {
	return v.pushConstant(0)
}

func (v *VmCode) null() string {
	return v.pushConstant(0)
}

func (v *VmCode) goTo(label string) string {
	return fmt.Sprintf("goto %s", label)
}

func (v *VmCode) ifGoTo(label string) string {
	return fmt.Sprintf("if-goto %s", label)
}

func (v *VmCode) label(label string) string {
	return fmt.Sprintf("label %s", label)
}

func (v *VmCode) newStr(in string) []string {
	// refs Sys.vm:69
	res := []string{
		v.push("constant", len(in)),
		v.call("String", "new", 1),
	}
	for _, r := range in {
		res = append(res, v.pushConstant(int(r)))
		res = append(res, v.call("String", "appendChar", 2))
	}
	return res
}
