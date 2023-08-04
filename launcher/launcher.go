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
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" CONFIG "), line.Render("No username setting found!"), check.Render(" OK ")))
	}
	if string(user) != "" {
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" CONFIG "), lipgloss.JoinHorizontal(lipgloss.Top, lin.Render("nick: "+config.username), lin.Align(lipgloss.Right).PaddingRight(1).PaddingLeft(16).Render("zmen nick ->"), check.Render(" username.txt "))))
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
			fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" CONFIG "), line.Render("Saving config error.."), check.Render(" ERR ")))
		}
		defer f.Close()
		nbytes, err := f.WriteString(nick)
		if err != nil {
			fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" CONFIG "), line.Render("Username not saved!"), check.Render(" ERR ")))
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

func checkUpdate(url string) {
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" STATUS "), line.Render("Update check.."), check.Render(" OK ")))

	resp, err := Client.Get(url)
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(body), &jsonMap)

	if jsonMap["name"] != config.version {
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" VER "), line.Render("DOSTUPNA NOVA VERZIA!"), check.Render(" OK ")))
	}

	if config.version != "error" {
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" VER "), line.Render("Aktualna verzia: "+config.version), check.Render(" OK ")))
	}
}

func runWin(nick string) {
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" RUN "), line.Width(50).Render("Client sa spusta.."), check.Render(" TY RAKLO! ")))
	comm := "jvm/java-runtime-gamma/bin/javaw.exe -XX:HeapDumpPath=MojangTricksIntelDriversForPerformance_javaw.exe_minecraft.exe.heapdump -Djava.library.path=.\\versions\\forge-1.20.1-47.1.0\\natives -Dminecraft.launcher.brand=mikiho-launcher -Djna.tmpdir=.\\versions\\forge-1.20.1-47.1.0\\natives -Dorg.lwjgl.system.SharedLibraryExtractPath=.\\versions\\forge-1.20.1-47.1.0\\natives -Dio.netty.native.workdir=.\\versions\\forge-1.20.1-47.1.0\\natives -Dminecraft.launcher.brand=minecraft-launcher-lib -Dminecraft.launcher.version=6.1 -cp .\\libraries\\cpw\\mods\\securejarhandler\\2.1.10\\securejarhandler-2.1.10.jar;.\\libraries\\org\\ow2\\asm\\asm\\9.5\\asm-9.5.jar;.\\libraries\\org\\ow2\\asm\\asm-commons\\9.5\\asm-commons-9.5.jar;.\\libraries\\org\\ow2\\asm\\asm-tree\\9.5\\asm-tree-9.5.jar;.\\libraries\\org\\ow2\\asm\\asm-util\\9.5\\asm-util-9.5.jar;.\\libraries\\org\\ow2\\asm\\asm-analysis\\9.5\\asm-analysis-9.5.jar;.\\libraries\\net\\minecraftforge\\accesstransformers\\8.0.4\\accesstransformers-8.0.4.jar;.\\libraries\\org\\antlr\\antlr4-runtime\\4.9.1\\antlr4-runtime-4.9.1.jar;.\\libraries\\net\\minecraftforge\\eventbus\\6.0.5\\eventbus-6.0.5.jar;.\\libraries\\net\\minecraftforge\\forgespi\\7.0.1\\forgespi-7.0.1.jar;.\\libraries\\net\\minecraftforge\\coremods\\5.0.1\\coremods-5.0.1.jar;.\\libraries\\cpw\\mods\\modlauncher\\10.0.9\\modlauncher-10.0.9.jar;.\\libraries\\net\\minecraftforge\\unsafe\\0.2.0\\unsafe-0.2.0.jar;.\\libraries\\net\\minecraftforge\\mergetool\\1.1.5\\mergetool-1.1.5-api.jar;.\\libraries\\com\\electronwill\\night-config\\core\\3.6.4\\core-3.6.4.jar;.\\libraries\\com\\electronwill\\night-config\\toml\\3.6.4\\toml-3.6.4.jar;.\\libraries\\org\\apache\\maven\\maven-artifact\\3.8.5\\maven-artifact-3.8.5.jar;.\\libraries\\net\\jodah\\typetools\\0.6.3\\typetools-0.6.3.jar;.\\libraries\\net\\minecrell\\terminalconsoleappender\\1.2.0\\terminalconsoleappender-1.2.0.jar;.\\libraries\\org\\jline\\jline-reader\\3.12.1\\jline-reader-3.12.1.jar;.\\libraries\\org\\jline\\jline-terminal\\3.12.1\\jline-terminal-3.12.1.jar;.\\libraries\\org\\spongepowered\\mixin\\0.8.5\\mixin-0.8.5.jar;.\\libraries\\org\\openjdk\\nashorn\\nashorn-core\\15.3\\nashorn-core-15.3.jar;.\\libraries\\net\\minecraftforge\\JarJarSelector\\0.3.19\\JarJarSelector-0.3.19.jar;.\\libraries\\net\\minecraftforge\\JarJarMetadata\\0.3.19\\JarJarMetadata-0.3.19.jar;.\\libraries\\cpw\\mods\\bootstraplauncher\\1.1.2\\bootstraplauncher-1.1.2.jar;.\\libraries\\net\\minecraftforge\\JarJarFileSystems\\0.3.19\\JarJarFileSystems-0.3.19.jar;.\\libraries\\net\\minecraftforge\\fmlloader\\1.20.1-47.1.0\\fmlloader-1.20.1-47.1.0.jar;.\\libraries\\net\\minecraftforge\\fmlearlydisplay\\1.20.1-47.1.0\\fmlearlydisplay-1.20.1-47.1.0.jar;.\\libraries\\com\\github\\oshi\\oshi-core\\6.2.2\\oshi-core-6.2.2.jar;.\\libraries\\com\\google\\code\\gson\\gson\\2.10\\gson-2.10.jar;.\\libraries\\com\\google\\guava\\failureaccess\\1.0.1\\failureaccess-1.0.1.jar;.\\libraries\\com\\google\\guava\\guava\\31.1-jre\\guava-31.1-jre.jar;.\\libraries\\com\\ibm\\icu\\icu4j\\71.1\\icu4j-71.1.jar;.\\libraries\\com\\mojang\\authlib\\4.0.43\\authlib-4.0.43.jar;.\\libraries\\com\\mojang\\blocklist\\1.0.10\\blocklist-1.0.10.jar;.\\libraries\\com\\mojang\\brigadier\\1.1.8\\brigadier-1.1.8.jar;.\\libraries\\com\\mojang\\datafixerupper\\6.0.8\\datafixerupper-6.0.8.jar;.\\libraries\\com\\mojang\\logging\\1.1.1\\logging-1.1.1.jar;.\\libraries\\com\\mojang\\patchy\\2.2.10\\patchy-2.2.10.jar;.\\libraries\\com\\mojang\\text2speech\\1.17.9\\text2speech-1.17.9.jar;.\\libraries\\commons-codec\\commons-codec\\1.15\\commons-codec-1.15.jar;.\\libraries\\commons-io\\commons-io\\2.11.0\\commons-io-2.11.0.jar;.\\libraries\\commons-logging\\commons-logging\\1.2\\commons-logging-1.2.jar;.\\libraries\\io\\netty\\netty-buffer\\4.1.82.Final\\netty-buffer-4.1.82.Final.jar;.\\libraries\\io\\netty\\netty-codec\\4.1.82.Final\\netty-codec-4.1.82.Final.jar;.\\libraries\\io\\netty\\netty-common\\4.1.82.Final\\netty-common-4.1.82.Final.jar;.\\libraries\\io\\netty\\netty-handler\\4.1.82.Final\\netty-handler-4.1.82.Final.jar;.\\libraries\\io\\netty\\netty-resolver\\4.1.82.Final\\netty-resolver-4.1.82.Final.jar;.\\libraries\\io\\netty\\netty-transport-classes-epoll\\4.1.82.Final\\netty-transport-classes-epoll-4.1.82.Final.jar;.\\libraries\\io\\netty\\netty-transport-native-unix-common\\4.1.82.Final\\netty-transport-native-unix-common-4.1.82.Final.jar;.\\libraries\\io\\netty\\netty-transport\\4.1.82.Final\\netty-transport-4.1.82.Final.jar;.\\libraries\\it\\unimi\\dsi\\fastutil\\8.5.9\\fastutil-8.5.9.jar;.\\libraries\\net\\java\\dev\\jna\\jna-platform\\5.12.1\\jna-platform-5.12.1.jar;.\\libraries\\net\\java\\dev\\jna\\jna\\5.12.1\\jna-5.12.1.jar;.\\libraries\\net\\sf\\jopt-simple\\jopt-simple\\5.0.4\\jopt-simple-5.0.4.jar;.\\libraries\\org\\apache\\commons\\commons-compress\\1.21\\commons-compress-1.21.jar;.\\libraries\\org\\apache\\commons\\commons-lang3\\3.12.0\\commons-lang3-3.12.0.jar;.\\libraries\\org\\apache\\httpcomponents\\httpclient\\4.5.13\\httpclient-4.5.13.jar;.\\libraries\\org\\apache\\httpcomponents\\httpcore\\4.4.15\\httpcore-4.4.15.jar;.\\libraries\\org\\apache\\logging\\log4j\\log4j-api\\2.19.0\\log4j-api-2.19.0.jar;.\\libraries\\org\\apache\\logging\\log4j\\log4j-core\\2.19.0\\log4j-core-2.19.0.jar;.\\libraries\\org\\apache\\logging\\log4j\\log4j-slf4j2-impl\\2.19.0\\log4j-slf4j2-impl-2.19.0.jar;.\\libraries\\org\\joml\\joml\\1.10.5\\joml-1.10.5.jar;.\\libraries\\org\\lwjgl\\lwjgl-glfw\\3.3.1\\lwjgl-glfw-3.3.1.jar;.\\libraries\\org\\lwjgl\\lwjgl-glfw\\3.3.1\\lwjgl-glfw-3.3.1-natives-windows.jar;.\\libraries\\org\\lwjgl\\lwjgl-glfw\\3.3.1\\lwjgl-glfw-3.3.1-natives-windows-arm64.jar;.\\libraries\\org\\lwjgl\\lwjgl-glfw\\3.3.1\\lwjgl-glfw-3.3.1-natives-windows-x86.jar;.\\libraries\\org\\lwjgl\\lwjgl-jemalloc\\3.3.1\\lwjgl-jemalloc-3.3.1.jar;.\\libraries\\org\\lwjgl\\lwjgl-jemalloc\\3.3.1\\lwjgl-jemalloc-3.3.1-natives-windows.jar;.\\libraries\\org\\lwjgl\\lwjgl-jemalloc\\3.3.1\\lwjgl-jemalloc-3.3.1-natives-windows-arm64.jar;.\\libraries\\org\\lwjgl\\lwjgl-jemalloc\\3.3.1\\lwjgl-jemalloc-3.3.1-natives-windows-x86.jar;.\\libraries\\org\\lwjgl\\lwjgl-openal\\3.3.1\\lwjgl-openal-3.3.1.jar;.\\libraries\\org\\lwjgl\\lwjgl-openal\\3.3.1\\lwjgl-openal-3.3.1-natives-windows.jar;.\\libraries\\org\\lwjgl\\lwjgl-openal\\3.3.1\\lwjgl-openal-3.3.1-natives-windows-arm64.jar;.\\libraries\\org\\lwjgl\\lwjgl-openal\\3.3.1\\lwjgl-openal-3.3.1-natives-windows-x86.jar;.\\libraries\\org\\lwjgl\\lwjgl-opengl\\3.3.1\\lwjgl-opengl-3.3.1.jar;.\\libraries\\org\\lwjgl\\lwjgl-opengl\\3.3.1\\lwjgl-opengl-3.3.1-natives-windows.jar;.\\libraries\\org\\lwjgl\\lwjgl-opengl\\3.3.1\\lwjgl-opengl-3.3.1-natives-windows-arm64.jar;.\\libraries\\org\\lwjgl\\lwjgl-opengl\\3.3.1\\lwjgl-opengl-3.3.1-natives-windows-x86.jar;.\\libraries\\org\\lwjgl\\lwjgl-stb\\3.3.1\\lwjgl-stb-3.3.1.jar;.\\libraries\\org\\lwjgl\\lwjgl-stb\\3.3.1\\lwjgl-stb-3.3.1-natives-windows.jar;.\\libraries\\org\\lwjgl\\lwjgl-stb\\3.3.1\\lwjgl-stb-3.3.1-natives-windows-arm64.jar;.\\libraries\\org\\lwjgl\\lwjgl-stb\\3.3.1\\lwjgl-stb-3.3.1-natives-windows-x86.jar;.\\libraries\\org\\lwjgl\\lwjgl-tinyfd\\3.3.1\\lwjgl-tinyfd-3.3.1.jar;.\\libraries\\org\\lwjgl\\lwjgl-tinyfd\\3.3.1\\lwjgl-tinyfd-3.3.1-natives-windows.jar;.\\libraries\\org\\lwjgl\\lwjgl-tinyfd\\3.3.1\\lwjgl-tinyfd-3.3.1-natives-windows-arm64.jar;.\\libraries\\org\\lwjgl\\lwjgl-tinyfd\\3.3.1\\lwjgl-tinyfd-3.3.1-natives-windows-x86.jar;.\\libraries\\org\\lwjgl\\lwjgl\\3.3.1\\lwjgl-3.3.1.jar;.\\libraries\\org\\lwjgl\\lwjgl\\3.3.1\\lwjgl-3.3.1-natives-windows.jar;.\\libraries\\org\\lwjgl\\lwjgl\\3.3.1\\lwjgl-3.3.1-natives-windows-arm64.jar;.\\libraries\\org\\lwjgl\\lwjgl\\3.3.1\\lwjgl-3.3.1-natives-windows-x86.jar;.\\libraries\\org\\slf4j\\slf4j-api\\2.0.1\\slf4j-api-2.0.1.jar;.\\versions\\forge-1.20.1-47.1.0\\forge-1.20.1-47.1.0.jar -Djava.net.preferIPv6Addresses=system -DignoreList=bootstraplauncher,securejarhandler,asm-commons,asm-util,asm-analysis,asm-tree,asm,JarJarFileSystems,client-extra,fmlcore,javafmllanguage,lowcodelanguage,mclanguage,forge-,forge-1.20.1-47.1.0.jar -DmergeModules=jna-5.10.0.jar,jna-platform-5.10.0.jar -DlibraryDirectory=.\\libraries -p .\\libraries/cpw/mods/bootstraplauncher/1.1.2/bootstraplauncher-1.1.2.jar;.\\libraries/cpw/mods/securejarhandler/2.1.10/securejarhandler-2.1.10.jar;.\\libraries/org/ow2/asm/asm-commons/9.5/asm-commons-9.5.jar;.\\libraries/org/ow2/asm/asm-util/9.5/asm-util-9.5.jar;.\\libraries/org/ow2/asm/asm-analysis/9.5/asm-analysis-9.5.jar;.\\libraries/org/ow2/asm/asm-tree/9.5/asm-tree-9.5.jar;.\\libraries/org/ow2/asm/asm/9.5/asm-9.5.jar;.\\libraries/net/minecraftforge/JarJarFileSystems/0.3.19/JarJarFileSystems-0.3.19.jar --add-modules ALL-MODULE-PATH --add-opens java.base/java.util.jar=cpw.mods.securejarhandler --add-opens java.base/java.lang.invoke=cpw.mods.securejarhandler --add-exports java.base/sun.security.util=cpw.mods.securejarhandler --add-exports jdk.naming.dns/com.sun.jndi.dns=java.naming cpw.mods.bootstraplauncher.BootstrapLauncher --username " + nick + " --version forge-1.20.1-47.1.0 --gameDir . --assetsDir .\\assets --assetIndex 5 --uuid - --accessToken - --clientId ${clientid} --xuid ${auth_xuid} --userType msa --versionType release --launchTarget forgeclient --fml.forgeVersion 47.1.0 --fml.mcVersion 1.20.1 --fml.forgeGroup net.minecraftforge --fml.mcpVersion 20230612.114412"
	args := strings.Split(comm, " ")
	go exec.Command(args[0], args[1:]...).Output()
	time.Sleep(5000 * time.Millisecond)
}

