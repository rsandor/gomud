package main

import (
    "io"
    "bufio"
)

var colorMap = map[byte][]byte {
    'D': []byte("\033[1;30m"),
    'r': []byte("\033[0;31m"),
    'R': []byte("\033[1;31m"),
    'g': []byte("\033[0;32m"),
    'G': []byte("\033[1;32m"),
    'y': []byte("\033[0;33m"),
    'Y': []byte("\033[1;33m"),
    'b': []byte("\033[0;34m"),
    'B': []byte("\033[1;34m"),
    'm': []byte("\033[0;35m"),
    'M': []byte("\033[1;35m"),
    'c': []byte("\033[0;36m"),
    'C': []byte("\033[1;36m"),
    'w': []byte("\033[0;37m"),
    'W': []byte("\033[1;37m"),
    'x': []byte("\033[0;37m"),
    '@': []byte("@"),
}

// Replaces human readable color codes with ANSI color escape sequences
// in a given byte slice and writes them to the given writer.
func ColorWrite(w io.Writer, data []byte) error {
    buf := bufio.NewWriter(w)
    defer buf.Flush()

    for i := 0; i < len(data); i++ {
        // If the byte is not a color control symbol (@) just write it
        if data[i] != '@' {
            err := buf.WriteByte(data[i])
            if err != nil {
                return err
            }
            continue
        }

        // Replace with color escape code if found in the map
        i++
        bytes := colorMap[data[i]]
        if bytes != nil {
            _, err := buf.Write(bytes)
            if err != nil {
                return err
            }
        }
    }

    return nil
}


// Strips human readable colors codes from the given byte slice
// and writes them to the given writer.
func ColorStripWrite(w io.Writer, data []byte) error {
    buf := bufio.NewWriter(w)
    defer buf.Flush()

    for i := 0; i < len(data); i++ {
        // If the byte is not a color control symbol (@) just write it
        if data[i] != '@' {
            err := buf.WriteByte(data[i])
            if err != nil {
                return err
            }
            continue
        }

        // Strip it out
        i++
        if data[i] == '@' {
            err := buf.WriteByte('@')
            if err != nil {
                return err
            }
        }
    }

    return nil
}


