package main

import (
	"bufio"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Env struct {
	EmailFrom     string
	EmailPassword string
}

var EnvInstance = new(Env)

func loadEnvData() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("loadEnvData() --> Can't Load .env File")
	}

	EnvInstance.EmailFrom = os.Getenv("EMAIL_FROM")
	EnvInstance.EmailPassword = os.Getenv("EMAIL_PASSWORD")
}

func send(_to []string, _subject, _body string) {
	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s\r\n", _subject, _body))

	auth := smtp.PlainAuth(
		"",
		EnvInstance.EmailFrom,
		EnvInstance.EmailPassword,
		"smtp.gmail.com",
	)

	for _, to := range _to {
		err := smtp.SendMail(
			"smtp.gmail.com:587",
			auth,
			EnvInstance.EmailFrom,
			[]string{to},
			message,
		)
		if err != nil {
			log.Fatalf("send() --> Can't Send Email To [%s]", to)
		}
	}
}

func init() {
	loadEnvData()
}

func main() {
	var (
		emailList []string
		option    uint8
	)
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print(`
		[1] Send a Text To Everyone
		[2] Send a Special To Each Person
		)> `,
	)
	fmt.Scanln(&option)

	switch option {
	case 1:
		var (
			subject, body string
		)

		fmt.Print(`
			Enter Emails, (EXAMPLE) A@A.com,B@B.com,C@C.com
			)> `,
		)
		scanner.Scan()
		emailList = strings.Split(scanner.Text(), ",")

		fmt.Print(`
			Enter The Subject of Email
			)> `,
		)
		scanner.Scan()
		subject = scanner.Text()

		fmt.Print(`
			Enter The Body of Email
			)> `,
		)
		scanner.Scan()
		body = scanner.Text()

		send(emailList, string(subject), string(body))

	case 2:
		var subject, body string

		for {
			fmt.Print(`
				Enter Emails, (EXAMPLE) A@A.com,B@B.com,C@C.com)> `,
			)
			scanner.Scan()
			emailList = strings.Split(scanner.Text(), ",")

			fmt.Print(`
				Enter The Subject of Email)> `,
			)
			scanner.Scan()
			subject = scanner.Text()

			fmt.Print(`
				Enter The Body of Email)> `,
			)
			scanner.Scan()
			body = scanner.Text()

			send(emailList, subject, body)
		}
	default:
		log.Fatal("Invalid Option")
	}
}
