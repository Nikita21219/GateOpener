# GateOpener

Телеграм бот для управления шлагбаумами в жилом комплексе.  
Предоставляет 4 кнопки:
- Открыть шлагбаум на въезд
- Открыть шлагбаум на выезд
- Открыть шлагбаумы на въезд и выезд на 5 минут
- Закрыть шлагбаумы раньше чем через 5 минут

Usage:  
``` bash
echo BOT_TOKEN=YOUR_BOT_TOKEN >> .env; \
echo SID=YOUR_USER_ID_TO_OPEN_GATES >> .env; \
docker-compose up -d
```
