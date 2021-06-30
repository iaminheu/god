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
{{.updatePartial}}
{{.txUpdate}}
{{.txUpdatePartial}}
{{.delete}}
{{.txDelete}}
`
