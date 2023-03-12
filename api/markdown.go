
package api

import (
	"bytes"
	"net/url"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
)

type Option struct{
	URLPrefix string
	DataURLPrefix string
	HeadingIDPrefix string
	HeadingIDSuffix string
}

func (o *Option)fixRelLink(prefix, src string)(string){
	if !strings.Contains(src, "://") && !strings.HasPrefix(src, "/") && !strings.HasPrefix(src, "#") {
		loger.Infof("src: %s", src, strings.HasPrefix(src, "#"))
		out, err := url.JoinPath(prefix, src)
		if err != nil {
			out = "about:blank#err:" + err.Error()
		}
		return out
	}
	return src
}

func (o *Option)exec(n ast.Node)(ast.Node){
	switch m := n.(type) {
	case *ast.Link:
		if len(o.URLPrefix) > 0 {
			m.Destination = ([]byte)(o.fixRelLink(o.URLPrefix, (string)(m.Destination)))
		}
	case *ast.Image:
		if len(o.DataURLPrefix) > 0 {
			m.Destination = ([]byte)(o.fixRelLink(o.DataURLPrefix, (string)(m.Destination)))
		}else if len(o.URLPrefix) > 0 {
			m.Destination = ([]byte)(o.fixRelLink(o.URLPrefix, (string)(m.Destination)))
		}
	case *ast.Heading:
		anchor := &ast.Link{
			Container: ast.Container{
				Parent: m,
				Children: []ast.Node{
					&ast.Text{
						Leaf: ast.Leaf{
							Literal: ([]byte)("+"),
						},
					},
				},
			},
			Destination: ([]byte)("#" + o.HeadingIDPrefix + m.HeadingID + o.HeadingIDSuffix),
			AdditionalAttributes: []string{`class="anchor"`},
		}
		m.Children = append(m.Children, anchor)
		return m
	case *ast.CodeBlock:
	}
	return n
}

func (o *Option)Walk(n ast.Node)(ast.Node){
	n = o.exec(n)
	children := n.GetChildren()
	for i, c := range children {
		children[i] = o.Walk(c)
	}
	return n
}

func (o *Option)SetHtmlOptions(opt *html.RendererOptions){
	opt.HeadingIDPrefix = o.HeadingIDPrefix
	opt.HeadingIDSuffix = o.HeadingIDSuffix
}

var langClassRe = regexp.MustCompile("language-[^ ]")

func (o *Option)NewPolicy()(p *bluemonday.Policy){
	p = bluemonday.UGCPolicy()
	p.AllowAttrs("class").Matching(langClassRe).OnElements("code")
	p.AllowAttrs("class").Matching(regexp.MustCompile("anchor")).OnElements("a")
	return
}

var DefaultOption = &Option{}

func RenderMarkdown(src []byte, opt ...*Option)(out []byte){
	var o = DefaultOption
	if len(opt) > 0 {
		o = opt[0]
	}

	src = bytes.ReplaceAll(src, ([]byte)("\r\n"), ([]byte)("\n"))
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.Mmark
	p := parser.NewWithExtensions(extensions)
	opts := html.RendererOptions{
		Flags: html.CommonFlags | html.LazyLoadImages,
	}
	a := p.Parse(src)
	a = o.Walk(a)
	o.SetHtmlOptions(&opts)
	r := html.NewRenderer(opts)
	out1 := markdown.Render(a, r)
	out = o.NewPolicy().SanitizeBytes(out1)
	return
}
