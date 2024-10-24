package mpv

import (
	"fmt"
	"os/exec"
)

func Play(id string) error {
	arg := fmt.Sprintf("https://www.youtube.com/watch?v=%v", id)
	mpv := exec.Command("mpv", arg)
	if err := mpv.Run(); err != nil {
		return err
	}
	return nil

	// ipcc := mpv.NewIPCClient("/tmp/mpvsocket") // Lowlevel client
	// c := mpv.NewClient(ipcc)                   // Highlevel client, can also use RPCClient
	//
	// c.Loadfile(arg, mpv.LoadFileModeReplace)
	// c.SetPause(true)
	// c.Seek(600, mpv.SeekModeAbsolute)
	// c.SetFullscreen(true)
	// c.SetPause(false)
	//
	// pos, err := c.Position()
	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("Position in Seconds: %.0f", pos)
	// return nil
}
