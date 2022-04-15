package db

type Commit struct {
	Digest    string `json:"digest"`
	Timestamp string `json:"timestamp"`
}

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Spec struct {
	Pkg  *Package `json:"package"`
	Data string   `json:"data"`
}

type InflectionPoint struct {
	Commit Commit    `json:"commit"`
	Tags   []*string `json:"tags"`
	Pkg    *Package  `json:"package"`
}

type Build struct {
	Spec             *Spec            `json:"spec"`
	Inflection_point *InflectionPoint `json:"inflection"`
	Status           string           `json:"status"`
}

//func ShowAllUsers() (au *AllUsers) {
//	file, err := os.OpenFile("list.json", os.O_RDWR|os.O_APPEND, 0666)
//	checkError(err)
//	b, err := ioutil.ReadAll(file)
//	var alUsrs AllUsers
//	json.Unmarshal(b, &alUsrs.Users)
//	checkError(err)
//	return &alUsrs
//}
