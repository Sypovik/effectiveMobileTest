package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
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

type enrichmentResult struct {
	Age     *int
	Gender  *string
	Country *string
}

func (r *enrichmentResult) setFromAPIResponse(apiResp interface{}, err error, field string) {
	if err != nil {
		log.Warn().Err(err).Msgf("API error for %s", field)
		return
	}

	switch v := apiResp.(type) {
	case *int:
		if v != nil {
			r.Age = v
		}
	case *string:
		if v != nil {
			switch field {
			case "gender":
				r.Gender = v
			case "country":
				r.Country = v
			}
		}
	}
}

func fetchAPI(ctx context.Context, url string, target interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}
	return nil
}

func getAge(ctx context.Context, name string) (*int, error) {
	type response struct{ Age int }
	var resp response

	url := "https://api.agify.io/?name=" + name
	if err := fetchAPI(ctx, url, &resp); err != nil {
		return nil, err
	}

	if resp.Age == 0 {
		return nil, nil
	}
	return &resp.Age, nil
}

func getGender(ctx context.Context, name string) (*string, error) {
	type response struct{ Gender string }
	var resp response

	url := "https://api.genderize.io/?name=" + name
	if err := fetchAPI(ctx, url, &resp); err != nil {
		return nil, err
	}

	if resp.Gender == "" {
		return nil, nil
	}
	return &resp.Gender, nil
}

func getCountry(ctx context.Context, name string) (*string, error) {
	type country struct {
		CountryID string  `json:"country_id"`
		Prob      float64 `json:"probability"`
	}
	type response struct{ Country []country }
	var resp response

	url := "https://api.nationalize.io/?name=" + name
	if err := fetchAPI(ctx, url, &resp); err != nil {
		return nil, err
	}

	if len(resp.Country) == 0 {
		return nil, nil
	}
	return &resp.Country[0].CountryID, nil
}

func enrichData(name string) (age *int, gender, country *string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result := &enrichmentResult{}
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		age, err := getAge(ctx, name)
		result.setFromAPIResponse(age, err, "age")
	}()

	go func() {
		defer wg.Done()
		gender, err := getGender(ctx, name)
		result.setFromAPIResponse(gender, err, "gender")
	}()

	go func() {
		defer wg.Done()
		country, err := getCountry(ctx, name)
		result.setFromAPIResponse(country, err, "country")
	}()

	wg.Wait()
	return result.Age, result.Gender, result.Country
}
