namespace go verification

struct SendVeriCodeToEmailReq {
    1: string email (api.body="email");      // 用户邮箱地址
    2: optional string purpose (api.body="purpose");   // 验证用途（如 signup / reset_password / bind_email）
}


struct SendVeriCodeToEmailResp {
    1: bool is_successful;       // 是否发送成功
    2: optional string error_message;  // 失败时的错误信息（如 “邮箱格式错误” / “发送过于频繁”）
    3: optional i64 expire_at;         // 验证码过期时间（Unix 时间戳）
}



struct VerifyEmailCodeReq {
    1: string email (api.body="email");           // 用户邮箱
    2: string verification_code (api.body="code");  // 用户输入的验证码
}


struct VerifyEmailCodeResp {
    1: bool is_successful       // 是否验证成功
    2: optional string error_message // 错误原因（如 “验证码错误” / “验证码已过期”）
    3: optional string jwt_token      // （可选）验证成功后颁发的 JWT（用于修改密码、注册等操作）
}
