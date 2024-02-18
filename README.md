# Rust server dockerized container, with oxide(umod) support

#### Configure server with in file:
./servers/[servername]/server_config.json

#### Add plugins to:
./servers/[servername]/oxide

#### Run with:
docker-compose up -d 

After server is started uptrace with server in-game statistics and events will be avialable on http://localhost:14318, with login/pass: gena@rust.ru/HermanFuLLer/