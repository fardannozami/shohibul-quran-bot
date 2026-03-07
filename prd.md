Baik, kita lanjutkan dengan **dokumen teknis lanjutan** yang akan sangat membantu sebelum coding dimulai. Saya buat 4 bagian penting:

1️⃣ **System Architecture Diagram**
2️⃣ **Flowchart Logic Bot**
3️⃣ **Database ERD**
4️⃣ **Struktur Project Go + Whatsmeow (production ready)**

Ini biasanya disebut **Technical Design / System Design Document**.

---

# 🏗 1. System Architecture

Arsitektur sistem Shohibul Qur'an Bot.

```
                WhatsApp Server
                       │
                       │
                Whatsmeow Client
                       │
        ┌──────────────┴──────────────┐
        │                             │
   Message Handler              Event Handler
 (chat message)                (join group)
        │                             │
        │                             │
        ▼                             ▼
     Parser                    Welcome System
 (detect laporan)
        │
        ▼
  Gamification Engine
 (XP, streak, stats)
        │
        ▼
      Database
        │
        ▼
     Scheduler
(reminder, motivasi)
```

### Penjelasan komponen

| Komponen            | Fungsi                |
| ------------------- | --------------------- |
| Whatsmeow           | koneksi ke WhatsApp   |
| Message Handler     | menerima pesan grup   |
| Parser              | mendeteksi laporan    |
| Gamification Engine | XP, streak, statistik |
| Scheduler           | reminder & motivasi   |
| Database            | menyimpan data        |

---

# 🔄 2. Flowchart — Laporan Mengaji

```
User kirim pesan
       │
       ▼
Apakah mengandung
"alhamdulillah" ?
       │
   ┌───┴────┐
   │        │
  TIDAK    YA
   │        │
   ▼        ▼
Ignore   Parse halaman/juz
             │
             ▼
        Hitung halaman
             │
             ▼
        Simpan laporan
             │
             ▼
        Update streak
             │
             ▼
          Tambah XP
             │
             ▼
       Kirim respon bot
```

---

# 🔄 3. Flowchart — Reminder

```
Scheduler 18:00
       │
       ▼
Ambil semua user
       │
       ▼
Cek siapa yang
belum laporan
       │
       ▼
Jika ada user
belum laporan
       │
       ▼
Kirim reminder
mention user
```

---

# 🔄 4. Flowchart — Motivasi Random

```
Scheduler start
       │
       ▼
Generate random time
(06:00 - 21:00)
       │
       ▼
Tunggu waktu tersebut
       │
       ▼
Ambil random ayat
atau hadist
       │
       ▼
Kirim pesan motivasi
ke grup
```

---

# 🔄 5. Flowchart — Welcome Member

```
User join group
       │
       ▼
Event GroupParticipants
       │
       ▼
Action = add ?
       │
   ┌───┴────┐
   │        │
  NO       YES
   │        │
 ignore   Kirim welcome
            │
            ▼
      Jelaskan rules
```

---

# 🗄 6. Database ERD

Struktur relasi database.

```
users
 ├── id
 ├── phone
 ├── name
 ├── xp
 ├── level
 ├── streak
 └── joined_at
       │
       │
       ▼
reports
 ├── id
 ├── user_id
 ├── pages
 ├── message
 └── date
       │
       │
       ▼
daily_progress
 ├── user_id
 ├── date
 ├── pages
 └── reports_count
       │
       │
       ▼
badges
 ├── user_id
 ├── badge
 └── created_at
```

Relasi utama:

```
users 1 --- n reports
users 1 --- n badges
users 1 --- n daily_progress
```

---

# 📦 7. Struktur Project Go

Struktur project yang scalable.

```
shohibul-quran-bot
│
├── cmd
│   └── bot
│       └── main.go
│
├── internal
│
│   ├── bot
│   │   ├── client.go
│   │   ├── message_handler.go
│   │   └── event_handler.go
│
│   ├── parser
│   │   └── report_parser.go
│
│   ├── gamification
│   │   ├── xp.go
│   │   ├── streak.go
│   │   └── badge.go
│
│   ├── motivation
│   │   ├── quran_api.go
│   │   ├── hadith_api.go
│   │   └── engine.go
│
│   ├── scheduler
│   │   └── cron.go
│
│   └── database
│       ├── db.go
│       ├── models.go
│       └── repository.go
│
├── config
│   └── config.yaml
│
└── migrations
```

---

# ⚙️ 8. Module Responsibilities

### bot/

menangani komunikasi dengan WhatsApp.

```
client.go
message_handler.go
event_handler.go
```

---

### parser/

mendeteksi laporan.

```
ParsePages()
DetectReport()
```

---

### gamification/

menghitung XP dan streak.

```
AddXP()
UpdateStreak()
GrantBadge()
```

---

### motivation/

mengambil ayat dan hadist dari API.

```
GetRandomAyat()
GetRandomHadith()
```

---

### scheduler/

menjalankan:

```
Reminder
Motivation
ResetDaily
```

---

# 🧠 9. Algoritma Random Motivasi

Saat bot start:

```
randomHour = random(6,21)
randomMinute = random(0,59)
```

Contoh hasil:

```
14:37
```

Scheduler akan kirim motivasi pada waktu tersebut.

Setiap hari waktu akan berbeda.

---

# 📊 10. Estimasi Kompleksitas

Jika semua fitur di PRD dibuat.

| Modul             | LOC kira-kira |
| ----------------- | ------------- |
| Bot Handler       | 400           |
| Parser            | 200           |
| Gamification      | 400           |
| Motivation Engine | 300           |
| Scheduler         | 200           |

Total kira-kira:

```
1500 — 2500 lines Go
```

Masih sangat manageable.

---

# 🚀 11. Development Plan

### Step 1

setup project

```
Go
Whatsmeow
Database
```

---

### Step 2

fitur dasar

* laporan
* parsing halaman
* database

---

### Step 3

gamification

* XP
* streak
* leaderboard

---

### Step 4

scheduler

* reminder
* motivasi

---

### Step 5

welcome & rules

---

# ⭐ Rekomendasi Penting

Agar bot ini **benar-benar kuat dan scalable**, sebaiknya tambahkan juga:

### Heatmap Konsistensi

seperti GitHub:

```
🟩🟩⬜🟩🟩
🟩🟩🟩🟩⬜
```

Ini sangat memotivasi.

---

✅ Dengan dokumen ini sebenarnya kamu sudah punya **blueprint lengkap untuk membangun Shohibul Qur'an Bot**.

---

Jika kamu mau, saya bisa lanjut membuat sesuatu yang **sangat penting untuk developer**, yaitu:

* **starter source code Go + Whatsmeow (siap run)**
* lengkap dengan **handler laporan, parser halaman, database, scheduler**

Biasanya dengan itu bot **sudah bisa running dalam 30 menit**.
