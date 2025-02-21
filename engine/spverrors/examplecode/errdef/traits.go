package errdef

import "github.com/joomcode/errorx"

var globalTraits []TraitDefinition

type TraitDefinition struct {
	Title string
	Trait errorx.Trait
}

func RegisterTrait(name, title string) errorx.Trait {
	trait := errorx.RegisterTrait(name)
	def := TraitDefinition{
		Title: title,
		Trait: trait,
	}
	globalTraits = append(globalTraits, def)
	return trait
}

var TraitConfig = RegisterTrait("config", "Server may be configured incorrectly")
var TraitIllegalArgument = RegisterTrait("illegal_argument", "Illegal Argument")
