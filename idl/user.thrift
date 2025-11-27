namespace go user


struct LoginReq {
    1: string email (api.body="email");
    2: string password (api.body="password");
}

struct LoginResp {
    1: bool isSuccessful;
    2: string errorMessage;
    3: string userName;
    4: string email;
}

struct SignUpReq {
    1: string username (api.body="username");
    2: string email    (api.body="email");
    3: string password (api.body="password");
}

struct SignUpResp {
    1: bool isSuccessful;
    2: string errorMessage;
    3: string userName;
    4: string email;
}



struct LogoutReq {
}

struct LogoutResp {
    1: bool isSuccessful;
    2: string errorMessage;
}


//get user by ID or name
//cannot pass both ID and name, at least one should be passed
struct GetUserReq {
    1: i64 id (api.query = "id");
    2: string name (api.query = "name");
}


struct GetUserResp {
    1: bool   isSuccessful;
    2: string errorMessage;
    3: string userName;
    4: i64    id;
    5: string avatarUrl;
}


struct ResetPasswordReq {
    1: string email (api.body="email");
    2: string newPassword (api.body="new_password");
}

struct ResetPasswordResp {
    1: bool isSuccessful;
    2: string errorMessage;
}



//place holder
//body里，avatar=要上传的文件
struct UpdateAvatarReq {
}

struct UpdateAvatarResp {
    1: bool isSuccessful;
    2: string errorMessage;
}
