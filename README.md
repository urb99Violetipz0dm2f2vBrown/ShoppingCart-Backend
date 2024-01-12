# Getting Started with BackEnd



## Available Scripts

In the project directory, you can run:

### `docker compose up`

This will run the app and run on localhost:3000

### `docker compose down`

This will tear down the docker container


### `Available endpoints`
# Books

api.Get("/books/:title/:author/:genre", handlers.ListBooks)

api.Post("/book", handlers.CreateBook)

# Shopping Cart

api.Get("/cart", handlers.ListCart)

api.Put("/cart/:id", handlers.AddToCart)

api.Delete("/cart/:id", handlers.RemoveFromCart)

