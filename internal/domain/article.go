package domain

import "time"

type Article struct {
	Id      int64
	Author  Author
	Topic   string
	Status  ArticleStatus
	Content string
	CTime   time.Time
	UTime   time.Time
}

const (
	ArticleStatusUnknown = iota
	ArticleStatusUnpublished
	ArticleStatusPublished
	ArticleStatusPrivate
)

type ArticleStatus uint8

type Author struct {
	Id   int64
	Name string
}

func (s ArticleStatus) ToUint8() uint8 {
	return uint8(s)
}
