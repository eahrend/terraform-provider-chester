# api

Package API is a client side SDK to the chester-api http server

## Types

### type [Client](/api/client.go#L17)

`type Client struct { ... }`

Client is the wrapper for all the things the
client API will need.

#### func [NewClient](/api/client.go#L29)

`func NewClient(host, user, pass, audience string) (*Client, error)`

NewClient creates a pointer to a Client struct with specific
configuration options.
On a failure it will return a nil object and a non-nil error.

#### func [NewClientWithOptions](/api/client.go#L64)

`func NewClientWithOptions(opts ...ClientOption) (*Client, error)`

NewClientWithOptions allows for users to create a client with their own options.

#### func (*Client) [AddDatabase](/api/database.go#L96)

`func (c *Client) AddDatabase(database models.AddDatabaseRequest) (models.AddDatabaseResponse, error)`

AddDatabase will add an instance group to chester.
On a successful response it will return a filled models.AddDatabaseResponse and a nil error.
On a failure it will return an empty models.AddDatabaseResponse and a non-nil error.

```go
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
```

#### func (*Client) [CreateUser](/api/database.go#L195)

`func (c *Client) CreateUser(user models.ProxySqlMySqlUser) error`

CreateUser isn't currently used, this will be added at a later date

#### func (*Client) [DeleteUser](/api/database.go#L189)

`func (c *Client) DeleteUser(username string) error`

DeleteUser shouldn't be used, deleting an instance group should delete all associated users

#### func (*Client) [GetDatabase](/api/database.go#L54)

`func (c *Client) GetDatabase(instanceName string) (models.InstanceData, error)`

GetDatabase returns a model.InstanceData based on the instanceName parameter.
On a successful call, it will return a models.InstanceData struct and a nil error.
On an unsuccessful call, it will return an empty models.InstanceData struct and a non-nil error.

```go
package main
import github.com/eahrend/terraform-provider-chester/api

db, err := client.GetDatabase("sql-instance")
if err != nil {
	// handle error here
}
fmt.Println(db.InstanceName)
```

#### func (*Client) [GetDatabases](/api/database.go#L27)

`func (c *Client) GetDatabases() ([]models.InstanceData, error)`

GetDatabases returns a list of models.InstanceData of all the instance
data. On a successful call it will return the list and a nil error.
On an unsuccessful call, it will return an empty list and a non-nil error.

```go
package main
import github.com/eahrend/terraform-provider-chester/api

dbs, err := client.GetDatabases()
if err != nil {
	// handle errors here
}
for _, db := range dbs {
	fmt.Println(db.InstanceName)
}
```

#### func (*Client) [GetUser](/api/database.go#L205)

`func (c *Client) GetUser(userName string) (models.ProxySqlMySqlUser, error)`

GetUser isn't currently used, these tweaks will have to be made at a later date

#### func (*Client) [ModifyDatabase](/api/database.go#L158)

`func (c *Client) ModifyDatabase(database models.ModifyDatabaseRequest) error`

ModifyDatabase sends a modify command to the http api.
On a successful call this returns a nil error.
On a failed call, this will return a non-nil error.

```go
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
```

#### func (*Client) [ModifyQueryRuleByID](/api/database.go#L168)

`func (c *Client) ModifyQueryRuleByID(queryRule models.ProxySqlMySqlQueryRule) error`

ModifyQueryRuleByID shouldn't be used. Query Rules need to be authoritive.

#### func (*Client) [ModifyUser](/api/database.go#L179)

`func (c *Client) ModifyUser(userData models.ModifyUserRequest) error`

ModifyUser isn't currently used, will have to add at a later date when it becomes necessary

#### func (*Client) [RemoveDatabase](/api/database.go#L130)

`func (c *Client) RemoveDatabase(database models.RemoveDatabaseRequest) error`

RemoveDatabase removes a database based on the models.RemoveDatabaseRequest object.
On success, it will return a nil error.
On a failure, it will return a non-nil error.

```go
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
```

#### func (*Client) [UpdateCert](/api/database.go#L232)

`func (c *Client) UpdateCert(certData, instanceGroup string) error`

TODO: this is irrelevant for the most part until I update proxysql

#### func (*Client) [UpdateKey](/api/database.go#L219)

`func (c *Client) UpdateKey(keyData, instanceGroup string) error`

TODO: this is irrelevant for the most part until I update proxysql

### type [ClientOption](/api/client.go#L13)

`type ClientOption func(*Client)`

ClientOption is an option wrapper for the client

#### func [WithAudience](/api/client.go#L123)

`func WithAudience(audience string) ClientOption`

WithAudience sets the client's audience and generates
a token using google's idtoken package.

!!! DO NOT USE IN CONJUNCTION WITH WithToken !!!

#### func [WithHTTPClient](/api/client.go#L131)

`func WithHTTPClient(client *http.Client) ClientOption`

WithHTTPClient creates a ClientOption that overrides the
default http client.

#### func [WithHost](/api/client.go#L87)

`func WithHost(host string) ClientOption`

WithHost creates a ClientOption that modifies the client's
base host url

#### func [WithPassword](/api/client.go#L103)

`func WithPassword(password string) ClientOption`

WithPassword creates a ClientOption that modifies
the client's basic auth password

#### func [WithToken](/api/client.go#L113)

`func WithToken(token *oauth2.Token) ClientOption`

WithToken creates a ClientOption that sets the
the client's oauth2 token.

!!! DO NOT USE IN CONJUNCTION WITH WithAudience!!!

#### func [WithUsername](/api/client.go#L95)

`func WithUsername(username string) ClientOption`

WithUsername creates a ClientOption that modifies the
client's basic auth username

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
