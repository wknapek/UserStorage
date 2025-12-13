Endpoints Implemented
#	Method	Endpoint	Description
- GET	/users	List users
-	GET	/users/:id	Get user by ID
-	POST	/users	Create user → publish UserCreated
-	PUT	/users/:id	Update user → publish UserUpdated
-	DELETE	/users/:id	Delete user → publish UserDeleted
-	GET	/users/:id/files	Get user files
-	POST	/users/:id/files	Add file
-	DELETE	/users/:id/files	Delete all files
