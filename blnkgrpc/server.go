package blnkgrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/jerry-enebeli/blnk"
	"github.com/jerry-enebeli/blnk/api/middleware"
	pb "github.com/jerry-enebeli/blnk/proto"

	"github.com/sirupsen/logrus"
)


type BlnkServer struct {
	pb.UnimplementedBlnkServiceServer
	ctx   context.Context
	header []*pb.HeaderArray
	body  string
	path []string
	includes []string
	blnk   *blnk.Blnk
	auth   *middleware.AuthMiddleware
}


func NewGRPC(ctx context.Context, b *blnk.Blnk) *BlnkServer {
	auth := middleware.NewAuthMiddleware(b)
	return &BlnkServer{ctx: ctx, blnk: b, body: "", includes: nil, auth: auth}
}

func (s *BlnkServer) GSender(ctx context.Context, req *pb.GRequest) (*pb.GResponse, error) {
	path := strings.Split(req.Path, "/")
	includes := strings.Split(path[len(path)-1], "?")
	rpcName := strings.ToUpper(req.Method + "_" + path[1] + "_" + strconv.Itoa(len(path)))
	s.path = path
	for _,v := range(includes) {
		s.includes =  append(s.includes, v)
	}
	
	s.body=req.Body
	s.header= req.Header

	reflection := reflect.ValueOf(s).MethodByName(rpcName)
	if reflection.Kind() == 0 {
		errMessage := fmt.Sprintf("rpc implementation for method %s and rpcname %s was not found", req.Method, rpcName)

		response, err := s.generateErrorResponse(errMessage)
		if err != nil {
			logrus.Warningf("Failed to generate error response: %s", err.Error())
			return nil, fmt.Errorf("failed to generate error response: %s", err.Error())
		}

		return response, nil
	}

	results := reflection.Call(nil)
	errResult := results[1].Interface()

	if errResult != nil {
		response, err := s.generateErrorResponse(fmt.Sprintf("%s", errResult))
		if err != nil {
			logrus.Warningf("Failed to generate error response: %v", err.Error())
			return nil, fmt.Errorf("failed to generate error response: %s", err.Error())
		}

		return response, nil
	}

	result, ok := results[0].Interface().(*pb.GResponse)
	if !ok {
		return nil, fmt.Errorf("failed to assert rpc call result as response, got type %T", results[0].Interface())
	}

	return result, nil
}

func (s *BlnkServer) generateErrorResponse(errMessage string) (*pb.GResponse, error) {
	errData := map[string]string{
		"resp_code": "02",
		"resp_msg":  errMessage,
	}

	errDataBytes, err := json.Marshal(errData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal value to JSON: %s", err.Error())
	}

	errResponse := &pb.GResponse{
		ResponseCode: 999,
		ResponseMsg:  "Failed",
		Body:         string(errDataBytes),
		Status:       "nok",
	}

	return errResponse, nil
}
