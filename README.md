# serApi

站点服务和服务项目

## 开发安装

- go install github.com/zeromicro/go-zero/tools/goctl@latest
- goctl env check --install --verbose --force

## 项目描述

- gateway 网关项目
-  si 站点服务
- ag 服务

## 配置文件

- 所有项目公用config.yaml文件
- go run main.go  -config path\config\config.yaml

## 编译

- CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/site  cmd/site/main.go



echo "# goapiserver" >> README.md
git init
git add README.md
git commit -m "first commit"
git branch -M main
git remote add origin https://github.com/server.git
git push -u origin main


git clone https://bong05:q1@gitea.com/bong05/goserver.git

git remote set-url origin https://bong05:q1@gitea.com/bong05/goserver.git

更改提交地址
git remote set-url origin git@gitea.com:bong0/gose.git

bong05
gitea.com
id_ed25519_onemore.pub

bong558
gitea.com
id_ed25519.pub

正确用法
# 👉 配置 bong05 用户（使用 id_ed25519_onemore）
Host gitea-bong05
HostName gitea.com
User git
IdentityFile ~/.ssh/id_ed25519_onemore
IdentitiesOnly yes

# 👉 配置 bong558 用户（使用 id_ed25519）
Host gitea-bong558
HostName gitea.com
User git
IdentityFile ~/.ssh/id_ed25519
IdentitiesOnly yes

git remote set-url origin git@gitea-bong05:bong05/goserver.git
public




