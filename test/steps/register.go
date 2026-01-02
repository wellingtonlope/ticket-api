package steps

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

type RegisterFeatureSuite struct{}

func (r *RegisterFeatureSuite) InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
	})

	ctx.AfterSuite(func() {
	})
}

func (r *RegisterFeatureSuite) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		testCtx := NewBDDTestContext(&testing.T{})
		return SetTestContext(ctx, testCtx), nil
	})

	ctx.Step(`^I have valid user data:$`, r.IHaveValidUserData)
	ctx.Step(`^I have invalid user data:$`, r.IHaveInvalidUserData)
	ctx.Step(`^the user "([^"]*)" is already registered$`, r.TheUserIsAlreadyRegistered)
	ctx.Step(`^I send a POST request to "([^"]*)"$`, r.ISendAPostRequestTo)
	ctx.Step(`^I send a POST request to "([^"]*)" with invalid JSON body$`, r.ISendAPostRequestToWithInvalidJSONBody)
	ctx.Step(`^I should receive a (\d+) status code$`, r.IShouldReceiveAStatusCode)
	ctx.Step(`^the response should contain a "([^"]*)" field$`, r.TheResponseShouldContainAField)
	ctx.Step(`^the response should contain error message "([^"]*)"$`, r.TheResponseShouldContainErrorMessage)
	ctx.Step(`^the response should not contain an error message$`, r.TheResponseShouldNotContainAnErrorMessage)
	ctx.Step(`^the response should contain a "token" field$`, r.TheResponseShouldContainATokenField)
}

func RegisterFeature(t *testing.T, featurePath string) {
	absPath, err := filepath.Abs(featurePath)
	if err != nil {
		t.Fatalf("failed to get absolute path: %v", err)
	}

	opts := godog.Options{
		Output: colors.Colored(os.Stdout),
		Format: "pretty",
		Paths:  []string{absPath},
	}

	feature := &RegisterFeatureSuite{}

	status := godog.TestSuite{
		Name:                 "register",
		TestSuiteInitializer: feature.InitializeTestSuite,
		ScenarioInitializer:  feature.InitializeScenario,
		Options:              &opts,
	}.Run()

	if status != 0 {
		t.Fatalf("non-zero status returned, failed to run feature tests")
	}
}

func (r *RegisterFeatureSuite) IHaveValidUserData(ctx context.Context, table *godog.Table) (context.Context, error) {
	testCtx := GetTestContext(ctx)
	if testCtx == nil {
		return ctx, fmt.Errorf("test context not found")
	}

	testCtx.ClearRequestBody()

	for _, row := range table.Rows {
		key := row.Cells[0].Value
		value := row.Cells[1].Value

		switch key {
		case "name", "email", "password":
			testCtx.SetRequestBody(key, value)
		}
	}

	return ctx, nil
}

func (r *RegisterFeatureSuite) IHaveInvalidUserData(ctx context.Context, table *godog.Table) (context.Context, error) {
	testCtx := GetTestContext(ctx)
	if testCtx == nil {
		return ctx, fmt.Errorf("test context not found")
	}

	testCtx.ClearRequestBody()

	for _, row := range table.Rows {
		key := row.Cells[0].Value
		value := row.Cells[1].Value

		switch key {
		case "name", "email", "password":
			testCtx.SetRequestBody(key, value)
		}
	}

	return ctx, nil
}

func (r *RegisterFeatureSuite) TheUserIsAlreadyRegistered(ctx context.Context, email string) (context.Context, error) {
	testCtx := GetTestContext(ctx)
	if testCtx == nil {
		return ctx, fmt.Errorf("test context not found")
	}

	err := testCtx.CreateTestUser(email, "Existing User", "password123")
	if err != nil {
		return ctx, fmt.Errorf("failed to create test user: %w", err)
	}

	return ctx, nil
}

func (r *RegisterFeatureSuite) ISendAPostRequestTo(ctx context.Context, path string) (context.Context, error) {
	testCtx := GetTestContext(ctx)
	if testCtx == nil {
		return ctx, fmt.Errorf("test context not found")
	}

	err := testCtx.SendPostRequest(path)
	if err != nil {
		return ctx, fmt.Errorf("failed to send POST request: %w", err)
	}

	return ctx, nil
}

func (r *RegisterFeatureSuite) ISendAPostRequestToWithInvalidJSONBody(ctx context.Context, path string) (context.Context, error) {
	testCtx := GetTestContext(ctx)
	if testCtx == nil {
		return ctx, fmt.Errorf("test context not found")
	}

	testCtx.ClearRequestBody()
	testCtx.SetRawBody(`{"name":test"email":"test@test.com","password":"password123"}`)

	err := testCtx.SendPostRequest(path)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (r *RegisterFeatureSuite) IShouldReceiveAStatusCode(ctx context.Context, expectedCode int) error {
	testCtx := GetTestContext(ctx)
	if testCtx == nil {
		return fmt.Errorf("test context not found")
	}

	actualCode := testCtx.GetStatusCode()
	if actualCode != expectedCode {
		return fmt.Errorf("expected status code %d, got %d", expectedCode, actualCode)
	}

	return nil
}

func (r *RegisterFeatureSuite) TheResponseShouldContainAField(ctx context.Context, field string) error {
	testCtx := GetTestContext(ctx)
	if testCtx == nil {
		return fmt.Errorf("test context not found")
	}

	body := testCtx.GetResponseBody()
	if body == nil {
		return fmt.Errorf("response body is empty")
	}

	if _, exists := body[field]; !exists {
		return fmt.Errorf("field '%s' not found in response", field)
	}

	return nil
}

func (r *RegisterFeatureSuite) TheResponseShouldContainErrorMessage(ctx context.Context, expectedMessage string) error {
	testCtx := GetTestContext(ctx)
	if testCtx == nil {
		return fmt.Errorf("test context not found")
	}

	body := testCtx.GetResponseBody()
	if body == nil {
		return fmt.Errorf("response body is empty")
	}

	message, ok := body["message"].(string)
	if !ok {
		return fmt.Errorf("'message' field not found in response or not a string")
	}

	if message != expectedMessage {
		return fmt.Errorf("expected message '%s', got '%s'", expectedMessage, message)
	}

	return nil
}

func (r *RegisterFeatureSuite) TheResponseShouldNotContainAnErrorMessage(ctx context.Context) error {
	testCtx := GetTestContext(ctx)
	if testCtx == nil {
		return fmt.Errorf("test context not found")
	}

	body := testCtx.GetResponseBody()
	if body == nil {
		return fmt.Errorf("response body is empty")
	}

	if _, exists := body["message"]; exists {
		return fmt.Errorf("unexpected error message in response")
	}

	return nil
}

func (r *RegisterFeatureSuite) TheResponseShouldContainATokenField(ctx context.Context) error {
	return r.TheResponseShouldContainAField(ctx, "token")
}
