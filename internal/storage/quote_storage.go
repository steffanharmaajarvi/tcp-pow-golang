package storage

import "math/rand"

var Quotes = []string{
	"I like to listen. I have learned a great deal from listening carefully. Most people never listen.",
	"I think, that if the world were a bit more like ComicCon, it would be a better place",
	"Quit being so hard on yourself. We are what we are; we love what we love. We don't need to justify it to anyone... not even to ourselves.",
	"Voice is not just the sound that comes from your throat, but the feelings that come from your words.",
}

func GetRandomQuote() string {
	var limit = len(Quotes) - 1

	return Quotes[rand.Intn(limit)]
}
