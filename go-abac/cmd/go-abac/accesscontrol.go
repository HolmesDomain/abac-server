package main

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	mongodbadapter "github.com/casbin/mongodb-adapter/v2"
)

// ABAC model string
var Text = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub_rule, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = eval(p.sub_rule) && r.obj == p.obj && r.act == p.act
`
var Model, _ = model.NewModelFromString(Text)
var Mongo = goDotEnvVariable("MONGO_ENDPOINT")

func abac(user, resource, action string) bool {
	// get JSON file in MongoDB using the proper adapter
	// looks for casbin_rule
	adapter, err := mongodbadapter.NewAdapter(Mongo + "/PolicyDatabase")
	enforcer, _ := casbin.NewEnforcer(Model, adapter)

	type subject struct {
		Label string `json:"label"`
	}

	sub := subject{user}

	if err != nil {
		panic(err)
	}

	q := enforcer.GetPolicy()

	log.Println(q)

	ok, err := enforcer.Enforce(sub, resource, action)

	if err != nil {
		// handle err
	}

	if ok == true {
		return true
	} else {
		return false
	}
}

func getThePolicy() [][]string {
	adapter, err := mongodbadapter.NewAdapter(Mongo + "/PolicyDatabase")
	enforcer, _ := casbin.NewEnforcer(Model, adapter)

	if err != nil {
		panic(err)
	}

	policy := enforcer.GetPolicy()

	return policy
}
