package ngrok

// Subscriber is a function that acts on a LogMessage when one is emitted by a Process.
type Subscriber func(message *LogMessage)
