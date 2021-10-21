package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/eahrend/chestermodels"
)

var (
	mux       *http.ServeMux
	server    *httptest.Server
	client    *Client
	username  string
	password  string
	databases []models.InstanceData
)

func testAuthenciateMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		basicuser, basicpass, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		} else if basicuser != username || basicpass != password {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func emptyResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	dbList := []models.InstanceData{}
	b, _ := json.Marshal(&dbList)
	w.Write(b)
}

func setupBadHost() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	username = "foo"
	password = "bar"
	// skipping token for IAP auth here.
	client, _ = NewClientWithOptions(WithHost("http://notarealplace.co.uk"), WithPassword(password), WithUsername(username))
	return func() {
		server.Close()
	}
}

func setupBadPassword() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	username = "foo"
	password = "bar"
	// skipping token for IAP auth here.
	client, _ = NewClientWithOptions(WithHost(server.URL), WithPassword("not_a_real_password"), WithUsername(username))
	return func() {
		server.Close()
	}
}

func setupBadUser() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	password = "bar"
	// skipping token for IAP auth here.
	client, _ = NewClientWithOptions(WithHost(server.URL), WithPassword(password), WithUsername("not_a_real_user"))
	return func() {
		server.Close()
	}
}

func setup() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	username = "foo"
	password = "bar"
	// reverting to empty list here
	databases = []models.InstanceData{}
	dbOneReadReplicas := []models.AddDatabaseRequestDatabaseInformation{}
	dbOneReadOne := models.AddDatabaseRequestDatabaseInformation{
		Name:      "foo-reader-one",
		IPAddress: "1.2.3.5",
	}
	dbOneReadTwo := models.AddDatabaseRequestDatabaseInformation{
		Name:      "foo-reader-two",
		IPAddress: "1.2.3.6",
	}
	dbOneReadReplicas = append(dbOneReadReplicas, dbOneReadOne, dbOneReadTwo)
	dbOneQueryRules := []models.ProxySqlMySqlQueryRule{}
	dbOneQueryRuleOne := models.ProxySqlMySqlQueryRule{
		RuleID:               1,
		Username:             "foo",
		Active:               1,
		MatchDigest:          "foo",
		DestinationHostgroup: 5,
		Apply:                1,
		Comment:              "bar",
	}
	dbOneQueryRuleTwo := models.ProxySqlMySqlQueryRule{
		RuleID:               2,
		Username:             "foo",
		Active:               10,
		MatchDigest:          "foo",
		DestinationHostgroup: 5,
		Apply:                1,
		Comment:              "bar",
	}
	dbOneQueryRules = append(dbOneQueryRules, dbOneQueryRuleOne, dbOneQueryRuleTwo)
	dbOne := models.InstanceData{
		InstanceName:   "foo",
		ReadHostGroup:  5,
		WriteHostGroup: 10,
		Username:       "foo",
		Password:       "bar",
		QueryRules:     dbOneQueryRules,
		MasterInstance: models.AddDatabaseRequestDatabaseInformation{
			Name:      "foo",
			IPAddress: "1.2.3.4",
		},
		ReadReplicas: dbOneReadReplicas,
		UseSSL:       0,
		ChesterMetaData: models.ChesterMetaData{
			InstanceGroup:       "foo",
			MaxChesterInstances: 2,
		},
	}
	dbTwoReadReplicas := []models.AddDatabaseRequestDatabaseInformation{}
	dbTwoReadOne := models.AddDatabaseRequestDatabaseInformation{
		Name:      "shaz-reader-one",
		IPAddress: "1.2.3.7",
	}
	dbTwoReadTwo := models.AddDatabaseRequestDatabaseInformation{
		Name:      "shaz-reader-two",
		IPAddress: "1.2.3.8",
	}
	dbTwoReadReplicas = append(dbTwoReadReplicas, dbTwoReadOne, dbTwoReadTwo)
	dbTwoQueryRules := []models.ProxySqlMySqlQueryRule{}
	dbTwoQueryRuleOne := models.ProxySqlMySqlQueryRule{
		RuleID:               1,
		Username:             "shaz",
		Active:               1,
		MatchDigest:          "shaz",
		DestinationHostgroup: 5,
		Apply:                1,
		Comment:              "bot",
	}
	dbTwoQueryRuleTwo := models.ProxySqlMySqlQueryRule{
		RuleID:               2,
		Username:             "shaz",
		Active:               1,
		MatchDigest:          "shaz",
		DestinationHostgroup: 10,
		Apply:                1,
		Comment:              "bot",
	}
	dbTwoQueryRules = append(dbTwoQueryRules, dbTwoQueryRuleOne, dbTwoQueryRuleTwo)
	dbTwo := models.InstanceData{
		InstanceName:   "shaz",
		ReadHostGroup:  5,
		WriteHostGroup: 10,
		Username:       "shaz",
		Password:       "bot",
		QueryRules:     dbTwoQueryRules,
		MasterInstance: models.AddDatabaseRequestDatabaseInformation{
			Name:      "shaz",
			IPAddress: "1.2.3.8",
		},
		ReadReplicas: dbTwoReadReplicas,
		UseSSL:       0,
		ChesterMetaData: models.ChesterMetaData{
			InstanceGroup:       "shaz",
			MaxChesterInstances: 3,
		},
	}
	databases = append(databases, dbOne, dbTwo)
	// skipping token for IAP auth here.
	client, _ = NewClientWithOptions(WithHost(server.URL), WithPassword(password), WithUsername(username))
	return func() {
		server.Close()
	}
}

