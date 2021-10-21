package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	models "github.com/eahrend/chestermodels"
)

/*
	GetDatabases returns a list of models.InstanceData of all the instance
	data. On a successful call it will return the list and a nil error.
	On an unsuccessful call, it will return an empty list and a non-nil error.

		package main
		import github.com/eahrend/terraform-provider-chester/api

		dbs, err := client.GetDatabases()
		if err != nil {
			// handle errors here
		}
		for _, db := range dbs {
			fmt.Println(db.InstanceName)
		}
*/
func (c *Client) GetDatabases() ([]models.InstanceData, error) {
	resp, err := c.makeRequest(nil, fmt.Sprintf("%s/databases?filter=true", c.HostURL), http.MethodGet)
	if err != nil {
		return nil, err
	}
	id := []models.InstanceData{}
	err = json.NewDecoder(bytes.NewBuffer(resp)).Decode(&id)
	if err != nil {
		return nil, err
	}
	return id, nil
}

/*
	GetDatabase returns a model.InstanceData based on the instanceName parameter.
	On a successful call, it will return a models.InstanceData struct and a nil error.
	On an unsuccessful call, it will return an empty models.InstanceData struct and a non-nil error.

		package main
		import github.com/eahrend/terraform-provider-chester/api

		db, err := client.GetDatabase("sql-instance")
		if err != nil {
			// handle error here
		}
		fmt.Println(db.InstanceName)
*/
func (c *Client) GetDatabase(instanceName string) (models.InstanceData, error) {
	resp, err := c.makeRequest(nil, fmt.Sprintf("%s/databases/%s?filter=true", c.HostURL, instanceName), http.MethodGet)
	if err != nil {
		return models.InstanceData{}, err
	}
	id := models.InstanceData{}
	err = json.NewDecoder(bytes.NewBuffer(resp)).Decode(&id)
	if err != nil {
		return models.InstanceData{}, err
	}
	return id, nil
}

/*
	AddDatabase will add an instance group to chester.
	On a successful response it will return a filled models.AddDatabaseResponse and a nil error.
	On a failure it will return an empty models.AddDatabaseResponse and a non-nil error.

		package main
		import github.com/eahrend/terraform-provider-chester/api

		addDataBaseRequest := models.AddDatabaseRequest{
			Action: "add",
			InstanceName: "database-name",
			Username: "foo",
			Password: "bar",
			MasterInstance: models.AddDatabaseRequestDatabaseInformation{
				IPAddress: "1.2.3.4",
				Name: "foo",
			},
			ReadReplicas: []models.AddDatabaseRequestDatabaseInformation{},
			ChesterMetaData: models.ChesterMetaData{
				InstanceGroup:"foo",
				MaxReadReplicas:5,
			},
		}
		addDataBaseResponse, err := client.AddDatabase(addDataBaseRequest)
		if err != nil {
			// handle error here
		}

*/
func (c *Client) AddDatabase(database models.AddDatabaseRequest) (models.AddDatabaseResponse, error) {
	b, err := json.Marshal(&database)
	if err != nil {
		return models.AddDatabaseResponse{}, err
	}
	resp, err := c.makeRequest(b, fmt.Sprintf("%s", c.HostURL), http.MethodPost)
	if err != nil {
		return models.AddDatabaseResponse{}, err
	}
	adr := models.AddDatabaseResponse{}
	err = json.NewDecoder(bytes.NewBuffer(resp)).Decode(&adr)
	if err != nil {
		return models.AddDatabaseResponse{}, err
	}
	return adr, nil
}

/*
	RemoveDatabase removes a database based on the models.RemoveDatabaseRequest object.
	On success, it will return a nil error.
	On a failure, it will return a non-nil error.

		package main
		import github.com/eahrend/terraform-provider-chester/api

		removeDatabaseRequest := models.RemoveDatabaseRequest{
			Action:"remove",
			InstanceName: "foo",
		}
		err := client.RemoveDatabase(removeDatabaseRequest)
		if err != nil {
			// handle error here
		}
*/
func (c *Client) RemoveDatabase(database models.RemoveDatabaseRequest) error {
	b, err := json.Marshal(&database)
	if err != nil {
		return err
	}
	_, err = c.makeRequest(b, fmt.Sprintf("%s", c.HostURL), http.MethodDelete)
	return err
}

