package util

import (
	"math/rand"
	"time"

	"github.com/easc01/mindo-server/internal/models"
)

var AllColors = []models.Color{
	models.ColorRed, models.ColorBlue, models.ColorGreen, models.ColorYellow,
	models.ColorOrange, models.ColorPurple, models.ColorPink, models.ColorBrown,
	models.ColorTeal,
}

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func GetRandomColor() models.Color {
	return AllColors[rnd.Intn(len(AllColors))]
}
