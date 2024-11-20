package zanzibar

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"github.com/cryptkeeperhq/cryptkeeper/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Zanzibar struct {
	client *authzed.Client
	ctx    context.Context
	logger *slog.Logger
}

func NewZanzibar(config *config.Config) (*Zanzibar, error) {
	client, err := getSpiceDbClient(config.Zanzibar.Endpoint, config.Zanzibar.ApiKey)
	if err != nil {
		return nil, err
	}

	return &Zanzibar{
		client: client,
		ctx:    context.Background(),
		logger: config.Logger,
	}, nil
}

func getSpiceDbClient(endpoint, presharedKey string) (*authzed.Client, error) {
	opts := []grpc.DialOption{
		grpcutil.WithInsecureBearerToken(presharedKey),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	return authzed.NewClient(endpoint, opts...)
}

func (z *Zanzibar) WriteSchema() error {
	schema := `
	definition user {}

	definition group {
		relation member: user
	}

	definition path {
		relation list: user | group#member
		relation read: user | group#member
		relation create: user | group#member
		relation update: user | group#member
		relation delete: user | group#member
		relation rotate: user | group#member

		permission list_secret = list
		permission read_secret = read
		permission create_secret = create
		permission update_secret = update
		permission delete_secret = delete
		permission rotate_secret = rotate
	}

	definition secret {
		relation parent_path: path
		relation deny_list: user | group#member
		relation deny_read: user | group#member
		relation deny_create: user | group#member
		relation deny_update: user | group#member
		relation deny_delete: user | group#member
		relation deny_rotate: user | group#member

		permission list_allowed = parent_path->list_secret - deny_list
		permission read_allowed = parent_path->read_secret - deny_read
		permission create_allowed = parent_path->create_secret - deny_create
		permission update_allowed = parent_path->update_secret - deny_update
		permission delete_allowed = parent_path->delete_secret - deny_delete
		permission rotate_allowed = parent_path->rotate_secret - deny_rotate
	}
	`

	_, err := z.client.WriteSchema(context.Background(), &v1.WriteSchemaRequest{Schema: schema})
	if err != nil {
		z.logger.Error(fmt.Sprintf("Failed to apply schema: %v", err))
	} else {
		z.logger.Info("Schema applied successfully")
	}

	return err
}

func (z *Zanzibar) CheckSecretDenyPermission(secretID, denyPermission, userID string) (bool, error) {
	z.logger.Debug(fmt.Sprintf("Checking [%s] permissions for [%s] on Secret [%s]", denyPermission, userID, secretID))

	checkResp, err := z.client.CheckPermission(z.ctx, &v1.CheckPermissionRequest{
		Resource:   &v1.ObjectReference{ObjectType: "secret", ObjectId: secretID},
		Permission: denyPermission,
		Subject:    &v1.SubjectReference{Object: &v1.ObjectReference{ObjectType: "user", ObjectId: userID}},
	})
	if err != nil {
		return false, err
	}

	hasDenyPermission := checkResp.Permissionship == v1.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION
	if hasDenyPermission {
		z.logger.Debug(fmt.Sprintf("User %s is denied %s permission on secret %s", userID, denyPermission, secretID))
	} else {
		z.logger.Debug(fmt.Sprintf("User %s is not denied %s permission on secret %s", userID, denyPermission, secretID))
	}

	return hasDenyPermission, nil
}

func (z *Zanzibar) CheckPathPermission(pathID, permission, userID string) (bool, error) {
	z.logger.Debug(fmt.Sprintf("Checking [%s] permissions for [%s] on Path [%s]", permission, userID, pathID))

	checkResp, err := z.client.CheckPermission(z.ctx, &v1.CheckPermissionRequest{
		Resource:   &v1.ObjectReference{ObjectType: "path", ObjectId: pathID},
		Permission: permission,
		Subject:    &v1.SubjectReference{Object: &v1.ObjectReference{ObjectType: "user", ObjectId: userID}},
	})
	if err != nil {
		return false, err
	}

	hasPermission := checkResp.Permissionship == v1.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION
	if hasPermission {
		z.logger.Info(fmt.Sprintf("User %s has %s permission on path %s", userID, permission, pathID))
	} else {
		z.logger.Info(fmt.Sprintf("User %s does not have %s permission on path %s", userID, permission, pathID))
	}

	return hasPermission, nil
}

func (z *Zanzibar) CheckPathAndSecretPermission(pathID, secretID, permission, userID string) (bool, error) {
	hasPathPermission, err := z.CheckPathPermission(pathID, permission, userID)
	if err != nil || !hasPathPermission {
		return hasPathPermission, err
	}

	if secretID != "" {
		hasDenyPermission, err := z.CheckSecretDenyPermission(secretID, "deny_"+permission, userID)
		if err != nil {
			return false, err
		}
		if hasDenyPermission {
			return false, nil // Denied explicitly
		}
	}

	return true, nil // Allowed
}

func (z *Zanzibar) AddPathPermissions(pathID, userID string, permissions []string) error {
	var updates []*v1.RelationshipUpdate
	for _, perm := range permissions {
		z.logger.Debug(fmt.Sprintf("Adding [%s] permissions to [%s] for user [%s]", perm, pathID, userID))

		updates = append(updates, &v1.RelationshipUpdate{
			Operation: v1.RelationshipUpdate_OPERATION_CREATE,
			Relationship: &v1.Relationship{
				Resource: &v1.ObjectReference{ObjectType: "path", ObjectId: pathID},
				Relation: perm,
				Subject:  &v1.SubjectReference{Object: &v1.ObjectReference{ObjectType: "user", ObjectId: userID}},
			},
		})
	}

	_, err := z.client.WriteRelationships(z.ctx, &v1.WriteRelationshipsRequest{Updates: updates})
	return err
}

func (z *Zanzibar) AddPathGroupPermissions(pathID, groupID string, permissions []string) error {
	fmt.Printf("Adding permissions to [%s] for group [%s]\n", pathID, groupID)
	var updates []*v1.RelationshipUpdate
	for _, perm := range permissions {
		updates = append(updates, &v1.RelationshipUpdate{
			Operation: v1.RelationshipUpdate_OPERATION_TOUCH,
			Relationship: &v1.Relationship{
				Resource: &v1.ObjectReference{ObjectType: "path", ObjectId: pathID},
				Relation: perm,
				Subject: &v1.SubjectReference{
					Object:           &v1.ObjectReference{ObjectType: "group", ObjectId: groupID},
					OptionalRelation: "member",
				},
			},
		})
	}

	resp, err := z.client.WriteRelationships(z.ctx, &v1.WriteRelationshipsRequest{Updates: updates})
	if err != nil {
		fmt.Printf("Failed to write relationships: %v\n", err)
	}
	fmt.Printf("WriteRelationships response: %+v\n", resp)
	return err
}

func (z *Zanzibar) AddUserToGroup(userID, groupID string) error {
	_, err := z.client.WriteRelationships(z.ctx, &v1.WriteRelationshipsRequest{
		Updates: []*v1.RelationshipUpdate{
			{
				Operation: v1.RelationshipUpdate_OPERATION_CREATE,
				Relationship: &v1.Relationship{
					Resource: &v1.ObjectReference{ObjectType: "group", ObjectId: groupID},
					Relation: "member",
					Subject:  &v1.SubjectReference{Object: &v1.ObjectReference{ObjectType: "user", ObjectId: userID}},
				},
			},
		},
	})
	if err != nil {
		z.logger.Error(fmt.Sprintf("Failed to add user to group: %v", err))
	}
	z.logger.Info(fmt.Sprintf("---> Added user %s to group %s", userID, groupID))
	return err
}

func (z *Zanzibar) VerifyGroupMembership(userID, groupID string) {
	checkResp, err := z.client.CheckPermission(z.ctx, &v1.CheckPermissionRequest{
		Resource:   &v1.ObjectReference{ObjectType: "group", ObjectId: groupID},
		Permission: "member",
		Subject:    &v1.SubjectReference{Object: &v1.ObjectReference{ObjectType: "user", ObjectId: userID}},
	})
	if err != nil {
		log.Fatalf("Failed to verify group membership: %v", err)
	}

	if checkResp.Permissionship == v1.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION {
		fmt.Printf("User %s is a member of group %s\n", userID, groupID)
	} else {
		fmt.Printf("User %s is NOT a member of group %s\n", userID, groupID)
	}
}

func (z *Zanzibar) VerifyGroupPermissions(pathID, groupID, permission string) {
	fmt.Printf("Verifying permissions for group [%s] on path [%s] with permission [%s]\n", groupID, pathID, permission)
	checkResp, err := z.client.CheckPermission(z.ctx, &v1.CheckPermissionRequest{
		Resource:   &v1.ObjectReference{ObjectType: "path", ObjectId: pathID},
		Permission: permission,
		Subject:    &v1.SubjectReference{Object: &v1.ObjectReference{ObjectType: "group", ObjectId: groupID}},
	})
	if err != nil {
		log.Fatalf("Failed to verify group permissions: %v", err)
	}

	fmt.Printf("CheckPermission response: %+v\n", checkResp)
	if checkResp.Permissionship == v1.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION {
		fmt.Printf("Group %s has %s permission on path %s\n", groupID, permission, pathID)
	} else {
		fmt.Printf("Group %s does NOT have %s permission on path %s\n", groupID, permission, pathID)
	}
}

func (z *Zanzibar) DenySecretPermission(secretID, userID, denyRelation string) error {
	fmt.Printf("Adding Deny Permissions to [%s] for user [%s]\n", secretID, userID)
	_, err := z.client.WriteRelationships(z.ctx, &v1.WriteRelationshipsRequest{
		Updates: []*v1.RelationshipUpdate{
			{
				Operation: v1.RelationshipUpdate_OPERATION_CREATE,
				Relationship: &v1.Relationship{
					Resource: &v1.ObjectReference{ObjectType: "secret", ObjectId: secretID},
					Relation: denyRelation,
					Subject:  &v1.SubjectReference{Object: &v1.ObjectReference{ObjectType: "user", ObjectId: userID}},
				},
			},
		},
	})
	if err != nil {
		log.Println("Failed to deny permission: %v", err)
	}

	fmt.Printf("Denied %s permission for user %s on secret %s\n", denyRelation, userID, secretID)
	return err
}

func (z *Zanzibar) AddSecretToPath(pathID, secretID string) error {
	_, err := z.client.WriteRelationships(z.ctx, &v1.WriteRelationshipsRequest{
		Updates: []*v1.RelationshipUpdate{
			{
				Operation: v1.RelationshipUpdate_OPERATION_CREATE,
				Relationship: &v1.Relationship{
					Resource: &v1.ObjectReference{ObjectType: "secret", ObjectId: secretID},
					Relation: "parent_path",
					Subject:  &v1.SubjectReference{Object: &v1.ObjectReference{ObjectType: "path", ObjectId: pathID}},
				},
			},
		},
	})
	return err
}

func (z *Zanzibar) VerifyPathPermissions(pathID, userID, permission string) {
	fmt.Printf("Verifying permissions for user [%s] on path [%s] with permission [%s]\n", userID, pathID, permission)
	checkResp, err := z.client.CheckPermission(z.ctx, &v1.CheckPermissionRequest{
		Resource:   &v1.ObjectReference{ObjectType: "path", ObjectId: pathID},
		Permission: permission,
		Subject:    &v1.SubjectReference{Object: &v1.ObjectReference{ObjectType: "user", ObjectId: userID}},
	})
	if err != nil {
		log.Fatalf("Failed to verify path permissions: %v", err)
	}

	fmt.Printf("CheckPermission response: %+v\n", checkResp)
	if checkResp.Permissionship == v1.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION {
		fmt.Printf("User %s has %s permission on path %s\n", userID, permission, pathID)
	} else {
		fmt.Printf("User %s does NOT have %s permission on path %s\n", userID, permission, pathID)
	}
}

func (z *Zanzibar) VerifyRelationshipUpdates(path string) {
	relationships, err := z.client.ReadRelationships(context.TODO(), &v1.ReadRelationshipsRequest{
		RelationshipFilter: &v1.RelationshipFilter{
			ResourceType:       "path",
			OptionalResourceId: path,
			OptionalSubjectFilter: &v1.SubjectFilter{
				SubjectType:       "user",
				OptionalSubjectId: "test1",
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to read relationships: %v", err)
	}

	for {
		item, err := relationships.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			z.logger.Error(err.Error())
			break
		}

		fmt.Println(
			item.GetRelationship().Resource.ObjectType, " ",
			item.GetRelationship().Resource.ObjectId, " ",
			item.GetRelationship().Subject.Object.ObjectType, " ",
			item.GetRelationship().Subject.Object.ObjectId, " ",
			item.GetRelationship().Relation, " ",
		)
	}
}

// func (z *Zanzibar) AddTuple(ctx context.Context, namespace, objectID, relation, userID string) error {
// 	// return nil
// 	_, err := z.client.WriteRelationships(z.ctx, &v1.WriteRelationshipsRequest{
// 		Updates: []*v1.RelationshipUpdate{
// 			{
// 				Operation: v1.RelationshipUpdate_OPERATION_CREATE,
// 				Relationship: &v1.Relationship{
// 					Resource: &v1.ObjectReference{
// 						ObjectType: "path",
// 						ObjectId:   "path1",
// 					},
// 					Relation: "list",
// 					Subject: &v1.SubjectReference{
// 						Object: &v1.ObjectReference{
// 							ObjectType: "user",
// 							ObjectId:   "user1",
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		log.Println("Failed to write relationship: %v", err)
// 	}

// 	_, err = z.client.WriteRelationships(z.ctx, &v1.WriteRelationshipsRequest{
// 		Updates: []*v1.RelationshipUpdate{
// 			{
// 				Operation: v1.RelationshipUpdate_OPERATION_CREATE,
// 				Relationship: &v1.Relationship{
// 					Resource: &v1.ObjectReference{
// 						ObjectType: "secret",
// 						ObjectId:   "secret1",
// 					},
// 					Relation: "create",
// 					Subject: &v1.SubjectReference{
// 						Object: &v1.ObjectReference{
// 							ObjectType: "user",
// 							ObjectId:   "user1",
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		log.Println("Failed to write relationship: %v", err)
// 	}

// 	// Example: Explicitly deny "user1" from creating "secret1"
// 	_, err = z.client.WriteRelationships(ctx, &v1.WriteRelationshipsRequest{
// 		Updates: []*v1.RelationshipUpdate{
// 			{
// 				Operation: v1.RelationshipUpdate_OPERATION_CREATE,
// 				Relationship: &v1.Relationship{
// 					Resource: &v1.ObjectReference{
// 						ObjectType: "secret",
// 						ObjectId:   "secret1",
// 					},
// 					Relation: "deny",
// 					Subject: &v1.SubjectReference{
// 						Object: &v1.ObjectReference{
// 							ObjectType: "user",
// 							ObjectId:   "user1",
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		log.Println("Failed to write relationship: %v", err)
// 	}

// 	// Check if "user1" has read permission on "secret1"
// 	checkResp, err := z.client.CheckPermission(ctx, &v1.CheckPermissionRequest{
// 		Resource: &v1.ObjectReference{
// 			ObjectType: "secret",
// 			ObjectId:   "secret1",
// 		},
// 		Permission: "read_allowed",
// 		Subject: &v1.SubjectReference{
// 			Object: &v1.ObjectReference{
// 				ObjectType: "user",
// 				ObjectId:   "user1",
// 			},
// 		},
// 	})
// 	if err != nil {
// 		log.Println("Failed to check permission: %v", err)
// 	}

// 	if checkResp.Permissionship == v1.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION {
// 		fmt.Println("User has read permission")
// 	} else {
// 		fmt.Println("User does not have read permission")
// 	}

// 	return err
// }

// func (z *Zanzibar) CheckPermission(ctx context.Context, namespace, objectID, relation, userID string) (bool, error) {
// 	// request := &v1.SubjectReference{Object: &v1.ObjectReference{
// 	// 	ObjectType: "user",
// 	// 	ObjectId:   "user1",
// 	// }}

// 	// firstPost := &v1.ObjectReference{
// 	// 	ObjectType: "secret",
// 	// 	ObjectId:   "secret1",
// 	// }

// 	// response, err := z.client.CheckPermission(context.Background(), &v1.CheckPermissionRequest{
// 	// 	Resource:   firstPost,
// 	// 	Permission: "read",
// 	// 	Subject:    request,
// 	// })

// 	response, err := z.client.CheckPermission(ctx, &v1.CheckPermissionRequest{
// 		Resource: &v1.ObjectReference{
// 			ObjectType: "secret",
// 			ObjectId:   "secret1",
// 		},
// 		Permission: "read_allowed",
// 		Subject: &v1.SubjectReference{
// 			Object: &v1.ObjectReference{
// 				ObjectType: "user",
// 				ObjectId:   "user1",
// 			},
// 		},
// 	})

// 	fmt.Println(response)
// 	// response, err := z.client.CheckPermission(ctx, request)
// 	if err != nil {
// 		fmt.Println(err)
// 		return false, err
// 	}

// 	return response.Permissionship == v1.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION, nil
// }

// // List all permissions for the given path and user
// checkResp, err := z.client.LookupResources(context.TODO(), &v1.LookupResourcesRequest{
// 	ResourceObjectType: "path",
// 	Permission:         permission,
// 	Subject: &v1.SubjectReference{
// 		Object: &v1.ObjectReference{
// 			ObjectType: "user",
// 			ObjectId:   userID,
// 		},
// 	},
// })

// if err != nil {
// 	logger.Error(fmt.Sprint("failed to lookup resources: %v", err))
// 	return
// }

// var subjects []string
// for {
// 	item, err := checkResp.Recv()
// 	if err == io.EOF {
// 		break
// 	}
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(item.GetResourceObjectId(), " ", item.ResourceObjectId, permission, item.Permissionship)
// 	subjects = append(subjects, item.GetResourceObjectId())
// }

// fmt.Println(subjects)
