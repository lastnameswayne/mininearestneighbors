package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Ip     string `json:"ip"`
	DbName string `json:"dbName"`
}

var mongocfg = MongoConfig{
	Ip:     "mongodb+srv://swayne:swayne@cluster0.85aqi48.mongodb.net/?retryWrites=true&w=majority",
	DbName: "mesure",
}

type Measurement struct {
	Name  string         `json:"name"`
	Sizes []SizeKeyValue `json:"sizes"`
}

type SizeKeyValue struct {
	K string `json:"k"`
	V struct {
		Cm   MeasurementValue `json:"cm"`
		Inch MeasurementValue `json:"inch"`
	} `json:"v"`
}

type MeasurementValue struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type ProductMeasurements struct {
	SSenseProductID string      `json:"ssenseproductId"`
	Chest           Measurement `json:"Chest"`
	Shoulder        Measurement `json:"Shoulder"`
	SleeveLength    Measurement `json:"SleeveLength"`
	Length          Measurement `json:"Length"`
}

func Read() {
	cursor := fetch()
	for cursor.Next(context.TODO()) {
		var result ProductMeasurements
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v\n", result)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}

func newClient(mongoConfig MongoConfig) (*mongo.Client, error) {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoConfig.Ip).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	// Send a ping to confirm a successful connection
	err = client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func fetch() *mongo.Cursor {
	// Set client options
	ctx := context.TODO()
	client, err := newClient(mongocfg)
	if err != nil {
		panic(err)
	}

	// Open an aggregation cursor
	coll := client.Database("mesure").Collection("measurements")
	res, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

var pipeline = bson.A{
	bson.D{
		{"$project",
			bson.D{
				{"ssenseproductId", 1},
				{"SizesArrayChest",
					bson.D{
						{"$filter",
							bson.D{
								{"input", "$measurements"},
								{"as", "measurement"},
								{"cond",
									bson.D{
										{"$eq",
											bson.A{
												"$$measurement.name",
												"Chest",
											},
										},
									},
								},
							},
						},
					},
				},
				{"SizesArrayShoulder",
					bson.D{
						{"$filter",
							bson.D{
								{"input", "$measurements"},
								{"as", "measurement"},
								{"cond",
									bson.D{
										{"$eq",
											bson.A{
												"$$measurement.name",
												"Shoulder",
											},
										},
									},
								},
							},
						},
					},
				},
				{"SizesArrayLength",
					bson.D{
						{"$filter",
							bson.D{
								{"input", "$measurements"},
								{"as", "measurement"},
								{"cond",
									bson.D{
										{"$eq",
											bson.A{
												"$$measurement.name",
												"Length",
											},
										},
									},
								},
							},
						},
					},
				},
				{"SizesArraySleeveLength",
					bson.D{
						{"$filter",
							bson.D{
								{"input", "$measurements"},
								{"as", "measurement"},
								{"cond",
									bson.D{
										{"$eq",
											bson.A{
												"$$measurement.name",
												"Sleeve Length",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
	bson.D{
		{"$project",
			bson.D{
				{"ssenseproductId", 1},
				{"name", 1},
				{"url", 1},
				{"image", 1},
				{"priceByCountry", 1},
				{"brand", 1},
				{"Chest",
					bson.D{
						{"name", "Chest"},
						{"sizes",
							bson.D{
								{"$objectToArray",
									bson.D{
										{"$arrayElemAt",
											bson.A{
												"$SizesArrayChest.sizes",
												0,
											},
										},
									},
								},
							},
						},
					},
				},
				{"Shoulder",
					bson.D{
						{"name", "Shoulder"},
						{"sizes",
							bson.D{
								{"$objectToArray",
									bson.D{
										{"$arrayElemAt",
											bson.A{
												"$SizesArrayShoulder.sizes",
												0,
											},
										},
									},
								},
							},
						},
					},
				},
				{"SleeveLength",
					bson.D{
						{"name", "Sleeve Length"},
						{"sizes",
							bson.D{
								{"$objectToArray",
									bson.D{
										{"$arrayElemAt",
											bson.A{
												"$SizesArraySleeveLength.sizes",
												0,
											},
										},
									},
								},
							},
						},
					},
				},
				{"Length",
					bson.D{
						{"name", "Length"},
						{"sizes",
							bson.D{
								{"$objectToArray",
									bson.D{
										{"$arrayElemAt",
											bson.A{
												"$SizesArrayLength.sizes",
												0,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}
