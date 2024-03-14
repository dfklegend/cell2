package common

import (
	"fmt"
	"testing"
)

// 测试 class 嵌入
// 匿名嵌入，更多类似函数的调用查找顺序
// 如果有本地实现，优先本地实现，否则查找匿名实现

// 如果想调用推迟实现，可以使用接口

type IAExt interface {
	BeCalled1()
}

type A struct {
	Ext IAExt
}

func NewA() *A {
	return &A{}
}

func (a *A) HelloA() {
	fmt.Println("helloA from A")
}

func (a *A) Hello() {
	fmt.Println("hello from A")
	a.BeCalled()
	a.Ext.BeCalled1()
}

func (a *A) BeCalled() {
	fmt.Println("BeCalled from A")
}

type B struct {
	*A
}

func NewB() *B {
	return &B{
		A: NewA(),
	}
}

func (b *B) HelloB() {
	fmt.Println("helloB from B")
}

func (b *B) BeCalled() {
	fmt.Println("BeCalled from B")
}

func (b *B) BeCalled1() {
	fmt.Println("BeCalled1 from B")
}

type C struct {
	*B
}

func NewC() *C {
	return &C{
		B: NewB(),
	}
}

func Test_EmbedClass(t *testing.T) {
	b := &B{A: &A{}}
	b.A.Ext = b

	b.Hello()
}

func Test_EmbedChainClass(t *testing.T) {
	c := NewC()
	c.HelloA()
}
