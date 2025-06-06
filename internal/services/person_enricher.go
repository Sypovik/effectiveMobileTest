package services

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Sypovik/effectiveMobileTest/internal/dto"
	"github.com/Sypovik/effectiveMobileTest/internal/models"
	"github.com/rs/zerolog/log"
)

// toDTO преобразует модель Person в DTO PersonResponse
// @param p модель Person для преобразования
// @return DTO PersonResponse
func toDTO(p *models.Person) *dto.PersonResponse {
	log.Info().Msgf("Преобразование модели Person в DTO для ID: %d", p.ID)

	response := &dto.PersonResponse{
		ID:         p.ID,
		Name:       p.Name,
		Surname:    p.Surname,
		Patronymic: p.Patronymic,
		Age:        p.Age,
		Gender:     p.Gender,
		Country:    p.Country,
		CreatedAt:  p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  p.UpdatedAt.Format(time.RFC3339),
	}

	log.Info().Msgf("Успешное преобразование модели Person в DTO для ID: %d", p.ID)
	return response
}

func enrichData(name string) (age *int, gender, country *string) {
	type ageResp struct {
		Age int `json:"age"`
	}
	type genderResp struct {
		Gender string `json:"gender"`
	}
	type countryResp struct {
		Country []struct {
			CountryID string  `json:"country_id"`
			Prob      float64 `json:"probability"`
		} `json:"country"`
	}

	var (
		aResp ageResp
		gResp genderResp
		cResp countryResp
	)

	client := &http.Client{Timeout: 3 * time.Second}

	// Helper function to make requests, read body, and log
	fetchAndDecode := func(url, logPrefix string, target interface{}) ([]byte, error) {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Error().Err(err).Msgf("%s: Ошибка создания запроса", logPrefix)
			return nil, err
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Error().Err(err).Msgf("%s: Ошибка выполнения запроса", logPrefix)
			return nil, err
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Msgf("%s: Ошибка чтения тела ответа", logPrefix)
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			log.Error().
				Str("name", name).
				Str("url", url).
				Int("status", resp.StatusCode).
				Bytes("response_body", bodyBytes).
				Msgf("%s: Получен неожиданный статус код", logPrefix)
			return bodyBytes, nil // Return body even on non-200 for debugging
		}

		if err := json.Unmarshal(bodyBytes, target); err != nil {
			log.Error().
				Err(err).
				Str("name", name).
				Str("url", url).
				Bytes("raw_response", bodyBytes).
				Msgf("%s: Ошибка декодирования JSON", logPrefix)
			return bodyBytes, err
		}
		return bodyBytes, nil
	}

	// Age
	ageURL := "https://api.agify.io/?name=" + name
	if body, err := fetchAndDecode(ageURL, "Age API", &aResp); err == nil {
		if aResp.Age != 0 { // Check if age was actually enriched
			age = &aResp.Age
			log.Info().
				Str("name", name).
				Int("age", *age).
				Bytes("response_body", body).
				Interface("decoded_data", aResp). // Log the decoded struct
				Msg("Age API: Успешное обогащение возраста")
		} else {
			log.Warn().
				Str("name", name).
				Bytes("response_body", body).
				Interface("decoded_data", aResp).
				Msg("Age API: Возраст не получен или равен 0")
		}
	}

	// Gender
	genderURL := "https://api.genderize.io/?name=" + name
	if body, err := fetchAndDecode(genderURL, "Gender API", &gResp); err == nil {
		if gResp.Gender != "" { // Check if gender was actually enriched
			gender = &gResp.Gender
			log.Info().
				Str("name", name).
				Str("gender", *gender).
				Bytes("response_body", body).
				Interface("decoded_data", gResp). // Log the decoded struct
				Msg("Gender API: Успешное обогащение пола")
		} else {
			log.Warn().
				Str("name", name).
				Bytes("response_body", body).
				Interface("decoded_data", gResp).
				Msg("Gender API: Пол не получен или пустой")
		}
	}

	// Country
	countryURL := "https://api.nationalize.io/?name=" + name
	if body, err := fetchAndDecode(countryURL, "Country API", &cResp); err == nil {
		if len(cResp.Country) > 0 { // Check if country was actually enriched
			country = &cResp.Country[0].CountryID
			log.Info().
				Str("name", name).
				Str("country", *country).
				Bytes("response_body", body).
				Interface("decoded_data", cResp). // Log the decoded struct
				Msg("Country API: Успешное обогащение национальности")
		} else {
			log.Warn().
				Str("name", name).
				Bytes("response_body", body).
				Interface("decoded_data", cResp).
				Msg("Country API: Национальность не получена или список стран пуст")
		}
	}

	return
}
