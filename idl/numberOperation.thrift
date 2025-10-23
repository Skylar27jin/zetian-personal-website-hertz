namespace go numberOperation

struct GetToBinaryReq {
    1: i32 number (api.query="number");
}

struct GetToBinaryResp {
    1: string res;
}