package evaluator

import (
	"z/object"
)

var builtins = map[string]*object.Builtin{
	"len":         object.GetBuiltinByName("len"),
	"first":       object.GetBuiltinByName("first"),
	"last":        object.GetBuiltinByName("last"),
	"rest":        object.GetBuiltinByName("rest"),
	"push":        object.GetBuiltinByName("push"),
	"puts":        object.GetBuiltinByName("puts"),
	"execute":     object.GetBuiltinByName("execute"),
	"mysql_query": object.GetBuiltinByName("mysql_query"),
	"mysql_init":  object.GetBuiltinByName("mysql_init"),
	"typeof":      object.GetBuiltinByName("typeof"),
}
