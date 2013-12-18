package main

import "io/ioutil"

type Message struct {
    Name    string
    Text    []byte
}

// Returns the full filename for the message
func (msg *Message) Path() string {
    return messagesPath + msg.Name
}

// Reads a message from disk
func (msg *Message) Read() error {
    text, err := ioutil.ReadFile(msg.Path())
    if err != nil {
        return err
    }
    msg.Text = text
    return nil
}

// Writes a message to disk
func (msg *Message) Write() error {
    err := ioutil.WriteFile(msg.Path(), msg.Text, 0644)
    return err
}


// Map of file names to file contents from files in the $ROOT/messags folder
var messages map[string]*Message

// Loads all of the game messages files from the $ROOT/messages folder into
// the globally accessable messages map.
func loadMessages() error {
    messages = make(map[string]*Message)

    files, err := ioutil.ReadDir(messagesPath)
    if err != nil {
        return err
    }

    for _, file := range files {
        msg := &Message{Name: file.Name()}
        if err := msg.Read(); err != nil {
            return err
        }
        messages[msg.Name] = msg
    }

    return nil
}


// Gets the text for a message with the given name
func getMessage(name string) []byte {
    return messages[name].Text
}
