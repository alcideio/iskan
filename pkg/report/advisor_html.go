package report

import (
	"bytes"
	"encoding/json"
	"text/template"

	"github.com/alcideio/iskan/pkg/advisor"
)

type HtmlReport struct {
	Report *advisor.AdvisorClusterReport
}

func (r *HtmlReport) generateData() string {
	data, err := json.Marshal(r.Report)
	if err != nil {
		return "{error}"
	}
	return string(data)
}

func (r *HtmlReport) newTemplateEngine(name string, data string) (*template.Template, error) {
	funcs := template.FuncMap{
		"generateData": r.generateData,
	}
	return template.New(name).Funcs(funcs).Parse(data)
}

func (r *HtmlReport) Generate() (out string, err error) {

	html := `
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <link rel="shortcut icon" href="https://www.alcide.io/wp-content/themes/alcide/favicon.ico" />
  <title>[Alcide] iSKan .. Kubernetes Native Container Image Scanner</title>
  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no" />
  <meta name="theme-color" content="#000000" />

  <!-- Fonts -->
  <link href="https://fonts.googleapis.com/css?family=Poppins:200,300,400,600,700,800" rel="stylesheet" />
  <link href='https://fonts.googleapis.com/css?family=Source Code Pro' rel='stylesheet' />
  <link href='https://fonts.googleapis.com/css?family=Roboto' rel='stylesheet' />
  <!-- Icons -->
  <script src="https://kit.fontawesome.com/14367b7099.js" crossorigin="anonymous"></script>

  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css"
    integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh" crossorigin="anonymous" />

	<link rel="stylesheet" href="https://unpkg.com/@alcideio/alcide-skan-viewer@0.1.3/dist/main.css"  />

</head>

<body>
  <noscript>
    You need to enable JavaScript to run this app.
  </noscript>
  	<script type='text/javascript'>
		window['skanReportData'] = {{ generateData }}
  	</script>
  	<div id="root"></div>
	<script src="https://unpkg.com/@alcideio/alcide-skan-viewer@0.1.3/dist/main.js" crossorigin="anonymous"></script>
</body>

</html>
  `

	tmpl, err := r.newTemplateEngine("full-report", html)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, r)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
