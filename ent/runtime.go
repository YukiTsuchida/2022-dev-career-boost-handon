// Code generated by ent, DO NOT EDIT.

package ent

import (
	"2022-dev-career-boost-handon/ent/organization"
	"2022-dev-career-boost-handon/ent/schema"
	"2022-dev-career-boost-handon/ent/user"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	organizationFields := schema.Organization{}.Fields()
	_ = organizationFields
	// organizationDescName is the schema descriptor for name field.
	organizationDescName := organizationFields[0].Descriptor()
	// organization.DefaultName holds the default value on creation for the name field.
	organization.DefaultName = organizationDescName.Default.(string)
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescName is the schema descriptor for name field.
	userDescName := userFields[0].Descriptor()
	// user.DefaultName holds the default value on creation for the name field.
	user.DefaultName = userDescName.Default.(string)
}
