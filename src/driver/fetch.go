package driver

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	h "github.com/lastnameswayne/mininearestneighbors/src/hnsw"
	v "github.com/lastnameswayne/mininearestneighbors/src/vector"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Ip     string `json:"ip"`
	DbName string `json:"dbName"`
}

type Measurement struct {
	Name  string         `json:"name"`
	Sizes []SizeKeyValue `json:"sizes"`
}

type SizeKeyValue struct {
	K string           `json:"k"`
	V MeasurementGroup `json:"v"`
}

type MeasurementGroup struct {
	Cm   MeasurementValue `json:"cm"`
	Inch MeasurementValue `json:"inch"`
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
	var mongoIP = os.Getenv("MONGO_IP")
	var mongoDB = os.Getenv("MONGO_DB")
	var mongoCfg = MongoConfig{
		Ip:     mongoIP,
		DbName: mongoDB,
	}

	fmt.Println("mongoIP", mongoIP)
	layerCount := 15
	M := 5
	mMax := 2 * M //recommended
	efSize := 50
	hnsw := h.ConstructHNSW(layerCount)
	cursor := fetch(mongoCfg)

	for cursor.Next(context.TODO()) {
		var result ProductMeasurements
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		resultMapped := mapToVector(result)
		for _, vector := range resultMapped {
			hnsw = hnsw.InsertVector(vector, efSize, M, mMax)
			fmt.Println("inserted", vector.Id, vector.Size)
		}
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done")
	ints := []int{63, 73, 55, 67}
	//chest, length, shoulder, sleeve length
	q := v.Vector{
		Id:     "123",
		Vector: ints,
	}
	res := hnsw.Search(q, efSize, 50)
	fmt.Println(res)
	for idx, vertex := range res {
		fmt.Println("item", idx, vertex.Id, "distance", v.Distance(vertex.Vector, q.Vector), "and sizes", q.Vector)
	}
}

func mapToVector(m ProductMeasurements) []v.Vector {
	ssenseID, err := strconv.ParseInt(m.SSenseProductID, 10, 64)
	if err != nil {
		panic(err)
	}

	measurements := []Measurement{m.Chest, m.Length, m.Shoulder, m.SleeveLength}

	return mapMeasurement(measurements, int(ssenseID))

}

func mapMeasurement(measurements []Measurement, ssenseID int) []v.Vector {
	res := []v.Vector{}

	sizeToInts := fillMap(measurements)
	for size, vals := range sizeToInts {
		vector := v.Vector{
			Id:     addsizetoId(ssenseID, size),
			Size:   size,
			Vector: vals,
		}
		res = append(res, vector)
	}

	return res
}

func addsizetoId(id int, size string) string {
	return strconv.Itoa(id) + size
}

func fillMap(measurements []Measurement) map[string][]int {
	sizeToInts := map[string][]int{}
	for _, m := range measurements {
		for _, size := range m.Sizes {
			sizeToInts[size.K] = append(sizeToInts[size.K], int(size.V.Cm.Value))
		}

	}

	return sizeToInts
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

func fetch(mongoCfg MongoConfig) *mongo.Cursor {
	// Set client options
	ctx := context.TODO()
	client, err := newClient(mongoCfg)
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
