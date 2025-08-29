package extractor

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type TypeAnnotation struct {
	TypeName string
	FilePath string
	Params   map[string]string
}

type TypeExtractor interface {
	Extract(filePath string, keyword string) ([]TypeAnnotation, error)
}

type typeExtractor struct{}

func NewTypeExtractor() TypeExtractor {
	return &typeExtractor{}
}

func (e *typeExtractor) Extract(path string, keyword string) ([]TypeAnnotation, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var results []TypeAnnotation
	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}

		for _, spec := range gen.Specs {
			ts := spec.(*ast.TypeSpec)
			if gen.Doc == nil {
				continue
			}

			for _, c := range gen.Doc.List {
				if strings.HasPrefix(c.Text, "// "+keyword+":") {
					line := strings.TrimSpace(strings.TrimPrefix(c.Text, "// "+keyword+":"))

					args := map[string]string{
						"path": path,
						"type": ts.Name.Name,
					}
					for _, part := range strings.Fields(line) {
						if strings.Contains(part, "=") {
							kv := strings.SplitN(part, "=", 2)
							args[kv[0]] = kv[1]
						}
					}

					results = append(results, TypeAnnotation{
						TypeName: ts.Name.Name,
						FilePath: path,
						Params:   args,
					})
				}
			}
		}
	}
	return results, nil
}
