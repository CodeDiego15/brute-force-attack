package attack

import (
	"fmt"
	"os/exec"
)

func main() {
	fmt.Println(".")
	exec.Command("ls")
}

// Tomorrow we start.
