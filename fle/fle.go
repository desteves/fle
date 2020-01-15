package fle

import (
	"context"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

var (
	kmsProviders map[string]map[string]interface{}
	masterKey    map[string]interface{}
	keydb        string
	keycoll      string
	namespace    string
)

func init() {

	// vars set up
	keydb = os.Getenv("MONGODB_KEY_VAULT_DB")
	keycoll = os.Getenv("MONGODB_KEY_VAULT_COLL")
	namespace = keydb + "." + keycoll
	kmsProviders = map[string]map[string]interface{}{
		"aws": {
			"accessKeyId":     os.Getenv("AWS_KEY"),
			"secretAccessKey": os.Getenv("AWS_SECRET_ACCESS_KEY"),
		},
	}
	masterKey = map[string]interface{}{
		"region": os.Getenv("AWS_KMS_REGION"),
		"key":    os.Getenv("AWS_KMS_ARN_ROLE"),
	}
	// clean up --  important to do this with { w: maj }
	client, err := CreateMongoClient(os.Getenv("MONGODB_ATLAS_URI"), writeconcern.New(writeconcern.WMajority()), false)
	if err != nil {
		log.Fatalf("! client - %v", err)
	}
	defer client.Disconnect(context.TODO())
	_ = client.Database(keydb).Collection(keycoll).Drop(context.TODO())

	// config keys
	clientOptions := options.ClientEncryption().SetKeyVaultNamespace(namespace).SetKmsProviders(kmsProviders)
	// verify
	keyVaultClient, err := mongo.NewClientEncryption(client, clientOptions)
	if err != nil {
		log.Fatalf("! key vault client - %v", err)
	}
	defer keyVaultClient.Close(context.TODO())
	dataKeyOptions := options.DataKey().SetMasterKey(masterKey).SetKeyAltNames([]string{"awskms"})
	result, err := keyVaultClient.CreateDataKey(context.TODO(), "aws", dataKeyOptions)
	if err != nil {
		log.Fatalf("! create data key - %v", err)
	}

	// Insert the data specified into the admin.datakeys
	log.Infof(" key valut create data key result is %+v", result)

}

// CreateMongoClient is
func CreateMongoClient(uri string, wc *writeconcern.WriteConcern, useEncryption bool) (client *mongo.Client, err error) {

	clientOptions := options.Client().ApplyURI(uri).SetWriteConcern(wc)

	/////////////////// adds encryption to the mongo client
	if useEncryption {
		schemaMap := map[string]interface{}{
			"tutorial.foobar": readJSONFile("foobarSchemaMap.json"),
		}
		autoEncOptions := options.AutoEncryption().
			SetKeyVaultNamespace(namespace).
			SetKmsProviders(kmsProviders).
			SetSchemaMap(schemaMap)

		clientOptions.SetAutoEncryptionOptions(autoEncOptions)
	}
	/////////////////////////////////////////////////////////

	client, err = mongo.Connect(context.TODO(), clientOptions)
	return
}

// helper private function
func readJSONFile(file string) bson.D {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("ReadFile error for %v: %v", file, err)
	}

	var fileDoc bson.D
	if err = bson.UnmarshalExtJSON(content, false, &fileDoc); err != nil {
		log.Fatalf("UnmarshalExtJSON error for file %v: %v", file, err)
	}
	return fileDoc
}
