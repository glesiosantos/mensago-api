package main

import (
	"mensago-api/internal/domain/campaign"
	"time"

	"github.com/go-playground/validator/v10"
)

func main() {
	contacts := []campaign.Contact{{Email: "teste@gmail.com"}}
	campaign := campaign.Campaign{
		Id:        "01",
		Name:      "Teste de Campanha",
		Content:   "Criando uma nova campanha",
		Contacts:  contacts,
		CreatedAt: time.Now(),
	}
	validate := validator.New()

	err := validate.Struct(campaign)

	if err == nil {
		println("Nenhum erro")
	} else {
		validatorErros := err.(validator.ValidationErrors)

		for _, v := range validatorErros {
			println(v.Error())
		}
	}
}
