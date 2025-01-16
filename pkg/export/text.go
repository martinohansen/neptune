package export

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/martinohansen/neptune/pkg/places"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type textModel struct {
	Categories []textCategory
}

type textCategory struct {
	Name   string
	Places []textPlace
}

type textPlace struct {
	Name    string
	Address string
	Note    string
	URL     string
}

const tmpl = `
{{range .Categories}}
### {{.Name | categoryTitle}}
{{range .Places -}}
* [{{.Name}}]({{.URL}}) {{with .Note}}({{.}}){{end}}
  {{.Address}}
{{end}}
{{end}}
`

func model(ps []*places.Place) (model textModel) {
	categories := places.DistinctCategories(ps)
	for _, cat := range categories {
		x := textCategory{Name: cat}
		for _, p := range ps {
			url := fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=Google&query_place_id=%s", p.GooglePlaceID)
			if cat == p.Categories[0] {
				// TODO: Add notes value to places.Place type
				v := textPlace{
					Name:    p.Name,
					Address: p.FormattedAddress,
					Note:    "",
					URL:     url,
				}
				x.Places = append(x.Places, v)
			}
		}
		model.Categories = append(model.Categories, x)
	}
	return model
}

func Text(ps []*places.Place, wr io.Writer) error {
	funcMap := template.FuncMap{
		"categoryTitle": categoryTitle,
	}

	t := template.Must(template.New("text").Funcs(funcMap).Parse(tmpl))

	err := t.Execute(wr, model(ps))
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}
	return nil
}

// categoryTitle takes a category string and returns it in title case
func categoryTitle(s string) string {
	title := cases.Title(language.English)
	return title.String(strings.Replace(s, "_", " ", -1))
}
