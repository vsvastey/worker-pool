package main

import "fmt"

func showProgress(wps []*WorkerAndProgress, done <-chan struct{}) {
	for i := 0; i < len(wps); i++ {
		fmt.Println(wps[i].pb.Draw())
	}

	updated := make(chan struct{})
	for i := 0; i < len(wps); i++ {
		go func(wp *WorkerAndProgress) {
			for st := range wp.worker.Status() {
				wp.pb.Set(fmt.Sprintf("%s - %s", st.ID, st.Task), st.Progress)
				updated <- struct{}{}
			}
		}(wps[i])
	}

	for {
		select {
		case <-updated:
			fmt.Print("\033[", len(wps), "F")
			for i := 0; i < len(wps); i++ {
				fmt.Println(wps[i].pb.Draw())
			}
		case <-done:
			for i := 0; i < len(wps); i++ {
				wps[i].worker.Stop()
			}
			return
		}
	}
}
