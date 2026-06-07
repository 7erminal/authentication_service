package functions

import (
	"authentication_service/api"
	"authentication_service/structs/requestsDTOs"
	"authentication_service/structs/responsesDTOs"
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/logs"
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
	host, _ := beego.AppConfig.String("customerBaseUrl")

	reqText, _ := json.Marshal(req)

	logs.Info("Request to get User: ", string(reqText))

	request := api.NewRequest(
		host,
		"/v1/users/"+req.UserId,
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		setControllerJSON(c, err.Error())
		return responsesDTOs.UserResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responsesDTOs.UserResponseDTO{}, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, read, "", "  "); err != nil {
		logs.Info("Raw response received is ", string(read))
	} else {
		logs.Info("Raw response received is \n", prettyJSON.String())
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responsesDTOs.UserResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responsesDTOs.UserResponseDTO{}, err
	}
	setControllerJSON(c, data)

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}

func GetUserWithUsername(c *beego.Controller, req requestsDTOs.GetUserWithUsernameRequest) (responsesDTOs.UserResponseDTO, error) {
	host, _ := beego.AppConfig.String("customerBaseUrl")

	reqText, _ := json.Marshal(req)

	logs.Info("Request to get User: ", string(reqText))

	request := api.NewRequest(
		host,
		"/v1/users/get-user-by-username/"+req.Username,
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		setControllerJSON(c, err.Error())
		return responsesDTOs.UserResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responsesDTOs.UserResponseDTO{}, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, read, "", "  "); err != nil {
		logs.Info("Raw response received is ", string(read))
	} else {
		logs.Info("Raw response received is \n", prettyJSON.String())
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responsesDTOs.UserResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responsesDTOs.UserResponseDTO{}, err
	}
	setControllerJSON(c, data)

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}

func GetCustomer(c *beego.Controller, req requestsDTOs.GetCustomerRequest) (responsesDTOs.CustomerResponseDTO, error) {
	host, _ := beego.AppConfig.String("customerBaseUrl")

	reqText, _ := json.Marshal(req)

	logs.Info("Request to get Customer: ", string(reqText))

	request := api.NewRequest(
		host,
		"/v1/customers/"+req.CustomerId,
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		setControllerJSON(c, err.Error())
		return responsesDTOs.CustomerResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responsesDTOs.CustomerResponseDTO{}, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, read, "", "  "); err != nil {
		logs.Info("Raw response received is ", string(read))
	} else {
		logs.Info("Raw response received is \n", prettyJSON.String())
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responsesDTOs.CustomerResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responsesDTOs.CustomerResponseDTO{}, err
	}
	setControllerJSON(c, data)

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}

func GetCustomerWithUsername(c *beego.Controller, req requestsDTOs.GetCustomerWithUsernameRequest) (responsesDTOs.CustomerResponseDTO, error) {
	host, _ := beego.AppConfig.String("customerBaseUrl")

	reqText, _ := json.Marshal(req)

	logs.Info("Request to get Customer with username: ", string(reqText))

	request := api.NewRequest(
		host,
		"/v1/customers/username/"+req.Username,
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		setControllerJSON(c, err.Error())
		return responsesDTOs.CustomerResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responsesDTOs.CustomerResponseDTO{}, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, read, "", "  "); err != nil {
		logs.Info("Raw response received is ", string(read))
	} else {
		logs.Info("Raw response received is \n", prettyJSON.String())
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responsesDTOs.CustomerResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responsesDTOs.CustomerResponseDTO{}, err
	}
	setControllerJSON(c, data)

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}

func UpdateUserPassword(c *beego.Controller, req requestsDTOs.UpdateUserPasswordRequest) (responsesDTOs.UserResponseDTO, error) {
	host, _ := beego.AppConfig.String("customerBaseUrl")

	reqText, _ := json.Marshal(req)

	logs.Info("Request to update User password: ", string(reqText))

	request := api.NewRequest(
		host,
		"/v1/users/password/"+strconv.FormatInt(req.UserId, 10),
		api.PUT)
	if request.InterfaceParams == nil {
		request.InterfaceParams = map[string]interface{}{}
	}

	request.InterfaceParams["UserId"] = req.UserId
	request.InterfaceParams["OldPassword"] = req.OldPassword
	request.InterfaceParams["NewPassword"] = req.NewPassword

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		setControllerJSON(c, err.Error())
		return responsesDTOs.UserResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responsesDTOs.UserResponseDTO{}, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, read, "", "  "); err != nil {
		logs.Info("Raw response received is ", string(read))
	} else {
		logs.Info("Raw response received is \n", prettyJSON.String())
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responsesDTOs.UserResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responsesDTOs.UserResponseDTO{}, err
	}
	setControllerJSON(c, data)

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}
