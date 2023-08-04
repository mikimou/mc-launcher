import minecraft_launcher_lib
import subprocess
import sys
import time

lv = "forge-1.20.1-47.1.43"
md = "."
options = {
    "username": "replaceme",
    "uuid": "-",
    "token": "-"
}

minecraft_command = minecraft_launcher_lib.command.get_minecraft_command(lv, md, options)
minecraft_command.insert(3, '-Dminecraft.launcher.brand=mikiho-skurveny-launcher')
minecraft_command[0] = 'jvm/java-runtime-gamma/bin/javaw.exe'

#subprocess.Popen(minecraft_command)

print(" ".join(minecraft_command))
time.sleep(30)