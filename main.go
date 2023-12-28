package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Message interface {
	formatMessage() string
}

type TextMessage struct {
	user string
	text string
	sentAt time.Time
}

func (m TextMessage) formatMessage() string {
	return fmt.Sprintf("%s [%s]: %s", m.user, m.sentAt.Format("15:04:05"), m.text)
}

type ActionMessage struct {
	user string
	text string
}

func (m ActionMessage) formatMessage() string {
	return fmt.Sprintf("%s %s", m.user, m.text)
}

type model struct {
	topic	  string
	name	  string
	client    mqtt.Client
	messages  []Message
	textInput textinput.Model
	cursor    int
}

var received = make(chan mqtt.Message)

func initialModel(topic string, name string, client mqtt.Client) model {
	ti := textinput.New()
	ti.Placeholder = "Type your message..."
	ti.Focus()
	ti.Width = 80

	m := model{
		topic: topic,
		name: name,
		client: client,
		messages: []Message{},
		textInput: ti,
		cursor:  0,
	}

	// message handler adds messages to a global channel
	var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		received <- msg
	}

	if token := client.Subscribe(topic + "#", 0, messagePubHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	return m
}

func (m model) Init() tea.Cmd {
    return textinput.Blink
}

type CmdDecorator func() tea.Cmd

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	commands := map[string]CmdDecorator {
		"/quit": func() tea.Cmd {
			return tea.Quit
		},
		"/hello": func() tea.Cmd {
			action := ActionMessage{
				user: m.name,
				text: "says hello!",
			}
			token := m.client.Publish(m.topic + m.name, 0, false, action)
			token.Wait()
			m.textInput.SetValue("")
			return nil
		},
	}

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			action := ActionMessage{
				user: m.name,
				text: "left the chat",
			}
			token := m.client.Publish(m.topic + m.name, 0, false, action)
			token.Wait()

			return m, tea.Quit
		case tea.KeyTab:
			// check if input field start with slash
			selectedCmd := ""
			if strings.HasPrefix(m.textInput.Value(), "/") {
				for k := range commands {
					if strings.HasPrefix(k, m.textInput.Value()) {
						if selectedCmd == "" {
							selectedCmd = k
						} else {
							selectedCmd = ""
							break
						}
					}
				}
			}
			if selectedCmd != "" {
				m.textInput.SetValue(selectedCmd)
				m.textInput.CursorEnd()
			}
			return m, nil
		case tea.KeyEnter:
			textInput := m.textInput.Value()
			if textInput[0] == '/' {
				if cmd, ok := commands[textInput]; ok {
					return m, cmd()
				} else {
					// Print "no such command" action message
					m.messages = append(m.messages, ActionMessage{
						user: "Interface:",
						text: "The command " + textInput[1:] + " does not exist!",
					})
					m.textInput.SetValue("")
					return m, nil
				}
			}

			// TODO: send Action Message/TextMessage instead of string as payload. Use encoding: JSON?
			token := m.client.Publish(m.topic + m.name, 0, false, m.textInput.Value())
			token.Wait()
			m.textInput.SetValue("")

			return m, nil
		}
    }

	select {
	case msg := <-received:
		m.messages = append(m.messages, TextMessage{
			user: strings.Split(msg.Topic(), "/")[1],
			text: string(msg.Payload()),
			sentAt: time.Now(),
		})
	default:
	}
	m.textInput, cmd = m.textInput.Update(msg)
    return m, cmd
}

func (m model) View() string {
    // The header
    s := fmt.Sprintf("--- %s ---\n\n", m.topic)

    // Iterate over messages
    for _, msg := range m.messages {
        // Render the row
        s += msg.formatMessage() + "\n"
    }

    // The footer
    s += "\n----" + strings.Repeat("-", len(m.topic)) + "----\n"
	s += m.textInput.View() + "\n"

    // Send the UI for rendering
    return s
}

func main() {
	var displayName string
	if len(os.Args) > 1 {
		displayName = os.Args[1]
	} else {
		displayName = "go-simple"
	}

	opts := mqtt.NewClientOptions().AddBroker("0.0.0.0:1883")
	opts.SetClientID(displayName)
	opts.SetCleanSession(true)
	
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	p := tea.NewProgram(initialModel("topicul_de_miercuri_seara/", displayName, c))
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error: %v", err)
        os.Exit(1)
    }
}
