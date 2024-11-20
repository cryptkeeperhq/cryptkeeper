package keycloak

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
)

var (
	KeycloakURL = "http://localhost:9999"
	ClientID    = "apiClient"
	ClientSecet = "4NZ9L23PRPEo51yZnclWCbXEyWU70rgU"
	Realm       = "myrealm"
)

type KAuth struct {
	Client    *gocloak.GoCloak
	Ctx       context.Context
	Token     *gocloak.JWT
	ExpiresAt time.Time
}

func Init() (*KAuth, error) {
	client := gocloak.NewClient(KeycloakURL)
	ctx := context.Background()

	token, err := client.LoginAdmin(ctx, "vdparikh", "password", Realm)
	if err != nil {
		panic(err)
	}

	return &KAuth{
		Client:    client,
		Token:     token,
		Ctx:       ctx,
		ExpiresAt: time.Now().Add(time.Duration(token.ExpiresIn) * time.Second),
	}, nil
}

func (k *KAuth) RefreshToken() {

	// fmt.Println(k.Token.ExpiresIn)
	// expiresIn := time.Now().Add(time.Duration(k.Token.ExpiresIn) * time.Second)

	fmt.Println(k.ExpiresAt)
	fmt.Println(time.Now().Add(5 * time.Second))
	if k.ExpiresAt.After(time.Now().Add(5 * time.Second)) {
		fmt.Println("---> TOKEN VALID")
		return
	}

	fmt.Println("---> REFRESHING TOKEN")
	k.Token, _ = k.Client.LoginAdmin(k.Ctx, "vdparikh", "password", Realm)
	k.ExpiresAt = time.Now().Add(time.Duration(k.Token.ExpiresIn) * time.Second)
}

func (k *KAuth) Login(username, password string) (*gocloak.JWT, error) {

	return k.Client.GetToken(k.Ctx, Realm, gocloak.TokenOptions{
		ClientID:     &ClientID,
		ClientSecret: &ClientSecet,
		GrantType:    gocloak.StringP("password"),
		Username:     &username,
		Password:     &password,
		Scope:        gocloak.StringP("openid"),
	})
}

func (k *KAuth) CreateUser(user models.User) error {
	k.RefreshToken()

	keycloackUser := gocloak.User{
		FirstName:     gocloak.StringP(user.Name),
		LastName:      gocloak.StringP(""),
		EmailVerified: gocloak.BoolP(true),
		Credentials: &[]gocloak.CredentialRepresentation{
			{
				Temporary: gocloak.BoolP(false),
				Type:      gocloak.StringP("password"),
				Value:     gocloak.StringP(user.Password),
			},
		},
		Email:    gocloak.StringP(user.Email),
		Enabled:  gocloak.BoolP(true),
		Username: gocloak.StringP(user.Username),
	}

	x, err := k.Client.CreateUser(k.Ctx, k.Token.AccessToken, Realm, keycloackUser)
	fmt.Println(x)

	return err

}

func (k *KAuth) GetUserInfo(token string) (*gocloak.UserInfo, error) {

	client := gocloak.NewClient(KeycloakURL)
	ctx := context.Background()
	rptResult, err := client.RetrospectToken(ctx, token, ClientID, ClientSecet, Realm)
	if err != nil {
		return nil, err
	}

	if !*rptResult.Active {
		return nil, errors.New("invalid token")
	}

	return client.GetUserInfo(ctx, token, Realm)
}

func (k *KAuth) GetUserGroupsAndRoles(groups []*gocloak.Group) map[string][]string {

	// groups, err := k.Client.GetGroups(k.Ctx, k.Token.AccessToken, Realm, gocloak.GetGroupsParams{})

	gR := make(map[string][]string)

	for _, group := range groups {
		roles, _ := k.Client.GetRoleMappingByGroupID(k.Ctx, k.Token.AccessToken, Realm, *group.ID)
		gR[*group.ID] = []string{}

		if roles.RealmMappings != nil {
			for _, role := range *roles.RealmMappings {
				gR[*group.ID] = append(gR[*group.ID], *role.Name)
			}
		}
	}
	return gR
}
func (k *KAuth) CreateGroup(groupName string) error {
	k.RefreshToken()

	_, err := k.Client.CreateGroup(k.Ctx, k.Token.AccessToken, Realm, gocloak.Group{
		Name: gocloak.StringP(groupName),
	})

	return err
}

func (k *KAuth) AddUserToGroup(userID, groupID string) error {
	k.RefreshToken()

	return k.Client.AddUserToGroup(k.Ctx, k.Token.AccessToken, Realm, userID, groupID)

}

func (k *KAuth) RemoveUserFromGroup(userID, groupID string) error {
	k.RefreshToken()

	return k.Client.DeleteUserFromGroup(k.Ctx, k.Token.AccessToken, Realm, userID, groupID)
}

func (k *KAuth) GetUsers(groupID string) ([]*gocloak.User, error) {
	k.RefreshToken()

	if groupID == "" {
		return k.Client.GetUsers(k.Ctx, k.Token.AccessToken, Realm, gocloak.GetUsersParams{})
	}

	return k.Client.GetGroupMembers(k.Ctx, k.Token.AccessToken, Realm, groupID, gocloak.GetGroupsParams{})

}

func (k *KAuth) GetGroups(userID string) ([]*gocloak.Group, error) {
	k.RefreshToken()

	if userID == "" {

		groups, err := k.Client.GetGroups(k.Ctx, k.Token.AccessToken, Realm, gocloak.GetGroupsParams{
			BriefRepresentation: gocloak.BoolP(false),
			Full:                gocloak.BoolP(true),
		})

		return groups, err
	}

	return k.Client.GetUserGroups(k.Ctx, k.Token.AccessToken, Realm, userID, gocloak.GetGroupsParams{})

}
