namespace go verification

struct SendVeriCodeToEmailReq {
    1: string email (api.body="email");
    2: string purpose (api.body="purpose");
}


struct SendVeriCodeToEmailResp {
    1: bool is_successful;
    2: string error_message;
    3: i64 expire_at;
}



struct VerifyEmailCodeReq {
    1: string email (api.body="email");
    2: string verification_code (api.body="code");
}


struct VerifyEmailCodeResp {
    1: bool is_successful
    2: string error_message
}
