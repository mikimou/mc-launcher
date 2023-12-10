package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/hashicorp/go-getter"
)

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#243433"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	textc     = lipgloss.Color("#243433")

	Client = &http.Client{Timeout: 5 * time.Second}

	t_heading  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4")).PaddingTop(2).PaddingLeft(4).PaddingBottom(1)
	t_normal   = lipgloss.NewStyle().Padding(1)
	title      = lipgloss.NewStyle().Margin(3, 8, 0).BorderStyle(lipgloss.NormalBorder()).BorderForeground(subtle).BorderBottom(true)
	titleDesc  = lipgloss.NewStyle().Margin(0, 8, 3)
	divider    = lipgloss.NewStyle().SetString("•").Padding(0, 1).Foreground(subtle).String()
	url        = lipgloss.NewStyle().Foreground(special).Render
	txt_subtle = lipgloss.NewStyle().Foreground(subtle).Render
	//status     = lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" STATUS "), line.Render("Update check.."), check.Render(" OK "))
	stat     = lipgloss.NewStyle().Height(1).Foreground(lipgloss.Color("#D9DCCF")).Background(lipgloss.Color("#107869")).MarginLeft(2).MarginTop(0).MarginBottom(1)
	line     = lipgloss.NewStyle().SetString(" ").Width(60).Height(1).Foreground(lipgloss.AdaptiveColor{Light: "#43433", Dark: "#C1C6B2"}).Background(lipgloss.AdaptiveColor{Dark: "#353533"}).MarginTop(0).MarginBottom(1)
	lin      = lipgloss.NewStyle().SetString(" ").Height(1).Foreground(lipgloss.AdaptiveColor{Light: "#43433", Dark: "#C1C6B2"}).Background(lipgloss.AdaptiveColor{Dark: "#353533"}).MarginTop(0).MarginBottom(1)
	check    = lipgloss.NewStyle().Height(1).Foreground(textc).Background(special).MarginRight(2).MarginTop(0).MarginBottom(1)
	helptext = lipgloss.NewStyle().SetString(" ").Width(60).Height(1).Foreground(lipgloss.AdaptiveColor{Light: "#43433", Dark: "#C1C6B2"}).Background(lipgloss.AdaptiveColor{Dark: "#353533"})
)

type Config struct {
	username string
	version  string
}

var config Config

func main() {
	checkVer()
	fmt.Println(title.Render("Mikiho custom launcher for moded/enhanced minecraft clients"))
	fmt.Println(titleDesc.Render("Michal Hicz" + divider + url("https://github.com/mikimou") + divider + txt_subtle(config.version)))
	if config.version == "error" {
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" VER "), line.Render("Error loading version"), check.Render(" ERR ")))
	}
	checkUpdate("https://api.github.com/repos/mikimou/mc-launcher/releases/latest")
	time.Sleep(500 * time.Millisecond)
	loadUser()
	time.Sleep(1500 * time.Millisecond)
	if runtime.GOOS != "windows" {
		log.Fatal("unsupported os")
	} else {
		runWin(config.username)
	}
}

func loadUser() {
	user, err := os.ReadFile("username.txt")
	if err != nil {
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" CONFIG "), line.Render("Nenasiel sa ziadny nick!"), check.Render(" OK ")))
	}
	if string(user) != "" {
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" CONFIG "), lipgloss.JoinHorizontal(lipgloss.Top, lin.Render("nick: "+string(user)), lin.Align(lipgloss.Right).PaddingRight(1).PaddingLeft(16).Render("zmen nick ->"), check.Render(" username.txt "))))
	} else {
		setNick()
	}
	config.username = string(user)
}

func setNick() {
	fmt.Print("  Zadaj nick -> ")
	var nick string
	n, err := fmt.Scanln(&nick)
	fmt.Println()
	if n != 1 {
		log.Fatal("zly nick!")
	}
	if err != nil {
		log.Fatal("zly nick!")
	}
	config.username = nick
	err = os.WriteFile("username.txt", []byte(nick), 0644)
	if err != nil {
		f, err := os.Create("username.txt")
		if err != nil {
			fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" CONFIG "), line.Render("Ukladanie nastaveni zlyhalo!"), check.Render(" ERR ")))
		}
		defer f.Close()
		nbytes, err := f.WriteString(nick)
		if err != nil {
			fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" CONFIG "), line.Render("Nick nebol ulozeny!"), check.Render(" ERR ")))
			if nbytes == 0 {
			}
		}
	}
}

