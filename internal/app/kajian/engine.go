package kajian

import (
	"fmt"
	"math/rand"
	"time"
)

type KajianItem struct {
	Title   string
	Excerpt string
	Link    string
}

type Engine struct {
	fallback []KajianItem
}

func NewEngine() *Engine {
	return &Engine{
		fallback: []KajianItem{
			// Parenting Islami
			{
				Title:   "Parenting Islami: Mendidik Anak dengan Cinta dan Akhlak Mulia",
				Excerpt: "Rasulullah SAW bersabda: 'Sebaik-baik kalian adalah yang paling baik terhadap anak-anaknya.' (HR. Tirmidzi). Mendidik anak dalam Islam bukan sekadar mengajarkan ibadah, tapi juga menanamkan akhlak mulia, kejujuran, dan kasih sayang. Mulailah dengan memberikan teladan, karena anak belajar lebih banyak dari apa yang mereka lihat daripada apa yang mereka dengar.",
			},
			{
				Title:   "Parenting Islami: Membangun Komunikasi yang Baik dengan Anak",
				Excerpt: "Allah SWT berfirman: 'Dan rendahkanlah dirimu terhadap mereka dengan penuh kasih sayang.' (QS. Al-Isra: 24). Komunikasi yang baik antara orang tua dan anak dimulai dengan mendengarkan. Berikan waktu berkualitas, dengarkan cerita mereka tanpa menyela, dan responi dengan penuh pengertian. Anak yang merasa didengar akan tumbuh menjadi pribadi yang percaya diri dan terbuka.",
			},
			{
				Title:   "Parenting Islami: Mengajarkan Shalat sejak Kecil",
				Excerpt: "Rasulullah SAW bersabda: 'Perintahkan anak-anakmu untuk shalat ketika mereka berusia tujuh tahun.' (HR. Abu Daud). Shalat adalah tiang agama, dan mengajarkannya sejak dini merupakan investasi akhlak terbesar. Mulai dengan mengajak anak ke masjid, memberikan contoh, dan menjadikan shalat sebagai kegiatan yang menyenangkan, bukan kewajiban yang memberatkan.",
			},
			{
				Title:   "Parenting Islami: Menanamkan Kecintaan kepada Al-Qur'an",
				Excerpt: "Rasulullah SAW bersabda: 'Sebaik-baik kalian adalah yang mempelajari Al-Qur'an dan mengajarkannya.' (HR. Bukhari). Mulailah mengajak anak membaca Al-Qur'an sejak kecil dengan cara yang menyenangkan. Bacakan surat pendek sebelum tidur, ajak mengaji ke majelis ilmu, dan berikan apresiasi atas setiap kemajuan mereka. Kecintaan kepada Al-Qur'an yang tertanam sejak kecil akan menjadi penuntun hidup mereka.",
			},
			{
				Title:   "Parenting Islami: Mendidik Anak tentang Pentingnya Berbagi",
				Excerpt: "Rasulullah SAW bersabda: 'Berinfaklah, dan jika tidak mampu, maka bantu orang yang membutuhkan.' (HR. Bukhari & Muslim). Ajarkan anak untuk berbagi sejak kecil, baik dengan teman, tetangga, maupun yang membutuhkan. Mulai dari hal sederhana seperti berbagi makanan atau mainan, hingga mengajak mereka ke panti asuhan. Anak yang terbiasa berbagi akan tumbuh menjadi pribadi yang dermawan dan peduli sosial.",
			},
			// Ibadah
			{
				Title:   "Keutamaan Shalat Dhuha dan Waktu Pelaksanaannya",
				Excerpt: "Shalat Dhuha dikerjakan pada waktu dhuha (ketika matahari mulai naik), sekitar pukul 08:00-10:30. Rasulullah SAW bersabda: 'Shalat orang-orang yang awwabin (kembali kepada Allah) adalah shalat dhuha.' (HR. Muslim). Shalat Dhuha 2 rakaat dapat menghapus dosa-dosa sebagaimana menghapus karat pada besi. Semakin banyak rakaat, semakin besar pahalanya.",
			},
			{
				Title:   "Shalat Sunnah Rawatib: Pahala dan Cara Pelaksanaannya",
				Excerpt: "Shalat rawatib adalah shalat sunnah yang dikerjakan sebelum dan sesudah shalat wajib. Rasulullah SAW bersabda: 'Barangsiapa mengerjakan dua rakaat sunnah sebelum subuh, maka tidak akan dimasukkan ke dalam neraka.' (HR. Tirmidzi). Ada 12 rakaat shalat rawatib: 4 rakaat sebelum Zuhur dan 2 sesudahnya, 2 rakaat sesudah Maghrib, 2 rakaat sesudah Isya, dan 2 rakaat sebelum Subuh.",
			},
			{
				Title:   "Tahajjud: Shalat Malam yang Penuh Keberkahan",
				Excerpt: "Rasulullah SAW bersabda: 'Shalat malam itu shalat sunnah yang paling utama setelah shalat fardhu.' (HR. Muslim). Waktu terbaik tahajjud adalah sepertiga malam yang akhir (sekitar pukul 02:00-03:00). Allah SWT turun ke langit dunia pada waktu itu dan berfirman: 'Siapa yang berdoa kepada-Ku, Aku akan mengabulkannya.' Tahajjud dapat menebus dosa, mengabulkan doa, dan mengangkat derajat.",
			},
			{
				Title:   "Puasa Sunnah: Keutamaan Puasa Senin-Kamis",
				Excerpt: "Rasulullah SAW bersabda: 'Setiap amal kebaikan anak Adam dilipatgandakan pahalanya sepuluh hingga tujuh ratus kali lipat. Allah berfirman: 'Kecuali puasa, karena puasa itu milik-Ku dan Aku yang membalasnya.' (HR. Bukhari & Muslim). Puasa Senin-Kamis adalah puasa sunnah yang paling dianjurkan, karena pada hari Senin Nabi SAW dilahirkan dan pada hari Kamis beliau diangkat (menerima wahyu pertama).",
			},
			{
				Title:   "Sedekah: Investasi Akhirat yang Tidak Akan Rugi",
				Excerpt: "Rasulullah SAW bersabda: 'Sedekah tidak akan mengurangi harta.' (HR. Muslim). Sedekah bukan sekadar memberikan uang, tapi bisa dalam bentuk senyuman, ilmu yang bermanfaat, atau membantu orang lain. Allah SWT berfirman: 'Perumpamaan (nafkah yang dikeluarkan oleh) orang-orang yang menafkahkan hartanya di jalan Allah adalah serupa dengan sebutir benih yang menumbuhkan tujuh bulir, pada tiap-tiap bulir seratus biji.' (QS. Al-Baqarah: 261)",
			},
			{
				Title:   "Istighfar: Cara Membersihkan Haji dan Menambah Rezeki",
				Excerpt: "Allah SWT berfirman: 'Mohonlah ampun kepada Tuhanmu, sesungguhnya Dia Maha Pengampun. Niscaya Dia akan mengirimkan hujan yang lebat (bagimu) dan membanyakkan harta dan anak-anakmu.' (QS. Nuh: 10-12). Istighfar adalah senjata mukmin untuk menghapus dosa dan membuka pintu rezeki. Minimal 100 kali sehari, atau setelah shalat fardhu. Rasulullah SAW sendiri melakukan istighfar lebih dari 100 kali sehari.",
			},
		},
	}
}

func (e *Engine) GetRandomKajian() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	item := e.fallback[r.Intn(len(e.fallback))]
	return e.formatKajian(&item)
}

func (e *Engine) formatKajian(item *KajianItem) string {
	if item == nil {
		return ""
	}

	msg := fmt.Sprintf("📚 *Kajian Hari Ini*\n\n*%s*\n\n%s", item.Title, item.Excerpt)
	if item.Link != "" {
		msg += fmt.Sprintf("\n\n🔗 Baca selengkapnya: %s", item.Link)
	}
	return msg
}
