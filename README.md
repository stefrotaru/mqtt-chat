# MQTT Chat client

## Description
This is a simple MQTT chat client that is used to demonstrate the use of the
MQTT protocol together with the Go programming language.

The client connects to the MQTT broker running on the same machine on port
1883 and subscribes to a default topic (you can change the topic name in the
code below). It also publishes to the same topic.

The client uses the bubbletea library to create a terminal UI. The UI
consists of a header, a footer and a text input field. The header shows the
topic name and the footer shows the text input field. The messages are
displayed in the middle of the screen. The client allows to send text
messages and action messages (e.g. user joined/left the chat) by typing them
in the text input field and pressing "Enter". It also allows to send
commands by typing them in the text input field and pressing "Enter". For now,
the following commands are supported:

- /quit
  The client disconnects from the MQTT broker and exits.

- /hello
  The client sends an action message to the chat room saying that the user
  has joined the chat.

The client also allows to use the "Tab" key to auto-complete commands.

## Usage
To build this program, you need to install the bubbletea library. To do this
run the following command:
```
  go get -u github.com/charmbracelet/bubbletea
```
To build the program, run the following command:
```
  go build -o mqtt-chat main.go
```
To run the program, run the following command:
```
  ./mqtt-chat
```
You can specify a display name as a command line argument:
```
  ./mqtt-chat <display-name>
```
The default display name is "go-simple".

## Deployment
The project also includes a docker compose file which deploys a compatible
eclipse MQTT broker for the project. You may use it for testing the chat client:
```
docker compose up -d
```
