# qq

A question. An answer. That's it.

```
$ qq what is the mass of the sun
The mass of the Sun is approximately 1.989 × 10³⁰ kg.
```

`qq` is a command-line tool for when you just need a quick answer and don't want to open a browser, start a chat session, or wait for a heavy client to load. Ask a question, get a response, get back to work.

Zero external dependencies. Single binary. No config files.

## Install

```bash
go install github.com/tickloop/qq@latest
```

Or build from source:

```bash
git clone https://github.com/tickloop/qq.git
cd qq
make build
```

## Setup

`qq` uses [OpenRouter](https://openrouter.ai), giving you access to hundreds of models through a single API key.

```bash
export OPENROUTER_API_KEY="your-key-here"
```

That's the only configuration. Add it to your shell profile and forget about it.

## Usage

```bash
qq how do I reverse a list in python
```

```bash
qq what port does postgres use
```

```bash
qq convert 72 fahrenheit to celsius
```

No quotes needed. Everything after `qq` is your question.

### Pick a model

The default model is `google/gemini-3-flash-preview` — fast and cheap, exactly right for quick questions. Override it when you need to:

```bash
qq -m anthropic/claude-sonnet-4 explain the CAP theorem briefly
```

Any model on [OpenRouter](https://openrouter.ai/models) works.


## Philosophy

Some tools try to do everything. `qq` does one thing: it takes a question from your terminal and gives you an answer. No interactive mode, no conversation history, no streaming, no plugins. Just a question and an answer.

If you need more than that, there are great tools for that. This isn't one of them.

## License

MIT
