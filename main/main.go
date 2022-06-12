package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/Syfaro/telegram-bot-api"
	gt "github.com/bas24/googletranslatefree"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const UrlExamples = "https://context.reverso.net/translation/"
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
				msg := tgbotapi.NewMessage(message.Chat.ID, text)
				telebot.Send(msg)
			default:
				textToTranslate := getWords(message.Text)
				doc := getPage(UrlExamples, "russian-english/"+textToTranslate)
				examplesWord := examples(doc.Find("#examples-content .example"))
				translatedText := getTranslate(textToTranslate, "ru", "en")
				msg := tgbotapi.NewMessage(message.Chat.ID, translatedText+examplesWord)

				telebot.Send(msg)
			}
		}

	}
}

func getTranslate(text, fromLang, toLang string) string {
	result, _ := gt.Translate(text, fromLang, toLang)
	return SimpleTranslate + DNL + result + DNL
}

func examples(s *goquery.Selection) string {
	text := ""
	if s.Size() == 0 {
		return "Прикладів немає:)"
	}
	for i := 0; i < s.Size() && i < 3; i++ {
		example := s.Eq(i).Children().Find(".text")

		text += NL + ExampleNumber + strconv.Itoa(i+1) + DNL +
			strings.TrimSpace(example.Eq(0).Text()) +
			NL + Brackets + NL
		text += strings.TrimSpace(example.Eq(1).Text()) + NL
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
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", "Mozialla -1.0")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	return doc
}
