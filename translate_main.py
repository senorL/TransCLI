import hashlib
import random
import requests
import config

# 你的百度翻译API的APPID和密钥
APPID = config.APPID
SECRET_KEY = config.SECRET_KEY


def generate_sign(query, salt):
    """生成签名"""
    sign_str = APPID + query + salt + SECRET_KEY
    sign = hashlib.md5(sign_str.encode()).hexdigest()
    return sign


def translate(query):
    salt = str(random.randint(32768, 65536))
    sign = generate_sign(query, salt)

    # 构建请求参数
    params = {
        'q': query,
        'from': 'en',
        'to': 'zh',
        'appid': APPID,
        'salt': salt,
        'sign': sign
    }

    # 发送请求
    url = 'https://fanyi-api.baidu.com/api/trans/vip/translate'
    response = requests.get(url, params=params)

    # 解析响应
    result = response.json()
    if 'trans_result' in result:
        return result['trans_result'][0]['dst']
    else:
        return f"Error: {result.get('error_msg', 'Unknown error')}"


def main():

    print("Enter 'exit' to quit.")
    while True:
        query = input("Enter text to translate: ")
        if query.lower() == 'exit':
            break


        translation = translate(query)
        print(f"Translation: {translation}\n")


if __name__ == "__main__":
    main()