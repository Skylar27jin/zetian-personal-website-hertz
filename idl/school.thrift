struct School {
    1: i64 id;
    2: string name;
    3: string short_name;
    4: list<string> aliases;
    5: string description;
    6: i64 created_at;
    7: i64 updated_at;

}


struct GetAllSchoolsReq {
}

struct GetAllSchoolsResp {
    1: bool isSuccessful;
    2: string errorMessage;
    3: list<School> Schools
}