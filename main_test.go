package main

import (
	"2022-dev-career-boost-handon/ent"
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func Test_Organizationとそれに紐づくUserの追加検索削除ができる(t *testing.T) {
	cli := prepareClient()
	defer cli.Close()

	ctx := context.Background()

	org, err := CreateOrganization(ctx, cli, "PLAID")
	assert.Nil(t, err)
	assert.Equal(t, org.Name, "PLAID")

	user, err := CreateUser(ctx, cli, "YukiTsuchida", 26, true, org)
	assert.Nil(t, err)
	assert.Equal(t, user.Name, "YukiTsuchida")
	assert.Equal(t, user.Age, 26)
	assert.Equal(t, user.Active, true)

	users, err := org.QueryUsers().All(ctx)
	assert.Nil(t, err)
	assert.Equal(t, len(users), 1)
	assert.Equal(t, users[0].Name, "YukiTsuchida")

	orgs, err := user.QueryOrganizations().All(ctx)
	assert.Nil(t, err)
	assert.Equal(t, len(orgs), 1)
	assert.Equal(t, orgs[0].Name, "PLAID")

	err = cli.User.DeleteOne(user).Exec(ctx)
	assert.Nil(t, err)

	err = cli.Organization.DeleteOne(org).Exec(ctx)
	assert.Nil(t, err)
}

func Test_Userのageが負の場合にエラーが起きる(t *testing.T) {
	cli := prepareClient()
	defer cli.Close()

	ctx := context.Background()

	org, err := CreateOrganization(ctx, cli, "PLAID")
	assert.Nil(t, err)
	assert.Equal(t, org.Name, "PLAID")

	_, err = CreateUser(ctx, cli, "YukiTsuchida", -100, true, org)
	assert.EqualError(t, err, "failed creating user: ent: validator failed for field \"User.age\": value out of range")
}

func Test_UserがOrganizationに紐づけられていない場合にエラーが起きる(t *testing.T) {
	cli := prepareClient()
	defer cli.Close()

	ctx := context.Background()

	_, err := cli.User.
		Create().
		SetAge(26).
		SetName("YukiTsuchida").
		SetActive(true).
		Save(ctx)
	assert.EqualError(t, err, "ent: missing required edge \"User.organizations\"")
}

func CreateUser(ctx context.Context, client *ent.Client, name string, age int, active bool, organization *ent.Organization) (*ent.User, error) {
	user, err := client.User.
		Create().
		SetAge(age).
		SetName(name).
		SetActive(active).
		AddOrganizations(organization).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	return user, nil
}

func CreateOrganization(ctx context.Context, client *ent.Client, name string) (*ent.Organization, error) {
	organization, err := client.Organization.
		Create().
		SetName(name).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating organization: %w", err)
	}
	return organization, nil
}

func prepareClient() *ent.Client {
	client, err := ent.Open("postgres", "postgresql://admin:admin@localhost:5432/db?sslmode=disable")
	if err != nil {
		log.Fatalf("creating client: %v", err)
	}
	return client
}
