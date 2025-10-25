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
