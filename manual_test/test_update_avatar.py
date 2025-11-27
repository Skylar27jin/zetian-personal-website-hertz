import requests
from pathlib import Path

# ============================
# 全局配置：你只需要改这里！
# ============================

API_URL = "http://localhost:8888"       # 或 https://api.skylar27.com
EMAIL = "skyjin0127@gmail.com"        # 测试账号
PASSWORD = "Jzt20050127!"              # 测试密码
AVATAR_FILE = (Path(__file__).resolve().parent / "me.jpg")  # 本地图片路径

# ============================


def login_and_get_jwt():
    url = f"{API_URL}/login"
    data = {
        "email": EMAIL,
        "password": PASSWORD,
    }

    print(f"\n>>> Logging in as {EMAIL} ...")
    resp = requests.post(url, data=data)

    print("Login Status:", resp.status_code)
    print("Login Body:", resp.text)

    if resp.status_code != 200:
        print("[ERROR] login request failed")
        return None

    jwt = resp.cookies.get("JWT")
    if not jwt:
        print("[ERROR] login succeeded but no JWT returned!")
        return None

    print(">>> JWT acquired!\n")
    return jwt


def update_avatar(jwt):
    url = f"{API_URL}/user/update-avatar"

    print(f">>> Uploading avatar: {AVATAR_FILE}")

    files = {
        "avatar": open(AVATAR_FILE, "rb")
    }
    cookies = {
        "JWT": jwt
    }

    resp = requests.post(url, files=files, cookies=cookies)
    files["avatar"].close()

    print("\nUpdate Status:", resp.status_code)
    try:
        print("Update Response:", resp.json())
    except:
        print("Update Response:", resp.text)


def test_me(jwt):
    url = f"{API_URL}/me"
    print("\n>>> Checking /me info ...")

    resp = requests.get(url, cookies={"JWT": jwt})

    print("Me Status:", resp.status_code)

    try:
        print("Me Response:", resp.json())
    except:
        print("Me Response:", resp.text)


if __name__ == "__main__":
    # 1. 登录拿 JWT
    jwt = login_and_get_jwt()
    if not jwt:
        exit(1)

    # 2. 上传头像
    update_avatar(jwt)

    # 3. 重新取 /me 验证结果
    test_me(jwt)
