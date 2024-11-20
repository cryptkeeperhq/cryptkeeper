
# Zanzibar / SpiceDB
Zanzibar is Google's purpose-built authorization system. It's a centralized authorization database built to take authorization queries from high-traffic apps and return authorization decisions. For CryptKeeper's Path/Secret authorization, we will be leveraging SpiceDB which is Google Zanzibar-inspired database for storing and computing permissions data.

https://authzed.com/spicedb

docker run -p 50051:50051 quay.io/authzed/spicedb serve --grpc-preshared-key "spicedbsecret"


## Schema Design

#### Namespaces:
- User: Represents the users in the system.
- Path: Represents the paths which contain secrets.
- Secret: Represents the secrets within paths.

#### Permissions:
- Users have permissions on the path level (list, create, update, delete, rotate) for all secrets within a path.
- Users can have explicit deny permissions on specific secrets, overriding the path-level permissions.


#### Permissions:
- Path: list permission.
- Secret: read, create, delete, update, and deny permissions. The deny relation is checked first to enforce explicit deny rules.

Schema

```
definition user {}

definition path {
  relation list: user
  relation create: user
  relation update: user
  relation delete: user
  relation rotate: user
}

definition secret {
  relation read: user | path#list
  relation create: user | path#create
  relation update: user | path#update
  relation delete: user | path#delete
  relation rotate: user | path#rotate
  relation deny: user

  permission read_allowed = read - deny
  permission create_allowed = create - deny
  permission update_allowed = update - deny
  permission delete_allowed = delete - deny
  permission rotate_allowed = rotate - deny
}
```

1. Namespaces:
- **user**: Represents a user in the system.
- **path**: Represents a path that contains secrets. It has relations for `list`, `create`, `update`, `delete`, and `rotate` permissions.
- **secret**: Represents a secret within a path. It inherits permissions from the path and has an additional `deny` relation to explicitly deny access.


2. Permissions:
- Permissions like `read_allowed`, `create_allowed`, etc., are computed by subtracting the deny relation from the respective permission relations.


## HCL
We will leverage HCL for admin/users to create/update policy. The policy updates will be synced with SpiceDB for permission checking.




```go

	// zanzibarClient.VerifyPathPermissions("/zanzibar-kv", "test1", "read")
	// zanzibarClient.VerifyRelationshipUpdates("/zanzibar-kv")
	// zanzibarClient.CheckPathAndSecretPermission("/zanzibar-kv", "/attribute/foo", "read", "test1")
	// return
	// // Example check permission
	// hasPermission, err := zanzibarClient.CheckPermission("secret-id", "read_allowed", "user-id")
	// if err != nil {
	// 	log.Printf("failed to check permission: %v\n", err)
	// }
	// fmt.Println("Has Permission:", hasPermission)

	// // Example add tuple
	// // err = zanzibarClient.AddTuple(ctx, "secret", "secret-id", "owner", "user-id")
	// zanzibarClient.AddSecretToPath("/zanzibar-kv", "foo")
	// zanzibarClient.AddSecretToPath("/zanzibar-kv", "/attribute/foo")
	// err = zanzibarClient.AddPathPermissions("path-id", "user-id")
	// if err != nil {
	// 	log.Println("failed to add tuple: %v\n", err)
	// }

	// hasPermission, err = zanzibarClient.CheckPermission("secret-id", "read_allowed", "user-id")
	// if err != nil {
	// 	log.Printf("failed to check permission: %v\n", err)
	// }
	// fmt.Println("Has Permission:", hasPermission)

	// return
  ```