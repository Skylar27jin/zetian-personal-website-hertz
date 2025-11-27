# manual_test/test_upload_photo.py
import json
from http.cookies import SimpleCookie
from pathlib import Path

import requests

BASE_URL = "http://localhost:8888"

EMAIL = "skyjin0127@gmail.com"
PASSWORD = "Jzt20050127!"

IMAGE_PATH = (Path(__file__).resolve().parent / "me.jpg")
FILE_FIELD_NAME = "images"


def pretty_print_response(resp: requests.Response, title: str = ""):
    print("=" * 40)
    if title:
        print(title)
        print("-" * 40)
    print("Status:", resp.status_code)
    try:
        data = resp.json()
        print(json.dumps(data, indent=2, ensure_ascii=False))
    except Exception:
        print(resp.text)
    print("=" * 40)


def login_and_get_session() -> requests.Session:
    session = requests.Session()

    login_url = f"{BASE_URL}/login"
    payload = {"email": EMAIL, "password": PASSWORD}

    print(f"[*] 登录: POST {login_url}")
    resp = session.post(login_url, json=payload)

    pretty_print_response(resp, "Login Response")
    resp.raise_for_status()

    # 打印响应头里的 Set-Cookie
    set_cookie_raw = resp.headers.get("Set-Cookie", "")
    print("[*] Set-Cookie 原始值:", set_cookie_raw)

    # 打印 requests 当前认为的 cookies
    print("[*] session.cookies 初始:", session.cookies.get_dict())

    # 尝试从 Set-Cookie 里手动解析 JWT，并强制塞到 session 里
    cookie = SimpleCookie()
    cookie.load(set_cookie_raw)

    if "JWT" in cookie:
        jwt_value = cookie["JWT"].value
        print("[*] 从 Set-Cookie 解析到 JWT =", jwt_value)

        # 不指定 domain/path，让它变成“当前主机”的 cookie
        session.cookies.set("JWT", jwt_value)
    else:
        print("[!] Set-Cookie 里没有找到 JWT 字段，检查一下后端是否真的在 login 里写 cookie。")

    print("[*] session.cookies 之后:", session.cookies.get_dict())
    return session


def upload_photo(session: requests.Session):
    if not IMAGE_PATH.exists():
        raise FileNotFoundError(f"找不到图片文件: {IMAGE_PATH}")

    url = f"{BASE_URL}/post/media/upload"

    print(f"[*] 上传图片: POST {url}")
    print(f"[*] 使用文件: {IMAGE_PATH}")
    print("[*] 本次请求会携带的 cookies:", session.cookies.get_dict())

    with open(IMAGE_PATH, "rb") as f:
        files = {
            FILE_FIELD_NAME: (IMAGE_PATH.name, f, "image/png"),
        }
        resp = session.post(url, files=files)

    pretty_print_response(resp, "UploadPostMedia Response")


def main():
    print("=== 手动测试：登录 -> 上传图片到 /post/media/upload ===")
    session = login_and_get_session()
    upload_photo(session)


if __name__ == "__main__":
    main()
