package extra

import (
	"fmt"
	"github.com/gookit/color"
)

func ShowLogo() {
	logo := "\n███████╗████████╗ █████╗ ██████╗     ██╗      █████╗ ██████╗ ███████╗\n██╔════╝╚══██╔══╝██╔══██╗██╔══██╗    ██║     ██╔══██╗██╔══██╗██╔════╝\n███████╗   ██║   ███████║██████╔╝    ██║     ███████║██████╔╝███████╗\n╚════██║   ██║   ██╔══██║██╔══██╗    ██║     ██╔══██║██╔══██╗╚════██║\n███████║   ██║   ██║  ██║██║  ██║    ███████╗██║  ██║██████╔╝███████║\n╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝    ╚══════╝╚═╝  ╚═╝╚═════╝ ╚══════╝\n                                                                     "
	color.Cyan.Printf(logo + "\n")

}

func ShowDevInfo() {
	fmt.Println("\033[36mVERSION: \033[33m1.0\033[33m")
	fmt.Println("\033[36mDEV: \033[33mhttps://t.me/StarLabsTech\033[33m")
	fmt.Println("\033[36mGitHub: \033[33mhttps://github.com/0xStarLabs/StarLabs-Twitter\033[33m")
	fmt.Println("\033[36mDONATION EVM ADDRESS: \033[33m0x620ea8b01607efdf3c74994391f86523acf6f9e1\033[0m")
	fmt.Println()
}
