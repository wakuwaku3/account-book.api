package env

import (
	"os"

	"github.com/wakuwaku3/account-book.api/src/1-application-business-rules/usecases"

	"github.com/joho/godotenv"
)

type (
	env struct {
		credentialsFilePath string
		secret              string
	}
)

// NewEnv は環境変数管理用のインスタンスを生成します
func NewEnv() usecases.Env {
	return &env{credentialsFilePath: ""}
}
func (env *env) Initialize() error {
	slice := make([]string, 0)
	if exists("./.user.env") {
		slice = append(slice, "./.user.env")
	}
	if exists("./.env") {
		slice = append(slice, "./.env")
	}
	err := godotenv.Load(slice...)
	env.credentialsFilePath = os.Getenv("CREDENTIALS_FILE_PATH")
	env.secret = os.Getenv("SECRET")
	return err
}
func (env *env) GetCredentialsFilePath() string {
	return env.credentialsFilePath
}
func (env *env) GetSecret() string {
	return env.secret
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
