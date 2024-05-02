# CRUD Backend Using MySQL in GoLang

# Get All Users
curl -X GET http://localhost:8000/users

# Get User by ID (replace {id} with actual user ID)
curl -X GET http://localhost:8000/users/{id}

# Create User
curl -X POST -H "Content-Type: application/json" -d '{"name":"John Doe","age":30}' http://localhost:8000/users

# Update User (replace {id} with actual user ID)
curl -X PUT -H "Content-Type: application/json" -d '{"name":"Jane Doe","age":35}' http://localhost:8000/users/{id}

# Delete User (replace {id} with actual user ID)
curl -X DELETE http://localhost:8000/users/{id}
