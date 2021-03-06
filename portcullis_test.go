package portcullis_test

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"

	"github.com/kubex/portcullis-go"
	"github.com/kubex/portcullis-go/keys"
)

const (
	testProject  = "this-is-a-test-project-id"
	testUserID   = "this-is-a-test-user-id"
	testUsername = "this-is-a-test-username"

	testAppID  = "test-app-id"
	testVendor = "test-vendor"
)

// TestAuthDataExtraction tests for valid transaction of portcullis meta data values
func TestAuthDataExtraction(t *testing.T) {
	metamap := map[string]string{}
	metamap[keys.GetProjectKey()] = testProject
	metamap[keys.GetUserIDKey()] = testUserID
	metamap[keys.GetUsernameKey()] = testUsername
	metamap[keys.GetAppIDKey()] = testAppID
	metamap[keys.GetAppVendorKey()] = testVendor

	meta := metadata.New(metamap)

	meta[keys.GetRolesKey()] = append(meta[keys.GetRolesKey()], "role2")
	meta[keys.GetRolesKey()] = append(meta[keys.GetRolesKey()], "role1")

	ctx := metadata.NewIncomingContext(context.Background(), meta)
	in := portcullis.FromContext(ctx)

	if in.GlobalAppID() != fmt.Sprintf("%s/%s", testVendor, testAppID) {
		t.Error("Global app ID does not contain expected value")
	}

	if !in.HasRole("role1") || !in.HasRole("role2") {
		t.Error("Roles do not contain expected values")
	}

	if in.ProjectID != testProject {
		t.Error("Project does not contain expected value")
	}

	if in.Username != testUsername {
		t.Error("Username does not contain expected value")
	}

	if in.UserID != testUserID {
		t.Error("userID does not contain expected value")
	}
}

// TestAuthDataExtractionWithMissingFields tests for valid extraction of portcullis meta with missing values
func TestAuthDataExtractionWithMissingFields(t *testing.T) {
	metamap := map[string]string{}
	metamap[keys.GetUsernameKey()] = testUsername
	meta := metadata.New(metamap)
	ctx := metadata.NewIncomingContext(context.Background(), meta)

	project := portcullis.FromContext(ctx).ProjectID
	username := portcullis.FromContext(ctx).Username
	userID := portcullis.FromContext(ctx).UserID
	roles := portcullis.FromContext(ctx).Roles

	if username != testUsername {
		t.Error("Username does not contain expected value")
	}

	if project != "" {
		t.Error("Project does not contain expected value")
	}

	if userID != "" {
		t.Error("userID does not contain expected value")
	}

	if len(roles) != 0 {
		t.Error("roles do not contain expected value")
	}
}

// TestExtractionWithInvalidContext tests extraction result with context contains no metadata
func TestExtractionWithInvalidContext(t *testing.T) {
	ctx := context.Background()
	project := portcullis.FromContext(ctx).ProjectID

	if project != "" {
		t.Error("Project does not contain expected value")
	}
}
