package APIs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

		var getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects?page=%c", j)

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

		var projectsMap = projectsJson.(map[string]interface{})

		var lenProjectsMap = len(projectsMap)

		i := 0
		for i < lenProjectsMap {
			if id, ok := projectsMap["id"].(float64); ok {
				ids = append(ids, id)
			}
			if path, ok1 := projectsMap["path_with_namespace"].(string); ok1 {
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

		var getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects?page=%c&private_token=%s", j, auth)

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

		var projectsMap = projectsJson.(map[string]interface{})

		var lenProjectsMap = len(projectsMap)

		i := 0
		for i < lenProjectsMap {
			if id, ok := projectsMap["id"].(float64); ok {
				ids = append(ids, id)
			}
			if path, ok1 := projectsMap["path_with_namespace"].(string); ok1 {
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

			var getArgument string
			if withAuth {
				getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects/%f/commits?page=%c&private_token=%s", projectsId[i], j, authRequest)
			} else {
				getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects/%f/commits?page=%c", projectsId[i], j)
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

			var commitsMap = commitsJson.(map[string]interface{})

			var lenCommitsMap = len(commitsMap)

			commitsCount += lenCommitsMap

			if lenCommitsMap == 0 {
				lastPageNotReached = false
			}
		}

		commitsCounts = append(commitsCounts, commitsCount)
	}

	// define the final limit
	if limitRequest > lenProjectsId {
		limitRequest = lenProjectsId
	}

	for i := 0; i < lenProjectsId; i++ {
		commitsAnswer.Repos = append(commitsAnswer.Repos, ReposInfo{paths[i], commitsCounts[i]})
	}

	// encoding the answer
	commitsAnswer.Auth = withAuth
	json.NewEncoder(w).Encode(commitsAnswer)
}
