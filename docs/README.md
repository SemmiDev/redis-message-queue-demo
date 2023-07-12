Distributed task queue dengan Golang Asynq

## Apa itu Message Queue?
Message Queue / Antrian pesan adalah komponen perangkat lunak yang memfasilitasi komunikasi antara aplikasi-aplikasi yang ada dalam infrastruktur microservices dan serverless. Protokol komunikasi asinkron mengirim dan menerima pesan yang diantri dan tidak memerlukan respons segera dari penerima (asynchronous communication). 

Antrian pesan sangat penting, karena mereka menyediakan komunikasi dan koordinasi antara aplikasi-aplikasi yang terdistribusi. Selain itu, mereka dapat secara signifikan nge simplify the coding of decoupled applications sambil meningkatkan keandalan, performa, dan skalabilitas.

## Emangnya kenapa sih kita perlu pake Message Queue?

Berikut adalah contoh usecases untuk setiap pertimbangan penggunaan Message Queue:
- Komunikasi Asynchronous
Contohnya: Pengiriman Email

Ketika pengguna melakukan pendaftaran dalam aplikasi, Anda dapat menggunakan Message Queue untuk mengirimkan email verifikasi kepada pengguna. Proses pengiriman email dapat menjadi tugas yang memakan waktu lama, dan dengan menggunakan Message Queue, Anda dapat mengantri tugas pengiriman email ke dalam antrian pesan. Ini memungkinkan aplikasi untuk melanjutkan operasi lain tanpa harus menunggu pengiriman email selesai, sehingga meningkatkan responsivitas dan pengalaman pengguna.

- Skalabilitas dan Ketersediaan
Contohnya: Pengolahan Gambar Skala Besar

Misalkan Anda memiliki aplikasi yang memungkinkan pengguna mengunggah gambar dan Anda perlu memproses gambar-gambar tersebut dalam skala besar. Dengan menggunakan Message Queue, Anda dapat mengirimkan tugas pemrosesan gambar ke antrian pesan dan menjalankan beberapa worker secara paralel untuk memproses tugas-tugas tersebut. Ini memungkinkan Anda untuk mengatasi lonjakan lalu lintas dan memastikan bahwa pemrosesan gambar dapat dilakukan secara efisien bahkan dalam situasi dengan beban kerja yang tinggi.

Pengiriman Data Real-time 
Antrian pesan memungkinkan pengiriman data secara real-time ke service lain, memungkinkan integrasi yang mulus antara berbagai aplikasi.

- Manajemen Kesalahan dan Retries
Contohnya: Integrasi Eksternal

Ketika berintegrasi dengan sistem eksternal, ada kemungkinan terjadi kesalahan sementara seperti gangguan jaringan atau masalah koneksi. Dalam hal ini, Message Queue memungkinkan Anda untuk mengatur kebijakan retry yang sesuai untuk mencoba kembali pemrosesan tugas secara otomatis. Misalnya, jika terjadi kesalahan saat berkomunikasi dengan sistem pembayaran eksternal, Anda dapat mengantri ulang tugas pembayaran ke antrian pesan dan memberikan kebijakan retry untuk mencoba kembali dengan interval waktu yang ditentukan.


Dan masih banyak lagi contoh penggunaan Message Queue lainnya

## Apa itu Asynq?

Asynq adalah sebuah library Go yang digunakan untuk mengantri tugas-tugas dan memprosesnya secara asynchronous dengan menggunakan worker. Library ini didukung oleh Redis dan dirancang agar mudah digunakan dan scalable.

Secara umum, Cara kerja Asynq sebagai berikut berikut:

- Klien menaruh/enqueue tugas-tugas pada sebuah antrian (queue).
- Server menarik tugas-tugas dari antrian dan memulai goroutine worker untuk setiap tugas.
- Tugas-tugas diproses secara concurrent oleh beberapa worker. 

Antrian tugas digunakan sebagai mekanisme untuk mendistribusikan pekerjaan ke berbagai mesin. Sistem dapat terdiri dari beberapa server worker dan broker, yang memungkinkan ketersediaan/ yang tinggi (high availability) dan horizontal scaling.

Berikut ilustrasi cara kerja Asynq:

Berikut fitur fitur dari asynq:
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
- Low latency to add a task since writes are fast in Redis
- De-duplication of tasks using unique option
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

Selengkapnya, temen2 bisa baca di sini https://github.com/hibiken/asynq/wiki/ ya hehe


- Endpoints
  - create post
  - get post

- Admin
  - When admin create a post, will enqueue to `email` queue

- Subscriber users
  - When user click on link in email, will enqueue to `analytic` queue
  - When user read the post, will enqueue to `analytic` queue

- Email Worker
  - Notify to all subscribers with sending email with Urchin Tracking Module (UTM) parameters (for analytics) when admin create a post

- Analytic Worker
  - Count user read post
  - Count user read post by UTM parameters