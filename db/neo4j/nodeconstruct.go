package neo4j

import "fmt"

// # This file has utility funcs for creating
// # cypher strings.

// # nodeSpec is meant to reflect what neo4j
// # node properties are called. It is used
// # several places in this pkg for refactoring
// # purposes.
var nodeSpec = struct {
	label string
	title string
	html  string
}{
	label: "WikiData",
	title: "title",
	html:  "html",
}

// # Simply convert a slice of strings into
// # a linear csv. Example: items=["a","b"],
// # -> "a, b, c".
func linearCSV(items []string) string {
	str := ""
	l := len(items)
	for i := 0; i < l; i++ {
		str += items[i]
		if i < l-1 {
			str += ", "
		}
	}
	return str
}

// # Create CQL bindings, using property ids.
// # Example: propIDs=["a", "b"] -> "a:$a, b:$b"
func nodePropBindStr(propIDs []string) string {
	keyValued := make([]string, 0, len(propIDs))
	for _, id := range propIDs {
		s := fmt.Sprintf("%s:$%s", id, id)
		keyValued = append(keyValued, s)
	}
	return "{" + linearCSV(keyValued) + "}"
}

// # Simply add a cypher accessor to a str.
// # Example: alias="a", s="b" -> "a.b".
func aliasedProp(alias, s string) string {
	return fmt.Sprintf("%s.%s", alias, s)
}

// # Aliased argument to neo4j id() func.
// # Example: alias="a" -> "id(a)"
func aliasedID(alias string) string {
	return fmt.Sprintf("id(%s)", alias)
}

// # Gives aliased node prop vals for db.WikiDataBrief.
// # alias="a" ->
// # 	"id(a)",
// # 	"a.(nodeSpec.title)",
// # 	csv of vals above.
func aliasedPropsBrief(alias string) (id, title, csv string) {
	id = aliasedID(alias)
	title = aliasedProp(alias, nodeSpec.title)
	csv = linearCSV([]string{id, title})
	return
}

// # Gives aliased node prop vals for db.WikiData.
// # alias="a" ->
// # 	"id(a)",
// # 	"a.(nodeSpec.title)",
// # 	"a.(nodeSpec.html)",
// # 	csv of vals above.
func aliasedProps(alias string) (id, title, html, csv string) {
	id, title, _ = aliasedPropsBrief(alias)
	html = aliasedProp(alias, nodeSpec.html)
	csv = linearCSV([]string{id, title, html})
	return
}

// # Constructs an aliased cql node with fromat:
// # 	let a,b,c = alias, nodeSpec.label, nodeSpec.title
// # 	(a:b {c:$c})
// # Note; meant to be used with bindings.
func cqlNode(alias string, addProps bool) string {
	props := ""
	if addProps {
		props = nodePropBindStr([]string{nodeSpec.title})
	}
	return fmt.Sprintf("(%s:%s %s)",
		alias, nodeSpec.label, props)
}
