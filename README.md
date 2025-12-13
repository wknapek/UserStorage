Endpoints Implemented
#	Method	Endpoint	Description
1	GET	/users	List users
2	GET	/users/:id	Get user by ID
3	POST	/users	Create user → publish UserCreated
4	PUT	/users/:id	Update user → publish UserUpdated
5	DELETE	/users/:id	Delete user → publish UserDeleted
6	GET	/users/:id/files	Get user files
7	POST	/users/:id/files	Add file
8	DELETE	/users/:id/files	Delete all files
