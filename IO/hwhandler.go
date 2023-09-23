package IO

//
//import (
//	"FireFighter/comm"
//)
//
//type HwHandler struct {
//	SensorsHandler  *SensorsHandler
//	MotorsHandler   *MotorsHandler
//	EncodersHandler *EncodersHandler
//}
//
//func New(handler *comm.Handler) *HwHandler {
//	h := &HwHandler{
//		SensorsHandler:  NewSensorsHandler(handler),
//		MotorsHandler:   NewMotorsHandler(handler),
//		EncodersHandler: NewEncodersHandler(handler),
//	}
//
//	go h.update()
//
//	return h
//}
//
//func (h *HwHandler) update() {
//	for {
//		h.SensorsHandler.Update()
//		h.MotorsHandler.Update()
//		h.EncodersHandler.Update()
//	}
//}
