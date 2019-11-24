package mail

//邮件模板

type Verify interface {
	authVerify()
}

type MailTemplete struct {
}
