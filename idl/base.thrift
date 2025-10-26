namespace go base 

include "user.thrift"
include "numberOperation.thrift"
include "verification.thrift"

service UserService {
    user.LoginResp Login(1: user.LoginReq request) (api.post="/login");
    user.SignUpResp SignUp(1: user.SignUpReq request) (api.post="/signup");
}

service NumberOperationService {
    numberOperation.GetToBinaryResp GetToBinary(1: numberOperation.GetToBinaryReq request) (api.get="/to_binary");
}

service VerificationService {
    verification.SendVeriCodeToEmailResp SendVeriCodeToEmail(1: verification.SendVeriCodeToEmailReq request) (api.post="/verification/email/send-code")
    //1.Generate a 6 bit verification code; 2. send the code to the email; 3.store it to the db
    verification.VerifyEmailCodeResp VerifyEmailCode(1: verification.VerifyEmailCodeReq request) (api.post="/verification/email/verify-code")
    //1. check if the code is correct; 2.if correct, disable this code and give the user a veriEmailJWT
}