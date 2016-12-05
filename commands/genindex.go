// Copyright 2015 The Hugo Authors. All rights reserved.
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

package commands

import (
	"encoding/json"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/spf13/cobra"
	"github.com/spf13/hugo/hugolib"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type Entries []algoliasearch.Object

type AlgoliaCfg struct {
	IndexName     string `yaml:"index_name"`
	ApplicationId string `yaml:"application_id"`
	AdminApiKey   string `yaml:"admin_api_key"`
}

func init() {
	genindexCmd.PersistentFlags().BoolVar(&genindexalgolia, "algolia", false, "Update Algolia Index (requires Algolia config in algolia.yaml)")
	genindexCmd.PersistentFlags().BoolVar(&genindexjson, "json", true, "Generate json index")
	genindexCmd.PersistentFlags().StringVar(&genindexfile, "file", "site_index.json", "json index output filename")
}

var genindexalgolia bool
var genindexjson bool
var genindexfile string
var genindexCmd = &cobra.Command{
	Use:   "index",
	Short: "Generate index of site pages to enable search",
	Long: `Generate an index of site pages to enable search

This command is used to create an index of the rendered pages
of your site. It operats in two modes. 

1. JSON file: by default, this command will write the search index to
a json file that can be used with something like lunr.js

2. Algolia: Providing the --algolia flag and algolia configuration details (see below)
will *update* the specified index on Algolia via API. Note, the index specified must
already exist on Algolia. The update algorithm will:
  1. Create a temporary index, 
  2. Copy settings from existing index and apply it to temporary index
  3. Upload new index data to temporary index
  4. Move temporary index to existing index.

algolia.yaml: if you are updating an index on algolia, you must create algolia.yaml in
the root of the hugo site with the following values:

  index_name: name_of_target_index
  application_id: ALGOLIA_APPLICATION_ID
  admin_api_key: algolia_admin_api_key
`,

	RunE: func(cmd *cobra.Command, args []string) error {
		var entries Entries
		var entriesJson []byte
		var algoliaCfg AlgoliaCfg

		if err := InitializeConfig(); err != nil {
			return err
		}

		sites, err := hugolib.NewHugoSitesFromConfiguration()

		if err != nil {
			return newSystemError("Error creating sites", err)
		}

		if err := sites.Build(hugolib.BuildCfg{SkipRender: true}); err != nil {
			return newSystemError("Error Processing Source Content", err)
		}

		for _, p := range sites.Pages() {
			if p.IsPage() {
				entries = append(entries, processPage(p))
			}
		}

		jww.FEEDBACK.Printf("\nGenerated %d Index Entries\n", len(entries))

		if genindexjson {
			outFile := viper.GetString("staticDir") + "/" + genindexfile
			if entriesJson, err = json.Marshal(entries); err != nil {
				return newSystemError("Error converting index entries to JSON", err)
			}
			if err = ioutil.WriteFile(outFile, entriesJson, 0644); err != nil {
				return newSystemError("Error writing index JSON file", err)
			}
			jww.FEEDBACK.Printf("Wrote %d Index Entries to %s\n", len(entries), outFile)
		}

		if genindexalgolia {
			// Read and process algolia.yaml:
			algoliaSettings, err := ioutil.ReadFile("algolia.yaml")
			if err != nil {
				return newSystemError("Error reading algolia.yaml. Does the file exist?", err)
			}
			if err = yaml.Unmarshal(algoliaSettings, &algoliaCfg); err != nil {
				return newSystemError("Error converting algolia.yaml content", err)
			}

			err = updateAlgoliaIndex(algoliaCfg, entries)
			if err != nil {
				return newSystemError("Error updating Algolia", err)
			}
			jww.FEEDBACK.Printf("Updated Algolia index %s\n", algoliaCfg.IndexName)
		}
		return nil
	},
}

func processPage(page *hugolib.Page) map[string]interface{} {

	entry := map[string]interface{}{
		"title":   page.LinkTitle(),
		"content": page.Plain(),
		"url":     page.Permalink(),
		"uri":     page.RelPermalink(),
		"tags":    page.GetParam("tags"),
	}

	jww.FEEDBACK.Printf("\nPage: %s - %d Words\n", entry["title"], page.WordCount())
	jww.FEEDBACK.Println("URL: ", entry["url"])
	if entry["tags"] != nil {
		jww.FEEDBACK.Printf("Tags: [%s]\n", strings.Join(entry["tags"].([]string), ", "))
	}
	return entry
}

func updateAlgoliaIndex(cfg AlgoliaCfg, entries Entries) (err error) {

	// Init the Algolia api:
	client := algoliasearch.NewClient(cfg.ApplicationId, cfg.AdminApiKey)
	index := client.InitIndex(cfg.IndexName)

	// get settings and synonyms from existing index
	settings, err := index.GetSettings()
	if err != nil {
		return newSystemError("Error retrieving Algolia settings for index: ", cfg.IndexName)
		return err
	}

	// Initialize the temporary index:
	tmpIndex := cfg.IndexName + "_tmp"
	indextmp := client.InitIndex(tmpIndex)

	// 1. Delete the temporary index (will not error if index doesn't exist):
	_, err = indextmp.Delete()
	if err != nil {
		return newSystemError("Error Deleting temporary index: ", err)
	}
	// 2. Set settings on the temp index which will create it
	_, err = indextmp.SetSettings(settings.ToMap())
	if err != nil {
		return newSystemError("Error on Setting Settings for temp index: ", err)
	}
	// 3. Add the entries to the temporary index:
	_, err = indextmp.AddObjects(entries)
	if err != nil {
		return newSystemError("Error Adding Algolia objects to temp index: ", err)
	}
	// 5. Move new index to existing index (this deletes the source(temp) index)
	_, err = client.MoveIndex(tmpIndex, cfg.IndexName)

	return err
}
