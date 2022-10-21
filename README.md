# Coinche (SERVER)

This is the server part, the client part can be found here https://github.com/Alounv/coinche-front.

## What

A simple multiplayer coinche game. 

You can play here: https://coinche.vercel.app/

<img width="832" alt="Capture d’écran 2022-10-21 à 15 12 42" src="https://user-images.githubusercontent.com/34238160/197204169-798c9637-d5c0-45c0-bfbf-4363017852ab.png">

<details>
    <summary>Bidding</summary>
    <img width="831" alt="Capture d’écran 2022-10-21 à 15 11 12" src="https://user-images.githubusercontent.com/34238160/197204166-1f02e273-9e3b-4a71-b209-0a08a38ca28b.png">
</details>

## Why

This is a personal project to learn:
- Golang
- Test Driven Development
- Clean Architecture
- PostgreSQL
- WebSockets

## Developing

```bash
go run main.go
```

## Building

To create a production version of your app:

```bash
go run build
```

You can preview the production build with `./coinche`.

## Environment variables
```bash
## Developing

```bash
bun run dev

# or start the server and open the app in a new browser tab
bun run dev -- --open
```

## Building

To create a production version of your app:

```bash
bun run build
```

You can preview the production build with `bun run preview`.

## Environment variables
```bash
PORT=:5000
SQLX_POSTGRES_INFO="host=localhost user=aloun password=ILovePostgres port=5432"
DB_NAME=coincheDb
```
