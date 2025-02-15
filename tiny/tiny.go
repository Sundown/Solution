package tiny

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sundown/solution/prism"
)

func Entry(env *prism.Environment) {
	header := "from tinygrad import Tensor"
	header += `
a = Tensor.empty(4, 4)
b = Tensor.empty(4, 4)
print((a+b).tolist())`

	// final
	os.WriteFile("main.py", []byte(header), 0644)

	cmd := exec.Command("python3", "main.py")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "DEBUG=4")
	output, _ := cmd.CombinedOutput()
	fmt.Println(string(output))
}
