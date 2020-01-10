package fle

import (
	"context"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	// clean up
	client, err := CreateMongoClient()
	if err != nil {
		log.Fatalf("! client - %v", err)
	}
	_ = client.Database(keydb).Collection(keycoll).Drop(context.TODO())
	defer client.Disconnect(context.TODO())

	// config keys
	clientOptions := options.ClientEncryption().SetKeyVaultNamespace(namespace).SetKmsProviders(kmsProviders)
	// verify
	keyVaultClient, err := mongo.NewClientEncryption(client, clientOptions)
	if err != nil {
		log.Fatalf("! key vault client - %v", err)
	}
	defer keyVaultClient.Close(context.TODO())
	dataKeyOptions := options.DataKey().SetMasterKey(masterKey).SetKeyAltNames([]string{"aws_altname"})
	_, err = keyVaultClient.CreateDataKey(context.TODO(), "aws", dataKeyOptions)
	if err != nil {
		log.Fatalf("! create data key - %v", err)
	}
}

// CreateMongoClient is the a mongo client sans encryption
func CreateMongoClient() (client *mongo.Client, err error) {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_ATLAS_URI"))
	client, err = mongo.Connect(context.TODO(), clientOptions)
	return
}

// CreateEncryptedMongoClient creates a client configured with encryption
func CreateEncryptedMongoClient() (client *mongo.Client, err error) {
	schemaMap := map[string]interface{}{
		"tutorial.foobar": readJSONFile("foobarSchemaMap.json"),
	}
	autoEncOptions := options.AutoEncryption().
		SetKeyVaultNamespace(namespace).
		SetKmsProviders(kmsProviders).
		SetSchemaMap(schemaMap)
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_ATLAS_URI")).SetAutoEncryptionOptions(autoEncOptions)
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
