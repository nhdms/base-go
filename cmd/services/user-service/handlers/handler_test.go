package handlers

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/mitchellh/mapstructure"
	"github.com/nhdms/base-go/pkg/dbtool"
	"github.com/nhdms/base-go/proto/exmsg/models"
	"github.com/nhdms/base-go/proto/exmsg/services"
	"github.com/nhdms/base-go/tests"
	"github.com/spf13/cast"
	"log"
	"testing"
)

var userHandler *UserHandler
var ctx = context.Background()

func init() {
	err := tests.LoadTestConfig()
	if err != nil {
		log.Fatal("Failed to load config", err)
	}

	psql, err := dbtool.NewConnectionManager(dbtool.DBTypePostgreSQL, nil)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	redis, err := dbtool.CreateRedisConnection(nil)

	userHandler = NewUserHandler(psql, redis)
}

func TestGetUserByID(t *testing.T) {
	// Test GetUserByID
	resp := services.UserResponse{}
	err := userHandler.GetUserByID(ctx, &services.UserRequest{
		UserId: 2385,
	}, &resp)
	if err != nil {
		t.Fatal("Failed to get user by ID", err)
	}
	t.Log("User: ", resp.User)
}

func TestGetProfiles(t *testing.T) {
	// Test GetProfiles
	resp := services.ProfileResponse{}
	err := userHandler.GetProfiles(ctx, &services.ProfileRequest{
		//UserId: 2385,
		AncestorDepartmentIds: []int64{104},
		Query: &models.Query{
			Limit: 250,
			Page:  0,
		},
	}, &resp)
	if err != nil {
		t.Fatal("Failed to get profiles", err)
	}
	t.Log("Profiles: ", resp.Profiles)
}

func TestUserHandler_GetDataSetScopes(t *testing.T) {
	// Test GetUserDataSetScopes
	resp := services.DataSetScopeResponse{}
	err := userHandler.GetDataSetScopes(ctx, &services.DataSetScopeRequest{
		DepartmentIds: []int64{249, 251},
	}, &resp)
	if err != nil {
		t.Fatal("Failed to get data set scopes", err)
	}
	t.Log("DataSetScopes: ", resp.Scopes)
}

func TestUserHandler_GetRoles(t *testing.T) {
	// Test GetRoles
	resp := services.RoleResponse{}
	err := userHandler.GetRoles(ctx, &services.RoleRequest{
		RoleIds: []int64{628, 537},
	}, &resp)
	if err != nil {
		t.Fatal("Failed to get roles", err)
	}
	t.Log("Roles: ", resp.Roles)
}

func TestNewUserHandler(t *testing.T) {
	type str struct {
		A string  `json:"a,omitempty"`
		B *string `json:"b,omitempty"`
	}

	a := str{
		A: "1",
		B: nil,
	}

	var m map[string]interface{}
	mapstructure.Decode(a, &m)
	b, _ := json.Marshal(m)
	t.Log(string(b))

}

func TestUserHandler_GetProjects(t *testing.T) {
	resp := services.ProjectResponse{}
	err := userHandler.GetProjects(ctx, &services.ProjectRequest{
		ProjectIds: []int64{355},
	}, &resp)
	if err != nil {
		t.Fatal("Failed to get project", err)
	}
	t.Log("Project: ", resp.Projects)
	enableAfterSales := false
	settings := resp.Projects[0].GetSettings()

	enableAfterSales = settings != nil && cast.ToBool(settings.AsMap()["enable_after_sale"])
	if !enableAfterSales {
		t.Fatal("Enable after sale is not enabled. drop message")
	}
	t.Log("Enable after sale")
}
