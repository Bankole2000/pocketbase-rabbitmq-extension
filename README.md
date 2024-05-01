## Pocketbase Extended with Go + RabbitMQ

<br />
<p float="left" align="middle">
  <a href="https://pocketbase.io">
    <img src="./media/pocketbaselogo.svg" alt="Logo" width="100" height="100">
  </a>
  <img src="./media/plus-icon.svg" alt="Logo" width="50" style="margin: 20px" height="50">
  <a href="https://www.rabbitmq.com/" target="_blank">
    <img src="./media/rabbitmq-logo.svg" alt="Logo" width="80" height="90">
  </a>

<!-- ABOUT THE PROJECT -->
### A simple Pocketbase RabbitMQ publisher extension template

![Publishing with Rabbit MQ][demo-gif]

### What it does

- Publish any and all events from pocketbase to rabbitmq exchange

- Select and customize published events

- Useful for multiple scenarios like data replication, horizontal scaling, single source of truth accross microservices, async data processing & event handling etc

### Prerequisites

- [Golang](https://go.dev/) >= v1.19
  
- [RabbitMQ](https://www.rabbitmq.com/) (Running locally or in the cloud)
- [NodeJS](https://nodejs.org) (Optional - for testing the listener)

### Setup & installation

1. Clone or fork this repository (or use this template)

    ```sh
    git clone https://github.com/bankole2000/pb-rmq-ext.git
    ```

2. init go dependencies in `/pocketbase-publisher` folder

    ```bash
      # in /pocketbase-publisher folder
    go mod init <app-name> && go mod tidy
      # Note: <app-name> can be whatever you want but
      # will be name of the executable file created on build
    ```

3. install node dependences in `/listener-demo` folder (optional)

    ```sh
      # in /listener-demo folder
    npm i
    ```

4. Following the example in `/pocketbase-publisher/.env.example` Create `.env` file in the `/pocketbase-publisher` folder and fill in environemental

    ```bash
     # in /pocketbase-publisher/.env file
    RABBITMQ_URL="amqp://<username>:<password>@<rabbitmq-host>:<rabbitmg-port>/" 
     # replace with your rabbitmq connection string
    RABBITMQ_EXCHANGE="my-exchange"
     # replace with any name you wish to give the exchange
    ```

#### No RabbitMQ? No problem

This step requires [docker](https://www.docker.com/) - To get a local instance of Rabbitmq running simply run

```docker
docker run -d --hostname localhost -p 5672:5672 -p 15672:15672 --name rabbitmq rabbitmq:3-management
```

Or run `docker compose up` using [this docker compose yml file](./media/docker-compose.yml)

### How To Use

1. Start the pocketbase application

   ```sh
     # in /pocketbase-publisher folder
   go run main.go serve
   ```

2. Start up the message/event listener

   ```sh
     # in /listener-demo
   npm run dev
   ```

3. Use the pocketbase api or sdk or admin (running on `http://127.0.0.1:8090/_/`) as you normally would (creating collections, records etc) and watch the events logged by the listener.

4. To build a statically linked single file application, you can run

   ```sh
     # create build file
   CGO_ENABLED=0 go build
     # run created executable
   ./my-app serve
   ```

> ⚠ Warning: If you run into any `gcc failed` errors running `go run main.go serve` , try setting  `CGO_ENABLED=0` environment variable according to the [pocketbase docs here](https://pocketbase.io/docs/go-overview/)

> 📝 Note: in the `/pocketbase-publisher/main.go` file, you can also change the values of `rmqUrl` and `exchange` variables on _lines 55_ and _56_ to the string values of your rabbitmq Url and the exchange name. This way the build won't rely on the `.env` file to run (be sure to delete or comment out the `goDotEnvVariable` function in _line 39_)

<!-- CONTRIBUTING -->
## Contributing

Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'feat: added some amazing feature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<!-- CONTACT -->
## Contact

📧 Shoot me an email - <techybanky@gmail.com>

🌎 My Website - [The Neon Coder](https://bankole2000.github.io/webpieces)

💼 Project Link: [https://github.com/bankole2000/pocketbase-rabbitmq-extension](https://github.com/bankole2000/pocketbase-rabbitmq-extension)

<!-- ACKNOWLEDGEMENTS
## Acknowledgements

- [GitHub Emoji Cheat Sheet](https://www.webpagefx.com/tools/emoji-cheat-sheet)
- [Img Shields](https://shields.io)
- [Choose an Open Source License](https://choosealicense.com)
- [GitHub Pages](https://pages.github.com)
- [Animate.css](https://daneden.github.io/animate.css)
- [Loaders.css](https://connoratherton.com/loaders)
- [Slick Carousel](https://kenwheeler.github.io/slick)
- [Smooth Scroll](https://github.com/cferdinandi/smooth-scroll)
- [Sticky Kit](http://leafo.net/sticky-kit)
- [JVectorMap](http://jvectormap.com)
- [Font Awesome](https://fontawesome.com) -->

[demo-gif]: ./media/publisher.gif
