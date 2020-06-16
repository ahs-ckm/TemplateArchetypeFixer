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
	"strconv"
	"strings"

	"aqwari.net/xml/xmltree"
	"github.com/beevik/etree"
)

func main() {

	argsWithoutProg := os.Args[1:]
	var numberFilesRead = 0
	//var numberFilesWritten = 0
	if len(argsWithoutProg) > 0 {
		fmt.Println(argsWithoutProg)
		fmt.Println("Starting as a job....")

		template_path := argsWithoutProg[0]
		output_folder := argsWithoutProg[1]

		log.Println("working on " + template_path)
		log.Println("output folder= " + output_folder)

		//Check if this is a single template file or a folder
		if strings.HasSuffix(strings.ToUpper(template_path), ".OET") {
			log.Println("\r\n\r\n" + "This is a single template file\r\n")
			//For a single file
			removeSurplusArchetypes(template_path, output_folder)
		} else {
			log.Println("\r\n\r\n" + "Going to look for templates in " + template_path + " \r\n")
			/*** START **/
			var err = filepath.Walk(template_path,
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					//fmt.Println("Recursive File="+path, info.Size())
					if fi, err := os.Stat(path); err == nil {
						if fi.Mode().IsDir() {
							//fmt.Println("Is a Directory")
						} else if strings.HasSuffix(strings.ToUpper(""+path), ".OET") {
							//fmt.Println("filename=" + fi.Name())
							//var fullpath = template_path + "\\" + path
							numberFilesRead++
							fmt.Println("Examining " + path)
							var outputFile = output_folder
							// + fi.Name()
							fmt.Println("Output Path = " + outputFile)
							removeSurplusArchetypes(path, outputFile)
						}
					}

					//fmt.Println("Is Directory=" + IsDirectory(path))
					return nil
				})
			if err != nil {
				log.Println(err)
			}

			/*** END **/
		}

	}
	fmt.Println("# Files Read = " + strconv.Itoa(numberFilesRead))

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

	log.Println(" Removing surplus archetypes for = " + path)
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
	log.Println("writing new " + newfile)
	return true
}
