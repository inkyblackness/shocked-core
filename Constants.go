package core

import (
	model "github.com/inkyblackness/shocked-model"
)

type localizedFiles struct {
	cybstrng string
}

var localized = [model.LanguageCount]localizedFiles{
	{
		cybstrng: "cybstrng.res"},
	{
		cybstrng: "frnstrng.res"},
	{
		cybstrng: "gerstrng.res"}}
