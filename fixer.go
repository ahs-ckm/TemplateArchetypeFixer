package main

// ArchetypeFixer:
// a utility to parse templates and remove any references in the definition section that are not part of the template structure.
// This was needed due to an issue involving the template generation transform including additional definitions which in some cases are invalid.

// takes a template
// loads contents in to XML object
// iterates through definition
//	-- for each integrity section archetype name
//  -- check if there is min 1 reference to it in the definition
//  --  xsi:type="CLUSTER" archetype_id="openEHR-EHR-CLUSTER.adhoc_cluster_heading.v1"
//  -- if there isn't a reference, write the file out, without the unnecessary integrity check, to a new folder for later upload

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"aqwari.net/xml/xmltree"
	"github.com/beevik/etree"
)

func main() {

	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) > 0 {
		fmt.Println(argsWithoutProg)
		fmt.Println("Starting as a job....")

		template_path := argsWithoutProg[0]
		output_folder := argsWithoutProg[1]

		log.Println("working on " + template_path)

		removeSurplusArchetypes(template_path, output_folder)
	}
	fmt.Println("Exiting.")
}

func removeSurplus(path string, surplus *[]string, output_folder string) []string {

	for idx := range *surplus {
		archetype_id := (*surplus)[idx]
		log.Println(path + ": removing " + archetype_id)
	}

	return nil
}

//func getSurplusArchetypes(path string, surplus *[]string) string {
func removeSurplusArchetypes(path string, newfolder string) bool {

	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		return false
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	root, _ := xmltree.Parse(content)
	xmltree_elDefinition := root.Search("", "definition")

	flattened := xmltree_elDefinition[0].Flatten()
	flatxml := xmltree_elDefinition[0].Attr("", "archetype_id")

	for idxflat := range flattened {
		flatxml = flatxml + flattened[idxflat].String()
	}

	var elTemplate *etree.Element

	// get template structure
	elTemplate = doc.SelectElement("template")
	if elTemplate != nil {
		elIntegrityCheck := elTemplate.SelectElements("integrity_checks")

		if elIntegrityCheck != nil {
			elDefinition := elTemplate.SelectElement("definition")
			elFindItmes := elDefinition.FindElements(".")

			log.Println(len(elFindItmes))

			for idxIC := range elIntegrityCheck {

				aCheck := elIntegrityCheck[idxIC]
				archetype_id := aCheck.SelectAttrValue("archetype_id", "")

				if archetype_id != "" {
					// check template structure for reference to this.

					if !strings.Contains(flatxml, archetype_id) {

						log.Println("removing " + archetype_id)

						elTemplate.RemoveChild(aCheck)
					}
				}
			}
		}
	}

	file := filepath.Base(path)
	newfile := newfolder + "" + file
	doc.WriteToFile(newfile)
	log.Println( "writing new " + newfile)
	return true
}
