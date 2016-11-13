// Copyright 2016-present The Hugo Authors. All rights reserved.
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
	"github.com/spf13/cast"
	"html/template"

	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

// TODO See the shortcode system to see the structure

// Data structures and methods
// ===========================

// A WidgetEntry represents a widget item defined
// in the site config.
// (TODO: See hugolib/menu.go for the data structure)
type Widget struct {
	Type       string
	Params     map[string]interface{}
	Identifier string
	Weight     int
	Template   *template.Template
}

func newWidget(widgetType string, options interface{}) (*Widget, error) {
	return &Widget{Type: widgetType, Params: options.(map[string]interface{})}, nil
}

type WidgetArea struct {
	Name     string
	Widgets  []*Widget
	Template *template.Template
}

func newWidgetArea(waname string) *WidgetArea {
	// TODO ??
	return &WidgetArea{Name: waname, Widgets: nil, Template: nil}
}

type Widgets map[string]*WidgetArea

// Internal widgets building
// =========================

// WidgetsConfig parses the widgets config variable
// (is a collection) and calls every widget configuration.
func getWidgetsFromConfig() Widgets {
	ret := Widgets{}

	if conf := viper.GetStringMap("widgets"); conf != nil {
		for waname, widgetarea := range conf {
			// wa is a widget area defined in the conf file
			wa, err := cast.ToSliceE(widgetarea)
			if err != nil {
				jww.ERROR.Printf("unable to process widgets in site config\n")
				jww.ERROR.Println(err)
			}

			// Instantiate a WidgetArea
			waobj := newWidgetArea(waname)

			// Retrieve all widgets
			for _, w := range wa {
				iw, err := cast.ToStringMapE(w)

				if err != nil {
					jww.ERROR.Printf("unable to process widget inside widget area in site config\n")
					jww.ERROR.Println(err)
				}

				// iw represents a widget inside a widget area
				wtype := cast.ToString(iw["type"])
				woptions, err := cast.ToStringMapE(iw["options"])
				wobj, err := newWidget(wtype, woptions)

				if err != nil {
					jww.ERROR.Printf("unable to instantiate widget: %s\n", iw)
					jww.ERROR.Println(err)
				}

				// then append it to the widget area object
				waobj.Widgets = append(waobj.Widgets, wobj)
			}

			// don't forget to append that widget area to the
			// Widgets object
			ret[waname] = waobj
		}
	}

	return ret
}

// instantiateWidget retrieves the widget's files
// and creates the templates
func instantiateWidget(s *Site, wa *WidgetArea, w *Widget) *Widget {
	// Load this widget's templates
	// using the site object's owner.tmpl
	s.owner.tmpl.LoadTemplatesWithPrefix(s.absWidgetDir()+"/"+w.Type+"/layouts", "widgets/"+w.Type)

	return w
}

// Main widgets entry point
// ========================

// This function adds the whole widgets' template code
// in the Site object. This is of type template.HTML.
// This function is called from hugo_sites.
func injectWidgets(s *Site) error {
	// Get widgets. This gives all information we need but
	// does not already read widget files.
	widgets := getWidgetsFromConfig()

	for _, widgetarea := range widgets {
		// _ is waname, if ever we need

		for _, w := range widgetarea.Widgets {
			w = instantiateWidget(s, widgetarea, w)
		}
	}

	// We now have all widgets with their templates.
	// Generate all widget areas with their templates
	// Now the template's content will be used inside
	// the main template files inside tpl/template_funcs
	// and in the templates using {{ widgets "mywidgetarea" }}
	s.Widgets = widgets

	return nil
}
