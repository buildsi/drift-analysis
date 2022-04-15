package views

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/buildsi/drift-server/config"
	"github.com/buildsi/drift-server/db"
)

// Upload a result, including commits, tags, and a package string
// Create Syntax: http://server/create/?jid=SOMETHING&mid=SOMETHINGELSE with authentication header credentials.
// Report Example: curl -F "upload=@FILENAME" http://localhost:8080/inflection-point/new/
func UploadInflectionPoint(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("upload")
	if err != nil {
		fmt.Println("Error Retrieving the upload Json file!")
		fmt.Println(err)
		return
	}
	defer file.Close()

	// read file and unmarshall json to inflection point
	b, err := ioutil.ReadAll(file)
	var point db.InflectionPoint
	err = json.Unmarshal(b, &point)
	if err != nil {
		fmt.Println(err)
	}

	// Prepare database and save resukt
	database := db.Open(config.Config.Database.Location)
	success, err := db.GetOrCreateInflectionPoint(database, point)
	fmt.Println(success)
}

//
//	result := token.Token{
//		Value:     token.GenValue(config.Global.Tokens.Length),
//		MachineID: args["mid"][0],
//		JobID:     args["jid"][0],
//	}
//	sucess, err := database.Add(db, result)
//	for !sucess {
//		if err != nil {
//			reportError(w, r, "database.Add", err)
//			return
//		}
//		result.Value = token.GenValue(config.Global.Tokens.Length)
//		sucess, err = database.Add(db, result)
//	}
//	fmt.Fprintln(w, result.Value)
//
