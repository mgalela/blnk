package blnkgrpc

import (
	"encoding/json"
	"errors"

	"github.com/sirupsen/logrus"

	model "github.com/jerry-enebeli/blnk/model"
	pb "github.com/jerry-enebeli/blnk/proto"
)

//CreateAccount
func (s *BlnkServer) POS_ACCOUNTS_2() (*pb.GResponse, error) {
	var newAccount model.Account
	var resp pb.GResponse
	
	if err := json.Unmarshal([]byte(s.body), &newAccount); err != nil {
		logrus.Warningf("Invalid Account JSON body: %v", err.Error())
		return nil, errors.New("Invalid Account JSON body: " + err.Error())
	}

	accountResp, err := s.blnk.CreateAccount(newAccount)
	if err != nil {
		logrus.Warningf("Failed to Created Account: %v", err.Error())
		resp.Status = "nok"
		resp.ResponseCode = 500
		resp.ResponseMsg = err.Error()
		resp.Body = ""

		return &resp, err
	}

	respbody, err := json.Marshal(accountResp)
    if err != nil {
    	logrus.Warningf("Failed to Marshal Account: %v", err.Error())
		resp.Status = "nok"
		resp.ResponseCode = 500
		resp.ResponseMsg = err.Error()
		resp.Body = ""
	   	return &resp, err
    }

	resp.Status = "ok"
	resp.ResponseCode = 200
	resp.ResponseMsg = "Account Created"
	resp.Body = string(respbody)

	return &resp, nil
}

//GetAllAccount
func (s *BlnkServer) GET_ACCOUNTS_2() (*pb.GResponse, error) {
	var resp pb.GResponse

	listaccount, err := s.blnk.GetAllAccounts()
	if err != nil {
		logrus.Warningf("Invalid Account JSON Responses: %v", err.Error())
		resp.Status = "nok"
		resp.ResponseCode = 500
		resp.ResponseMsg = err.Error()
		resp.Body = ""
		return &resp, errors.New("Invalid Account JSON Responses: " + err.Error())
	}

	respbody, err := json.Marshal(listaccount)
	if err != nil {
		logrus.Warningf("Failed to Marshal Accounts: %v", err.Error())
		resp.Status = "nok"
		resp.ResponseCode = 500
		resp.ResponseMsg = err.Error()
		resp.Body = ""
		return &resp, err
	}

	resp.Status = "ok"
	resp.ResponseCode = 200
	resp.ResponseMsg = "Get Accounts Successfully"
	resp.Body = string(respbody)

	return &resp, nil

}

//GetAccount
func (s *BlnkServer) GET_ACCOUNTS_3() (*pb.GResponse, error) {
	var resp pb.GResponse

	listaccount, err := s.blnk.GetAccount(s.path[2], s.includes)
	if err != nil {
		logrus.Warningf("Invalid Account JSON Response: %v", err.Error())
		resp.Status = "nok"
		resp.ResponseCode = 500
		resp.ResponseMsg = err.Error()
		resp.Body = ""
		return &resp, errors.New("Invalid Account JSON Response: " + err.Error())
	}

	respbody, err := json.Marshal(listaccount)
	if err != nil {
		logrus.Warningf("Failed to Marshal Account: %v", err.Error())
		resp.Status = "nok"
		resp.ResponseCode = 500
		resp.ResponseMsg = err.Error()
		resp.Body = ""
		return &resp, err
	}

	resp.Status = "ok"
	resp.ResponseCode = 200
	resp.ResponseMsg = "Get Account Successfully"
	resp.Body = string(respbody)

	return &resp, nil

}
