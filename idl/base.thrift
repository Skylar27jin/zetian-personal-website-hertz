namespace go base 

include "user.thrift"
include "numberOperation.thrift"
include "verification.thrift"
include "post.thrift"

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

    verification.MeResp Me(1: verification.MeReq request) (api.get="/me")
    //1, get "JWT" from Cookie; 2, verify whether cookie is not expired and JWT is not expired; 3, if both are not expired, return id, name, and email
}


service PostService {
    post.GetPostByIDResp GetPostByID(1: post.GetPostByIDReq request) (api.get="/post/get")

    post.CreatePostResp CreatePost(1: post.CreatePostReq request) (api.post="/post/create")

    post.EditPostResp EditPost(1: post.EditPostReq request) (api.post="/post/edit")

    post.DeletePostResp DeletePost(1: post.DeletePostReq request) (api.post="/post/delete")

    post.GetSchoolRecentPostsResp GetSchoolRecentPosts(1: post.GetSchoolRecentPostsReq request) (api.get="/post/school/recent")

    post.GetPersonalRecentPostsReq GetPersonalRecentPosts(1: post.GetPersonalRecentPostsResp request) (api.get="/post/personal")

    //like, unlike, fav, unfav will authorize user based on the JWT
    //user_id likes post_id, but will check whether user_id == JWT

    post.LikePostReq LikePost(1: post.UserFlagPostResq request) (api.post="/post/like")
    post.UnlikePostReq UnlikePost(1: post.UserFlagPostResq request) (api.post="/post/unlike")
    post.FavPostReq FavPost(1: post.UserFlagPostResq request) (api.post="/post/fav")
    post.UnfavPostReq UnfavPost(1: post.UserFlagPostResq request) (api.post="/post/unfav")
}
