# Control Panel

Frontend каркас для системы Web App Security.

## Быстрый старт

```bash
npm install
npm run dev
```

## Команды

- `npm run dev` — локальный запуск (Vite).
- `npm run build` — сборка для продакшена.
- `npm run preview` — предпросмотр сборки.
- `npm run lint` — запуск линтера.
- `npm run generate-api` — генерация API клиента по OpenAPI.

## Docker

```bash
docker build -t web-app-security-control-panel .
docker run --rm -p 8080:80 web-app-security-control-panel
```
