// Декодер BBcode
// За основу был взят данный пакет:
// https://pkg.go.dev/zxq.co/ripple/hanayo/modules/bbcode
package bbcode

import (
	"strings"

	"github.com/frustra/bbcode"
)

type BBcode struct {
	Compiler bbcode.Compiler
}

// Создание нового компи
func New() BBcode {
	compiler := bbcode.NewCompiler(true, true)
	compiler.SetTag("list", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "ul"
		style := node.GetOpeningTag().Value
		switch style {
		case "a":
			out.Attrs["style"] = "list-style-type: lower-alpha;"
		case "A":
			out.Attrs["style"] = "list-style-type: upper-alpha;"
		case "i":
			out.Attrs["style"] = "list-style-type: lower-roman;"
		case "I":
			out.Attrs["style"] = "list-style-type: upper-roman;"
		case "1":
			out.Attrs["style"] = "list-style-type: decimal;"
		default:
			out.Attrs["style"] = "list-style-type: disc;"
		}

		if len(node.Children) == 0 {
			out.AppendChild(bbcode.NewHTMLTag(""))
		} else {
			node.Info = []*bbcode.HTMLTag{out, out}
			tags := node.Info.([]*bbcode.HTMLTag)
			for _, child := range node.Children {
				curr := tags[1]
				curr.AppendChild(node.Compiler.CompileTree(child))
			}
			if len(tags[1].Children) > 0 {
				last := tags[1].Children[len(tags[1].Children)-1]
				if len(last.Children) > 0 && last.Children[len(last.Children)-1].Name == "br" {
					last.Children[len(last.Children)-1] = bbcode.NewHTMLTag("")
				}
			} else {
				tags[1].AppendChild(bbcode.NewHTMLTag(""))
			}
		}
		return out, false
	})

	compiler.SetTag("*", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		parent := node.Parent
		for parent != nil {
			if parent.ID == bbcode.OPENING_TAG && parent.GetOpeningTag().Name == "list" {
				out := bbcode.NewHTMLTag("")
				out.Name = "li"
				tags := parent.Info.([]*bbcode.HTMLTag)
				if len(tags[1].Children) > 0 {
					last := tags[1].Children[len(tags[1].Children)-1]
					if len(last.Children) > 0 && last.Children[len(last.Children)-1].Name == "br" {
						last.Children[len(last.Children)-1] = bbcode.NewHTMLTag("")
					}
				} else {
					tags[1].AppendChild(bbcode.NewHTMLTag(""))
				}
				tags[1] = out
				tags[0].AppendChild(out)

				if len(parent.Children) == 0 {
					out.AppendChild(bbcode.NewHTMLTag(""))
				} else {
					for _, child := range node.Children {
						curr := tags[1]
						curr.AppendChild(node.Compiler.CompileTree(child))
					}
				}
				if node.ClosingTag != nil {
					tag := bbcode.NewHTMLTag(node.ClosingTag.Raw)
					bbcode.InsertNewlines(tag)
					out.AppendChild(tag)
				}
				return nil, false
			}
			parent = parent.Parent
		}
		return bbcode.DefaultTagCompiler(node)
	})

	compiler.SetTag("hr", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "div"
		out.Attrs["class"] = "ui divider"
		out.AppendChild(nil)
		return out, false
	})

	compiler.SetTag("p", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "p"
		return out, true
	})

	compiler.SetTag("table", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "table"
		return out, true
	})

	compiler.SetTag("tr", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "tr"
		return out, true
	})

	compiler.SetTag("td", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "td"
		return out, true
	})

	return BBcode{
		Compiler: compiler,
	}
}

// Компилятор
func (c *BBcode) Compile(s string) string {
	s = strings.TrimSpace(s)
	return c.Compiler.Compile(s)
}
