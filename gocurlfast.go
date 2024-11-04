//go:build gocurlfast
// +build gocurlfast

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	// Определяем флаги
	followRedirects := flag.Bool("L", false, "Следовать по редиректам")
	onlyHeaders := flag.Bool("I", false, "Запрашивать только заголовки")
	outputFile := flag.String("o", "", "Файл для записи тела ответа")
	showHTTPCode := flag.Bool("w", false, "Показать только HTTP-код")
	silent := flag.Bool("s", false, "Тихий режим")
	flag.Parse()

	// Получаем URL из аргументов командной строки
	if len(flag.Args()) < 1 {
		log.Fatal("URL не указан")
	}
	siteURL := flag.Arg(0)

	// Парсим URL
	_, err := url.ParseRequestURI(siteURL)
	if err != nil {
		log.Fatalf("Некорректный URL: %v", err)
	}

	// Оптимизированный HTTP-клиент с тайм-аутом
	client := &http.Client{
		Timeout: 5 * time.Second, // Максимальное время ожидания ответа
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     30 * time.Second,
			DisableCompression:  true,            // Отключение сжатия может иногда ускорить обработку
			TLSHandshakeTimeout: 5 * time.Second, // Тайм-аут для TLS-рукопожатия
		},
	}
	if !*followRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	// Создаем запрос
	req, err := http.NewRequest("GET", siteURL, nil)
	if err != nil {
		log.Fatalf("Ошибка при создании запроса: %v", err)
	}
	if *onlyHeaders {
		req.Method = "HEAD"
	}

	// Выполняем запрос
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	// Если нужен только HTTP-код
	if *showHTTPCode {
		fmt.Println(resp.StatusCode)
		return
	}

	// Если тихий режим, отключаем вывод ошибок
	if *silent {
		log.SetOutput(io.Discard)
	}

	// Обрабатываем вывод
	var writer io.Writer = os.Stdout
	if *outputFile != "" {
		file, err := os.Create(*outputFile)
		if err != nil {
			log.Fatalf("Ошибка при создании файла: %v", err)
		}
		defer file.Close()
		writer = file
	}

	if *onlyHeaders {
		for key, values := range resp.Header {
			for _, value := range values {
				fmt.Fprintf(writer, "%s: %s\n", key, value)
			}
		}
	} else {
		// Увеличенный буфер для оптимизации копирования тела ответа
		buffer := make([]byte, 32*1024) // 32 KB
		_, err = io.CopyBuffer(writer, resp.Body, buffer)
		if err != nil {
			log.Fatalf("Ошибка при записи тела ответа: %v", err)
		}
	}
}
