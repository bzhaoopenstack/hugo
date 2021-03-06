// Copyright 2019 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

import (
	"fmt"
	"testing"
)

func TestPaginator(t *testing.T) {
	configFile := `
baseURL = "https://example.com/foo/"
paginate = 3
paginatepath = "thepage"

[languages.en]
weight = 1
contentDir = "content/en"

[languages.nn]
weight = 2
contentDir = "content/nn"

`
	b := newTestSitesBuilder(t).WithConfigFile("toml", configFile)
	var content []string
	for i := 0; i < 9; i++ {
		for _, contentDir := range []string{"content/en", "content/nn"} {
			content = append(content, fmt.Sprintf(contentDir+"/blog/page%d.md", i), fmt.Sprintf(`---
title: Page %d
---

Content.
`, i))
		}

	}

	b.WithContent(content...)

	pagTemplate := `
{{ $pag := $.Paginator }}
Total: {{ $pag.TotalPages }}
First: {{ $pag.First.URL }}
Page Number: {{ $pag.PageNumber }}
URL: {{ $pag.URL }}
{{ with $pag.Next }}Next: {{ .URL }}{{ end }}
{{ with $pag.Prev }}Prev: {{ .URL }}{{ end }}
{{ range $i, $e := $pag.Pagers }}
{{ printf "%d: %d/%d  %t" $i $pag.PageNumber .PageNumber (eq . $pag) -}}
{{ end }}
`

	b.WithTemplatesAdded("index.html", pagTemplate)
	b.WithTemplatesAdded("index.xml", pagTemplate)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html",
		"Page Number: 1",
		"0: 1/1  true")

	b.AssertFileContent("public/thepage/2/index.html",
		"Total: 3",
		"Page Number: 2",
		"URL: /foo/thepage/2/",
		"Next: /foo/thepage/3/",
		"Prev: /foo/",
		"1: 2/2  true",
	)

	b.AssertFileContent("public/index.xml",
		"Page Number: 1",
		"0: 1/1  true")
	b.AssertFileContent("public/thepage/2/index.xml",
		"Page Number: 2",
		"1: 2/2  true")

	b.AssertFileContent("public/nn/index.html",
		"Page Number: 1",
		"0: 1/1  true")

	b.AssertFileContent("public/nn/index.xml",
		"Page Number: 1",
		"0: 1/1  true")

}
