Distributed task queue dengan Golang Asynq

## Apa itu Tasks Queue?
Task Queue adalah sebuah mekanisme yang digunakan untuk mengatur dan menjalankan tugas atau pekerjaan yang memerlukan eksekusi terpisah dalam sebuah sistem (asynchronous). Dengan menggunakan Task Queue, kita dapat memisahkan tugas-tugas tersebut dari proses utama dan membagikannya ke beberapa proses atau thread untuk nge handle task tersebut. Hal ini memungkinkan sistem untuk secara efisien mengelola workload yang kompleks dan meningkatkan kinerja serta skalabilitas aplikasi.

## Contoh Kasus Penggunaan Task Queue
- Pemrosesan Gambar
Ketika kita memiliki aplikasi atau layanan yang memerlukan pemrosesan gambar, seperti thumbnail generation atau filter aplikasi foto, Task Queue dapat digunakan untuk membagi tugas pemrosesan gambar menjadi bagian-bagian kecil yang dapat dijalankan secara paralel.

- Pengiriman Email Massal
Ketika kita perlu mengirim email massal ke ribuan atau bahkan jutaan pelanggan, Task Queue dapat digunakan untuk mengatur dan menjalankan tugas pengiriman email secara terpisah. Setiap email dapat diwakilkan sebagai tugas dalam antrian, dan proses atau worker dapat mengambil tugas tersebut dan mengirim email dengan efisien. Dengan menggunakan Task Queue, kita dapat menghindari pemblokiran atau penundaan dalam pengiriman email, mengelola antrian pengiriman dengan baik, sehingga meningkatkan responsivitas dan pengalaman pengguna.

