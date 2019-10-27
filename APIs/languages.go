package APIs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type LanguageInfo struct {
	Name  string
	Count int
}

// to sort by count at the end
type ByCount []LanguageInfo

func (a ByCount) Len() int           { return len(a) }
func (a ByCount) Less(i, j int) bool { return a[i].Count > a[j].Count } // which is More for our purpose
func (a ByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type LanguagesAnswer struct {
	Languages []string `json:"languages"`
	Auth      bool     `json:"auth"`
}

func HandlerLanguages(w http.ResponseWriter, r *http.Request) {
	var languagesAnswer = LanguagesAnswer{}

	http.Header.Add(w.Header(), "content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 4 || parts[3] != "languages" {
		http.Error(w, "Malformed URL", http.StatusBadRequest)
		return
	}

	limitRequest := 5

	// get the limit number if it's not the default one
	limit, ok0 := r.URL.Query()["limit"]

	if ok0 {
		limit, err := strconv.Atoi(limit[0])

		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		limitRequest = limit
	}

	authRequest := ""
	withAuth := false

	//get the authentification if there is one
	auth, ok1 := r.URL.Query()["auth"]

	if ok1 {
		authRequest = auth[0]
		withAuth = true
	}

	// webhook call
	var parameters []string
	if ok0 {
		parameters = append(parameters, "limit")
	}
	if withAuth {
		parameters = append(parameters, "authentification")
	}

	WebhookChecking(w, "languages", parameters)

	// get the payload if there is one
	var projectsName []string
	err := json.NewDecoder(r.Body).Decode(&projectsName)

	// determine the ids we will work with
	var projectsId []float64
	if err != nil { // without playload
		if withAuth {
			projectsId = GetProjectsIdWithAuth(w, authRequest)
		} else {
			projectsId = GetProjectsId(w)
		}
	} else { // with playload
		var nameIdMap map[string]float64
		nameIdMap = make(map[string]float64)

		if withAuth {
			nameIdMap = GetProjectsNameIdMapWithAuth(w, authRequest)
		} else {
			nameIdMap = GetProjectsNameIdMap(w)
		}

		for i := 0; i < len(projectsName); i++ {
			projectsId = append(projectsId, nameIdMap[projectsName[i]])
		}
	}

	var languageCountMap map[string]int
	languageCountMap = make(map[string]int)

	// for each project, we get the languages use and update languageCountMap
	lenProjectsId := len(projectsId)
	for i := 0; i < lenProjectsId; i++ {
		id := strconv.Itoa(int(projectsId[i]))

		var getArgument string
		if withAuth {
			getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects/%s/languages?private_token=%s", id, authRequest)
		} else {
			getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects/%s/languages", id)
		}

		resp, err := http.Get(getArgument)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var languagesJson interface{}

		err = json.Unmarshal(body, &languagesJson)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		languagesMap, ok := languagesJson.(map[string]interface{})

		if ok {
			for language := range languagesMap {
				languageCountMap[language] += languageCountMap[language] + 1
			}
		} else {
			continue
		}
	}

	// define the final limit
	if limitRequest > len(languageCountMap) {
		limitRequest = len(languageCountMap)
	}

	var languagesInfo []LanguageInfo

	for language := range languageCountMap {
		languagesInfo = append(languagesInfo, LanguageInfo{language, languageCountMap[language]})
	}

	sort.Sort(ByCount(languagesInfo))

	for i := 0; i < limitRequest; i++ {
		languagesAnswer.Languages = append(languagesAnswer.Languages, languagesInfo[i].Name)
	}

	// encoding the answer
	languagesAnswer.Auth = withAuth
	json.NewEncoder(w).Encode(languagesAnswer)
}
