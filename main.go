package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// Широта и долгота Москвы
const latMos string = "55.7522"
const lonMos string = "37.6156"

// Широта и долгота  Петербурга
const latSp string = "59.8944"
const lonSp string = "30.2642"

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("Не найден .env файл")
	}

}

type jsonResp struct {
	Main struct {
		Feels_like float64 `json:"feels_like"`
		Temp       float64 `json:"temp`
	} `json:"main"`
}

// getTemp принимает географические широту и долготу и возвращает текущую температуру
// Также принимаетс ключ к api
func getTemp(lat, lon, key string) {
	cityCoor := "https://api.openweathermap.org/data/2.5/weather?lat=" + lat + "&lon=" + lon + "&appid=" + key + "&units=metric"
	res, err := http.Get(cityCoor)
	if err != nil {
		fmt.Println("Ошибка запроса")
	}
	bodyJson, err1 := io.ReadAll(res.Body)
	if err1 != nil {
		fmt.Println("Ошибка ответа")
	}

	var body jsonResp

	err2 := json.Unmarshal(bodyJson, &body)
	if err2 != nil {
		fmt.Println("Ошибка кодировки", err2)
	}

	fmt.Printf("В выбранном городе %v градусов по цельсию", body.Main.Temp)
}

func main() {
	// Здесь указывается ключ, получаемый при регистрации
	// Хранится в .env файле
	var apiKey, _ = os.LookupEnv("API_KEY")

	// Команда предполагающая ввод гео координат интересуемого города
	coords := &cobra.Command{
		Use:     "coords [широта] [долгота]",
		Version: "v.1.1",
		Args:    cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			getTemp(args[0], args[1], apiKey)
		},
	}

	// Команда для вызова погоды в Москве
	moscow := &cobra.Command{
		Use:     "mos",
		Version: "v.1.1",
		Run: func(cmd *cobra.Command, args []string) {
			getTemp(latMos, lonMos, apiKey)
		},
	}
	// Команда для вызова погоды в Петербурге
	spb := &cobra.Command{
		Use:     "spb",
		Version: "v.1.1",
		Run: func(cmd *cobra.Command, args []string) {
			getTemp(latSp, latSp, apiKey)
		},
	}
	mainComm := &cobra.Command{Use: "main"}
	mainComm.AddCommand(coords, moscow, spb)
	mainComm.Execute()
}