- Notifikasi
- PDF Report
- Backup
- daan masih banyak lagi contoh kasus penggunaan Task Queue [lainnya](https://www.google.com/search?q=tasks+queues+usecases&oq=tasks+queues+usecases)

## Apa itu Asynq?

Asynq adalah sebuah library Go yang digunakan untuk mengantri tugas-tugas (tasks queues) dan memprosesnya secara asynchronous dengan menggunakan worker. Di balik layar Library ini menggunakan Redis dan dirancang agar mudah digunakan dan scalable.

Secara umum, Cara kerja Asynq sebagai berikut berikut:

- Klien enqueue tugas-tugas pada sebuah antrian (queue).
- Server menarik tugas-tugas dari antrian dan memulai goroutine worker untuk setiap tugas.
- Tugas-tugas diproses secara concurrent oleh beberapa worker. 

Antrian tugas digunakan sebagai mekanisme untuk mendistribusikan pekerjaan ke berbagai mesin ataupun concurrent process pada mesin yang sama. Sistem dapat terdiri dari beberapa server worker dan broker, yang memungkinkan ketersediaan/ yang tinggi (high availability) dan horizontal scaling.

Berikut ilustrasi cara kerja Asynq:
![Asynq](https://user-images.githubusercontent.com/11155743/116358505-656f5f80-a806-11eb-9c16-94e49dab0f99.jpg)

Fitur fitur yang ditawarkan asynq sebagai berikut:
- Guaranteed at least one execution of a task
Jika terjadi failure yang menyebabkan hilangnya pesan atau membutuhkan waktu terlalu lama untuk recovery, pesan akan dikirim ulang untuk memastikan pesan terkirim setidaknya satu kali.

- Scheduling of tasks
Dapat menjadwalkan tugas untuk dieksekusi di masa depan dengan menggunakan asynq.ProcessIn(time)

- Retries of failed tasks
Dapat mengatur jumlah maksimal retry untuk tugas yang gagal, dengan menggunakan asynq.MaxRetry(10)

- Automatic recovery of tasks in the event of a worker crash
Jika worker crash, tugas yang sedang diproses akan dikembalikan ke antrian untuk diproses ulang.

- Weighted priority queues
Dapat mengatur prioritas tugas, misalnya menggunakan asynq.Queue("critical") untuk tugas yang kritis atau penting.

- Strict priority queues
asynq mendukung antrian prioritas ketat, yang berarti tugas dengan prioritas lebih tinggi akan selalu dieksekusi terlebih dahulu dibandingkan dengan tugas dengan prioritas lebih rendah.

- Low latency to add a task since writes are fast in Redis
Penulisan data pada Redis sangat cepat, sehingga asynq memiliki latensi rendah dalam menambahkan tugas ke dalam antrian.

- De-duplication of tasks using unique option
Anda dapat menghindari duplikasi tugas dengan menggunakan opsi unik pada tugas-tugas yang sama. Jika sebuah tugas dengan opsi unik sudah ada dalam antrian, tugas baru dengan opsi yang sama akan diabaikan.

- Allow timeout and deadline per task
- Allow aggregating group of tasks to batch multiple successive operations
- Flexible handler interface with support for middlewares
- Ability to pause queue to stop processing tasks from the queue
- Periodic Tasks
- Support Redis Cluster for automatic sharding and high availability
- Support Redis Sentinels for high availability
- Integration with Prometheus to collect and visualize queue metrics
- Web UI to inspect and remote-control queues and tasks
- CLI to inspect and remote-control queues and tasks

Supeer duper komplit kan? Selengkapnya, temen2 bisa baca di sini https://github.com/hibiken/asynq/wiki/ ya hehe

## Nyobain Asynq dengan studi kasus sederhana
Mari kita jelajahi penerapan Asynq pada sebuah aplikasi blog sederhana yang melibatkan admin, subscriber, dan tugas-tugas terkait notifikasi dan analitik. Dalam aplikasi ini, admin dapat membuat post, dan para pembaca dapat berlangganan newsletter. Setiap kali admin membuat post, subscriber akan menerima email notifikasi. Selain itu, kita ingin melacak berapa banyak subscriber yang membaca post melalui tautan di email notifikasi, serta berapa banyak yang membaca langsung di website.

- Admin membuat post: Ketika admin membuat post baru, kita akan menerapkan konsep tasks queue. Setelah post dibuat, kita akan mengirimkan tugas "kirim notifikasi email" ke dalam antrian "email". Tugas ini akan dikonsumsi oleh worker email yang akan mengirimkan email notifikasi kepada semua subscriber.

- Pembaca membaca post melalui tautan di email: Ketika pembaca membuka tautan di email notifikasi, kita akan menerapkan konsep tasks queue sekali lagi. Setiap kali ada klik tautan, kita akan menambahkan tugas "lacak pembacaan melalui tautan" ke dalam antrian "analytic". Tugas ini akan dikonsumsi oleh worker analitik yang akan mencatat bahwa pembaca tersebut membaca post melalui tautan. 

- Pembaca membaca post langsung di website: Ketika pembaca membaca post langsung di website, kita juga akan menerapkan konsep tasks queue. Setiap kali ada aksi membaca post, kita akan menambahkan tugas "lacak pembacaan langsung" ke dalam antrian "analytic". Tugas ini akan dikonsumsi oleh worker analitik yang akan mencatat bahwa pembaca tersebut membaca post langsung di website.

Untuk mendeteksi pembaca membuka link dari email atau website nya langsung kita tambahkan parameter UTM (Urchin Tracking Module) pada link yang dikirimkan ke email. Parameter UTM ini akan kita gunakan untuk melacak pembacaan post.

## Apa aja yang kita butuhkan?
- Go, kita akan menggunakan Go sebagai bahasa pemrograman
- Redis, kita akan menggunakan Redis sebagai broker untuk Asynq
- Docker & Docker Compose, kita akan menjalankan Redis menggunakan Docker
- Dan package2 go seperti chi router, dll

## Constrains
- Datastore (penyimpanan post, dll) kita hanya akan menggunakan memory (map), jadi ketika server mati, data akan hilang hehe
- Kita ga bikin tests (unittest,apitest,dll) krna lama bet bikinnya, jadi kita akan langsung coba run aja (blackbox test)

## Show me the code
Eitz, sabar dulu, berikut folder struktur project kita

```
├── cmd
   └── main.go
├── article_analytic.go
├── article_analytic_handler.go
├── article.go
├── article_handler.go
├── article_subscriber.go
├── article_subscriber_handler.go
├── http_server.go
├── mail_sender.go
├── task_send_new_article_notification_email.go
├── task_send_new_subscriber_welcome_email.go
├── task_track_read_article.go
├── utm_generator.go
├── utm_generator_test.go
├── worker_distributor.go
├── worker_logger.go
├── worker_processor.go
├── go.mod
├── go.sum
├── docker-compose.yaml
└── playground.http

```
