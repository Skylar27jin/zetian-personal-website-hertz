namespace go base 

include "user.thrift"
include "numberOperation.thrift"

service UserService {
    user.LoginResp Login(1: user.LoginReq request) (api.post="/login");
    user.SignUpResp SignUp(1: user.SignUpReq request) (api.post="/signup");
}

service NumberOperationService {
    numberOperation.GetToBinaryResp GetToBinary(1: numberOperation.GetToBinaryReq request) (api.get="/to_binary");
}