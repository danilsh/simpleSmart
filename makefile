all: _ventService _historyCaptureService
	@echo "***** Successfull Build *****"

_ventService:
	@$(MAKE) -C ventService

_historyCaptureService:
	@$(MAKE) -C historyCaptureService
