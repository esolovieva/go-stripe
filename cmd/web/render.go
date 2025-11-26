package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

// templateData represents a collection of data that can be passed
// to 'HTML templates' during rendering. It provides several typed maps
// for flexible data transport as well as common fields used across
// web pages such as flash messages, CSRF tokens, authentication flags,
// and version information.
type templateData struct {
	StringMap       map[string]string  //Universal container for transfering string data
	IntMap          map[string]int     //Universal container for transfering int data
	FloatMap        map[string]float32 //Universal cotainer for transfering float data
	Data            map[string]any     //Any data common container
	CSRFToken       string             //String with CSRFToken
	Flash           string             //Message for the user
	Warning         string             //Message for the user
	Error           string             //Message for the user
	IsAuthenticated int                //Flag of the authentication
	API             string             //URL API
	CSSVersion      string             //Version of CSS for cache-busting
}

var functions = template.FuncMap{} //An empty map of functions that can be connected in HTML templates (very useful for helper functions in templates).

//go:embed templates
var templateFS embed.FS //Embeds the templates folder inside the binary using go:embed

// addDefaultData injects default values into the provided templateData
// structure before the template is rendered. This typically includes
// global data shared across all templates, such as CSRF tokens or
// authentication information.
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	return td
}

// renderTemplate renders the specified page template to the provided
// ResponseWriter. It loads the template from cache when available, or
// parses it from the embedded filesystem otherwise. Optional partial
// templates may be included. The method returns an error if parsing or
// execution fails.
func (app *application) renderTemplate(w http.ResponseWriter, r *http.Request, page string, td *templateData, partials ...string) error {
	var t *template.Template
	var err error

	templateToRender := fmt.Sprintf("templates/%s.page.gohtml", page) //sets path to a template

	_, templateInMap := app.templateCache[templateToRender] //bool - if teplate exists in the template cache

	if app.config.env == "production" && templateInMap {
		t = app.templateCache[templateToRender] //if this is production env, get template from cache (if it exists there)
	} else { //otherwise parce template from scratch
		t, err = app.parseTemplate(partials, page, templateToRender)
		if err != nil {
			app.errorLog.Println(err)
			return err
		}
	}

	if td == nil {
		td = &templateData{} //set template data if it is empty
	}

	td = app.addDefaultData(td, r) //add default data to the template data

	err = t.Execute(w, td) //applies a parsed template to the specified data object, writing the output to w
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	return nil
}

// parseTemplate builds and parses a complete template by combining a base
// layout, optional partial templates, and the main page template. It reads
// templates from the embedded filesystem and stores the resulting parsed
// template in the application's template cache. It returns the constructed
// template or an error if parsing fails.
func (app *application) parseTemplate(partials []string, page, templateToRender string) (*template.Template, error) {
	var t *template.Template

	var err error

	// build partials (names => paths)
	if len(partials) > 0 {
		for i, x := range partials {
			partials[i] = fmt.Sprintf("templates/%s.partial.gohtml", x)
		}
	}

	//forms the list of template files: base layout, partials, template to render
	if len(partials) > 0 {
		args := make([]string, 0, len(partials)+2)
		args = append(args, "templates/base.layout.gohtml")
		args = append(args, partials...)
		args = append(args, templateToRender)
		t, err = template.New(fmt.Sprintf("%s.page.gohtml", page)).Funcs(functions).ParseFS(templateFS, args...)
	} else {
		t, err = template.New(fmt.Sprintf("%s.page.gohtml", page)).Funcs(functions).ParseFS(templateFS, "templates/base.layout.gohtml", templateToRender)
	}
	if err != nil {
		app.errorLog.Println(err)
		return nil, err
	}

	app.templateCache[templateToRender] = t //Save ready template into templates cache
	return t, nil
}
