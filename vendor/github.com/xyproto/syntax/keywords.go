package syntax

var keywords = map[string]struct{}{
	"BEGIN":            {},
	"END":              {},
	"False":            {},
	"Infinity":         {},
	"NaN":              {},
	"None":             {},
	"True":             {},
	"abstract":         {},
	"alias":            {},
	"align_union":      {},
	"alignof":          {},
	"and":              {},
	"append":           {},
	"as":               {},
	"asm":              {},
	"assert":           {},
	"auto":             {},
	"axiom":            {},
	"begin":            {},
	"bool":             {},
	"boolean":          {},
	"break":            {},
	"byte":             {},
	"caller":           {},
	"case":             {},
	"catch":            {},
	"char":             {},
	"class":            {},
	"concept":          {},
	"concept_map":      {},
	"const":            {},
	"const_cast":       {},
	"constexpr":        {},
	"continue":         {},
	"debugger":         {},
	"decltype":         {},
	"def":              {},
	"default":          {},
	"defined":          {},
	"del":              {},
	"delegate":         {},
	"delete":           {},
	"die":              {},
	"do":               {},
	"double":           {},
	"dump":             {},
	"dynamic_cast":     {},
	"elif":             {},
	"else":             {},
	"elsif":            {},
	"end":              {},
	"ensure":           {},
	"enum":             {},
	"eval":             {},
	"except":           {},
	"exec":             {},
	"exit":             {},
	"explicit":         {},
	"export":           {},
	"extends":          {},
	"extern":           {},
	"false":            {},
	"final":            {},
	"finally":          {},
	"float":            {},
	"float32":          {},
	"float64":          {},
	"for":              {},
	"foreach":          {},
	"friend":           {},
	"from":             {},
	"func":             {},
	"function":         {},
	"generic":          {},
	"get":              {},
	"global":           {},
	"goto":             {},
	"if":               {},
	"implements":       {},
	"import":           {},
	"in":               {},
	"inline":           {},
	"instanceof":       {},
	"int":              {},
	"int8":             {},
	"int16":            {},
	"int32":            {},
	"int64":            {},
	"interface":        {},
	"is":               {},
	"lambda":           {},
	"last":             {},
	"late_check":       {},
	"local":            {},
	"long":             {},
	"make":             {},
	"map":              {},
	"module":           {},
	"mutable":          {},
	"my":               {},
	"namespace":        {},
	"native":           {},
	"new":              {},
	"next":             {},
	"nil":              {},
	"no":               {},
	"nonlocal":         {},
	"not":              {},
	"null":             {},
	"nullptr":          {},
	"operator":         {},
	"or":               {},
	"our":              {},
	"package":          {},
	"pass":             {},
	"print":            {},
	"private":          {},
	"property":         {},
	"protected":        {},
	"public":           {},
	"raise":            {},
	"redo":             {},
	"register":         {},
	"reinterpret_cast": {},
	"require":          {},
	"rescue":           {},
	"retry":            {},
	"return":           {},
	"self":             {},
	"set":              {},
	"short":            {},
	"signed":           {},
	"sizeof":           {},
	"static":           {},
	"static_assert":    {},
	"static_cast":      {},
	"strictfp":         {},
	"struct":           {},
	"sub":              {},
	"super":            {},
	"switch":           {},
	"synchronized":     {},
	"template":         {},
	"then":             {},
	"this":             {},
	"throw":            {},
	"throws":           {},
	"transient":        {},
	"true":             {},
	"try":              {},
	"type":             {},
	"typedef":          {},
	"typeid":           {},
	"typename":         {},
	"typeof":           {},
	"undef":            {},
	"undefined":        {},
	"union":            {},
	"unless":           {},
	"unsigned":         {},
	"until":            {},
	"use":              {},
	"using":            {},
	"var":              {},
	"virtual":          {},
	"void":             {},
	"volatile":         {},
	"wantarray":        {},
	"when":             {},
	"where":            {},
	"while":            {},
	"with":             {},
	"yield":            {},
}
