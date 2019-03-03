# account-book.api

## Required

go(v1.11 以降)

## Command

#### Run

```sh
go run ./main.go
```

#### Build

```sh
go build -o dist/main ./main.go
```

#### Release

```sh
# 手動デプロイの場合
gcloud app deploy

# CIでデプロイする場合
git tag 1.0.0
git push origin --tags
```
