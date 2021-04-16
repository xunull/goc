package rand_fake

import "github.com/brianvoe/gofakeit/v6"

func FakeGameName() string {
	return gofakeit.Gamertag()
}
