package fle

import (
	"context"
	"encoding/base64"
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
	encRulesFile string
)

func init() {

	// vars set up
	keydb = os.Getenv("MONGODB_KEY_VAULT_DB")
	keycoll = os.Getenv("MONGODB_KEY_VAULT_COLL")
	encRulesFile = os.Getenv("ENCRYPTION_RULES_JSON_FILE")
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
	dataKeyOptions := options.DataKey().SetMasterKey(masterKey).SetKeyAltNames([]string{"altname", "aws_dk", "awskms"})
	result, err := keyVaultClient.CreateDataKey(context.TODO(), "aws", dataKeyOptions)
	if err != nil {
		log.Fatalf("! create data key - %v", err)
	}

	// Insert the data specified into the admin.datakeys
	log.Infof(" data key is {\"$binary\": { \"base64\": \"%+v\", \"subType\": \"%+v\" } } ", base64.StdEncoding.EncodeToString(result.Data), result.Subtype)
	log.Infof(" For Deterministic algorithm, you ****MUST**** provide the above result as part of the \"keyId\" in %+v", encRulesFile)
}

// CreateMongoClient is
func CreateMongoClient(uri string, wc *writeconcern.WriteConcern, useEncryption bool) (client *mongo.Client, err error) {

	clientOptions := options.Client().ApplyURI(uri).SetWriteConcern(wc)

	/////////////////// adds encryption to the mongo client
	if useEncryption {
		var sm bson.D
		var content []byte
		content, err = ioutil.ReadFile(encRulesFile)
		if err != nil {
			log.Errorf("! reading %v:%v", encRulesFile, err)
			return
		}

		if err = bson.UnmarshalExtJSON(content, false, &sm); err != nil {
			log.Errorf("! unmarshal %v: %v", encRulesFile, err)
			return
		}

		schemaMap := map[string]interface{}{
			"tutorial.foobar": sm,
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
