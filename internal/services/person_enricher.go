package services

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Sypovik/effectiveMobileTest/internal/dto"
	"github.com/Sypovik/effectiveMobileTest/internal/models"
)

func toDTO(p *models.Person) *dto.PersonResponse {
	return &dto.PersonResponse{
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

	// Age
	if resp, err := client.Get("https://api.agify.io/?name=" + name); err == nil {
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&aResp)
		age = &aResp.Age
	}

	// Gender
	if resp, err := client.Get("https://api.genderize.io/?name=" + name); err == nil {
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&gResp)
		gender = &gResp.Gender
	}

	// Country
	if resp, err := client.Get("https://api.nationalize.io/?name=" + name); err == nil {
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&cResp)
		if len(cResp.Country) > 0 {
			country = &cResp.Country[0].CountryID
		}
	}

	return
}
