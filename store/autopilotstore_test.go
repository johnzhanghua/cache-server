package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

const (
	apiHost = "https://api2.autopilothq.com"
)

type TestContact struct {
	ID string `json:"contact_id,omitempty"`
}

func TestGetContact(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()

	gock.New(apiHost).
		MatchHeader("autopilotapikey", apiKey).
		MatchHeader("Content-Type", "application/json").
		Get("/v1/contact/person_9EAF39E4-9AEC-4134-964A-D9D8D54162E7").
		Reply(200).
		JSON(map[string]string{"contact_id": "person_9EAF39E4-9AEC-4134-964A-D9D8D54162E7"})

	store, err := NewAutoPilotStore()
	assert.NoError(err)

	key := "contact:person_9EAF39E4-9AEC-4134-964A-D9D8D54162E7"
	var contact TestContact
	resp := store.Get(context.Background(), key, &contact)
	assert.Equal(resp.StatusCode, 200)

	assert.Equal(contact.ID, "person_9EAF39E4-9AEC-4134-964A-D9D8D54162E7")
}

func TestUpsertContact(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()

	gock.New(apiHost).
		MatchHeader("autopilotapikey", apiKey).
		MatchHeader("Content-Type", "application/json").
		Post("/v1/contact").
		Reply(200).
		JSON(map[string]string{"contact_id": "person_9EAF39E4-9AEC-4134-964A-D9D8D54162E7"})

	store, err := NewAutoPilotStore()
	assert.NoError(err)

	key := "contact:person_9EAF39E4-9AEC-4134-964A-D9D8D54162E7"
	var resp TestContact
	data := TestContact{
		ID: "person_9EAF39E4-9AEC-4134-964A-D9D8D54162E7",
	}
	rs := store.Upsert(context.Background(), key, data, &resp)
	assert.Equal(rs.StatusCode, 200)

	assert.Equal(resp.ID, "person_9EAF39E4-9AEC-4134-964A-D9D8D54162E7")
}
