package domain

import "time"

type Article struct {
	Id      int64
	Author  Author
	Topic   string
	Status  ArticleStatus
	Content string
	Ctime   time.Time
	Utime   time.Time
}

func (a Article) Abstract() string {
	var abstractLimit = 140
	// rune counts non-ascii chars as 1, bytes won't
	str := []rune(a.Content)
	if len(str) > abstractLimit {
		str = str[:abstractLimit]
	}
	return string(str)
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

type PublishedArticle Article

func (a PublishedArticle) Abstract() string {
	var abstractLimit = 140
	// rune counts non-ascii chars as 1, bytes won't
	str := []rune(a.Content)
	if len(str) > abstractLimit {
		str = str[:abstractLimit]
	}
	return string(str)
}
