/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package properties

import (
	"context"
	"io"
	"net/http"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/galasaapi"
)

// DeleteProperty - performs all the logic to implement the `galasactl properties delete` command,
// but in a unit-testable manner.
func DeleteProperty(
	namespace string,
	name string,
	apiClient *galasaapi.APIClient,
) error {
	var err error
	err = validateInputsAreNotEmpty(namespace, name)
	if err == nil {
		err = deleteCpsProperty(namespace, name, apiClient)
	}
	return err
}

func deleteCpsProperty(namespace string,
	name string,
	apiClient *galasaapi.APIClient,
) error {
	var err error = nil
	var resp *http.Response
	var context context.Context = nil
	var responseBody []byte

	apicall := apiClient.ConfigurationPropertyStoreAPIApi.DeleteCpsProperty(context, namespace, name)
	_, resp, err = apicall.Execute()

	if (resp != nil) && (resp.StatusCode != http.StatusOK) {
		defer resp.Body.Close()
		
		responseBody, err = io.ReadAll(resp.Body)
		if err == nil{
			var apiError *galasaErrors.GalasaAPIError
			apiError, err = galasaErrors.GetApiErrorFromResponse(responseBody)

			if err == nil {
				//return galasa api error, because status code is not 200 (OK)
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_DELETE_PROPERTY_FAILED, name, apiError.Message)
			} else{
				//error occurred when trying to retrieve the api error
				//unable to retrieve galasa apiError
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_UNABLE_TO_GET_API_ERROR, err)
			}

		} else {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_UNABLE_TO_READ_RESPONSE_BODY, err)
		}
	}
	return err
}
