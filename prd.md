Berikut **update PRD Shohibul Qur’an Bot** sesuai perubahan yang kamu minta:

* ❌ **Anti Abuse System dihapus**
* ✅ **Laporan boleh berkali-kali**
* ✅ **Halaman tidak dibatasi**
* ✅ **Motivasi hanya 1x sehari**
* ✅ **Waktu motivasi random setiap hari**

Saya tulis bagian yang berubah agar PRD tetap konsisten untuk development.

---

# 📄 PRD — Shohibul Qur'an Bot (Updated)

## 1. Overview

**Nama Produk:** Shohibul Qur'an Bot
**Platform:** WhatsApp Group
**Teknologi:** Golang + Whatsmeow

Shohibul Qur’an Bot adalah bot WhatsApp yang membantu komunitas membangun kebiasaan membaca Al-Qur’an melalui sistem:

* laporan membaca
* gamifikasi
* reminder
* motivasi ayat & hadist
* challenge komunitas

Bot dirancang agar membaca Al-Qur’an menjadi **habit harian yang menyenangkan dan konsisten**.

---

# 🎯 2. Tujuan Produk

Tujuan utama:

* membangun kebiasaan membaca Al-Qur’an setiap hari
* meningkatkan konsistensi membaca
* menciptakan motivasi komunitas

---

# 👥 3. Target User

Komunitas muslim seperti:

* grup halaqah
* keluarga
* komunitas kajian
* komunitas Qur’an

---

# ⚙️ 4. Core Features

---

# 📖 4.1 Laporan Membaca Qur’an

Bot mendeteksi pesan yang mengandung kata:

```
alhamdulillah
alhamdulillahirabbilalamin
```

Contoh laporan:

```
Alhamdulillah sudah baca 2 halaman
alhamdulillah 1 juz
hari ini alhamdulillah 5 halaman
```

---

## Behavior Bot

Ketika laporan diterima bot akan:

1️⃣ mencatat laporan
2️⃣ menghitung halaman
3️⃣ menambah XP
4️⃣ update statistik

---

# 📊 4.2 Perhitungan Bacaan

Standar mushaf Madinah:

```
1 juz = 20 halaman
```

Contoh parsing:

| Input     | Hasil |
| --------- | ----- |
| 3 halaman | 3     |
| 5 hlm     | 5     |
| 1 juz     | 20    |
| 0.5 juz   | 10    |

---

# 🔁 4.3 Multiple Reports per Day

User **boleh laporan berkali-kali dalam satu hari**.

Contoh:

```
Alhamdulillah 3 halaman
Alhamdulillah 5 halaman
Alhamdulillah 1 juz
```

Total dihitung:

```
3 + 5 + 20 = 28 halaman
```

Tidak ada batasan jumlah laporan.

---

# 🔥 4.4 Daily Streak

Streak dihitung jika user **melaporkan bacaan minimal sekali dalam sehari**.

Jika satu hari tidak ada laporan → streak reset.

### Level Streak

| Streak | Title             |
| ------ | ----------------- |
| 3      | 🌱 Pemula Qur'an  |
| 7      | 🌿 Sahabat Qur'an |
| 14     | 🌳 Pecinta Qur'an |
| 30     | 🕌 Ahlul Qur'an   |
| 100    | 👑 Penjaga Qur'an |

---

# 🎮 4.5 XP System

XP diberikan berdasarkan aktivitas.

| Aktivitas       | XP  |
| --------------- | --- |
| laporan membaca | +10 |
| 1 halaman       | +2  |
| streak 7 hari   | +20 |
| quiz benar      | +5  |

---

# 🎖 4.6 Level System

Level berdasarkan total XP.

| Level | Title          |
| ----- | -------------- |
| 1     | Pemula         |
| 5     | Sahabat Qur'an |
| 10    | Penjaga Ayat   |
| 20    | Ahlul Qur'an   |

---

# 🏆 4.7 Leaderboard

Leaderboard mingguan berdasarkan total halaman.

Contoh:

```
🏆 Ranking Mingguan

🥇 Ahmad — 45 halaman
🥈 Fatimah — 38 halaman
🥉 Ali — 32 halaman
```

---

# 📊 4.8 Statistik User

Command:

```
!stats
```

Output:

```
📊 Statistik Ahmad

Streak: 6 hari 🔥
Hari ini: 12 halaman
Bulan ini: 3.5 juz
XP: 320
Level: 7
```

---

# 👋 4.9 Welcome Message

Ketika ada member baru join.

