package prayer

import (
	"fmt"
	"math/rand"
	"time"
)

type PrayerNotification struct {
	Name      string
	Time      string // HH:MM format
	Message   string
	Keutamaan []string
}

type Engine struct {
	notifications []PrayerNotification
}

func NewEngine() *Engine {
	return &Engine{
		notifications: []PrayerNotification{
			{
				Name:    "Shalat Dhuha",
				Time:    "09:00",
				Message: "Selamat pagi, kaum Muslimin! ☀️\nSaatnya menunaikan Shalat Dhuha.",
				Keutamaan: []string{
					"Shalat Dhuha dapat menghapus dosa-dosa sebagaimana menghapus karat pada besi. (HR. Muslim)",
					"Allah SWT akan membangunkan baginya istana di surga. (HR. Tirmidzi)",
					"Allah SWT memberikan rezeki yang lapang (cukup) bagi yang mengerjakannya. (HR. Tirmidzi)",
					"Shalat Dhuha 2 rakaat dapat menuliskan 360 kebaikan dan menghapus 360 keburukan. (HR. Tirmidzi)",
					"Shalat Dhuha 4 rakaat dapat menuliskan 1.200 kebaikan dan menghapus 1.200 keburukan. (HR. Tirmidzi)",
				},
			},
			{
				Name:    "Shalat Tahajjud",
				Time:    "03:00",
				Message: "Selamat malam, kaum Muslimin! 🌙\nSaatnya menunaikan Shalat Tahajjud (shalat malam).",
				Keutamaan: []string{
					"Allah SWT turun ke langit dunia pada sepertiga malam yang akhir dan berfirman: 'Siapa yang berdoa kepada-Ku, Aku akan mengabulkannya.' (HR. Bukhari & Muslim)",
					"Shalat Tahajjud dapat menebus dosa-dosa dan menghapus keburukan. (HR. Tirmidzi)",
					"Allah SWT akan membangunkan baginya istana yang mulia di surga. (HR. Tirmidzi)",
					"Shalat Tahajjud adalah tanda kedekatan hamba dengan Allah SWT. (HR. Tirmidzi)",
					"Allah SWT akan mengangkat derajat orang yang mengerjakan Shalat Tahajjud. (HR. Tirmidzi)",
				},
			},
			{
				Name:    "Shalat Qabliyyah Subuh",
				Time:    "04:15",
				Message: "Selamat subuh, kaum Muslimin! 🌄\nSaatnya menunaikan Shalat Qabliyyah Subuh (shalat sebelum Subuh).",
				Keutamaan: []string{
					"Shalat 2 rakaat sebelum Subuh lebih baik daripada dunia dan segala isinya. (HR. Muslim)",
					"Allah SWT memberikan keberkahan bagi yang mengerjakannya. (HR. Tirmidzi)",
					"Shalat Qabliyyah Subuh dapat menjaga iman dan mendekatkan diri kepada Allah SWT. (HR. Tirmidzi)",
				},
			},
		},
	}
}

func (e *Engine) GetNotificationByTime(hour, minute int) *PrayerNotification {
	currentTime := fmt.Sprintf("%02d:%02d", hour, minute)

	for _, n := range e.notifications {
		if n.Time == currentTime {
			return &n
		}
	}
	return nil
}

func (e *Engine) GetAllNotifications() []PrayerNotification {
	return e.notifications
}

func (e *Engine) GetRandomKeutamaan(notification *PrayerNotification) string {
	if notification == nil || len(notification.Keutamaan) == 0 {
		return ""
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return notification.Keutamaan[r.Intn(len(notification.Keutamaan))]
}
