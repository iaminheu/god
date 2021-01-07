package tpl

var Model = `package model
{{.imports}}
{{.vars}}
{{.types}}
{{.new}}
{{.insert}}
{{.txInsert}}
{{.find}}
{{.update}}
{{.txUpdate}}
{{.delete}}
{{.txDelete}}
`
