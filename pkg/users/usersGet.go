/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package users

import (
	"context"
	"fmt"

	"github.com/galasa-dev/cli/pkg/embedded"
	"github.com/galasa-dev/cli/pkg/galasaapi"
)

func GetUsers(loginId string, apiClient *galasaapi.APIClient) error {
	var err error
	var outputText []galasaapi.UserData

	err = validateLoginIdFlag(loginId)
	if err == nil {

		outputText, err = getUserDataFromRestApi(loginId, apiClient)
		fmt.Print(outputText)
	}

	return err
}

func getUserDataFromRestApi(
	loginId string,
	apiClient *galasaapi.APIClient,
) ([]galasaapi.UserData, error) {

	var err error
	var context context.Context = nil

	var restApiVersion string

	var userProperties = make([]galasaapi.UserData, 0)

	restApiVersion, err = embedded.GetGalasactlRestApiVersion()

	if err == nil {

		apiCall := apiClient.UsersAPIApi.GetUserByLoginId(context).LoginId(loginId).ClientApiVersion(restApiVersion)
		userProperties, _, err = apiCall.Execute()

	}

	return userProperties, err
}