func runUnix(nick string) {
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stat.Render(" RUN "), line.Width(50).Render("Client sa spusta.."), check.Render(" TY RAKLO! ")))
	// FINISH
	time.Sleep(300000 * time.Millisecond)
	comm := "java -XX:HeapDumpPath=MojangTricksIntelDriversForPerformance_javaw.exe_minecraft.exe.heapdump -Djava.library.path=.\\versions\\forge-1.20.1-47.1.0\\natives -Dminecraft.launcher.brand=mikiho-launcher -Djna.tmpdir=.\\versions\\forge-1.20.1-47.1.0\\natives -Dorg.lwjgl.system.SharedLibraryExtractPath=.\\versions\\forge-1.20.1-47.1.0\\natives -Dio.netty.native.workdir=.\\versions\\forge-1.20.1-47.1.0\\natives -Dminecraft.launcher.brand=minecraft-launcher-lib -Dminecraft.launcher.version=6.1 -cp .\\libraries\\cpw\\mods\\securejarhandler\\2.1.10\\securejarhandler-2.1.10.jar;.\\libraries\\org\\ow2\\asm\\asm\\9.5\\asm-9.5.jar;.\\libraries\\org\\ow2\\asm\\asm-commons\\9.5\\asm-commons-9.5.jar;.\\libraries\\org\\ow2\\asm\\asm-tree\\9.5\\asm-tree-9.5.jar;.\\libraries\\org\\ow2\\asm\\asm-util\\9.5\\asm-util-9.5.jar;.\\libraries\\org\\ow2\\asm\\asm-analysis\\9.5\\asm-analysis-9.5.jar;.\\libraries\\net\\minecraftforge\\accesstransformers\\8.0.4\\accesstransformers-8.0.4.jar;.\\libraries\\org\\antlr\\antlr4-runtime\\4.9.1\\antlr4-runtime-4.9.1.jar;.\\libraries\\net\\minecraftforge\\eventbus\\6.0.5\\eventbus-6.0.5.jar;.\\libraries\\net\\minecraftforge\\forgespi\\7.0.1\\forgespi-7.0.1.jar;.\\libraries\\net\\minecraftforge\\coremods\\5.0.1\\coremods-5.0.1.jar;.\\libraries\\cpw\\mods\\modlauncher\\10.0.9\\modlauncher-10.0.9.jar;.\\libraries\\net\\minecraftforge\\unsafe\\0.2.0\\unsafe-0.2.0.jar;.\\libraries\\net\\minecraftforge\\mergetool\\1.1.5\\mergetool-1.1.5-api.jar;.\\libraries\\com\\electronwill\\night-config\\core\\3.6.4\\core-3.6.4.jar;.\\libraries\\com\\electronwill\\night-config\\toml\\3.6.4\\toml-3.6.4.jar;.\\libraries\\org\\apache\\maven\\maven-artifact\\3.8.5\\maven-artifact-3.8.5.jar;.\\libraries\\net\\jodah\\typetools\\0.6.3\\typetools-0.6.3.jar;.\\libraries\\net\\minecrell\\terminalconsoleappender\\1.2.0\\terminalconsoleappender-1.2.0.jar;.\\libraries\\org\\jline\\jline-reader\\3.12.1\\jline-reader-3.12.1.jar;.\\libraries\\org\\jline\\jline-terminal\\3.12.1\\jline-terminal-3.12.1.jar;.\\libraries\\org\\spongepowered\\mixin\\0.8.5\\mixin-0.8.5.jar;.\\libraries\\org\\openjdk\\nashorn\\nashorn-core\\15.3\\nashorn-core-15.3.jar;.\\libraries\\net\\minecraftforge\\JarJarSelector\\0.3.19\\JarJarSelector-0.3.19.jar;.\\libraries\\net\\minecraftforge\\JarJarMetadata\\0.3.19\\JarJarMetadata-0.3.19.jar;.\\libraries\\cpw\\mods\\bootstraplauncher\\1.1.2\\bootstraplauncher-1.1.2.jar;.\\libraries\\net\\minecraftforge\\JarJarFileSystems\\0.3.19\\JarJarFileSystems-0.3.19.jar;.\\libraries\\net\\minecraftforge\\fmlloader\\1.20.1-47.1.0\\fmlloader-1.20.1-47.1.0.jar;.\\libraries\\net\\minecraftforge\\fmlearlydisplay\\1.20.1-47.1.0\\fmlearlydisplay-1.20.1-47.1.0.jar;.\\libraries\\com\\github\\oshi\\oshi-core\\6.2.2\\oshi-core-6.2.2.jar;.\\libraries\\com\\google\\code\\gson\\gson\\2.10\\gson-2.10.jar;.\\libraries\\com\\google\\guava\\failureaccess\\1.0.1\\failureaccess-1.0.1.jar;.\\libraries\\com\\google\\guava\\guava\\31.1-jre\\guava-31.1-jre.jar;.\\libraries\\com\\ibm\\icu\\icu4j\\71.1\\icu4j-71.1.jar;.\\libraries\\com\\mojang\\authlib\\4.0.43\\authlib-4.0.43.jar;.\\libraries\\com\\mojang\\blocklist\\1.0.10\\blocklist-1.0.10.jar;.\\libraries\\com\\mojang\\brigadier\\1.1.8\\brigadier-1.1.8.jar;.\\libraries\\com\\mojang\\datafixerupper\\6.0.8\\datafixerupper-6.0.8.jar;.\\libraries\\com\\mojang\\logging\\1.1.1\\logging-1.1.1.jar;.\\libraries\\com\\mojang\\patchy\\2.2.10\\patchy-2.2.10.jar;.\\libraries\\com\\mojang\\text2speech\\1.17.9\\text2speech-1.17.9.jar;.\\libraries\\commons-codec\\commons-codec\\1.15\\commons-codec-1.15.jar;.\\libraries\\commons-io\\commons-io\\2.11.0\\commons-io-2.11.0.jar;.\\libraries\\commons-logging\\commons-logging\\1.2\\commons-logging-1.2.jar;.\\libraries\\io\\netty\\netty-buffer\\4.1.82.Final\\netty-buffer-4.1.82.Final.jar;.\\libraries\\io\\netty\\netty-codec\\4.1.82.Final\\netty-codec-4.1.82.Final.jar;.\\libraries\\io\\netty\\netty-common\\4.1.82.Final\\netty-common-4.1.82.Final.jar;.\\libraries\\io\\netty\\netty-handler\\4.1.82.Final\\netty-handler-4.1.82.Final.jar;.\\libraries\\io\\netty\\netty-resolver\\4.1.82.Final\\netty-resolver-4.1.82.Final.jar;.\\libraries\\io\\netty\\netty-transport-classes-epoll\\4.1.82.Final\\netty-transport-classes-epoll-4.1.82.Final.jar;.\\libraries\\io\\netty\\netty-transport-native-unix-common\\4.1.82.Final\\netty-transport-native-unix-common-4.1.82.Final.jar;.\\libraries\\io\\netty\\netty-transport\\4.1.82.Final\\netty-transport-4.1.82.Final.jar;.\\libraries\\it\\unimi\\dsi\\fastutil\\8.5.9\\fastutil-8.5.9.jar;.\\libraries\\net\\java\\dev\\jna\\jna-platform\\5.12.1\\jna-platform-5.12.1.jar;.\\libraries\\net\\java\\dev\\jna\\jna\\5.12.1\\jna-5.12.1.jar;.\\libraries\\net\\sf\\jopt-simple\\jopt-simple\\5.0.4\\jopt-simple-5.0.4.jar;.\\libraries\\org\\apache\\commons\\commons-compress\\1.21\\commons-compress-1.21.jar;.\\libraries\\org\\apache\\commons\\commons-lang3\\3.12.0\\commons-lang3-3.12.0.jar;.\\libraries\\org\\apache\\httpcomponents\\httpclient\\4.5.13\\httpclient-4.5.13.jar;.\\libraries\\org\\apache\\httpcomponents\\httpcore\\4.4.15\\httpcore-4.4.15.jar;.\\libraries\\org\\apache\\logging\\log4j\\log4j-api\\2.19.0\\log4j-api-2.19.0.jar;.\\libraries\\org\\apache\\logging\\log4j\\log4j-core\\2.19.0\\log4j-core-2.19.0.jar;.\\libraries\\org\\apache\\logging\\log4j\\log4j-slf4j2-impl\\2.19.0\\log4j-slf4j2-impl-2.19.0.jar;.\\libraries\\org\\joml\\joml\\1.10.5\\joml-1.10.5.jar;.\\libraries\\org\\lwjgl\\lwjgl-glfw\\3.3.1\\lwjgl-glfw-3.3.1.jar;.\\libraries\\org\\lwjgl\\lwjgl-glfw\\3.3.1\\lwjgl-glfw-3.3.1-natives-windows.jar;.\\libraries\\org\\lwjgl\\lwjgl-glfw\\3.3.1\\lwjgl-glfw-3.3.1-natives-windows-arm64.jar;.\\libraries\\org\\lwjgl\\lwjgl-glfw\\3.3.1\\lwjgl-glfw-3.3.1-natives-windows-x86.jar;.\\libraries\\org\\lwjgl\\lwjgl-jemalloc\\3.3.1\\lwjgl-jemalloc-3.3.1.jar;.\\libraries\\org\\lwjgl\\lwjgl-jemalloc\\3.3.1\\lwjgl-jemalloc-3.3.1-natives-windows.jar;.\\libraries\\org\\lwjgl\\lwjgl-jemalloc\\3.3.1\\lwjgl-jemalloc-3.3.1-natives-windows-arm64.jar;.\\libraries\\org\\lwjgl\\lwjgl-jemalloc\\3.3.1\\lwjgl-jemalloc-3.3.1-natives-windows-x86.jar;.\\libraries\\org\\lwjgl\\lwjgl-openal\\3.3.1\\lwjgl-openal-3.3.1.jar;.\\libraries\\org\\lwjgl\\lwjgl-openal\\3.3.1\\lwjgl-openal-3.3.1-natives-windows.jar;.\\libraries\\org\\lwjgl\\lwjgl-openal\\3.3.1\\lwjgl-openal-3.3.1-natives-windows-arm64.jar;.\\libraries\\org\\lwjgl\\lwjgl-openal\\3.3.1\\lwjgl-openal-3.3.1-natives-windows-x86.jar;.\\libraries\\org\\lwjgl\\lwjgl-opengl\\3.3.1\\lwjgl-opengl-3.3.1.jar;.\\libraries\\org\\lwjgl\\lwjgl-opengl\\3.3.1\\lwjgl-opengl-3.3.1-natives-windows.jar;.\\libraries\\org\\lwjgl\\lwjgl-opengl\\3.3.1\\lwjgl-opengl-3.3.1-natives-windows-arm64.jar;.\\libraries\\org\\lwjgl\\lwjgl-opengl\\3.3.1\\lwjgl-opengl-3.3.1-natives-windows-x86.jar;.\\libraries\\org\\lwjgl\\lwjgl-stb\\3.3.1\\lwjgl-stb-3.3.1.jar;.\\libraries\\org\\lwjgl\\lwjgl-stb\\3.3.1\\lwjgl-stb-3.3.1-natives-windows.jar;.\\libraries\\org\\lwjgl\\lwjgl-stb\\3.3.1\\lwjgl-stb-3.3.1-natives-windows-arm64.jar;.\\libraries\\org\\lwjgl\\lwjgl-stb\\3.3.1\\lwjgl-stb-3.3.1-natives-windows-x86.jar;.\\libraries\\org\\lwjgl\\lwjgl-tinyfd\\3.3.1\\lwjgl-tinyfd-3.3.1.jar;.\\libraries\\org\\lwjgl\\lwjgl-tinyfd\\3.3.1\\lwjgl-tinyfd-3.3.1-natives-windows.jar;.\\libraries\\org\\lwjgl\\lwjgl-tinyfd\\3.3.1\\lwjgl-tinyfd-3.3.1-natives-windows-arm64.jar;.\\libraries\\org\\lwjgl\\lwjgl-tinyfd\\3.3.1\\lwjgl-tinyfd-3.3.1-natives-windows-x86.jar;.\\libraries\\org\\lwjgl\\lwjgl\\3.3.1\\lwjgl-3.3.1.jar;.\\libraries\\org\\lwjgl\\lwjgl\\3.3.1\\lwjgl-3.3.1-natives-windows.jar;.\\libraries\\org\\lwjgl\\lwjgl\\3.3.1\\lwjgl-3.3.1-natives-windows-arm64.jar;.\\libraries\\org\\lwjgl\\lwjgl\\3.3.1\\lwjgl-3.3.1-natives-windows-x86.jar;.\\libraries\\org\\slf4j\\slf4j-api\\2.0.1\\slf4j-api-2.0.1.jar;.\\versions\\forge-1.20.1-47.1.0\\forge-1.20.1-47.1.0.jar -Djava.net.preferIPv6Addresses=system -DignoreList=bootstraplauncher,securejarhandler,asm-commons,asm-util,asm-analysis,asm-tree,asm,JarJarFileSystems,client-extra,fmlcore,javafmllanguage,lowcodelanguage,mclanguage,forge-,forge-1.20.1-47.1.0.jar -DmergeModules=jna-5.10.0.jar,jna-platform-5.10.0.jar -DlibraryDirectory=.\\libraries -p .\\libraries/cpw/mods/bootstraplauncher/1.1.2/bootstraplauncher-1.1.2.jar;.\\libraries/cpw/mods/securejarhandler/2.1.10/securejarhandler-2.1.10.jar;.\\libraries/org/ow2/asm/asm-commons/9.5/asm-commons-9.5.jar;.\\libraries/org/ow2/asm/asm-util/9.5/asm-util-9.5.jar;.\\libraries/org/ow2/asm/asm-analysis/9.5/asm-analysis-9.5.jar;.\\libraries/org/ow2/asm/asm-tree/9.5/asm-tree-9.5.jar;.\\libraries/org/ow2/asm/asm/9.5/asm-9.5.jar;.\\libraries/net/minecraftforge/JarJarFileSystems/0.3.19/JarJarFileSystems-0.3.19.jar --add-modules ALL-MODULE-PATH --add-opens java.base/java.util.jar=cpw.mods.securejarhandler --add-opens java.base/java.lang.invoke=cpw.mods.securejarhandler --add-exports java.base/sun.security.util=cpw.mods.securejarhandler --add-exports jdk.naming.dns/com.sun.jndi.dns=java.naming cpw.mods.bootstraplauncher.BootstrapLauncher --username " + nick + " --version forge-1.20.1-47.1.0 --gameDir . --assetsDir .\\assets --assetIndex 5 --uuid - --accessToken - --clientId ${clientid} --xuid ${auth_xuid} --userType msa --versionType release --launchTarget forgeclient --fml.forgeVersion 47.1.0 --fml.mcVersion 1.20.1 --fml.forgeGroup net.minecraftforge --fml.mcpVersion 20230612.114412"
	args := strings.Split(comm, " ")
	out, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
	}
	fmt.Println(string(out))

}