// TestNewClient_AuthPass is the basic functionality of our basic auth setup
func TestNewClient_AuthPass(t *testing.T) {
	teardown := setup()
	defer teardown()
	emptyHandler := http.HandlerFunc(emptyResponse)
	mux.Handle("/databases", testAuthenciateMiddleWare(emptyHandler))
	_, err := client.GetDatabases()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// TestNewClient_AuthFailBadUser checks for a bad username
// if no error is returned, this fails.
func TestNewClient_AuthFailBadUser(t *testing.T) {
	teardown := setupBadUser()
	defer teardown()
	emptyHandler := http.HandlerFunc(emptyResponse)
	mux.Handle("/databases", testAuthenciateMiddleWare(emptyHandler))
	_, err := client.GetDatabases()
	if err == nil {
		t.FailNow()
	}
}

// TestNewClient_AuthFailBadPassword checks for a bad password
// if no error is returned, this fails.
func TestNewClient_AuthFailBadPassword(t *testing.T) {
	teardown := setupBadPassword()
	defer teardown()
	emptyHandler := http.HandlerFunc(emptyResponse)
	mux.Handle("/databases", testAuthenciateMiddleWare(emptyHandler))
	_, err := client.GetDatabases()
	if err == nil {
		t.FailNow()
	}
}

// TestNewClient_AuthFailBadHost checks for a bad password
// if no error is returned, this fails.
func TestNewClient_AuthFailBadHost(t *testing.T) {
	teardown := setupBadHost()
	defer teardown()
	emptyHandler := http.HandlerFunc(emptyResponse)
	mux.Handle("/databases", testAuthenciateMiddleWare(emptyHandler))
	_, err := client.GetDatabases()
	if err == nil {
		t.FailNow()
	}
}

func getDatabasesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(&databases)
	if err != nil {
		http.Error(w, "failed to marshal json", http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

// TestClient_GetDatabases checks that we're able to query all
// affected databases.
func TestClient_GetDatabases(t *testing.T) {
	teardown := setup()
	defer teardown()
	mux.HandleFunc("/databases", http.HandlerFunc(getDatabasesHandler))
	dbs, err := client.GetDatabases()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	// doing some checking that we're good on our end
	for _, db := range dbs {
		if db.MasterInstance.Name != db.InstanceName {
			t.FailNow()
		}
	}
}

func TestClient_GetDatabaseByName(t *testing.T) {
	teardown := setup()
	defer teardown()
	mux.HandleFunc("/databases/foo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		for _, db := range databases {
			if db.InstanceName == "foo" {
				b, err := json.Marshal(db)
				if err != nil {
					t.Fatalf("failed to marshal response from mock server %s", err.Error())
				}
				w.Write(b)
			}
		}
	})
	db, err := client.GetDatabase("foo")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if db.InstanceName != "foo" {
		t.FailNow()
	}
}

// TestClient_GetDatabaseByNameFail tests that we're returning a 404
// on a databases that doesn't exist
func TestClient_GetDatabaseByNameFail(t *testing.T) {
	teardown := setup()
	defer teardown()
	mux.HandleFunc("/databases/foo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "not found", http.StatusNotFound)
	})
	_, err := client.GetDatabase("foo")
	if err == nil {
		t.FailNow()
	}
}

