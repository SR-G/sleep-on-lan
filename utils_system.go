package main

import (
    "os/exec"
    "strings"
)

func Execute(command string) (bool, string, error) {
    // splitting head => g++ parts => rest of the command
    parts := strings.Fields(command)
    head := parts[0]
    parts = parts[1:len(parts)]

    out, err := exec.Command(head, parts...).Output()
    if err != nil {
        return false, "", err
    }
    return true, string(out), nil
}