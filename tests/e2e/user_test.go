package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"guilhermefaleiros/go-service-template/internal/api"
	"guilhermefaleiros/go-service-template/internal/api/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type TestHelper struct {
	Server *httptest.Server
	Client *http.Client
	API    *api.API
}

func setup(t *testing.T) *TestHelper {
	t.Helper()

	api, err := api.NewAPI("e2e")
	assert.NoError(t, err)

	server := httptest.NewServer(api.Router)

	return &TestHelper{
		Server: server,
		Client: server.Client(),
		API:    api,
	}
}

func teardown(helper *TestHelper, t *testing.T) {
	t.Helper()

	helper.Server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := helper.API.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestUserControllerIntegration(t *testing.T) {
	helper := setup(t)
	defer teardown(helper, t)

	createUserEndpoint := fmt.Sprintf("%s/users", helper.Server.URL)

	t.Run("CreateUser_Success", func(t *testing.T) {
		requestBody := model.CreateUserRequest{
			Name:  faker.FirstName(),
			Email: faker.Email(),
			Phone: faker.Phonenumber(),
		}

		body, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		resp, err := helper.Client.Post(createUserEndpoint, "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()

		var response model.CreateUserResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		assert.NotEmpty(t, response.ID)
	})

	t.Run("CreateUser_FailWithSameEmail", func(t *testing.T) {
		requestBody := model.CreateUserRequest{
			Name:  faker.FirstName(),
			Email: faker.Email(),
			Phone: faker.Phonenumber(),
		}

		body, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		resp, err := helper.Client.Post(createUserEndpoint, "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()

		var response model.CreateUserResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		assert.NotEmpty(t, response.ID)

		resp, err = helper.Client.Post(createUserEndpoint, "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