// TestClient_AddDatabase tests the functionality of adding a new instance
// group, as well as the ability to retreive that instance group
func TestClient_AddDatabase(t *testing.T) {
	teardown := setup()
	defer teardown()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		addDbRequest := models.AddDatabaseRequest{}
		err := json.NewDecoder(r.Body).Decode(&addDbRequest)
		if err != nil {
			http.Error(w, "failed to parse json", http.StatusBadRequest)
		}
		newDatabase := models.InstanceData{
			InstanceName:    addDbRequest.InstanceName,
			ReadHostGroup:   10,
			WriteHostGroup:  5,
			Username:        addDbRequest.Username,
			Password:        addDbRequest.Password,
			QueryRules:      addDbRequest.QueryRules,
			MasterInstance:  addDbRequest.MasterInstance,
			ReadReplicas:    addDbRequest.ReadReplicas,
			UseSSL:          0,
			ChesterMetaData: addDbRequest.ChesterMetaData,
		}
		databases = append(databases, newDatabase)
		addDataBaseResponse := models.AddDatabaseResponse{
			Action:          "add",
			QueryRules:      addDbRequest.QueryRules,
			InstanceName:    addDbRequest.InstanceName,
			Username:        addDbRequest.Username,
			Password:        addDbRequest.Password,
			WriteHostGroup:  5,
			ReadHostGroup:   10,
			SSLEnabled:      0,
			ChesterMetaData: addDbRequest.ChesterMetaData,
		}
		respBytes, err := json.Marshal(&addDataBaseResponse)
		if err != nil {
			http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		}
		w.WriteHeader(200)
		w.Write(respBytes)
	})
	mux.HandleFunc("/databases/temp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		dbCheck := false
		for _, db := range databases {
			if db.InstanceName == "temp" {
				dbCheck = true
				b, err := json.Marshal(db)
				if err != nil {
					t.Fatalf("failed to marshal response from mock server %s", err.Error())
				}
				w.Write(b)
			}
		}
		if !dbCheck {
			http.Error(w, "instance not found", http.StatusNotFound)
		}
	})
	mux.HandleFunc("/databases", http.HandlerFunc(getDatabasesHandler))

	addDbReadReplicaOne := models.AddDatabaseRequestDatabaseInformation{
		Name:      "foo-read-1",
		IPAddress: "2.3.4.6",
	}
	addDbReadReplicaTwo := models.AddDatabaseRequestDatabaseInformation{
		Name:      "foo-read-2",
		IPAddress: "2.3.4.7",
	}
	addDbReadReplicas := []models.AddDatabaseRequestDatabaseInformation{
		addDbReadReplicaOne, addDbReadReplicaTwo,
	}
	addDbQueryRuleOne := models.ProxySqlMySqlQueryRule{
		RuleID:               1,
		Username:             "foo",
		Active:               1,
		MatchDigest:          "bar",
		DestinationHostgroup: 5,
		Apply:                1,
		Comment:              "baz",
	}
	addDbQueryRuleTwo := models.ProxySqlMySqlQueryRule{
		RuleID:               2,
		Username:             "foo",
		Active:               1,
		MatchDigest:          "barzoople",
		DestinationHostgroup: 10,
		Apply:                1,
		Comment:              "baz",
	}
	addDbQueryRules := []models.ProxySqlMySqlQueryRule{
		addDbQueryRuleOne, addDbQueryRuleTwo,
	}
	addDbRequest := models.AddDatabaseRequest{
		Action:       "add",
		InstanceName: "temp",
		Username:     "foo",
		Password:     "bar",
		MasterInstance: models.AddDatabaseRequestDatabaseInformation{
			Name:      "foo",
			IPAddress: "2.3.4.5",
		},
		ReadReplicas: addDbReadReplicas,
		QueryRules:   addDbQueryRules,
		ChesterMetaData: models.ChesterMetaData{
			InstanceGroup:       "temp",
			MaxChesterInstances: 3,
		},
		KeyData:   "",
		CertData:  "",
		CAData:    "",
		EnableSSL: 0,
	}
	_, err := client.AddDatabase(addDbRequest)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	dbs, err := client.GetDatabases()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(dbs) != 3 {
		t.Errorf("wtf %v", len(dbs))
		t.FailNow()
	}
	db, err := client.GetDatabase("temp")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if db.InstanceName != "temp" {
		t.Errorf("wrong instance %s", db.InstanceName)
		t.FailNow()
	}
}

