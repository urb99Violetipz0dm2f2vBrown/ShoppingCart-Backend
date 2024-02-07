# Getting Started with BackEnd



## Available Scripts

In the project directory, you can run:

### `docker compose up`

This will run the app and run on localhost:3000

### `docker compose down`

This will tear down the docker container

### `./send_requests.sh`

This will run the send_requests.sh which will generate book data for you. You have to make it executable with `chmod +x send_requests.sh` and then run it using `./send_requests.sh`

### `Available endpoints`
# Books

api.Get("/books/:title/:author/:genre", handlers.ListBooks)

api.Post("/book", handlers.CreateBook)

api.Put("/book/:id", handlers.EditBook)

# Shopping Cart

api.Get("/cart", handlers.ListCart)

api.Put("/cart/:id", handlers.AddToCart)

api.Delete("/cart/:id", handlers.RemoveFromCart)

