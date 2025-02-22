# 2425a大作业
基于区块链的智能博客系统, 基于PoW, 加密算法为sha256, 可渲染markdown类型的博客文档

## 技术栈
后端: go  

前端: python -flask  
     js -marked  

## 构建

```
cd docker 
./build_image.sh # 构建基础镜像pythonpkgs, 安装python相关包, 并以此作为后续python镜像的基础, 加快构建速度
cd ..
./run.fish # 运行
```


停止
```
./stop.sh
```