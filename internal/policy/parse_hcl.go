package policy

import (
	"fmt"
	"log"

	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

/*
path "path.com" {
		permissions = ["list", "read", "create", "update", "delete", "rotate"]
		deny_permissions = []
		users       = []
		certificates = []
		apps = ["app-50fee5e4-6eb6-49be-81b3-ef478f44d259"]
		groups = ["admins"]
}
*/

/*
	{
		"keyId": "key123",
		"statements": [
		  {
			"principal": "user:alice",
			"actions": ["encrypt", "decrypt"],
			"effect": "allow",
			"conditions": {
			  "ipAddress": "192.168.1.0/24"
			}
		  },
		  {
			"principal": "role:admin",
			"actions": ["rotate", "delete"],
			"effect": "allow"
		  }
		]
	  }
*/
func ParseHCLPolicy(hclString string) (models.Policy, error) {

	parser := hclparse.NewParser()

	fmt.Println(hclString)

	file, diags := parser.ParseHCL([]byte(hclString), "policy.hcl")
	if diags.HasErrors() {
		log.Printf("failed to parse HCL file: %v\n", diags)
		return models.Policy{}, diags
	}

	var policy models.Policy
	diags = gohcl.DecodeBody(file.Body, nil, &policy)
	if diags.HasErrors() {
		log.Printf("failed to parse HCL file: %v\n", diags)
		return models.Policy{}, diags
	}

	policy.HCL = hclString

	return policy, nil
}

func GetDefaultHCLPolicy(path string) (models.Policy, error) {
	hclString := fmt.Sprintf(`path "%s" {
		permissions = ["list", "read", "create", "update", "delete", "rotate"]
		deny_permissions = []
		users       = []
		certificates = []
		apps = []
		groups = ["admins"]
	}`, path)

	return ParseHCLPolicy(hclString)

}
