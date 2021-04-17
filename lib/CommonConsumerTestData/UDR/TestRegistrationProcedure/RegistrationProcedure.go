package TestRegistrationProcedure

import (
	"free5gc/lib/openapi/models"
	"free5gc/lib/MongoDBLibrary/logger"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"encoding/json"
	"strings"
	"strconv"
	"time"
)


// Here are my own commands
type buildData struct {
	RannnodeName		string 	`json:"rannnodeName" bson:"rannnodeName"`
	Plmnid				string 	`json:"plmnid" bson:"plmnid"`
	Tac					string 	`json:"tac" bson:"tac"`
	Snssai_sst			string 	`json:"snssai_sst" bson:"snssai_sst"`
	Snssai_sd			string 	`json:"snssai_sd" bson:"snssai_sd"`
	GNBId				string 	`json:"gNBId" bson:"gNBId"`
	Imsi				string 	`json:"imsi" bson:"imsi"`
	Dnn					string 	`json:"dnn" bson:"dnn"`
	AuthenticationSubs 	string 	`json:"authenticationSubs" bson:"authenticationSubs"`
	ServingPlmnId		string 	`json:"servingPlmnId" bson:"servingPlmnId"`
}

func stringToByteArray(str string) ([]byte){
	// use format "x02xf8x39" for representation of "\x02\xf8\x39"
	var outputArr []byte
	strSplit := strings.Split(str, "x")
	for i := 1; i < len(strSplit); i++ {
		hx,_ := strconv.ParseUint("0x"+strSplit[i], 0, 64)
		outputArr = append(outputArr,byte(hx))
	}
	return outputArr
}

var Client *mongo.Client = nil

func getbuildDataRegistrationProcedureFromMongoDB() (string, string, string) {
	url := "mongodb://localhost:27017"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	defer cancel()
	if err != nil {
		//defer cancel()
		logger.MongoDBLog.Panic(err.Error())
	}
	Client = client
	collection := Client.Database("build").Collection("config")

	var result map[string]interface{}
	filter := bson.M{"_id": 2}
	collection.FindOne(context.TODO(), filter).Decode(&result)

	tmp, err := json.Marshal(result)
	var buildDatainstance *buildData
	buildDatainstance = new(buildData)
	err = json.Unmarshal(tmp, buildDatainstance)

	var dnn = buildDatainstance.Dnn
	var Snssai_sst = strings.Replace(buildDatainstance.Snssai_sst, "x", "", -1)
	var Snssai_sd = strings.Replace(buildDatainstance.Snssai_sd, "x", "", -1)

	return dnn, Snssai_sst, Snssai_sd
}


const (
	FREE5GC_CASE = "hexoutNet"
)

var TestAmDataTable = make(map[string]models.AccessAndMobilitySubscriptionData)
var TestSmfSelDataTable = make(map[string]models.SmfSelectionSubscriptionData)
var TestSmSelDataTable = make(map[string]models.SessionManagementSubscriptionData)
var TestAmPolicyDataTable = make(map[string]models.AmPolicyData)
var TestSmPolicyDataTable = make(map[string]models.SmPolicyData)

func init() {

	dnnFromMongoDB, sstFromMongoDB, sdFromMongoDB  := getbuildDataRegistrationProcedureFromMongoDB()
	ix,_ := strconv.Atoi(sstFromMongoDB)
	sstFromMongoDBInt := int32(ix)
	
	TestAmDataTable[FREE5GC_CASE] = models.AccessAndMobilitySubscriptionData{
		Gpsis: []string{
			"msisdn-0900000000",
		},
		Nssai: &models.Nssai{
			DefaultSingleNssais: []models.Snssai{
				// {
				// 	Sst: 1,
				// 	Sd:  "010203",
				// },
				{
					Sst: sstFromMongoDBInt,
					Sd:  sdFromMongoDB,
				},
				{
					Sst: 1,
					Sd:  "112233",
				},
			},
			SingleNssais: []models.Snssai{
				// {
				// 	Sst: 1,
				// 	Sd:  "010203",
				// },
				{
					Sst: sstFromMongoDBInt,
					Sd:  sdFromMongoDB,
				},
				{
					Sst: 1,
					Sd:  "112233",
				},
			},
		},
		SubscribedUeAmbr: &models.AmbrRm{
			Uplink:   "1000 Kbps",
			Downlink: "1000 Kbps",
		},
	}

	// SMF Slicing data to be used here!
	TestSmfSelDataTable[FREE5GC_CASE] = models.SmfSelectionSubscriptionData{
		SubscribedSnssaiInfos: map[string]models.SnssaiInfo{
			// "01010203": { // sst:1, sd:010203
			sstFromMongoDB + sdFromMongoDB: { // sst:1, sd:010203
				DnnInfos: []models.DnnInfo{
					{
						// Dnn: "internetx",
						Dnn: dnnFromMongoDB,
					},
				},
			},
			"01112233": { // sst:1, sd:112233
				DnnInfos: []models.DnnInfo{
					{
						// Dnn: "internetx",
						Dnn: dnnFromMongoDB,
					},
				},
			},
		},
	}

	TestAmPolicyDataTable[FREE5GC_CASE] = models.AmPolicyData{
		SubscCats: []string{
			"hexoutNet",
		},
	}

	// SMF Slicing data to be used here!
	TestSmPolicyDataTable[FREE5GC_CASE] = models.SmPolicyData{
		SmPolicySnssaiData: map[string]models.SmPolicySnssaiData{
			// "01010203": {
			sstFromMongoDB + sdFromMongoDB: {
				Snssai: &models.Snssai{
					Sd:  sdFromMongoDB,
					Sst: sstFromMongoDBInt,
				},
				SmPolicyDnnData: map[string]models.SmPolicyDnnData{
					// "internetx": {
					// 	Dnn: "internetx",
					// },
					dnnFromMongoDB: {
						Dnn: dnnFromMongoDB,
					},
				},
			},
			"01112233": {
				Snssai: &models.Snssai{
					Sd:  "112233",
					Sst: 1,
				},
				SmPolicyDnnData: map[string]models.SmPolicyDnnData{
					// "internetx": {
					// 	Dnn: "internetx",
					// },
					dnnFromMongoDB: {
						Dnn: dnnFromMongoDB,
					},					
				},
			},
		},
	}

	TestSmSelDataTable[FREE5GC_CASE] = models.SessionManagementSubscriptionData{
		SingleNssai: &models.Snssai{
			Sst: sstFromMongoDBInt,
			Sd:  sdFromMongoDB,
		},
		DnnConfigurations: map[string]models.DnnConfiguration{
			// "internetx": {
			dnnFromMongoDB: {
				SscModes: &models.SscModes{
					DefaultSscMode:  models.SscMode__1,
					AllowedSscModes: []models.SscMode{models.SscMode__1, models.SscMode__2, models.SscMode__3},
				},
				PduSessionTypes: &models.PduSessionTypes{DefaultSessionType: models.PduSessionType_IPV4,
					AllowedSessionTypes: []models.PduSessionType{models.PduSessionType_IPV4},
				},
				SessionAmbr: &models.Ambr{
					Uplink:   "1000 Kbps",
					Downlink: "1000 Kbps",
				},
				Var5gQosProfile: &models.SubscribedDefaultQos{
					Var5qi: 9,
					Arp: &models.Arp{
						PriorityLevel: 8,
					},
					PriorityLevel: 8,
				},
			},
		},
	}
}
