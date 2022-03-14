package main

import (
	"fmt"
	"os"
	"text/template"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/google/uuid"
)

var qs = []*survey.Question{
	{
		Name: "profile_name",
		Prompt: &survey.Input{
			Message: "Enter a short name for the profile (e.g tfcreds)",
		},
		Validate: survey.Required,
	},
	{
		Name: "username_varname",
		Prompt: &survey.Input{
			Message: "Enter the Environment Variable Name for the Username (e.g. MY_USER)",
		},
		Validate: survey.Required,
	},
	{
		Name: "username_value",
		Prompt: &survey.Input{
			Message: "Enter the Username value (e.g. jdoe@bigcorp.com)",
		},
		Validate: survey.Required,
	},
	{
		Name: "password_varname",
		Prompt: &survey.Input{
			Message: "Enter the Environment Variable Name for the Password (e.g. MY_PASSWORD)",
		},
		Validate: survey.Required,
	},
	{
		Name: "password_value",
		Prompt: &survey.Password{
			Message: "Enter the Password (e.g. correcthorsebatterystaple)",
		},
		Validate: survey.Required,
	},
}

const tpl string = `#!/bin/bash

op_uuid="{{.OpUUID}}"

{{.UsernameEnvVar}}="$(op get item "$op_uuid" --fields username)"
{{.PasswordEnvVar}}="$(op get item "$op_uuid" --fields password)"
`

var scriptTemplate *template.Template = template.Must(template.New("script").Parse(tpl))

type PasswordServiceInput struct {
	UsernameEnvVar, PasswordEnvVar, UsernameValue, PasswordValue string
}
type PasswordServiceOutput struct {
	UsernameEnvVar, PasswordEnvVar, OpUUID string
}
type PasswordService interface {
	Save(input *PasswordServiceInput) (*PasswordServiceOutput, error)
}

type dummyPasswordService struct{}

func (d *dummyPasswordService) Save(input *PasswordServiceInput) (*PasswordServiceOutput, error) {
	return &PasswordServiceOutput{
		UsernameEnvVar: input.UsernameEnvVar,
		PasswordEnvVar: input.PasswordEnvVar,
		OpUUID:         uuid.New().String(),
	}, nil
}

func main() {
	answers := struct {
		Profile_Name     string
		Username_Varname string
		Username_Value   string
		Password_Varname string
		Password_Value   string
	}{}

	survey.Ask(qs, &answers)

	op := &dummyPasswordService{}
	res, err := op.Save(&PasswordServiceInput{
		UsernameEnvVar: answers.Username_Varname,
		UsernameValue:  answers.Username_Value,
		PasswordEnvVar: answers.Password_Varname,
		PasswordValue:  answers.Password_Value,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	if err := scriptTemplate.Execute(os.Stdout, res); err != nil {
		panic(err)
	}
}
