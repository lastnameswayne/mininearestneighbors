package fetch

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapToVector(t *testing.T) {
	productMeasurements := ProductMeasurements{
		SSenseProductID: "13408241",
		Chest: Measurement{
			Name: "Chest",
			Sizes: []SizeKeyValue{
				{K: "S", V: MeasurementGroup{Cm: MeasurementValue{Value: 54, Unit: "cm"}, Inch: MeasurementValue{Value: 21.5, Unit: "inch"}}},
				{K: "M", V: MeasurementGroup{Cm: MeasurementValue{Value: 57, Unit: "cm"}, Inch: MeasurementValue{Value: 22.75, Unit: "inch"}}},
				{K: "L", V: MeasurementGroup{Cm: MeasurementValue{Value: 60, Unit: "cm"}, Inch: MeasurementValue{Value: 23.75, Unit: "inch"}}},
				{K: "XL", V: MeasurementGroup{Cm: MeasurementValue{Value: 62, Unit: "cm"}, Inch: MeasurementValue{Value: 24.75, Unit: "inch"}}},
			},
		},
		Shoulder: Measurement{
			Name: "Shoulder",
			Sizes: []SizeKeyValue{
				{K: "S", V: MeasurementGroup{
					Cm: MeasurementValue{Value: 45, Unit: "cm"}, Inch: MeasurementValue{Value: 18, Unit: "inch"}}},
				{K: "M", V: MeasurementGroup{Cm: MeasurementValue{Value: 46, Unit: "cm"}, Inch: MeasurementValue{Value: 18.25, Unit: "inch"}}},
				{K: "L", V: MeasurementGroup{Cm: MeasurementValue{Value: 48, Unit: "cm"}, Inch: MeasurementValue{Value: 19, Unit: "inch"}}},
				{K: "XL", V: MeasurementGroup{Cm: MeasurementValue{Value: 48, Unit: "cm"}, Inch: MeasurementValue{Value: 19, Unit: "inch"}}},
			},
		},
		SleeveLength: Measurement{
			Name: "Sleeve Measurement",
			Sizes: []SizeKeyValue{
				{K: "S", V: MeasurementGroup{Cm: MeasurementValue{Value: 64, Unit: "cm"}, Inch: MeasurementValue{Value: 25.5, Unit: "inch"}}},
				{K: "M", V: MeasurementGroup{Cm: MeasurementValue{Value: 65, Unit: "cm"}, Inch: MeasurementValue{Value: 25.75, Unit: "inch"}}},
				{K: "L", V: MeasurementGroup{Cm: MeasurementValue{Value: 65, Unit: "cm"}, Inch: MeasurementValue{Value: 25.75, Unit: "inch"}}},
				{K: "XL", V: MeasurementGroup{Cm: MeasurementValue{Value: 65, Unit: "cm"}, Inch: MeasurementValue{Value: 25.75, Unit: "inch"}}},
			},
		},
		Length: Measurement{
			Name: "Length",
			Sizes: []SizeKeyValue{
				{K: "S", V: MeasurementGroup{Cm: MeasurementValue{Value: 70, Unit: "cm"}, Inch: MeasurementValue{Value: 27.75, Unit: "inch"}}},
				{K: "M", V: MeasurementGroup{Cm: MeasurementValue{Value: 70, Unit: "cm"}, Inch: MeasurementValue{Value: 27.75, Unit: "inch"}}},
				{K: "L", V: MeasurementGroup{Cm: MeasurementValue{Value: 71, Unit: "cm"}, Inch: MeasurementValue{Value: 28.25, Unit: "inch"}}},
				{K: "XL", V: MeasurementGroup{Cm: MeasurementValue{Value: 73, Unit: "cm"}, Inch: MeasurementValue{Value: 29, Unit: "inch"}}},
			},
		},
	}

	//I want output (select cm vals, group by sizes)
	var expectedSVector = []int{54, 45, 64, 70}
	var expectedMVector = []int{57, 46, 65, 70}
	var expectedLVector = []int{60, 48, 65, 71}
	var expectedXLVector = []int{62, 48, 65, 73}
	slices.Sort(expectedSVector)
	slices.Sort(expectedMVector)
	slices.Sort(expectedLVector)
	slices.Sort(expectedXLVector)

	t.Run("maps to one vector for each size", func(t *testing.T) {
		vectors := mapToVector(productMeasurements)
		assert.Len(t, vectors, 4) //4 measurements

		// Assert that each vector has the correct size and measurements
		for _, vector := range vectors {
			assert.Equal(t, 4, len(vector.Vector)) // Each vector should have 4 measurements

			res := vector.Vector
			slices.Sort(res)

			if vector.Size == "S" {
				assert.Equal(t, res, expectedSVector)
			}
			if vector.Size == "M" {
				assert.Equal(t, res, expectedMVector)
			}
			if vector.Size == "L" {
				assert.Equal(t, res, expectedLVector)
			}
			if vector.Size == "XL" {
				assert.Equal(t, res, expectedXLVector)
			}
		}
	})
}
