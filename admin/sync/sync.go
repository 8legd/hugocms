package admin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Syncs creation and update events for a page with Hugo
func WritePage(model interface{}, path string, slug string) (err error) {
	var filename = "data" + path + slug + ".json"
	err = write(model, filename)
	if err == nil {
		// If doesn't exist create page ref in Hugo
		fmt.Printf("\nTODO create content file for %s \n", filename)
	}
	return
}

// Syncs rename or deletion of a page with Hugo
func RemovePageRef(path string, slug string) (err error) {
	// TODO just remove content file from Hugo (data files can remain)
	// this way the page will in effect be un-published
	// TODO if after removing the content file the section directory is empty then
	// also remove it to clear up (again the sections will remain in the data dir)
	var filename = "content" + path + slug + ".json"
	fmt.Printf("\nTODO remove %s \n", filename)
	fmt.Printf("TODO and if required remove empty section dir %s \n", path)
	return
}

func write(model interface{}, filename string) error {
	output, err := json.MarshalIndent(model, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, output, 0644)
	return err
}

func handleError(err error) {
	fmt.Println(err)
	os.Exit(1)
	//TODO report to user through QOR Admin
}
