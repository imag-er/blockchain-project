# 预先安装python依赖 避免每次构建都要重新安装
FROM pythonpkgs:latest

# 工作目录
WORKDIR /app/python

# 复制文件到工作目录
COPY python/app.py .
COPY python/templates templates/
COPY python/static static/

# 暴露端口
EXPOSE 80

# 启动应用
CMD ["python", "app.py"]