// TestClient_RemoveDatabase tests the functionality of removing an
// existing instance group.
func TestClient_RemoveDatabase(t *testing.T) {
	teardown := setup()
	defer teardown()
	mux.HandleFunc("/databases", http.HandlerFunc(getDatabasesHandler))
	mux.HandleFunc("/databases/foo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		dbCheck := false
		for _, db := range databases {
			if db.InstanceName == "foo" {
				dbCheck = true
				b, err := json.Marshal(&db)
				if err != nil {
					http.Error(w, "failed to marshal response", http.StatusInternalServerError)
					return
				}
				w.Write(b)
				return
			}
		}
		if !dbCheck {
			http.Error(w, "database not found", http.StatusNotFound)
			return
		}
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		removeDbRequest := models.RemoveDatabaseRequest{}
		err := json.NewDecoder(r.Body).Decode(&removeDbRequest)
		if err != nil {
			http.Error(w, "failed to parse json", http.StatusBadRequest)
			return
		}
		for index, db := range databases {
			if db.InstanceName == removeDbRequest.InstanceName {
				databases = append(databases[:index], databases[index+1:]...)
			}
		}
	})
	removeDataBaseReq := models.RemoveDatabaseRequest{
		Action:       "remove",
		InstanceName: "foo",
		Username:     "foo",
	}
	err := client.RemoveDatabase(removeDataBaseReq)
	if err != nil {
		t.Errorf("failed to remove database: %s", err.Error())
		t.FailNow()
	}
	dbs, err := client.GetDatabases()
	if err != nil {
		t.Errorf("failed to get databases: %s", err.Error())
		t.FailNow()
	}
	if len(dbs) != 1 {
		t.Errorf("too many databases")
		t.FailNow()
	}
	_, err = client.GetDatabase("foo")
	if err == nil {
		t.Errorf("should have failed")
		t.FailNow()
	}
}

