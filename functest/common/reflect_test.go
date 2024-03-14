package common

// 测试反射 implememts
// 结论:

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Model interface {
	m()
}

func HasModels(m Model) {
	s := reflect.ValueOf(m).Elem()
	t := s.Type()

	// 接口类型
	modelType := reflect.TypeOf((*Model)(nil)).Elem()

	for i := 0; i < s.NumField(); i++ {
		f := t.Field(i)
		fmt.Printf("%d: %s %s -> %t\n", i, f.Name, f.Type, f.Type.Implements(modelType))
	}
}

func HasImpl(t reflect.Type, it reflect.Type) bool {
	fmt.Printf("%s %v -> %t\n", t.Name(), t, t.Implements(it))
	return t.Implements(it)
}

type Company struct{}

func (Company) m() {}

type Department struct{}

func (*Department) m() {}

type User struct {
	CompanyA    Company
	CompanyB    *Company
	DepartmentA Department
	DepartmentB *Department
}

func (User) m() {}

// 不同接受者
func Test(t *testing.T) {
	u := &User{}
	HasModels(u)

	modelType := reflect.TypeOf((*Model)(nil)).Elem()
	assert.Equal(t, true, HasImpl(reflect.TypeOf(u.CompanyA), modelType))
	assert.Equal(t, true, HasImpl(reflect.TypeOf(u.CompanyB), modelType))
	//
	assert.Equal(t, false, HasImpl(reflect.TypeOf(u.DepartmentA), modelType))
	assert.Equal(t, true, HasImpl(reflect.TypeOf(u.DepartmentB), modelType))
}

// 测试函数
type ClassA struct {
}

func (a *ClassA) Func1(args ...any) {}

func Func1(args ...any) {

}

// 如何区分Method和Func
func TestFunc(t *testing.T) {
	a := &ClassA{}

	f1 := reflect.ValueOf(a.Func1)
	t1 := f1.Type()
	log.Println(t1.NumIn())

	f2 := reflect.ValueOf(Func1)
	t2 := f2.Type()
	log.Println(t2.NumIn())
}

func toAny(i any) any {
	return i
}

func TestEqual(t *testing.T) {
	a1 := &ClassA{}
	a2 := a1
	b := a1 == a2

	i1 := toAny(a1)
	i2 := toAny(a2)
	b1 := i1 == i2
	log.Println(b)
	log.Println(b1)
}
