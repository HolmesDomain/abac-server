package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func getPolicyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	results := getThePolicyV2()

	if results == nil {
		var message = "Policy rule database is empty."
		json.NewEncoder(w).Encode(message)
	} else {
		json.NewEncoder(w).Encode(results)
	}
}

func querySinglePolicyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	var post Post

	json.NewDecoder(r.Body).Decode(&post)

	obj := post.Object
	sub := post.Label
	act := post.Action

	results := getSinglePolicy(sub, obj, act)

	if results == nil {
		var message = "No results found."
		json.NewEncoder(w).Encode(message)
	} else {
		json.NewEncoder(w).Encode(results)
	}
}

func postPolicyHandler(w http.ResponseWriter, r *http.Request) {
	var policy PolicyRule
	json.NewDecoder(r.Body).Decode(&policy)

	success := AccessResponse{"00", "success", getLog()}
	fail := AccessResponse{"01", "error", getLog()}

	if postPolicy(policy.PType, policy.Evaluator, policy.Object, policy.Action) == true {
		json.NewEncoder(w).Encode(success)
	} else {
		json.NewEncoder(w).Encode(fail)
	}
}

func queryPolicyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	var post Post

	json.NewDecoder(r.Body).Decode(&post)

	sub := post.Label

	results := getMatchingPolicy(sub)

	if results == nil {
		var message = "No results found."
		json.NewEncoder(w).Encode(message)
	} else {
		json.NewEncoder(w).Encode(results)
	}
}

func postEnforce(w http.ResponseWriter, r *http.Request) {
	var post Post
	json.NewDecoder(r.Body).Decode(&post)

	obj := post.Object
	sub := post.Label
	act := post.Action

	//adding backslash delineation matching policy-rule file
	newObj := "/" + obj

	//Error code struct
	errorMssg := ErrResponse{"Error has occured"}

	//Initialize struct for ABAC JSON responses
	success := AccessResponse{"00", "granted", getAccessLog(sub, obj, act, true)}
	fail := AccessResponse{"01", "denied", getAccessLog(sub, obj, act, false)}

	result := abac(sub, newObj, act)

	if result == true {
		json.NewEncoder(w).Encode(success)
	} else if result == false {
		json.NewEncoder(w).Encode(fail)
	} else {
		err := json.NewEncoder(w).Encode(errorMssg)
		panic(err)
	}
}

func paramEnforce(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	//Store URL parameters
	vars := mux.Vars(r)
	obj := vars["obj"]
	sub := vars["sub"]
	act := vars["act"]

	//adding backslash delineation matching policy-rule file
	newObj := "/" + obj

	//Error code struct
	errorMssg := ErrResponse{"Error has occured"}

	//Initialize struct for ABAC JSON responses
	successResp := AccessResponse{"00", "granted", getAccessLog(sub, obj, act, true)}
	failResp := AccessResponse{"01", "denied", getAccessLog(sub, obj, act, false)}

	result := abac(sub, newObj, act)

	if result == true {
		json.NewEncoder(w).Encode(successResp)
	} else if result == false {
		json.NewEncoder(w).Encode(failResp)
	} else {
		err := json.NewEncoder(w).Encode(errorMssg)
		panic(err)
	}
}

func removePolicyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	var delete Delete
	json.NewDecoder(r.Body).Decode(&delete)

	id := delete.ID

	//Initialize struct for ABAC JSON responses
	success := AccessResponse{"00", "success", getLog()}
	fail := AccessResponse{"01", "error", getLog()}

	if removePolicy(id) == true {
		json.NewEncoder(w).Encode(success)
	} else {
		json.NewEncoder(w).Encode(fail)
	}
}

type AccessResponse struct {
	Code   string `json:"code"`
	Result string `json:"result"`
	Log    string `json:"log"`
}

type Post struct {
	Label  string `json:"label"`
	Object string `json:"object"`
	Action string `json:"action"`
}

type Delete struct {
	ID string `json:"_id,omitempty"`
}

type ErrResponse struct {
	Error string `json:"error"`
}

type PolicyRule struct {
	PType     string `json:ptype,omitempty”`
	Evaluator string `json:evaluator,omitempty”`
	Object    string `json:object,omitempty”`
	Action    string `json:action,omitempty”`
}

type CasbinModel struct {
	PType string `json:PType,omitempty”`
	V0    string `json:V0,omitempty”`
	V1    string `json:V1,omitempty”`
	V2    string `json:V2,omitempty”`
}
