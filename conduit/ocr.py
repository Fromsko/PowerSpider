import os
import argparse
from ddddocr import DdddOcr
import json
import base64


def ocr(image_data: bytes):
    json_msg = {'result': "", 'error': ''}
    try:
        result = DdddOcr(show_ad=False).classification(image_data)
        json_msg['result'] = result
    except Exception:
        json_msg['error'] = '识别失败'
    return json.dumps(json_msg)  # 返回JSON格式的数据


def main():
    import sys
    if len(args := sys.argv) == 2:
        image_data = args[1]
        # 将参数传递给ocr函数处理数据并获取结果
        result = ocr(base64.b64decode(image_data))
    else:
        sys.exit()
    print(result)


if __name__ == "__main__":
    main()
