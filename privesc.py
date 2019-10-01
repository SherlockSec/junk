#!/usr/env python

import sys

try:
    import subprocess
    version = 0
except ImportError:
    import os
    version = 1

try:
    import privEscCommands
except ImportError:
    print("privEscCommands.py file not found.")
    sys.exit(1)

#Functions

def execute(command_dict):
    for item in command_dict:
        command = command_dict[item]["command"]
        if version == 0:
            out, error = subprocess.Popen([command], stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True).communicate()
            results = out.split('\n')
        else:
            echo_stdout = os.popen(command, 'r')
            results = echo_stdout.read().split('\n')
        command_dict[item]["results"] = results
    return command_dict

def showResults(command_dict):
    for item in command_dict:
        msg = command_dict[item]["msg"]
        results = command_dict[item]["results"]
        print("[+] " + msg)
        for result in results:
            if result.strip() != "":
                print("    " + result.strip())
        print
    return

def writeResults(msg, results):
    f = open("privcheckout.txt", "a")
    f.write("[+] " + str( len( results )-1 ) + " " + msg)
    for result in results:
        if result.strip() != "":
            f.write("    " + result.strip())
    f.close()
    return


