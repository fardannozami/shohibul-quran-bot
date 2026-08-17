package sunnah

import (
	"math/rand"
	"time"
)

type SleepSunnah struct {
	Title   string
	Message string
}

type SpecialDay struct {
	Name    string
	Message string
	Items   []string
}

type Engine struct {
	sleepSunnahs []SleepSunnah
	specialDays  []SpecialDay
}

func NewEngine() *Engine {
	return &Engine{
		sleepSunnahs: []SleepSunnah{
			{
				Title:   "Baca Surat Al-Mulk",
				Message: "Rasulullah SAW bersabda: 'Barangsiapa membaca Surat Al-Mulk setiap malam, Allah akan melindunginya dari siksa kubur.' (HR. Tirmidzi & Abu Daud)",
			},
			{
				Title:   "Baca Ayat Kursi",
				Message: "Rasulullah SAW bersabda: 'Barangsiapa membaca Ayat Kursi setelah shalat fardhu, tidak ada yang menghalanginya masuk surga kecuali mati.' (HR. An-Nasa'i & Ibnu Majah)",
			},
			{
				Title:   "Baca 3 Qul (Al-Ikhlas, Al-Falaq, An-Nas)",
				Message: "Rasulullah SAW bersabda: 'Bacalah Surat Al-Ikhlas, Al-Falaq, dan An-Nas setiap malam tiga kali. Itu akan melindungimu dari segala sesuatu.' (HR. Abu Daud & At-Tirmidzi)",
			},
			{
				Title:   "Baca Surat Al-Kahfi (Malam Jum'at)",
				Message: "Rasulullah SAW bersabda: 'Barangsiapa membaca Surat Al-Kahfi pada malam Jum'at, akan disinari cahaya antara dirinya dan Ka'bah.' (HR. Ad-Darimi)",
			},
			{
				Title:   "Baca Surat Ad-Dukhan",
				Message: "Rasulullah SAW bersabda: 'Barangsiapa membaca Surat Ad-Dukhan pada malam Jum'at, Allah akan memudahkan baginya di hari kiamat.' (HR. Ahmad)",
			},
			{
				Title:   "Baca Surat As-Sajdah & Al-Insan",
				Message: "Rasulullah SAW biasa membaca Surat As-Sajdah dan Al-Insan pada dua rakaat setelah Isya dan sebelum tidur. (HR. Muslim)",
			},
			{
				Title:   "Tidur dalam Keadaan Berwudhu",
				Message: "Rasulullah SAW bersabda: 'Tidaklah seorang Muslim berbaring untuk tidur dalam keadaan berwudhu, melainkan malaikat akan berada di sampingnya. Dan jika ia bangun di tengah malam, malaikat akan berdoa: Ya Allah, ampunilah hamba-Mu ini.' (HR. Bukhari & Muslim)",
			},
			{
				Title:   "Tidur Menghadap Kanan (Sebelah Timur)",
				Message: "Rasulullah SAW bersabda: 'Apabila seorang dari kalian hendak tidur, hendaklah ia berbaring di sisi kanannya, dan hendaklah membaca: Allahumma Rabbi Samawati Wa Ardhi, Rabbi Man 'Ala Syarika Lahu, Rabbi Khuliqni Walimka Ukhlik (Tafsir Ringkas). Kemudian setelah itu, hendaklah ia menghadap kiri.' (HR. Bukhari)",
			},
			{
				Title:   "Membaca Doa Sebelum Tidur",
				Message: "Rasulullah SAW bersabda: 'Apabila kalian hendak tidur, maka ucapkanlah: Bismika Allahumma amutu wa ahya. (Dengan nama-Mu ya Allah, aku mati dan aku hidup.)' (HR. Bukhari & Muslim)",
			},
			{
				Title:   "Meniup Tiga Kali dengan Doa dan Usapkan ke Seluruh Tubuh",
				Message: "Rasulullah SAW bersabda: 'Apabila kalian hendak tidur, maka tiuplah ke kedua tangan kalian dengan membaca: Bismikal Allahumma amutu wa ahya, sebanyak tiga kali. Kemudian usapkanlah ke seluruh tubuh kalian. Mulailah dengan mengusap kepala dan wajah kalian.' (HR. Bukhari & Muslim)",
			},
		},
		specialDays: []SpecialDay{
			{
				Name:    "Malam Jum'at",
				Message: "Malam Jum'at adalah malam yang penuh keberkahan. Berikut amalan sunnah yang dianjurkan:",
				Items: []string{
					"Membaca Surat Al-Kahfi. Rasulullah SAW bersabda: 'Barangsiapa membaca Surat Al-Kahfi pada malam Jum'at, akan disinari cahaya antara dirinya dan Ka'bah.' (HR. Ad-Darimi)",
					"Membaca Surat Yasin. Rasulullah SAW bersabda: 'Barangsiapa membaca Yasin di malam Jum'at, diampuni dosa-dosanya.' (HR. An-Nasa'i, Al-Hakim, & At-Tirmidzi)",
					"Banyak bershalawat kepada Nabi Muhammad SAW. Rasulullah SAW bersabda: 'Sebaik-baik hari kalian adalah Jum'at, maka perbanyaklah shalawat kepadaku pada hari itu.' (HR. An-Nasa'i)",
					"Membaca Surat Ad-Dukhan. 'Barangsiapa membaca Surat Ad-Dukhan pada malam Jum'at, akan memudahkan baginya di hari kiamat.' (HR. Ahmad)",
					"Membaca 3 Qul (Al-Ikhlas, Al-Falaq, An-Nas) sebanyak 3 kali setiap selesai shalat fardhu.",
					"Bersedekah. Rasulullah SAW bersabda: 'Sedekah pada hari Jum'ah itu lebih afdhal daripada sedekah pada hari-hari lainnya.' (HR. Al-Baihaqi)",
				},
			},
			{
				Name:    "Malam Sabtu",
				Message: "Malam Sabtu juga memiliki keutamaan tersendiri:",
				Items: []string{
					"Membaca Surat Al-Mulk. Rasulullah SAW bersabda: 'Barangsiapa membaca Surat Al-Mulk setiap malam, Allah akan melindunginya dari siksa kubur.' (HR. Tirmidzi & Abu Daud)",
					"Membaca Surat As-Sajdah dan Al-Insan sebagai amalan malam.",
					"Membaca 3 Qul (Al-Ikhlas, Al-Falaq, An-Nas) sebanyak 3 kali setiap selesai shalat fardhu.",
				},
			},
		},
	}
}

func (e *Engine) GetRandomSleepSunnah() *SleepSunnah {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	item := e.sleepSunnahs[r.Intn(len(e.sleepSunnahs))]
	return &item
}

func (e *Engine) GetSpecialDay(dayName string) *SpecialDay {
	for _, d := range e.specialDays {
		if d.Name == dayName {
			return &d
		}
	}
	return nil
}
