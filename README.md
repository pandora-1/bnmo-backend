# BNMO - BACKEND

## Prasyarat
1. Go 1.18.3 atau terbaru
2. MacOS, Windows (termasuk WSL), dan / atau Linux
3. Docker 20.10.17 atau terbaru
4. MySQL

## Cara menjalankan
1. Pada file .env, isi bagian dibawah ini sesuai dengan lingkungan pengguna
DB_USER=
DB_PASSWORD=
DB_HOST=
DB_PORT=
DB_NAME=
2. Pada ./database/database.go di fungsi NewDatabase, isi bagian dibawah ini sesuai dengan lingkungan pengguna
    USER := 
    PASS := 
    HOST :=
    DBNAME := 
3. Menjalankan command berikut

```bash
docker-compose up --build    
```

Buka [http://localhost:8000](http://localhost:8000) untuk melihat hasil backend yang dibuat

## Teknologi yang digunakan

- Go 1.18.3
- Docker 20.10.17
- Gin
- MySQL

=======
### Marcellus Michael Herman Kahari - 13520057

