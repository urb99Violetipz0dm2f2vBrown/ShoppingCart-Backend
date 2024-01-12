#!/bin/bash

# Loop to send 10 requests
#if you want to send more change {1..10} to whatever you wish
for i in {1..10}
do
    curl --location 'http://localhost:3000/api/book' \
    --header 'Content-Type: application/json' \
    --data '{
        "title": "Book'$i'",
        "author": "Author'$i'",
        "description": "Description'$i'",
        "price": '$((i * 10))',
        "genre": "Genre'$i'"
    }'
    
    # Add a delay between requests (optional)
    sleep 1
done