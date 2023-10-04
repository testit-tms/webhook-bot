# Telegram bot for sending messages to users

## Description

This bot is used to send messages to users or channels from webhooks.

## Usage

### Our installation

How to use our installation of the bot you can read [here](https://docs.testit.software/user-guide/work-with-projects/set-up-webhooks/set-up-telegram-notifications-using-webhooks.html#%D1%81%D0%BE%D0%B7%D0%B4%D0%B0%D0%BD%D0%B8%D0%B5-%D1%80%D0%B5%D0%B4%D0%B0%D0%BA%D1%82%D0%B8%D1%80%D0%BE%D0%B2%D0%B0%D0%BD%D0%B8%D0%B5-%D0%B2%D0%B5%D0%B1%D1%85%D1%83%D0%BA%D0%B0-%D0%B4%D0%BB%D1%8F-%D1%87%D0%B0%D1%82-%D0%B1%D0%BE%D1%82%D0%B0-telegram)

### Your own installation

You can run the bot locally or deploy it to your server.

Just follow the instructions below.

1. Create a bot in Telegram and get a token. You can read more about it [here](https://core.telegram.org/bots#6-botfather).
2. Clone this repository.
3. Add your bot token, image name and tag to the .env file.
4. Build the image

```bash
docker build -f deploy/Dockerfile -t webhook-bot:latest .
```

5. Run the containers

```bash
docker-compose -f deploy/docker-compose.yml -p webhook-bot up -d
```

## Contributing

Contributions are welcome! If you find a bug or have a feature request, please open an issue on the [issue tracker](https://github.com/testit-tms/webhook-bot/issues).

If you want to contribute code, please follow these steps:

1. Fork the repository.
2. Create a new branch for your changes.
3. Make your changes and write tests if necessary.
4. Run the tests and make sure they pass.
5. Commit your changes and push them to your fork.
6. Open a pull request to the `main` branch of the original repository.

We will review your changes and merge them if they meet our quality standards. Thank you for contributing!

## License

Distributed under the Apache-2.0 License. See LICENSE for more information.
