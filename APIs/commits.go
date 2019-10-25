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

func floatInSlice(a float64, list []float64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

type ReposInfo struct {
	Repository string `json:"repository"`
	Commits    int    `json:"commits"`
}

// to sort by commits count at the end
type ByCommitsCount []ReposInfo

func (a ByCommitsCount) Len() int           { return len(a) }
func (a ByCommitsCount) Less(i, j int) bool { return a[i].Commits > a[j].Commits } // which is More for our purpose
func (a ByCommitsCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type CommitsAnswer struct {
	Repos []ReposInfo `json:"repos"`
	Auth  bool        `json:"auth"`
}

func GetProjectsInfo(w http.ResponseWriter) ([]float64, []string) {
	var ids []float64
	var paths []string

	lastPageNotReached := true
	j := 0

	for lastPageNotReached {
		j++

		jString := strconv.Itoa(j)
		var getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects?page=%s", jString)

		resp, err := http.Get(getArgument)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return ids, paths
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return ids, paths
		}

		var projectsJson interface{}

		err = json.Unmarshal(body, &projectsJson)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return ids, paths
		}

		var projectsResult = projectsJson.([]interface{})

		var lenProjects = len(projectsResult)

		projects := make(map[int]map[string]interface{})

		for i := 0; i < lenProjects; i++ {
			projects[i] = projectsResult[i].(map[string]interface{})
		}

		i := 0
		for i < lenProjects {
			if id, ok := projects[i]["id"].(float64); ok {
				ids = append(ids, id)
			}
			if path, ok1 := projects[i]["path_with_namespace"].(string); ok1 {
				paths = append(paths, path)
			}

			i++
		}

		if i == 0 {
			lastPageNotReached = false
		}
	}

	return ids, paths
}

func GetProjectsInfoWithAuth(w http.ResponseWriter, auth string) ([]float64, []string) {
	var ids []float64
	var paths []string

	lastPageNotReached := true
	j := 0

	for lastPageNotReached {
		j++

		jString := strconv.Itoa(j)
		var getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects?page=%s&private_token=%s", jString, auth)

		resp, err := http.Get(getArgument)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return ids, paths
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return ids, paths
		}

		var projectsJson interface{}

		err = json.Unmarshal(body, &projectsJson)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return ids, paths
		}

		var projectsResult = projectsJson.([]interface{})

		var lenProjects = len(projectsResult)

		projects := make(map[int]map[string]interface{})

		for i := 0; i < lenProjects; i++ {
			projects[i] = projectsResult[i].(map[string]interface{})
		}

		i := 0
		for i < lenProjects {
			if id, ok := projects[i]["id"].(float64); ok {
				ids = append(ids, id)
			}
			if path, ok1 := projects[i]["path_with_namespace"].(string); ok1 {
				paths = append(paths, path)
			}

			i++
		}

		if i == 0 {
			lastPageNotReached = false
		}
	}

	return ids, paths
}

var commitsAnswer = CommitsAnswer{}

func HandlerCommits(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 4 || parts[3] != "commits" {
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

	var projectsId []float64
	var paths []string
	if withAuth {
		projectsId, paths = GetProjectsInfoWithAuth(w, authRequest)
	} else {
		projectsId, paths = GetProjectsInfo(w)
	}

	var commitsCounts []int

	// for each project, we get the number of commits
	lenProjectsId := len(projectsId)
	for i := 0; i < lenProjectsId; i++ {
		commitsCount := 0

		lastPageNotReached := true
		j := 0

		for lastPageNotReached {
			j++

			jString := strconv.Itoa(j)
			id := strconv.Itoa(int(projectsId[i]))

			var getArgument string
			if withAuth {
				getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects/%s/repository/commits?page=%s&private_token=%s", id, jString, authRequest)
			} else {
				getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects/%s/repository/commits?page=%s", id, jString)
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

			var commitsJson interface{}

			err = json.Unmarshal(body, &commitsJson)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			commitsMap, ok := commitsJson.([]interface{})

			if ok {
				var lenCommitsMap = len(commitsMap)

				commitsCount += lenCommitsMap

				if lenCommitsMap == 0 {
					lastPageNotReached = false
				}
			} else {
				lastPageNotReached = false
			}
		}

		commitsCounts = append(commitsCounts, commitsCount)
	}

	// define the final limit
	if limitRequest > lenProjectsId {
		limitRequest = lenProjectsId
	}

	var repos []ReposInfo

	for i := 0; i < lenProjectsId; i++ {
		repos = append(repos, ReposInfo{paths[i], commitsCounts[i]})
	}

	sort.Sort(ByCommitsCount(repos))

	commitsAnswer.Repos = repos[0:limitRequest]

	// encoding the answer
	commitsAnswer.Auth = withAuth
	json.NewEncoder(w).Encode(commitsAnswer)
}
