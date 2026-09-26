class ButlerVoice {
  constructor() {
    this.synth = window.speechSynthesis;
    this.isListening = false;
    this.backendUrl = "http://localhost:8080/api/assistant";

    this.initSpeechRecognition();
    this.createUI();
    this.bindEvents();
    this.loadSavedPreferences();
  }

  initSpeechRecognition() {
    const SpeechRecognition =
      window.SpeechRecognition || window.webkitSpeechRecognition;

    if (!SpeechRecognition) {
      this.hasSpeech = false;
      return;
    }

    this.hasSpeech = true;
    this.recognition = new SpeechRecognition();
    this.recognition.continuous = false;
    this.recognition.lang = "en-US";
    this.recognition.interimResults = false;

    this.recognition.onstart = () => {
      this.isListening = true;
      this.updateStatus("listening", "Listening...");
      this.setButtonsListening(true);
    };

    this.recognition.onend = () => {
      this.isListening = false;
      this.setButtonsListening(false);
      this.updateStatus("ready", "Ready");
    };

    this.recognition.onresult = (event) => {
      const transcript = event.results[0][0].transcript;
      this.addMessage(transcript, "user");
      this.handleCommand(transcript);
    };

    this.recognition.onerror = () => {
      this.isListening = false;
      this.setButtonsListening(false);
      this.updateStatus("ready", "Ready");
    };
  }

  startListening() {
    if (!this.hasSpeech) {
      this.openPanel();
      const input = document.getElementById("butler-text-input");
      if (input) input.focus();
      return;
    }

    try {
      if (this.isListening) {
        this.recognition.stop();
      } else {
        this.openPanel();
        this.recognition.start();
      }
    } catch {
      this.openPanel();
    }
  }

  async handleCommand(transcript) {
    this.updateStatus("processing", "Thinking...");

    try {
      const res = await fetch(this.backendUrl, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ transcript }),
      });

      if (!res.ok) throw new Error("Network response error");

      const data = await res.json();
      this.addMessage(data.reply, "assistant");

      if (data.action) {
        this.executeAction(data.action);
      }

      this.speak(data.reply);
    } catch {
      const msg = "Unable to connect to IRIS backend. Please make sure the server is running on :8080.";
      this.addMessage(msg, "assistant");
      this.speak(msg);
    }
  }

  executeAction(action) {
    const preview = document.getElementById("preview");

    switch (action.type) {
      case "SET_THEME":
        if (action.payload === "baby-pink") {
          document.body.classList.toggle("baby-pink");
          if (preview) preview.classList.toggle("baby-pink");
        }
        break;

      case "SET_READING_MODE":
        document.body.classList.toggle("reading-mode");
        if (preview) preview.classList.toggle("reading-mode");
        break;

      case "SET_TEXT_SIZE":
        document.body.classList.remove("text-size-small", "text-size-large");
        if (action.payload === "large") {
          document.body.classList.add("text-size-large");
        } else if (action.payload === "small") {
          document.body.classList.add("text-size-small");
        }
        break;

      case "HIGH_CONTRAST":
        document.body.classList.toggle("high-contrast");
        if (preview) preview.classList.toggle("high-contrast");
        break;

      case "SIMPLIFY":
        document.body.classList.toggle("simplified-layout");
        break;

      case "REDUCE_MOTION":
        document.body.classList.toggle("reduce-motion");
        if (preview) preview.classList.toggle("reduce-motion");
        break;

      case "READ_PAGE":
        this.readPageContent();
        break;

      case "RESET":
        document.body.classList.remove(
          "baby-pink",
          "high-contrast",
          "reading-mode",
          "reduce-motion",
          "simplified-layout",
          "text-size-small",
          "text-size-large"
        );
        if (preview) {
          preview.classList.remove("baby-pink", "high-contrast", "reading-mode", "reduce-motion");
        }
        localStorage.removeItem("iris_theme");
        break;
    }
  }

  readPageContent() {
    const preview = document.getElementById("previewText");
    const text = preview
      ? preview.textContent
      : document.querySelector("main p")?.textContent || "Welcome to IRIS.";
    this.speak(text);
  }

  speak(text) {
    if (!this.synth) return;
    this.synth.cancel();

    const utterance = new SpeechSynthesisUtterance(text);
    utterance.lang = "en-GB";
    utterance.rate = 0.95;

    const voices = this.synth.getVoices();
    const gbVoice = voices.find((v) => v.lang.includes("en-GB"));
    if (gbVoice) utterance.voice = gbVoice;

    this.synth.speak(utterance);
  }

  createUI() {
    if (document.getElementById("butler-floating-widget")) return;

    const widget = document.createElement("div");
    widget.id = "butler-floating-widget";
    widget.className = "butler-floating-widget";
    widget.innerHTML = `
      <div id="butler-panel" class="butler-panel" role="dialog" aria-label="IRIS Butler Dialog">
        <div class="butler-panel-header">
          <div class="butler-panel-title">
            <span>🎙️</span>
            <span>IRIS Butler</span>
            <span id="butler-status" class="butler-status-badge">Ready</span>
          </div>
          <button id="butler-close-btn" class="butler-close-btn" type="button" aria-label="Close assistant">✕ Close</button>
        </div>

        <div id="butler-dialog-box" class="butler-dialog-box">
          <div class="butler-bubble assistant">
            Good day. I am IRIS Butler. How may I assist your browsing today?
          </div>
        </div>

        <div class="butler-chips">
          <button type="button" class="butler-chip" data-cmd="Change the page to baby pink">🌸 Baby Pink</button>
          <button type="button" class="butler-chip" data-cmd="Make this easier to read">📖 Reading Mode</button>
          <button type="button" class="butler-chip" data-cmd="Increase the text size">🔍 Larger Text</button>
          <button type="button" class="butler-chip" data-cmd="Turn on text to speech">🔊 Read Page</button>
          <button type="button" class="butler-chip" data-cmd="Make the page simpler">🧘 Simplify</button>
          <button type="button" class="butler-chip" data-cmd="Turn on high contrast">🌓 High Contrast</button>
          <button type="button" class="butler-chip" data-cmd="Reset to default">↺ Reset</button>
        </div>

        <form id="butler-input-form" class="butler-input-row">
          <input
            id="butler-text-input"
            class="butler-text-input"
            type="text"
            placeholder="Ask Butler (e.g. 'baby pink')..."
            aria-label="Butler text input"
          />
          <button type="submit" class="butler-send-btn">Send</button>
        </form>
      </div>

      <button id="butler-fab" class="butler-fab" aria-label="Open IRIS Butler" type="button">
        🎙️
      </button>
    `;

    document.body.appendChild(widget);
  }

  bindEvents() {
    const inlineBtn = document.getElementById("butler-btn");
    const fab = document.getElementById("butler-fab");
    const closeBtn = document.getElementById("butler-close-btn");
    const form = document.getElementById("butler-input-form");
    const panel = document.getElementById("butler-panel");

    const toggle = (e) => {
      e.stopPropagation();
      if (panel && panel.classList.contains("open")) {
        this.closePanel();
      } else {
        this.openPanel();
        this.startListening();
      }
    };

    if (inlineBtn) inlineBtn.addEventListener("click", toggle);
    if (fab) fab.addEventListener("click", toggle);

    if (closeBtn) {
      closeBtn.addEventListener("click", (e) => {
        e.stopPropagation();
        this.closePanel();
      });
    }

    if (panel) {
      panel.addEventListener("click", (e) => e.stopPropagation());
    }

    document.addEventListener("click", (e) => {
      if (panel && panel.classList.contains("open")) {
        if (!panel.contains(e.target) && e.target !== inlineBtn && e.target !== fab) {
          this.closePanel();
        }
      }
    });

    document.addEventListener("keydown", (e) => {
      if (e.key === "Escape") {
        this.closePanel();
      }
    });

    if (form) {
      form.addEventListener("submit", (e) => {
        e.preventDefault();
        const input = document.getElementById("butler-text-input");
        const val = input ? input.value.trim() : "";
        if (val) {
          input.value = "";
          this.addMessage(val, "user");
          this.handleCommand(val);
        }
      });
    }

    document.querySelectorAll(".butler-chip").forEach((chip) => {
      chip.addEventListener("click", (e) => {
        e.stopPropagation();
        const cmd = chip.getAttribute("data-cmd");
        if (cmd) {
          this.addMessage(cmd, "user");
          this.handleCommand(cmd);
        }
      });
    });
  }

  openPanel() {
    const panel = document.getElementById("butler-panel");
    if (panel) panel.classList.add("open");
  }

  closePanel() {
    const panel = document.getElementById("butler-panel");
    if (panel) panel.classList.remove("open");

    if (this.recognition && this.isListening) {
      this.recognition.stop();
    }
    if (this.synth) {
      this.synth.cancel();
    }

    this.setButtonsListening(false);
    this.updateStatus("ready", "Ready");
  }

  updateStatus(statusKey, text) {
    const badge = document.getElementById("butler-status");
    if (badge) badge.textContent = text;
  }

  setButtonsListening(listening) {
    const inlineBtn = document.getElementById("butler-btn");
    const fab = document.getElementById("butler-fab");
    if (inlineBtn) inlineBtn.classList.toggle("listening", listening);
    if (fab) fab.classList.toggle("listening", listening);
  }

  addMessage(text, sender) {
    const box = document.getElementById("butler-dialog-box");
    if (!box) return;

    const bubble = document.createElement("div");
    bubble.className = `butler-bubble ${sender}`;
    bubble.textContent = text;
    box.appendChild(bubble);
    box.scrollTop = box.scrollHeight;
  }

  loadSavedPreferences() {
    const saved = localStorage.getItem("iris_theme");
    if (saved) document.body.classList.add(saved);
  }
}

window.addEventListener("DOMContentLoaded", () => {
  new ButlerVoice();
});
