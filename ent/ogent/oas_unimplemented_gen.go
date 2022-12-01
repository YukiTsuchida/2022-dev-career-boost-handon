// Code generated by ogen, DO NOT EDIT.

package ogent

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// AddUser implements addUser operation.
//
// PATCH /organizations/{id}/addUser
func (UnimplementedHandler) AddUser(ctx context.Context, params AddUserParams) (r AddUserOK, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateOrganization implements createOrganization operation.
//
// Creates a new Organization and persists it to storage.
//
// POST /organizations
func (UnimplementedHandler) CreateOrganization(ctx context.Context, req CreateOrganizationReq) (r CreateOrganizationRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateUser implements createUser operation.
//
// Creates a new User and persists it to storage.
//
// POST /users
func (UnimplementedHandler) CreateUser(ctx context.Context, req CreateUserReq) (r CreateUserRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteOrganization implements deleteOrganization operation.
//
// Deletes the Organization with the requested ID.
//
// DELETE /organizations/{id}
func (UnimplementedHandler) DeleteOrganization(ctx context.Context, params DeleteOrganizationParams) (r DeleteOrganizationRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteUser implements deleteUser operation.
//
// Deletes the User with the requested ID.
//
// DELETE /users/{id}
func (UnimplementedHandler) DeleteUser(ctx context.Context, params DeleteUserParams) (r DeleteUserRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListOrganization implements listOrganization operation.
//
// List Organizations.
//
// GET /organizations
func (UnimplementedHandler) ListOrganization(ctx context.Context, params ListOrganizationParams) (r ListOrganizationRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListOrganizationUsers implements listOrganizationUsers operation.
//
// List attached Users.
//
// GET /organizations/{id}/users
func (UnimplementedHandler) ListOrganizationUsers(ctx context.Context, params ListOrganizationUsersParams) (r ListOrganizationUsersRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListUser implements listUser operation.
//
// List Users.
//
// GET /users
func (UnimplementedHandler) ListUser(ctx context.Context, params ListUserParams) (r ListUserRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListUserOrganizations implements listUserOrganizations operation.
//
// List attached Organizations.
//
// GET /users/{id}/organizations
func (UnimplementedHandler) ListUserOrganizations(ctx context.Context, params ListUserOrganizationsParams) (r ListUserOrganizationsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadOrganization implements readOrganization operation.
//
// Finds the Organization with the requested ID and returns it.
//
// GET /organizations/{id}
func (UnimplementedHandler) ReadOrganization(ctx context.Context, params ReadOrganizationParams) (r ReadOrganizationRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadUser implements readUser operation.
//
// Finds the User with the requested ID and returns it.
//
// GET /users/{id}
func (UnimplementedHandler) ReadUser(ctx context.Context, params ReadUserParams) (r ReadUserRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateOrganization implements updateOrganization operation.
//
// Updates a Organization and persists changes to storage.
//
// PATCH /organizations/{id}
func (UnimplementedHandler) UpdateOrganization(ctx context.Context, req UpdateOrganizationReq, params UpdateOrganizationParams) (r UpdateOrganizationRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateUser implements updateUser operation.
//
// Updates a User and persists changes to storage.
//
// PATCH /users/{id}
func (UnimplementedHandler) UpdateUser(ctx context.Context, req UpdateUserReq, params UpdateUserParams) (r UpdateUserRes, _ error) {
	return r, ht.ErrNotImplemented
}