func TestClient_ModifyDatabase(t *testing.T) {
	teardown := setup()
	defer teardown()
	mux.HandleFunc("/databases", http.HandlerFunc(getDatabasesHandler))
	mux.HandleFunc("/databases/foo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		dbCheck := false
		for _, db := range databases {
			if db.InstanceName == "foo" {
				dbCheck = true
				b, err := json.Marshal(&db)
				if err != nil {
					http.Error(w, "failed to marshal response", http.StatusInternalServerError)
					return
				}
				w.Write(b)
				return
			}
		}
		if !dbCheck {
			http.Error(w, "database not found", http.StatusNotFound)
			return
		}
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, "invalid method", http.StatusMethodNotAllowed)
			return
		}
		modifyRequest := models.ModifyDatabaseRequest{}
		err := json.NewDecoder(r.Body).Decode(&modifyRequest)
		if err != nil {
			http.Error(w, "failed to decode", http.StatusInternalServerError)
		}
		// get the DB
		database := models.InstanceData{}
		for _, db := range databases {
			if db.InstanceName == "foo" {
				database = db
				break
			}
		}
		if modifyRequest.ReadReplicas != nil {
			database.ReadReplicas = modifyRequest.ReadReplicas
			for index, db := range databases {
				if db.InstanceName == modifyRequest.InstanceName {
					// remove the database
					databases = append(databases[:index], databases[index+1:]...)
					// add the same one back
					databases = append(databases, database)
					break
				}
			}
		}
		if modifyRequest.RemoveQueryRules != nil {
			newQueryRules := []models.ProxySqlMySqlQueryRule{}
			for _, qr := range database.QueryRules {
				if qr.RuleID != modifyRequest.RemoveQueryRules[0] {
					newQueryRules = append(newQueryRules, qr)
				}
			}
			database.QueryRules = newQueryRules
			for index, db := range databases {
				if db.InstanceName == modifyRequest.InstanceName {
					// remove the database
					databases = append(databases[:index], databases[index+1:]...)
					// add the same one back
					databases = append(databases, database)
					break
				}
			}
		}
		if modifyRequest.AddQueryRules != nil {
			database.QueryRules = append(database.QueryRules, modifyRequest.AddQueryRules...)
			for index, db := range databases {
				if db.InstanceName == modifyRequest.InstanceName {
					// remove the database
					databases = append(databases[:index], databases[index+1:]...)
					// add the same one back
					databases = append(databases, database)
					break
				}
			}
		}
		if database.ChesterMetaData.MaxChesterInstances != modifyRequest.ChesterMetaData.MaxChesterInstances {
			database.ChesterMetaData = modifyRequest.ChesterMetaData
			// just replacing the database here since there isn't any "backend" to update.
			for index, db := range databases {
				if db.InstanceName == modifyRequest.InstanceName {
					// remove the database
					databases = append(databases[:index], databases[index+1:]...)
					// add the same one back
					databases = append(databases, database)
					break
				}
			}
		}

	})
	modifyDatabaseRequest := models.ModifyDatabaseRequest{
		Action:           "modify",
		InstanceName:     "foo",
		AddQueryRules:    nil,
		RemoveQueryRules: nil,
		ReadReplicas:     nil,
		ChesterMetaData: models.ChesterMetaData{
			InstanceGroup:       "foo",
			MaxChesterInstances: 20,
		},
	}
	err := client.ModifyDatabase(modifyDatabaseRequest)
	if err != nil {
		t.Errorf("failed to modify chester metadata")
		t.FailNow()
	}
	db, err := client.GetDatabase("foo")
	if err != nil {
		t.Errorf("failed to get metadata")
		t.FailNow()
	}
	if db.ChesterMetaData.MaxChesterInstances != 20 {
		t.Errorf("failed to modify chester instance")
		t.FailNow()
	}
	modifyDatabaseRequest.RemoveQueryRules = []int{1}
	err = client.ModifyDatabase(modifyDatabaseRequest)
	if err != nil {
		t.Errorf("failed to modify database: %s", err.Error())
		t.FailNow()
	}
	db, err = client.GetDatabase("foo")
	if err != nil {
		t.Errorf("failed to get metadata")
		t.FailNow()
	}
	if len(db.QueryRules) != 1 {
		t.Error("failed to remove query rule")
		t.FailNow()
	}
	newQueryRuleOne := models.ProxySqlMySqlQueryRule{
		RuleID:               3,
		Username:             "rick",
		Active:               1,
		MatchDigest:          "never",
		DestinationHostgroup: 5,
		Apply:                1,
		Comment:              "gonna",
	}
	newQueryRuleTwo := models.ProxySqlMySqlQueryRule{
		RuleID:               4,
		Username:             "astley",
		Active:               1,
		MatchDigest:          "give",
		DestinationHostgroup: 10,
		Apply:                1,
		Comment:              "you",
	}
	modifyDatabaseRequest.RemoveQueryRules = nil
	newQueryRules := []models.ProxySqlMySqlQueryRule{newQueryRuleOne, newQueryRuleTwo}
	modifyDatabaseRequest.AddQueryRules = newQueryRules
	err = client.ModifyDatabase(modifyDatabaseRequest)
	if err != nil {
		t.Errorf("failed to add query rules %s", err.Error())
		t.FailNow()
	}
	db, err = client.GetDatabase("foo")
	if err != nil {
		t.Errorf("failed to get database %s", err.Error())
		t.FailNow()
	}
	if len(db.QueryRules) != 3 {
		t.Error("failed to add more query rules")
		t.FailNow()
	}
	readReplica := models.AddDatabaseRequestDatabaseInformation{
		Name:      "reader-three",
		IPAddress: "2.3.4.6",
	}
	modifyDatabaseRequest.ReadReplicas = []models.AddDatabaseRequestDatabaseInformation{readReplica}
	err = client.ModifyDatabase(modifyDatabaseRequest)
	if err != nil {
		t.Errorf("failed to modify database: %s", err.Error())
		t.FailNow()
	}
	db, err = client.GetDatabase("foo")
	if err != nil {
		t.Errorf("failed to get database: %s", err.Error())
		t.FailNow()
	}
	if len(db.ReadReplicas) != 1 {
		t.Error("failed to reset read replicas")
		t.FailNow()
	}
}