/*
	ModifyDatabase sends a modify command to the http api.
	On a successful call this returns a nil error.
	On a failed call, this will return a non-nil error.


		package main
		import github.com/eahrend/terraform-provider-chester/api

		modifyDatabaseRequest := models.ModifyDatabaseRequest{
			Action:"add",
			InstanceName: "foo",
			ChesterMetaData: models.ChesterMetaData{
				InstanceGroup:"foo",
				MaxReadReplicas:5,
			},
		}
		err := client.ModifyDatabase(modifyDatabaseRequest)
*/
func (c *Client) ModifyDatabase(database models.ModifyDatabaseRequest) error {
	b, err := json.Marshal(&database)
	if err != nil {
		return err
	}
	_, err = c.makeRequest(b, fmt.Sprintf("%s", c.HostURL), http.MethodPatch)
	return err
}

// ModifyQueryRuleByID shouldn't be used. Query Rules need to be authoritive.
func (c *Client) ModifyQueryRuleByID(queryRule models.ProxySqlMySqlQueryRule) error {
	queryRuleID := queryRule.RuleID
	b, err := json.Marshal(&queryRule)
	if err != nil {
		return err
	}
	_, err = c.makeRequest(b, fmt.Sprintf("%s/queryrules/%v", c.HostURL, queryRuleID), http.MethodPatch)
	return err
}

// ModifyUser isn't currently used, will have to add at a later date when it becomes necessary
func (c *Client) ModifyUser(userData models.ModifyUserRequest) error {
	b, err := json.Marshal(&userData)
	if err != nil {
		return err
	}
	_, err = c.makeRequest(b, fmt.Sprintf("%s/users", c.HostURL), http.MethodPatch)
	return err
}

// DeleteUser shouldn't be used, deleting an instance group should delete all associated users
func (c *Client) DeleteUser(username string) error {
	_, err := c.makeRequest(nil, fmt.Sprintf("%s/users/%s", c.HostURL, username), http.MethodDelete)
	return err
}

// CreateUser isn't currently used, this will be added at a later date
func (c *Client) CreateUser(user models.ProxySqlMySqlUser) error {
	b, err := json.Marshal(&user)
	if err != nil {
		return err
	}
	_, err = c.makeRequest(b, fmt.Sprintf("%s/users", c.HostURL), http.MethodPost)
	return err
}

// GetUser isn't currently used, these tweaks will have to be made at a later date
func (c *Client) GetUser(userName string) (models.ProxySqlMySqlUser, error) {
	b, err := c.makeRequest(nil, fmt.Sprintf("%s/users/%s", c.HostURL, userName), http.MethodGet)
	if err != nil {
		return models.ProxySqlMySqlUser{}, err
	}
	user := models.ProxySqlMySqlUser{}
	err = json.NewDecoder(bytes.NewReader(b)).Decode(&user)
	if err != nil {
		return models.ProxySqlMySqlUser{}, err
	}
	return user, nil
}

// TODO: this is irrelevant for the most part until I update proxysql
func (c *Client) UpdateKey(keyData, instanceGroup string) error {
	keyMap := map[string]string{
		"key": keyData,
	}
	b, err := json.Marshal(&keyMap)
	if err != nil {
		return err
	}
	_, err = c.makeRequest(b, fmt.Sprintf("%s/key/%s", c.HostURL, instanceGroup), http.MethodPatch)
	return err
}

// TODO: this is irrelevant for the most part until I update proxysql
func (c *Client) UpdateCert(certData, instanceGroup string) error {
	certMap := map[string]string{
		"cert": certData,
	}
	b, err := json.Marshal(&certMap)
	if err != nil {
		return err
	}
	_, err = c.makeRequest(b, fmt.Sprintf("%s/cert/%s", c.HostURL, instanceGroup), http.MethodPatch)
	return err
}
