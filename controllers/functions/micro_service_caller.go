package functions

import (
	"authentication_service/api"
	"authentication_service/structs/requestsDTOs"
	"authentication_service/structs/responsesDTOs"
	"context"
	"encoding/json"
	"io"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

// containsIgnoreCase checks if substr is in s, case-insensitive.
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) &&
		strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func setControllerJSON(c *beego.Controller, value interface{}) {
	if c == nil {
		return
	}
	if c.Data == nil {
		c.Data = map[interface{}]interface{}{}
	}
	c.Data["json"] = value
}

func GetUser(c *beego.Controller, req requestsDTOs.GetUserRequest) (responsesDTOs.UserResponseDTO, error) {
	return getUser(controllerContext(c), c, req)
}

func GetUserContext(ctx context.Context, req requestsDTOs.GetUserRequest) (responsesDTOs.UserResponseDTO, error) {
	return getUser(ctx, nil, req)
}

func getUser(ctx context.Context, c *beego.Controller, req requestsDTOs.GetUserRequest) (responsesDTOs.UserResponseDTO, error) {
	host, _ := beego.AppConfig.String("customerBaseUrl")

	request := api.NewRequest(
		host,
		"/v1/users/"+req.UserId,
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request:  request,
		Type_:    "params",
		Internal: true,
	}
	res, err := client.SendRequestContext(ctx)
	if err != nil {

		setControllerJSON(c, err.Error())
		return responsesDTOs.UserResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responsesDTOs.UserResponseDTO{}, err
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responsesDTOs.UserResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responsesDTOs.UserResponseDTO{}, err
	}
	setControllerJSON(c, data)

	return data, nil
}

func GetUserWithUsername(c *beego.Controller, req requestsDTOs.GetUserWithUsernameRequest) (responsesDTOs.UserResponseDTO, error) {
	host, _ := beego.AppConfig.String("customerBaseUrl")

	request := api.NewRequest(
		host,
		"/v1/users/get-user-by-username/"+req.Username,
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request:  request,
		Type_:    "params",
		Internal: true,
	}
	res, err := client.SendRequestContext(controllerContext(c))
	if err != nil {

		setControllerJSON(c, err.Error())
		return responsesDTOs.UserResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responsesDTOs.UserResponseDTO{}, err
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responsesDTOs.UserResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responsesDTOs.UserResponseDTO{}, err
	}
	setControllerJSON(c, data)

	return data, nil
}

func GetRole(c *beego.Controller, roleId string) (responsesDTOs.RoleApiResponseDTO, error) {
	host, _ := beego.AppConfig.String("customerBaseUrl")

	request := api.NewRequest(
		host,
		"/v1/roles/"+roleId,
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request:  request,
		Type_:    "params",
		Internal: true,
	}
	res, err := client.SendRequestContext(controllerContext(c))
	if err != nil {

		setControllerJSON(c, err.Error())
		return responsesDTOs.RoleApiResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responsesDTOs.RoleApiResponseDTO{}, err
	}

	// var prettyJSON bytes.Buffer
	// if err := json.Indent(&prettyJSON, read, "", "  "); err != nil {
	// } else {
	// }
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responsesDTOs.RoleApiResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responsesDTOs.RoleApiResponseDTO{}, err
	}
	setControllerJSON(c, data)

	return data, nil
}

func GetCustomer(c *beego.Controller, req requestsDTOs.GetCustomerRequest) (responsesDTOs.CustomerResponseDTO, error) {
	host, _ := beego.AppConfig.String("customerBaseUrl")

	request := api.NewRequest(
		host,
		"/v1/customers/"+req.CustomerId,
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request:  request,
		Type_:    "params",
		Internal: true,
	}
	res, err := client.SendRequestContext(controllerContext(c))
	if err != nil {

		setControllerJSON(c, err.Error())
		return responsesDTOs.CustomerResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responsesDTOs.CustomerResponseDTO{}, err
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responsesDTOs.CustomerResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responsesDTOs.CustomerResponseDTO{}, err
	}
	setControllerJSON(c, data)

	return data, nil
}

func GetCustomerWithUsername(c *beego.Controller, req requestsDTOs.GetCustomerWithUsernameRequest) (responsesDTOs.CustomerResponseDTO, error) {
	host, _ := beego.AppConfig.String("customerBaseUrl")

	request := api.NewRequest(
		host,
		"/v1/customers/username/"+req.Username,
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request:  request,
		Type_:    "params",
		Internal: true,
	}
	res, err := client.SendRequestContext(controllerContext(c))
	if err != nil {

		setControllerJSON(c, err.Error())
		return responsesDTOs.CustomerResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responsesDTOs.CustomerResponseDTO{}, err
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responsesDTOs.CustomerResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responsesDTOs.CustomerResponseDTO{}, err
	}
	setControllerJSON(c, data)

	return data, nil
}

func UpdateUserPassword(c *beego.Controller, req requestsDTOs.UpdateUserPasswordRequest) (responsesDTOs.UserResponseDTO, error) {
	host, _ := beego.AppConfig.String("customerBaseUrl")

	request := api.NewRequest(
		host,
		"/v1/users/password/"+req.UserId,
		api.PUT)
	if request.InterfaceParams == nil {
		request.InterfaceParams = map[string]interface{}{}
	}

	request.InterfaceParams["UserId"] = req.UserId
	request.InterfaceParams["OldPassword"] = req.OldPassword
	request.InterfaceParams["NewPassword"] = req.NewPassword

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request:  request,
		Type_:    "params",
		Internal: true,
	}
	res, err := client.SendRequestContext(controllerContext(c))
	if err != nil {

		setControllerJSON(c, err.Error())
		return responsesDTOs.UserResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responsesDTOs.UserResponseDTO{}, err
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responsesDTOs.UserResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responsesDTOs.UserResponseDTO{}, err
	}
	setControllerJSON(c, data)

	return data, nil
}

// controllerContext supports legacy non-request callers without manufacturing IDs.
func controllerContext(c *beego.Controller) context.Context {
	if c != nil && c.Ctx != nil && c.Ctx.Request != nil {
		return c.Ctx.Request.Context()
	}
	return context.Background()
}
