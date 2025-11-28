namespace go category

struct Category {
    1: i64 id,
    2: string name,
    3: string key,
    4: list<string> aliases,
    5: string description,
}

struct GetAllCategoriesReq {
}

struct GetAllCategoriesResp {
    1: bool isSuccessful,
    2: string errorMessage,
    3: list<Category> categories,
}
