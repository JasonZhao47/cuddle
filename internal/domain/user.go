package domain

import "time"

type User struct {
	// 可以用int64
	// 敏感信息PII
	Id       int64
	Nickname string
	Password string
	Email    string

	Birthday time.Time
	AboutMe  string

	// UTC
	Ctime time.Time
}
