package main

import (
	"log"

	"golang.org/x/sys/windows/svc"
)

func main() {
	winSvc, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("Failed to determine if we were a windows service")
	}

	if !winSvc {
		log.Fatalf("This program must be run as a windows service")
	}

	err = svc.Run("", &windowsService{})
	if err != nil {
		log.Fatalf("Failed to run service: %s", err)
	}

}

type windowsService struct{}

func (sh *windowsService) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (bool, uint32) {
	s <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}
	for {
		req := <-r
		switch req.Cmd {
		case svc.Interrogate:
			s <- req.CurrentStatus
		case svc.Stop, svc.Shutdown:
			return false, 0
		default:
			return false, 1052
		}
	}
}
