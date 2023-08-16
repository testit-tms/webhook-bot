# webhook-bot

To run the bot locally in docker compose:
1. Build the image and run start project with docker-compose.yml:    
```bash
docker build -f deploy/Dockerfile -t webhook-bot:latest .
```
2. Add your telgram-bot token to deploy/.env and customize other variables if needed.
3. Deploy the project in docker compose:
```bash
docker-compose -f deploy/docker-compose.yml -p webhook-bot up -d
```
