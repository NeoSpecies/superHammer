import socket
import struct
import json
import uuid

SOCKET_PATH = '/www/wwwroot/runtime/mainSocket.sock'  # 修改为实际的socket路径
VERSION = 0x0101  # 协议版本v1.1
MSG_TYPE = 0x01   # 同步消息类型

def send_request():
    # 创建Unix domain socket连接
    sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
    try:
        sock.connect(SOCKET_PATH)
        
        # 构造请求数据
        request = {
            "id": str(uuid.uuid4()),
            "method": "test_method",
            "params": {"key": "value"}
        }
        payload = json.dumps(request).encode('utf-8')
        
        # 封装协议头
        header = struct.pack('>HBI', VERSION, MSG_TYPE, len(payload))
        
        # 发送完整消息
        sock.sendall(header + payload)
        
        # 接收响应
        response_header = sock.recv(7)
        if len(response_header) != 7:
            raise Exception("Invalid response header")
            
        version, msg_type, length = struct.unpack('>HBI', response_header)
        response_data = sock.recv(length)
        
        # 解析响应
        response = json.loads(response_data.decode('utf-8'))
        print("Received response:", response)
        
    finally:
        sock.close()

if __name__ == '__main__':
    send_request()