mocks:
	mockgen -source=queueHandler/queueHandler.go -destination=queueHandler/mock_queueHandler.go -package=queueHandler
	mockgen -source=dbhandler/handler.go -destination=dbhandler/mock_handler.go -package=dbhandler
