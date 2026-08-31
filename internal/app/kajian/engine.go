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
				Title:   "Parenting Islami: Menjadi Teladan bagi Anak",
				Excerpt: "Anak belajar bukan hanya dari nasihat, tetapi terutama dari apa yang mereka lihat setiap hari. Orang tua yang ingin anaknya jujur, sabar, rajin beribadah, dan santun perlu terlebih dahulu memberikan contoh. Keteladanan adalah salah satu cara terbaik menanamkan nilai Islam dalam kehidupan anak.",
			},
			{
				Title:   "Parenting Islami: Mengajarkan Anak Mengucapkan Salam",
				Excerpt: "Salam adalah doa dan bentuk kasih sayang kepada sesama Muslim. Ajarkan anak untuk membiasakan mengucapkan assalamu'alaikum ketika bertemu, masuk rumah, atau berkunjung. Kebiasaan sederhana ini dapat menanamkan kesopanan sekaligus mempererat hubungan dengan orang lain.",
			},
			{
				Title:   "Parenting Islami: Mengajarkan Anak Bersikap Jujur",
				Excerpt: "Kejujuran merupakan akhlak mulia yang sangat ditekankan dalam Islam. Ajarkan anak untuk berkata benar meskipun sedang melakukan kesalahan. Ketika anak jujur, jangan langsung memarahinya, tetapi berikan apresiasi agar ia belajar bahwa berkata benar adalah sesuatu yang aman dan mulia.",
			},
			{
				Title:   "Parenting Islami: Mengajarkan Anak Menghormati Orang Tua",
				Excerpt: "Allah SWT memerintahkan manusia untuk berbuat baik kepada kedua orang tua. Ajarkan anak berbicara dengan lembut, membantu pekerjaan sederhana di rumah, dan menghargai orang tua. Menghormati orang tua bukan hanya tentang kata-kata, tetapi juga tercermin dari sikap dan perilaku sehari-hari.",
			},
			{
				Title:   "Parenting Islami: Mengajarkan Anak Berbakti kepada Orang Tua",
				Excerpt: "Berbakti kepada orang tua merupakan amal yang memiliki kedudukan besar dalam Islam. Orang tua dapat mengenalkan makna birrul walidain melalui kebiasaan sederhana seperti membantu, mendoakan, dan berbicara dengan sopan. Pendidikan tentang berbakti sebaiknya dimulai sejak anak masih kecil.",
			},
			{
				Title:   "Parenting Islami: Mengajarkan Anak Bersyukur",
				Excerpt: "Ajarkan anak bahwa setiap nikmat berasal dari Allah SWT. Biasakan mengucapkan alhamdulillah setelah mendapatkan makanan, hadiah, kesehatan, atau pengalaman menyenangkan. Anak yang dibiasakan bersyukur akan belajar melihat nikmat dalam hal-hal sederhana.",
			},
			{
				Title:   "Parenting Islami: Mengajarkan Anak Meminta Maaf",
				Excerpt: "Meminta maaf bukan tanda kelemahan, tetapi tanda keberanian untuk mengakui kesalahan. Ajarkan anak untuk meminta maaf ketika menyakiti orang lain dan belajar memperbaiki kesalahannya. Orang tua juga dapat memberikan contoh dengan meminta maaf ketika melakukan kesalahan kepada anak.",
			},
			{
				Title:   "Parenting Islami: Mengajarkan Anak Memaafkan",
				Excerpt: "Islam mengajarkan umatnya untuk menjadi pribadi yang pemaaf. Ketika anak mengalami konflik dengan teman atau saudaranya, bantu mereka memahami kesalahan tanpa menanamkan dendam. Ajarkan bahwa memaafkan bukan berarti membenarkan kesalahan, tetapi memilih untuk tidak menyimpan kebencian.",
			},
			{
				Title:   "Parenting Islami: Membiasakan Anak Berdoa",
				Excerpt: "Doa mengajarkan anak untuk bergantung kepada Allah SWT dalam setiap keadaan. Biasakan anak berdoa sebelum makan, sebelum tidur, ketika keluar rumah, dan ketika membutuhkan pertolongan. Dengan kebiasaan ini, anak belajar bahwa Allah selalu menjadi tempat meminta dan bersandar.",
			},
			{
				Title:   "Parenting Islami: Mengajarkan Anak Adab Sebelum Makan",
				Excerpt: "Ajarkan anak untuk mencuci tangan, membaca basmalah, menggunakan tangan kanan, makan makanan yang ada di hadapannya, dan menghindari berlebihan. Adab makan yang diajarkan sejak kecil akan menjadi kebiasaan yang terbawa hingga dewasa.",
			},
			{
				Title:   "Parenting Islami: Mengajarkan Anak Menjaga Kebersihan",
				Excerpt: "Kebersihan merupakan bagian penting dalam kehidupan seorang Muslim. Ajarkan anak merapikan mainan, membuang sampah pada tempatnya, menjaga kebersihan tubuh, dan merawat lingkungan. Kebiasaan menjaga kebersihan dapat menjadi bagian dari pendidikan akhlak sejak dini.",
			},
			{
				Title:   "Parenting Islami: Mengajarkan Anak Menyayangi Hewan",
				Excerpt: "Islam mengajarkan kasih sayang kepada seluruh makhluk Allah. Orang tua dapat mengajarkan anak untuk tidak menyakiti hewan, memberikan makanan ketika diperlukan, dan memperlakukannya dengan lembut. Dari sini anak belajar tentang rahmat, kasih sayang, dan tanggung jawab.",
			},
			{
				Title:   "Parenting Islami: Mengajarkan Anak Berbagi dengan Saudara",
				Excerpt: "Anak perlu belajar bahwa tidak semua hal harus dimiliki sendiri. Biasakan mereka berbagi makanan, mainan, atau kesempatan dengan saudara dan teman. Kebiasaan berbagi dapat membantu membentuk karakter dermawan dan mengurangi sifat egois.",
			},
			{
				Title:   "Parenting Islami: Mengajarkan Anak Menjaga Amanah",
				Excerpt: "Amanah adalah salah satu akhlak mulia yang harus ditanamkan sejak kecil. Mulailah dengan memberikan tanggung jawab sederhana seperti menjaga barang miliknya, menyelesaikan tugas, atau mengembalikan sesuatu yang dipinjam. Dari kebiasaan kecil tersebut anak belajar arti tanggung jawab.",
			},
			{
				Title:   "Parenting Islami: Mengajarkan Anak Tidak Sombong",
				Excerpt: "Ajarkan anak bahwa kepintaran, kecantikan, kekayaan, dan berbagai kelebihan adalah nikmat dari Allah SWT. Anak perlu belajar menghargai orang lain dan tidak merasa lebih tinggi karena memiliki sesuatu yang lebih. Rendah hati membuat seseorang disukai dan dihormati.",
			},

			// Akhlak & Kehidupan Muslim
			{
				Title:   "Akhlak Islami: Menjaga Lisan dari Perkataan Buruk",
				Excerpt: "Rasulullah SAW mengajarkan bahwa seorang Muslim hendaknya berkata yang baik atau diam. Menjaga lisan berarti menghindari kebohongan, hinaan, fitnah, dan perkataan yang menyakiti orang lain. Sebelum berbicara, biasakan bertanya kepada diri sendiri apakah perkataan tersebut membawa kebaikan.",
			},
			{
				Title:   "Akhlak Islami: Menghindari Ghibah",
				Excerpt: "Allah SWT melarang ghibah dan menggambarkannya sebagai perbuatan yang sangat buruk dalam QS. Al-Hujurat ayat 12. Ghibah berarti membicarakan keburukan seseorang yang tidak ada di hadapan kita. Menjaga lisan dari ghibah merupakan bentuk menjaga kehormatan sesama Muslim.",
			},
			{
				Title:   "Akhlak Islami: Belajar Bersabar dalam Menghadapi Ujian",
				Excerpt: "Kesabaran adalah akhlak yang sangat mulia dalam Islam. Ketika menghadapi kesulitan, seorang Muslim dianjurkan untuk tetap berusaha, berdoa, dan tidak berputus asa dari rahmat Allah. Kesabaran bukan berarti pasrah tanpa usaha, tetapi tetap teguh sambil mencari jalan keluar.",
			},
			{
				Title:   "Akhlak Islami: Belajar Rendah Hati",
				Excerpt: "Rendah hati bukan berarti merendahkan diri, tetapi tidak merasa lebih tinggi daripada orang lain. Ingatlah bahwa segala kelebihan adalah karunia Allah SWT. Semakin banyak nikmat yang diterima, seharusnya semakin besar pula rasa syukur dan kerendahan hati.",
			},
			{
				Title:   "Akhlak Islami: Menjaga Janji",
				Excerpt: "Menepati janji merupakan bagian dari akhlak seorang Muslim. Jangan mudah berjanji jika belum yakin dapat menunaikannya. Ketika sudah berjanji, berusahalah untuk memenuhi apa yang telah disampaikan karena kepercayaan dibangun dari konsistensi.",
			},
			{
				Title:   "Akhlak Islami: Menghormati Tetangga",
				Excerpt: "Islam memberikan perhatian besar terhadap hubungan dengan tetangga. Jangan mengganggu kenyamanan mereka, bantu ketika membutuhkan, dan jaga hubungan dengan sikap yang baik. Lingkungan yang dipenuhi kepedulian akan menciptakan kehidupan yang lebih damai.",
			},
			{
				Title:   "Akhlak Islami: Membalas Keburukan dengan Kebaikan",
				Excerpt: "Membalas keburukan dengan kebaikan bukanlah sesuatu yang mudah, tetapi merupakan akhlak yang mulia. Ketika seseorang menyakiti kita, Islam mengajarkan agar kita tidak membalas dengan kezaliman. Berbuat baik dapat menjadi jalan untuk menjaga hati dari kebencian.",
			},
			{
				Title:   "Akhlak Islami: Menjaga Hati dari Iri dan Dengki",
				Excerpt: "Iri dan dengki dapat membuat hati sulit menikmati nikmat yang telah diberikan Allah. Ketika melihat keberhasilan orang lain, biasakan mendoakan keberkahan untuk mereka dan meminta Allah memberikan kebaikan kepada kita. Setiap orang memiliki rezeki dan perjalanan hidup masing-masing.",
			},

			// Ibadah
			{
				Title:   "Keutamaan Shalat Lima Waktu",
				Excerpt: "Shalat lima waktu merupakan kewajiban utama bagi setiap Muslim yang telah memenuhi syarat. Shalat menjadi sarana seorang hamba untuk mengingat Allah dan memohon pertolongan-Nya. Jagalah shalat meskipun sedang sibuk karena waktu terus berjalan dan kesempatan tidak selalu kembali.",
			},
			{
				Title:   "Keutamaan Shalat Berjamaah",
				Excerpt: "Shalat berjamaah memiliki keutamaan yang besar dibandingkan shalat sendirian. Selain mendapatkan pahala, berjamaah juga mengajarkan kedisiplinan, persaudaraan, dan kepedulian terhadap sesama Muslim. Luangkan waktu untuk menghadiri masjid ketika memiliki kesempatan.",
			},
			{
				Title:   "Membiasakan Shalat Tepat Waktu",
				Excerpt: "Salah satu bentuk menjaga shalat adalah berusaha menunaikannya di awal waktu. Jangan menjadikan kesibukan sebagai alasan untuk terus menunda. Ketika waktu shalat tiba, berhentilah sejenak dari aktivitas dunia dan penuhi panggilan untuk menghadap Allah SWT.",
			},
			{
				Title:   "Keutamaan Membaca Al-Qur'an Setiap Hari",
				Excerpt: "Al-Qur'an adalah petunjuk bagi kehidupan manusia. Membaca beberapa ayat setiap hari lebih baik daripada tidak membacanya sama sekali. Jadikan Al-Qur'an sebagai bagian dari rutinitas, misalnya setelah shalat Subuh, Maghrib, atau sebelum tidur.",
			},
			{
				Title:   "Tadabbur Al-Qur'an: Tidak Hanya Membaca",
				Excerpt: "Membaca Al-Qur'an merupakan ibadah, tetapi memahami dan merenungkan maknanya juga sangat penting. Luangkan waktu untuk membaca terjemahan dan memahami pesan dari ayat yang dibaca. Dengan tadabbur, Al-Qur'an tidak hanya berada di lisan, tetapi juga memengaruhi kehidupan.",
			},
			{
				Title:   "Dzikir Pagi dan Petang: Mengingat Allah Setiap Hari",
				Excerpt: "Dzikir pagi dan petang merupakan amalan yang dapat membantu seorang Muslim mengingat Allah di tengah kesibukan. Luangkan beberapa menit setelah Subuh dan menjelang Maghrib untuk membaca dzikir yang diajarkan dalam sunnah. Hati yang selalu mengingat Allah akan lebih tenang.",
			},
			{
				Title:   "Keutamaan Membaca Ayat Kursi",
				Excerpt: "Ayat Kursi terdapat dalam QS. Al-Baqarah ayat 255 dan merupakan salah satu ayat yang agung dalam Al-Qur'an. Banyak Muslim membiasakan membacanya setelah shalat wajib dan sebelum tidur sebagai bagian dari dzikir dan doa.",
			},
			{
				Title:   "Membaca Surat Al-Kahfi di Hari Jumat",
				Excerpt: "Hari Jumat merupakan hari yang istimewa bagi umat Islam. Salah satu amalan yang banyak dilakukan adalah membaca Surat Al-Kahfi. Jadikan hari Jumat sebagai momentum untuk memperbanyak membaca Al-Qur'an, bershalawat, berdoa, dan mengingat Allah.",
			},
			{
				Title:   "Memperbanyak Shalawat di Hari Jumat",
				Excerpt: "Hari Jumat merupakan waktu yang baik untuk memperbanyak shalawat kepada Nabi Muhammad SAW. Shalawat dapat menjadi bagian dari rutinitas seorang Muslim untuk mengingat Rasulullah dan menunjukkan kecintaan kepada beliau. Luangkan waktu khusus pada hari Jumat untuk memperbanyak shalawat.",
			},
			{
				Title:   "Keutamaan Shalat Witir",
				Excerpt: "Shalat Witir merupakan shalat sunnah yang sangat dianjurkan dan menjadi penutup shalat malam. Witir dapat dikerjakan setelah Isya hingga sebelum masuk waktu Subuh. Jika khawatir tidak dapat bangun malam, Witir dapat dikerjakan sebelum tidur.",
			},
			{
				Title:   "Shalat Istikharah: Memohon Petunjuk Allah",
				Excerpt: "Ketika menghadapi pilihan penting, seorang Muslim dapat melakukan shalat istikharah dan berdoa memohon pilihan terbaik kepada Allah. Istikharah bukan berarti berhenti berusaha, tetapi menyerahkan hasil kepada Allah setelah mempertimbangkan pilihan dan melakukan ikhtiar.",
			},
			{
				Title:   "Shalat Taubat: Kembali kepada Allah",
				Excerpt: "Setiap manusia dapat melakukan kesalahan. Islam membuka pintu taubat selama seseorang masih hidup. Selain meninggalkan dosa dan menyesalinya, seorang Muslim hendaknya bertekad untuk tidak mengulangi kesalahan serta memperbanyak amal kebaikan.",
			},
			{
				Title:   "Keutamaan Membaca Istighfar",
				Excerpt: "Istighfar merupakan bentuk pengakuan seorang hamba atas kelemahan dan kesalahannya di hadapan Allah. Perbanyaklah membaca Astaghfirullah dalam keseharian, bukan hanya ketika merasa berdosa. Istighfar juga menjadi pengingat agar hati selalu kembali kepada Allah.",
			},
			{
				Title:   "Doa Setelah Shalat: Momen untuk Mendekat kepada Allah",
				Excerpt: "Setelah menyelesaikan shalat, luangkan waktu untuk berdzikir dan berdoa. Jangan terburu-buru kembali kepada aktivitas dunia. Gunakan momen tersebut untuk memohon ampunan, kesehatan, rezeki yang halal, keteguhan iman, dan kebaikan bagi keluarga.",
			},

			// Sedekah & Rezeki
			{
				Title:   "Rezeki dalam Islam: Tidak Hanya Berupa Uang",
				Excerpt: "Rezeki tidak selalu berbentuk harta. Kesehatan, keluarga, waktu, ilmu, teman yang baik, dan kesempatan berbuat baik juga merupakan nikmat dari Allah. Dengan memahami hal ini, kita belajar untuk lebih banyak bersyukur daripada terus membandingkan diri dengan orang lain.",
			},
			{
				Title:   "Sedekah Secara Diam-Diam",
				Excerpt: "Tidak semua kebaikan perlu diketahui orang lain. Sedekah yang dilakukan secara sembunyi-sembunyi dapat membantu menjaga keikhlasan dan menjauhkan diri dari keinginan mendapatkan pujian. Yang terpenting bukan seberapa banyak yang diberikan, tetapi ketulusan di balik pemberian tersebut.",
			},
			{
				Title:   "Sedekah dengan Senyuman",
				Excerpt: "Kebaikan tidak selalu membutuhkan uang. Senyuman yang tulus, membantu seseorang membawa barang, memberikan jalan, atau mengucapkan perkataan yang baik juga dapat menjadi bentuk kebaikan. Biasakan mencari kesempatan untuk berbuat baik setiap hari.",
			},
			{
				Title:   "Mencari Rezeki yang Halal",
				Excerpt: "Islam mengajarkan umatnya untuk mencari rezeki dengan cara yang halal dan baik. Jangan mengorbankan kejujuran demi keuntungan sesaat. Rezeki yang sedikit tetapi halal dan penuh keberkahan lebih baik daripada harta banyak yang diperoleh melalui cara yang tidak benar.",
			},
			{
				Title:   "Tawakal Setelah Berikhtiar",
				Excerpt: "Tawakal bukan berarti tidak melakukan apa-apa. Seorang Muslim tetap diperintahkan untuk berusaha, kemudian menyerahkan hasilnya kepada Allah. Lakukan bagian kita dengan sungguh-sungguh dan percayakan hasil akhirnya kepada Allah SWT.",
			},
			{
				Title:   "Bersyukur atas Rezeki yang Sedikit",
				Excerpt: "Jangan menunggu kaya untuk bersyukur. Nikmat kecil yang sering dianggap biasa seperti makanan, kesehatan, tempat tinggal, dan keluarga merupakan karunia yang sangat berharga. Dengan bersyukur, hati menjadi lebih tenang dan tidak mudah merasa kekurangan.",
			},

			// Doa & Dzikir
			{
				Title:   "Doa Ketika Menghadapi Kesulitan",
				Excerpt: "Ketika menghadapi masalah, jangan hanya mengandalkan kekuatan diri sendiri. Berdoalah kepada Allah, kemudian lakukan usaha terbaik untuk mencari jalan keluar. Kesulitan dapat menjadi sarana untuk semakin dekat kepada Allah dan belajar bersabar.",
			},
			{
				Title:   "Membiasakan Dzikir di Tengah Kesibukan",
				Excerpt: "Dzikir tidak harus selalu dilakukan dalam waktu yang panjang. Di sela perjalanan, menunggu pekerjaan, atau setelah menyelesaikan aktivitas, kita dapat membiasakan membaca tasbih, tahmid, takbir, dan istighfar. Sedikit tetapi rutin dapat membantu hati tetap terhubung dengan Allah.",
			},
			{
				Title:   "Doa untuk Kedua Orang Tua",
				Excerpt: "Salah satu bentuk bakti kepada orang tua adalah mendoakan mereka. Allah SWT mengajarkan doa: 'Rabbighfir li wa liwalidayya warhamhuma kama rabbayani shaghira.' Biasakan mendoakan kedua orang tua, baik ketika mereka masih hidup maupun setelah mereka wafat.",
			},
			{
				Title:   "Doa Sebelum Tidur",
				Excerpt: "Menutup hari dengan mengingat Allah adalah kebiasaan yang baik. Sebelum tidur, biasakan berwudhu jika memungkinkan, membaca dzikir dan doa yang diajarkan Rasulullah SAW, serta memaafkan kesalahan orang lain. Tidur pun menjadi lebih tenang dengan mengingat Allah.",
			},
			{
				Title:   "Bangun Tidur dan Mengingat Allah",
				Excerpt: "Ketika membuka mata di pagi hari, jangan langsung tenggelam dalam notifikasi dan aktivitas dunia. Mulailah dengan bersyukur kepada Allah karena masih diberikan kesempatan untuk hidup. Jadikan pagi sebagai awal yang baik dengan doa, dzikir, dan niat melakukan kebaikan.",
			},

			// Kehidupan Sehari-hari
			{
				Title:   "Islam Mengajarkan Keseimbangan dalam Kehidupan",
				Excerpt: "Islam tidak mengajarkan seseorang untuk hanya mengejar dunia atau meninggalkan dunia sepenuhnya. Seorang Muslim perlu bekerja, belajar, menjaga keluarga, beribadah, dan membantu sesama secara seimbang. Kehidupan dunia dapat menjadi jalan menuju kebaikan akhirat.",
			},
			{
				Title:   "Menggunakan Waktu dengan Bijak",
				Excerpt: "Waktu adalah salah satu nikmat yang sering baru disadari ketika telah berlalu. Gunakan waktu untuk hal-hal yang bermanfaat seperti bekerja, belajar, beribadah, menjaga keluarga, dan membantu orang lain. Hindari menghabiskan sebagian besar waktu untuk aktivitas yang tidak memberikan manfaat.",
			},
			{
				Title:   "Menjaga Amanah dalam Pekerjaan",
				Excerpt: "Bekerja bukan hanya tentang mendapatkan penghasilan, tetapi juga tentang menjalankan amanah. Kerjakan tugas dengan jujur, jangan mengambil hak orang lain, dan jangan sengaja mengabaikan tanggung jawab. Profesionalitas dan kejujuran merupakan bagian dari akhlak seorang Muslim.",
			},
			{
				Title:   "Islam Mengajarkan untuk Terus Belajar",
				Excerpt: "Menuntut ilmu merupakan bagian penting dalam kehidupan seorang Muslim. Belajarlah ilmu agama sekaligus ilmu yang bermanfaat bagi kehidupan. Jadikan proses belajar sebagai bentuk syukur atas kemampuan berpikir yang Allah berikan.",
			},
			{
				Title:   "Menjaga Silaturahmi",
				Excerpt: "Silaturahmi dapat dilakukan dengan banyak cara, seperti mengunjungi keluarga, menanyakan kabar, membantu ketika membutuhkan, atau sekadar mengirim pesan. Jangan menunggu momen besar untuk menjaga hubungan. Perhatian kecil yang dilakukan secara rutin dapat mempererat persaudaraan.",
			},
			{
				Title:   "Menghindari Berlebihan dalam Kehidupan",
				Excerpt: "Islam mengajarkan keseimbangan dan tidak berlebih-lebihan. Baik dalam makanan, belanja, hiburan, maupun penggunaan waktu, seorang Muslim perlu belajar mengendalikan diri. Hidup sederhana bukan berarti kekurangan, tetapi mampu membedakan antara kebutuhan dan keinginan.",
			},
			{
				Title:   "Menjaga Hati Ketika Mendapat Pujian",
				Excerpt: "Pujian dapat menjadi ujian bagi hati. Ketika mendapatkan apresiasi, bersyukurlah kepada Allah dan jangan sampai pujian membuat kita merasa lebih tinggi daripada orang lain. Ingat bahwa setiap kelebihan adalah karunia yang dapat berubah kapan saja.",
			},
			{
				Title:   "Ketika Doa Belum Dikabulkan",
				Excerpt: "Terkadang kita telah berdoa berkali-kali tetapi belum melihat hasil yang diharapkan. Jangan langsung menganggap Allah tidak mendengar doa kita. Tetaplah berdoa, berusaha, dan percaya bahwa Allah mengetahui waktu serta sesuatu yang paling baik bagi hamba-Nya.",
			},
			{
				Title:   "Jangan Berputus Asa dari Rahmat Allah",
				Excerpt: "Sebesar apa pun kesalahan seseorang, pintu taubat dan rahmat Allah tetap terbuka selama hidup. Jangan biarkan dosa membuat kita menjauh dari Allah. Justru ketika jatuh, segera bangkit, memohon ampun, dan kembali memperbaiki diri.",
			},
			{
				Title:   "Memulai Perubahan dari Hal Kecil",
				Excerpt: "Perubahan menjadi lebih mudah ketika dimulai dari kebiasaan sederhana. Mulailah dengan menjaga satu waktu shalat, membaca beberapa ayat Al-Qur'an, memperbanyak istighfar, atau bersedekah secara rutin. Jangan menunggu menjadi sempurna untuk mulai menjadi lebih baik.",
			},
			{
				Title:   "Istiqamah dalam Kebaikan",
				Excerpt: "Menjadi baik selama satu hari mungkin terasa mudah, tetapi mempertahankannya setiap hari membutuhkan istiqamah. Fokuslah pada amal yang mampu dilakukan secara konsisten. Sedikit tetapi rutin dapat menjadi kebiasaan yang membentuk karakter dan mendekatkan diri kepada Allah.",
			},
			{
				Title:   "Muhasabah Diri Sebelum Tidur",
				Excerpt: "Sebelum tidur, luangkan beberapa menit untuk mengevaluasi hari yang telah dilalui. Tanyakan kepada diri sendiri: kebaikan apa yang sudah dilakukan, kesalahan apa yang perlu diperbaiki, dan apa yang bisa dilakukan lebih baik besok? Muhasabah membantu kita terus berkembang menjadi pribadi yang lebih baik.",
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
