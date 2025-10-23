namespace go user


struct LoginReq {
    1: string email (api.query="email");
    2: string password (api.query="password");
}

struct LoginResp {
    1: bool isSuccessful;
    2: string errorMessage;
    3: string userName;
    4: string email;
}

struct SignUpReq {
    1: string username (api.query="username");
    2: string email    (api.query="email");
    3: string password (api.query="password");
}

struct SignUpResp {
    1: bool isSuccessful;
    2: string errorMessage;
    3: string userName;
    4: string email;
}
