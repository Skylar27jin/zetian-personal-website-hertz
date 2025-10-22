namespace go numberOperation

struct GetToBinaryReq {
    1: i32 number (api.query="number");
}

struct GetToBinaryResp {
    1: string res;
}


struct DecodeJWTReq {
}

struct DecodeJWTResp {
    1: bool isValid;
    2: map<string, string> payLoad;
}

service NumberOperationService {
    GetToBinaryResp GetToBinary(1: GetToBinaryReq request) (api.get="/to_binary");
    DecodeJWTResp DecodeJWT(1: DecodeJWTReq request) (api.get="/decode_jwt");
}