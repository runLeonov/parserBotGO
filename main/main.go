package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/Syfaro/telegram-bot-api"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const UrlExamples = "https://context.reverso.net/translation/"
const UrlTranslation = "https://www.reverso.net/text-translation/"
const TOKEN = "5436943485:AAG1Bnft74nScMPUXEJlCkKHSNsFPnBOtdE"

func main() {
	telebot, err := tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
		panic(err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := telebot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		message := update.Message
		if reflect.TypeOf(message.Text).Kind() == reflect.String && message.Text != "" {
			switch message.Text {
			case "/start":
				text := "Введіть слово або речення (поки тільки рос)"
				if message.Chat.UserName == "parodyOfLife" {
					msg := tgbotapi.NewMessage(message.Chat.ID, "Привіт Саша \n Надіюсь в тебе все добре")
					telebot.Send(msg)
				}
				println(message.Chat.UserName)
				msg := tgbotapi.NewMessage(message.Chat.ID, text)
				telebot.Send(msg)
			default:
				textToTranslate := getWords(message.Text)
				doc := getPage(UrlExamples, "russian-english/"+textToTranslate)

				examples := examples(doc.Find("#examples-content .example"))
				doc = getPage(UrlTranslation, "#sl=rus&tl=eng&text="+message.Text)
				//translatedText := dumpTranslation(doc.Find(".translation-inputs").Find("textarea").Eq(1))
				//fmt.Println(doc.Find("textarea").Html())
				msg := tgbotapi.NewMessage(message.Chat.ID, examples)

				telebot.Send(msg)
			}
		}

	}
}

func dumpTranslation(s *goquery.Selection) string {
	text := "Перекладений текст: \n"
	text += s.Find(".trg .ltr .text").Text()
	return text
}

func examples(s *goquery.Selection) string {
	text := ""
	if s.Size() == 0 {
		return "Прикладів немає:)"
	}
	for i := 0; i < s.Size() && i < 3; i++ {
		example := s.Eq(i).Children().Find(".text")
		text += "\nПриклад №" + strconv.Itoa(i+1) + "\n\n" + strings.TrimSpace(example.Eq(0).Text()) + "\n ---- > \n"
		text += strings.TrimSpace(example.Eq(1).Text()) + "\n"
	}
	return strings.TrimSpace(text)
}

func getWords(text string) string {
	strings.TrimSpace(text)
	strings.ReplaceAll(text, " ", "+")
	return text
}

func getPage(url, text string) *goquery.Document {
	client := &http.Client{}
	reqUrl := url + text
	req, err := http.NewRequest("GET", reqUrl, nil)
	fmt.Println(req.RequestURI)
	if err != nil {
		log.Fatalln(err)
	}
	//req.Header.Set("User-Agent", "Golang_Spider_Bot/3.0")
	req.Header.Set("User-Agent", "Mozialla -1.0")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	//fmt.Println(ioutil.ReadAll(req.Body))
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	return doc
}
