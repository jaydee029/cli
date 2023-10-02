/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package properties

import (
	"testing"

	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/stretchr/testify/assert"
)

func createCpsPropertyForConverter(name string, value string) galasaapi.CpsProperty {
	// name = ""
	// value = ""

	cpsStructure := galasaapi.CpsProperty{
		Name:  &name,
		Value: &value,
	}

	return cpsStructure
}

func TestGasalsaapiPropertyWithNoRecordsreturnsNoRecord(t *testing.T) {
	//Given
	properties := make([]galasaapi.CpsProperty, 0)

	//When
	output := FormattablePropertyFromGalasaApi(properties)

	//Then
	assert.Equal(t, 0, len(output), "The input record is empty and so should be the output record")
}

func TestGalasaapiPropertyWithRecordsReturnsSameAmountOfRecords(t *testing.T) {
	//Given
	properties := make([]galasaapi.CpsProperty, 0)

	property1 := createCpsPropertyForConverter("framework.name1", "value1")
	property2 := createCpsPropertyForConverter("multi.name1", "multValue1")
	property3 := createCpsPropertyForConverter("multi.name2", "multValue2")
	property4 := createCpsPropertyForConverter("framework.name2", "value2")
	property5 := createCpsPropertyForConverter("multi.name3", "multValue3")
	properties = append(properties, property1, property2, property3, property4, property5)

	//When
	output := FormattablePropertyFromGalasaApi(properties)

	//Then
	assert.Equal(t, len(properties), len(output), "The input record has a length of %v whilst the output has length of %v", len(properties), len(output))
	assert.Equal(t, "framework", output[0].Namespace)
	assert.Equal(t, "name1", output[0].Name)
	assert.Equal(t, "value1", output[0].Value)
	assert.Equal(t, "multi", output[2].Namespace)
	assert.Equal(t, "name2", output[2].Name)
	assert.Equal(t, "multValue2", output[2].Value)
}

func TestGalasaapiPropertyWithRecordsReturnsRecordsInOrder(t *testing.T) {
	//Given
	properties := make([]galasaapi.CpsProperty, 0)

	property1 := createCpsPropertyForConverter("framework.name1", "value1")
	property2 := createCpsPropertyForConverter("framework.name2", "value2")
	property3 := createCpsPropertyForConverter("framework.jindex", "jindexValue")
	property4 := createCpsPropertyForConverter("framework.custard", "custardValue")
	properties = append(properties, property1, property2, property3, property4)

	//When
	output := FormattablePropertyFromGalasaApi(properties)

	//Then
	assert.Equal(t, len(properties), len(output), "The input record has a length of %v whilst the output has length of %v", len(properties), len(output))
	assert.Equal(t, "custard", output[0].Name)
	assert.Equal(t, "jindex", output[1].Name)
	assert.Equal(t, "name1", output[2].Name)
	assert.Equal(t, "name2", output[3].Name)
}
