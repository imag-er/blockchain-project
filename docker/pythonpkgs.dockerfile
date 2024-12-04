# 使用官方 Python 运行时作为父镜像
FROM python:latest

# 安装 Flask 和 requests
RUN pip install Flask requests
