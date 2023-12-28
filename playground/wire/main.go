package wire

import "fmt"

// 不使用wireinject标签
// build时候不会进来

func UseRepository() {
	repo := InitUserRepository()
	fmt.Println(repo)
}
