package motivation

import (
	"math/rand"
	"time"
)

type Engine struct {
	quotes []string
}

func NewEngine() *Engine {
	return &Engine{
		quotes: []string{
			"“Sebaik-baik kalian adalah yang mempelajari Al-Qur'an dan mengajarkannya.” (HR. Bukhari)",
			"“Bacalah Al-Qur'an, karena sesungguhnya ia akan datang pada hari kiamat sebagai pemberi syafaat bagi pembacanya.” (HR. Muslim)",
			"“Orang yang lancar membaca Al-Qur'an akan bersama para malaikat yang mulia dan senantiasa taat...” (HR. Bukhari & Muslim)",
			"“Barangsiapa membaca satu huruf dari Kitabullah, maka baginya satu kebaikan...” (HR. Tirmidzi)",
			"Allah berfirman: “Ingatlah, hanya dengan mengingat Allah-lah hati menjadi tenteram.” (QS. Ar-Ra'd: 28)",
			"Allah berfirman: “Dan sesungguhnya telah Kami mudahkan Al-Qur'an untuk pelajaran, maka adakah orang yang mengambil pelajaran?” (QS. Al-Qamar: 17)",
			"“Sesungguhnya Allah mengangkat derajat sebuah kaum dengan kitab ini (Al-Qur'an), dan merendahkan kaum yang lain dengannya.” (HR. Muslim)",
		},
	}
}

// GetRandomMotivation picks a random motivation quote from the predefined list
func (e *Engine) GetRandomMotivation() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Intn(len(e.quotes))
	return e.quotes[index]
}
