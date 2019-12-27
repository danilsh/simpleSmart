all: _ventService _historyCaptureService

_ventService:
	$(MAKE) -C ventService

_historyCaptureService:
	$(MAKE) -C historyCaptureService
