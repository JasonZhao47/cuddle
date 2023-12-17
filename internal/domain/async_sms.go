package domain

type AsyncSms struct {
	Id        int64
	TplId     string
	RetryMax  int
	PhoneNums []string
	Args      []string
}