Bot mengirim pesan:

```
Assalamu'alaikum 👋

Selamat datang di Shohibul Qur'an 📖

Cara laporan membaca:

Alhamdulillah 2 halaman
Alhamdulillah 1 juz

Bot akan mencatat laporan otomatis.
```

---

# 📜 4.10 Group Rules

Rules disampaikan saat welcome.

1️⃣ niatkan membaca karena Allah
2️⃣ laporan dengan kata "Alhamdulillah"
3️⃣ boleh laporan berkali-kali
4️⃣ minimal membaca beberapa ayat

---

# 🌙 4.11 Random Ayat Qur’an

Bot mengirim ayat random dari API.

API:

```
https://api.quran.com/api/v4/verses/random
```

Pesan bot:

```
📖 Ayat Qur'an

اللَّهُ نُورُ السَّمَاوَاتِ وَالْأَرْضِ

Allah adalah cahaya langit dan bumi

(QS An-Nur 35)
```

---

# 📚 4.12 Random Hadist

API:

```
https://api.hadith.gading.dev
```

Pesan bot:

```
🌙 Hadist

"Sebaik-baik kalian adalah yang belajar Al-Qur'an dan mengajarkannya."

(HR Bukhari)
```

---

# ⏰ 5. Motivasi Harian

Motivasi dikirim **1 kali setiap hari**.

Waktu pengiriman **random**.

Contoh range waktu:

```
06:00 — 21:00
```

Contoh pesan:

```
📖 Motivasi Qur'an

"Bacalah Al-Qur'an karena ia akan datang memberi syafaat bagi pembacanya."

(HR Muslim)
```

---

# 🧠 6. Smart Qur'an Motivation Engine

Fitur ini membuat bot terasa hidup.

Bot akan menganalisis aktivitas grup.

---

## Metric yang dipantau

| Metric             | Fungsi             |
| ------------------ | ------------------ |
| last_report_time   | aktivitas terakhir |
| reports_today      | jumlah laporan     |
| active_users_today | user aktif         |

---

## Behavior Engine

### Kondisi 1 — Grup Sepi

Jika tidak ada laporan >6 jam.

Bot mengirim motivasi.

```
📖 Jangan lupa membaca Al-Qur'an hari ini 🤍
Walaupun hanya beberapa ayat.
```

---

### Kondisi 2 — Banyak laporan

Jika laporan tinggi.

Bot memberi apresiasi.

```
MasyaAllah 🔥

Hari ini sudah ada 20 laporan membaca Qur'an.
Semoga Allah memberkahi kita semua.
```

---

### Kondisi 3 — Target hampir tercapai

```
🎯 Target hampir tercapai

Tinggal 2 juz lagi untuk mencapai target minggu ini.
```

---

# 🤝 7. Challenge Grup

Target bacaan komunitas.

Contoh:

```
🎯 Target Mingguan
50 Juz
```

Progress:

```
Progress: 27 / 50 Juz
```

---

# 🗄 8. Database Schema

### users

```
id
phone
name
xp
level
streak
last_read_date
joined_at
```

---

### reports

```
id
user_id
pages
message
date
created_at
```

---

### daily_progress

```
user_id
date
pages
reports_count
```

---

### badges

```
user_id
badge
created_at
```

---

# 🏗 9. System Architecture

```
WhatsApp
   ↓
Whatsmeow Client
   ↓
Message Handler
   ↓
Parser
   ↓
Gamification Engine
   ↓
Database
```

---

# ⏰ 10. Scheduler

Jobs yang dijalankan:

| Job             | Waktu     |
| --------------- | --------- |
| reset harian    | 00:00     |
| motivasi random | 1x sehari |
| reminder        | 18:00     |

---

# 🚀 11. Deployment

Server minimal:

```
1 vCPU
1GB RAM
Ubuntu
```

Service:

```
systemd
docker
```

---

# ⭐ Future Vision

Shohibul Qur'an Bot dapat berkembang menjadi:

* Qur'an habit tracker
* komunitas khatam Qur'an global
* platform Qur'an community

---

Kalau kamu mau, langkah berikutnya yang paling penting sebelum coding adalah membuat:

1️⃣ **System Architecture Diagram yang detail**
2️⃣ **Flowchart seluruh logic bot (laporan, streak, motivasi, reminder)**
3️⃣ **Database ERD diagram**
4️⃣ **Struktur project Go + Whatsmeow yang production ready**

Biasanya ini membuat development **3–5x lebih cepat dan minim bug**.
