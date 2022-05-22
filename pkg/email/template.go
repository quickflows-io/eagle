package email

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"time"

	"github.com/go-eagle/eagle/pkg/log"
)

// NewActivationEmail Send activation email
func NewActivationEmail(username, activateURL string) (subject string, body string) {
	return "Account activation link", "Hi, " + username + "<br>Please activate your account： <a href = '" + activateURL + "'>" + activateURL + "</a>"
}

// ActiveUserMailData 激活用户模板数据
type ActiveUserMailData struct {
	HomeURL       string `json:"home_url"`
	WebsiteName   string `json:"website_name"`
	WebsiteDomain string `json:"website_domain"`
	ActivateURL   string `json:"activate_url"`
	Year          int    `json:"year"`
}

// NewActivationHTMLEmail send activation email html
func NewActivationHTMLEmail(username, activateURL string) (subject string, body string) {
	mailData := ActiveUserMailData{
		//HomeURL:       conf.Conf.Web.Domain,
		//WebsiteName:   conf.Conf.Web.Name,
		//WebsiteDomain: conf.Conf.Web.Domain,
		ActivateURL: activateURL,
		Year:        time.Now().Year(),
	}
	mailTplContent := getEmailHTMLContent("./templates/active-mail.html", mailData)
	return "Account activation link", mailTplContent
}

// ResetPasswordMailData Activate user template data
type ResetPasswordMailData struct {
	HomeURL       string `json:"home_url"`
	WebsiteName   string `json:"website_name"`
	WebsiteDomain string `json:"website_domain"`
	ResetURL      string `json:"reset_url"`
	Year          int    `json:"year"`
}

// NewResetPasswordEmail Send reset password email
func NewResetPasswordEmail(username, resetURL string) (subject string, body string) {
	return "reset Password", "Hi, " + username + "<br>Your reset link is： <a href = '" + resetURL + "'>" + resetURL + "</a>"
}

// NewResetPasswordHTMLEmail send reset password email html
func NewResetPasswordHTMLEmail(username, resetURL string) (subject string, body string) {
	mailData := ResetPasswordMailData{
		//HomeURL:       conf.Conf.Web.Domain,
		//WebsiteName:   conf.Conf.Web.Name,
		//WebsiteDomain: conf.Conf.Web.Domain,
		ResetURL: resetURL,
		Year:     time.Now().Year(),
	}
	mailTplContent := getEmailHTMLContent("./templates/reset-mail.html", mailData)
	return "reset Password", mailTplContent
}

// getEmailHTMLContent Get email template
func getEmailHTMLContent(tplPath string, mailData interface{}) string {
	b, err := ioutil.ReadFile(tplPath)
	if err != nil {
		log.Warnf("[util.email] read file err: %v", err)
		return ""
	}
	mailTpl := string(b)
	tpl, err := template.New("email tpl").Parse(mailTpl)
	if err != nil {
		log.Warnf("[util.email] template new err: %v", err)
		return ""
	}
	buffer := new(bytes.Buffer)
	err = tpl.Execute(buffer, mailData)
	if err != nil {
		fmt.Println("exec err", err)
		log.Warnf("[util.email] execute template err: %v", err)
	}
	return buffer.String()
}
