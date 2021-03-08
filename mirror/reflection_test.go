package mirror

import (
	"testing"
)

type Child struct {
	Name string
}

type Parent struct {
	Name     string
	Children []Child
}

type GrandParent struct {
	Name  string
	Child Parent
}

type GreatGrandParent struct {
	Child GrandParent
	Name  string
}

func (ggp *GreatGrandParent) GetGreatGrandChildren() []Child {
	return ggp.Child.Child.Children
}

func NewSingleChild() *GreatGrandParent {
	return &GreatGrandParent{
		Child: GrandParent{
			Name: "GiPa",
			Child: Parent{
				Name:     "Pa",
				Children: []Child{{Name: "Susy"}},
			},
		},
		Name: "GiGiPa",
	}
}

func TestReflection_GetPath(t *testing.T) {
	test := NewSingleChild()
	reflector := Reflect(test)
	children := reflector.GetPath("/Child/Child/Children").PanicIfErr().Value().Interface().([]Child)
	if len(children) != 1 || children[0].Name != "Susy" {
		t.Error("Failed to get path")
	}
}

func TestReflection_SetPath(t *testing.T) {
	test := NewSingleChild()
	reflector := Reflect(test)
	setChildren := reflector.SetPath(
		"/Child/Child/Children",
		[]Child{{Name: "Jerry"}},
	).
		PanicIfErr().
		Value().
		Interface().([]Child)

	gotChildren := reflector.GetPath("/Child/Child/Children").PanicIfErr().Value().Interface().([]Child)

	newReflector := Reflect(test)
	newReflectorGotChildren := newReflector.GetPath("/Child/Child/Children").PanicIfErr().Value().Interface().([]Child)

	if len(setChildren) != 1 || setChildren[0].Name != "Jerry" {
		t.Error("Failed to set path")
	}

	if len(gotChildren) != 1 || gotChildren[0].Name != "Jerry" {
		t.Error("Failed to set path")
	}

	if len(newReflectorGotChildren) != 1 || newReflectorGotChildren[0].Name != "Jerry" {
		t.Error("Failed to set path")
	}
}

func TestReflection_Exec(t *testing.T) {
	test := NewSingleChild()
	greatGrandChildren := Reflect(test).GetPath("/GetGreatGrandChildren").Exec().UnwrapResult()[0].([]Child)
	if len(greatGrandChildren) != 1 || greatGrandChildren[0].Name != "Susy" {
		t.Error("Failed to get path")
	}
}
