package APIs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

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

func GetProjectsId(w http.ResponseWriter) []float64 {
	var ids []float64

	lastPageNotReached := true
	j := 0

	for lastPageNotReached {
		j++

		jString := strconv.Itoa(j)
		var getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects?page=%s", jString)

		resp, err := http.Get(getArgument)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return ids
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return ids
		}

		var projectsJson interface{}

		err = json.Unmarshal(body, &projectsJson)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return ids
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

			i++
		}

		if i == 0 {
			lastPageNotReached = false
		}
	}

	return ids
}

func GetProjectsIdWithAuth(w http.ResponseWriter, auth string) []float64 {
	var ids []float64

	lastPageNotReached := true
	j := 0

	for lastPageNotReached {
		j++

		jString := strconv.Itoa(j)
		var getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects?page=%s&private_token=%s", jString, auth)

		resp, err := http.Get(getArgument)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return ids
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return ids
		}

		var projectsJson interface{}

		err = json.Unmarshal(body, &projectsJson)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return ids
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

			i++
		}

		if i == 0 {
			lastPageNotReached = false
		}
	}

	return ids
}

func GetProjectsNameIdMapWithAuth(w http.ResponseWriter, auth string) map[string]float64 {
	var nameIdMap map[string]float64
	nameIdMap = make(map[string]float64)

	lastPageNotReached := true
	j := 0

	for lastPageNotReached {
		j++

		jString := strconv.Itoa(j)
		var getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects?page=%s&private_token=%s", jString, auth)

		resp, err := http.Get(getArgument)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nameIdMap
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nameIdMap
		}

		var projectsJson interface{}

		err = json.Unmarshal(body, &projectsJson)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nameIdMap
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
				if path, ok1 := projects[i]["path_with_namespace"].(string); ok1 {
					nameIdMap[path] = id
				}
			}

			i++
		}

		if i == 0 {
			lastPageNotReached = false
		}
	}

	return nameIdMap
}

func GetProjectsNameIdMap(w http.ResponseWriter) map[string]float64 {
	var nameIdMap map[string]float64
	nameIdMap = make(map[string]float64)

	lastPageNotReached := true
	j := 0

	for lastPageNotReached {
		j++

		jString := strconv.Itoa(j)
		var getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects?page=%s", jString)

		resp, err := http.Get(getArgument)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nameIdMap
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nameIdMap
		}

		var projectsJson interface{}

		err = json.Unmarshal(body, &projectsJson)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nameIdMap
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
				if path, ok1 := projects[i]["path_with_namespace"].(string); ok1 {
					nameIdMap[path] = id
				}
			}

			i++
		}

		if i == 0 {
			lastPageNotReached = false
		}
	}

	return nameIdMap
}