func checkVer() {
	var version string
	ver, err := os.ReadFile("launcher.version")
	if ver != nil {
		version = string(ver)
	} else {
		version = "error"
	}
	if err != nil {
		version = "error"
	}
	config.version = version
}

func setUpdatedVersion(ver string) {
	err := os.WriteFile("launcher.version", []byte(ver), 0644)
	if err != nil {
	}
}

func checkUpdate(url string) {
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" STATUS "), line.Render("Update check.."), check.Render(" OK ")))
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" STATUS "), line.Render("Ziadny internet!"), check.Render(" ERR ")))
		}
	}()
	resp, err := Client.Get(url)
	if err != nil {
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" STATUS "), line.Render("Server neodpoveda!"), check.Render(" ERR ")))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	var jsonMap map[string]string
	json.Unmarshal([]byte(body), &jsonMap)

	if config.version != "error" {
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" VER "), line.Render("Aktualna verzia: "+config.version), check.Render(" OK ")))
	}

	if jsonMap["name"] != config.version {
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" VER "), line.Render("DOSTUPNA NOVA VERZIA!"), check.Render(" OK ")))
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" VER "), line.Render("Stahujem a instalujem"), check.Render(" OK ")))
		installUpdate(jsonMap["name"])
		setUpdatedVersion(jsonMap["name"])
	}
}

func installUpdate(ver string) {
	client := getter.Client{DisableSymlinks: true}
	client.Dst = "."
	client.Dir = true
	client.Src = "https://github.com/mikimou/mc-launcher/releases/download/" + ver + "/update.zip"
	err := client.Get()
	if err != nil {
	}
}

func runWin(nick string) {
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" RUN "), line.Width(50).Render("Client sa spusta.."), check.Render(" MODED ")))
	comm, err := os.ReadFile("launcher.command")
	if err != nil {
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" RUN "), line.Render("Chyba v spustaci, preinstaluj alebo stiahni novu verziu!"), check.Render(" ERR ")))
		time.Sleep(4000 * time.Millisecond)
	}
	if comm != nil {
		args := strings.Split(strings.Replace(string(comm), "replaceme", nick, 1), " ")
		go exec.Command(args[0], args[1:]...).Output()
		time.Sleep(6000 * time.Millisecond)
	} else {
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" RUN "), line.Render("Chyba v spustaci, preinstaluj alebo stiahni novu verziu!"), check.Render(" ERR ")))
		time.Sleep(4000 * time.Millisecond)
	}
}

func runUnix(nick string) {
	comm, err := os.ReadFile("launcher.command")
	if err != nil {
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" RUN "), line.Render("Chyba v spustaci, preinstaluj alebo stiahni novu verziu!"), check.Render(" ERR ")))
		time.Sleep(4000 * time.Millisecond)
	}
	if comm != nil {
		javaPath, err := exec.LookPath("java")
		if err != nil {
			fmt.Print("  Zadaj umiestnenie javy (full path) -> ")
			n, err := fmt.Scanln(&javaPath)
			fmt.Println()
			if n != 1 {
				log.Fatal("zla cesta!")
				javaPath = "java"
			}
			if err != nil {
				log.Fatal("zla cesta!")
				javaPath = "java"
			}
		}

		args := strings.Split(strings.Replace(string(comm), "replaceme", nick, 1), " ")
		args[0] = javaPath
		go exec.Command(args[0], args[1:]...).Output()
		time.Sleep(6000 * time.Millisecond)
	} else {
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" RUN "), line.Render("Chyba v spustaci, preinstaluj alebo stiahni novu verziu!"), check.Render(" ERR ")))
		time.Sleep(4000 * time.Millisecond)
	}
}
