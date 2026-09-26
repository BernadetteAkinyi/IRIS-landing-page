/**
 * IRIS Butler - Adaptive Voice & Accessibility Companion
 * Connects speech recognition and accessible controls to the Go backend.
 */
class ButlerVoice {
  constructor() {
    this.synth = window.speechSynthesis;
    this.isListening = false;
    this.backendUrl = "http://localhost:8080/api/assistant";

    this.initSpeechRecognition();
    this.createFloatingUI();
    this.bindEvents();
    this.loadSavedPreferences();
  }

  initSpeechRecognition() {
    const SpeechRecognition =
      window.SpeechRecognition || window.webkitSpeechRecognition;

    if (!SpeechRecognition) {
      console.warn("Speech recognition is not supported in this browser. Butler will operate via text input.");
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
      if (this.currentStatus === "listening") {
        this.updateStatus("ready", "Ready");
      }
    };

    this.recognition.onresult = (event) => {
      const transcript = event.results[0][0].transcript;
      this.addMessage(transcript, "user");
      this.handleCommand(transcript);
    };

    this.recognition.onerror = (e) => {
      console.warn("Speech recognition notice:", e.error);
      this.isListening = false;
      this.setButtonsListening(false);
      this.updateStatus("ready", "Ready");
      if (e.error === "not-allowed") {
        this.addMessage("Microphone permission was denied. You may type commands to me below.", "assistant");
      }
    };
  }

  startListening() {
    if (!this.hasSpeech) {
      this.openPanel();
      this.focusInput();
      return;
    }

    try {
      if (this.isListening) {
        this.recognition.stop();
      } else {
        this.openPanel();
        this.recognition.start();
      }
    } catch (err) {
      console.error("Failed to start speech recognition:", err);
      this.openPanel();
      this.focusInput();
    }
  }

  async handleCommand(transcript) {
    this.updateStatus("processing", "Thinking...");

    const currentPrefs = this.getCurrentPreferences();
    const storedUser = localStorage.getItem("irisRememberMe") || "user_dyslexia";

    try {
      const res = await fetch(this.backendUrl, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          transcript: transcript,
          user_id: storedUser,
          current_preferences: currentPrefs,
        }),
      });

      if (!res.ok) {
        throw new Error(`Server returned ${res.status}`);
      }

      const data = await res.json();

      this.addMessage(data.reply, "assistant");
      this.updateStatus("speaking", "Speaking...");

      // Execute returned accessibility actions
      if (data.actions && data.actions.length > 0) {
        data.actions.forEach((action) => this.executeAction(action));
      }

      this.speak(data.reply, () => {
        this.updateStatus("ready", "Ready");
      });
    } catch (err) {
      console.error("Butler assistant backend error:", err);
      const fallbackMsg = "My apologies, I had trouble reaching the IRIS service. Please ensure the backend server is running.";
      this.addMessage(fallbackMsg, "assistant");
      this.speak(fallbackMsg, () => {
        this.updateStatus("ready", "Ready");
      });
    }
  }

  executeAction(action) {
    if (!action || !action.type) return;

    switch (action.type) {
      case "SET_COLOR_THEME":
        this.applyColorTheme(action.payload.theme);
        break;

      case "SET_READING_MODE":
        this.applyReadingMode(action.payload.enabled, action.payload.font, action.payload.line_spacing);
        break;

      case "SET_TEXT_SIZE":
        this.applyTextSize(action.payload.size);
        break;

      case "SET_HIGH_CONTRAST":
        this.applyHighContrast(action.payload.enabled);
        break;

      case "SET_SIMPLIFIED_LAYOUT":
        this.applySimplifiedLayout(action.payload.enabled);
        break;

      case "SET_REDUCE_MOTION":
        this.applyReduceMotion(action.payload.enabled);
        break;

      case "SET_TTS":
        if (action.payload.action === "stop") {
          if (this.synth) this.synth.cancel();
        } else if (action.payload.action === "read") {
          this.readCurrentContent();
        }
        break;

      case "APPLY_USER_PREFERENCES":
        this.applyAllPreferences(action.payload.preferences);
        break;

      case "RESET_INTERFACE":
        this.resetAll();
        break;

      case "NAVIGATE":
        if (action.payload && action.payload.page) {
          setTimeout(() => {
            window.location.href = action.payload.page;
          }, 1200);
        }
        break;

      default:
        console.log("Unhandled Butler action:", action);
    }
  }

  applyColorTheme(theme) {
    const preview = document.getElementById("preview");
    const themeClasses = [
      "theme-baby-pink",
      "theme-lavender",
      "theme-mint-green",
      "theme-peach",
      "theme-sky-blue",
      "baby-pink",
      "lavender",
      "mint-green",
      "peach",
      "sky-blue",
    ];

    document.body.classList.remove(...themeClasses);
    if (preview) {
      preview.classList.remove(...themeClasses);
    }

    if (theme && theme !== "default" && theme !== "high-contrast") {
      const cls = `theme-${theme}`;
      document.body.classList.add(cls);
      if (preview) {
        preview.classList.add(cls);
        this.updatePreviewText("Color Theme Applied", `IRIS interface adjusted to ${theme.replace("-", " ")} palette.`);
      }
    } else if (theme === "default") {
      if (preview) {
        this.updatePreviewText("Default Theme", "Restored standard warm palette.");
      }
    }

    this.savePreference("color_theme", theme);
  }

  applyReadingMode(enabled, font = "Lexend", lineSpacing = "relaxed") {
    const preview = document.getElementById("preview");

    if (enabled) {
      document.body.classList.add("reading-mode");
      if (preview) {
        preview.classList.add("reading-mode");
        this.updatePreviewText(
          "Focus Reading Mode",
          `Clear typography (${font}) and ${lineSpacing} spacing applied to ease visual stress and support reading.`
        );
        this.highlightButton("reading-mode-btn");
      }
    } else {
      document.body.classList.remove("reading-mode");
      if (preview) {
        preview.classList.remove("reading-mode");
      }
    }

    this.savePreference("reading_mode", enabled);
  }

  applyTextSize(size) {
    const sizeClasses = ["text-size-small", "text-size-medium", "text-size-large", "text-size-x-large"];
    document.body.classList.remove(...sizeClasses);

    if (size && size !== "medium") {
      document.body.classList.add(`text-size-${size}`);
    }

    const preview = document.getElementById("preview");
    if (preview) {
      this.updatePreviewText("Text Size Adjusted", `Font size scaled to ${size} for optimal comfort.`);
    }

    this.savePreference("text_size", size);
  }

  applyHighContrast(enabled) {
    const preview = document.getElementById("preview");

    if (enabled) {
      document.body.classList.add("high-contrast");
      if (preview) {
        preview.classList.add("high-contrast");
        this.updatePreviewText(
          "High Contrast Enabled",
          "Dark palette with high contrast borders and text activated for visual definition."
        );
        this.highlightButton("high-contrast-btn");
      }
    } else {
      document.body.classList.remove("high-contrast");
      if (preview) {
        preview.classList.remove("high-contrast");
      }
    }

    this.savePreference("high_contrast", enabled);
  }

  applySimplifiedLayout(enabled) {
    if (enabled) {
      document.body.classList.add("simplified-layout");
    } else {
      document.body.classList.remove("simplified-layout");
    }

    const preview = document.getElementById("preview");
    if (preview) {
      this.updatePreviewText(
        enabled ? "Simplified View" : "Full Layout Restored",
        enabled
          ? "Unnecessary clutter and decorative elements minimized to maintain focus."
          : "Standard comprehensive layout displayed."
      );
    }

    this.savePreference("simplified_layout", enabled);
  }

  applyReduceMotion(enabled) {
    const preview = document.getElementById("preview");

    if (enabled) {
      document.body.classList.add("reduce-motion");
      if (preview) {
        preview.classList.add("reduce-motion");
        this.updatePreviewText(
          "Reduce Motion Enabled",
          "Animations and transitions are paused to reduce visual fatigue."
        );
        this.highlightButton("reduce-motion-btn");
      }
    } else {
      document.body.classList.remove("reduce-motion");
      if (preview) {
        preview.classList.remove("reduce-motion");
      }
    }

    this.savePreference("reduce_motion", enabled);
  }

  applyAllPreferences(prefs) {
    if (!prefs) return;
    if (prefs.color_theme) this.applyColorTheme(prefs.color_theme);
    if (prefs.reading_mode !== undefined) this.applyReadingMode(prefs.reading_mode, prefs.reading_font, prefs.line_spacing);
    if (prefs.text_size) this.applyTextSize(prefs.text_size);
    if (prefs.high_contrast !== undefined) this.applyHighContrast(prefs.high_contrast);
    if (prefs.simplified_layout !== undefined) this.applySimplifiedLayout(prefs.simplified_layout);
    if (prefs.reduce_motion !== undefined) this.applyReduceMotion(prefs.reduce_motion);
  }

  resetAll() {
    const themeClasses = [
      "theme-baby-pink",
      "theme-lavender",
      "theme-mint-green",
      "theme-peach",
      "theme-sky-blue",
      "high-contrast",
      "reading-mode",
      "reduce-motion",
      "simplified-layout",
      "text-size-small",
      "text-size-medium",
      "text-size-large",
      "text-size-x-large",
    ];

    document.body.classList.remove(...themeClasses);
    const preview = document.getElementById("preview");
    if (preview) {
      preview.classList.remove(...themeClasses);
      this.updatePreviewText("Welcome to IRIS", "Technology that adapts to you, making every digital experience more comfortable.");
    }

    localStorage.removeItem("iris_preferences");
  }

  readCurrentContent() {
    const preview = document.getElementById("previewText");
    const textToRead = preview
      ? preview.textContent
      : document.querySelector("main p")?.textContent || "Welcome to IRIS, your adaptive accessibility companion.";

    this.speak(textToRead);
  }

  speak(text, onEnd) {
    if (!this.synth) {
      if (onEnd) onEnd();
      return;
    }

    this.synth.cancel(); // cancel any active speech
    const utterance = new SpeechSynthesisUtterance(text);
    utterance.lang = "en-GB";
    utterance.rate = 0.95;
    utterance.pitch = 0.92;

    // Pick British voice if available
    const voices = this.synth.getVoices();
    const britishVoice = voices.find(
      (v) => v.lang.includes("en-GB") || v.name.includes("UK") || v.name.includes("George") || v.name.includes("Oliver")
    );
    if (britishVoice) {
      utterance.voice = britishVoice;
    }

    utterance.onend = () => {
      if (onEnd) onEnd();
    };

    utterance.onerror = () => {
      if (onEnd) onEnd();
    };

    this.synth.speak(utterance);
  }

  updatePreviewText(title, text) {
    const t = document.getElementById("previewTitle");
    const p = document.getElementById("previewText");
    if (t) t.textContent = title;
    if (p) p.textContent = text;
  }

  highlightButton(id) {
    document.querySelectorAll(".feature-btn").forEach((btn) => btn.classList.remove("active"));
    const target = document.getElementById(id);
    if (target) target.classList.add("active");
  }

  savePreference(key, value) {
    try {
      const cur = JSON.parse(localStorage.getItem("iris_preferences") || "{}");
      cur[key] = value;
      localStorage.setItem("iris_preferences", JSON.stringify(cur));
    } catch (e) {
      console.warn("Storage error:", e);
    }
  }

  getCurrentPreferences() {
    try {
      return JSON.parse(localStorage.getItem("iris_preferences") || "{}");
    } catch {
      return {};
    }
  }

  loadSavedPreferences() {
    const prefs = this.getCurrentPreferences();
    if (Object.keys(prefs).length > 0) {
      this.applyAllPreferences(prefs);
    }
  }

  // ==========================================
  // Butler Accessible UI & Floating Widget
  // ==========================================
  createFloatingUI() {
    if (document.getElementById("butler-floating-widget")) return;

    const widget = document.createElement("div");
    widget.id = "butler-floating-widget";
    widget.className = "butler-floating-widget";
    widget.innerHTML = `
      <div id="butler-panel" class="butler-panel" role="dialog" aria-labelledby="butler-title">
        <div class="butler-panel-header">
          <div class="butler-panel-title">
            <span>🎙️</span>
            <span id="butler-title">IRIS Butler</span>
            <span id="butler-status" class="butler-status-badge">Ready</span>
          </div>
          <button id="butler-close-btn" class="butler-close-btn" aria-label="Close Butler panel">&times;</button>
        </div>

        <div id="butler-dialog-box" class="butler-dialog-box" aria-live="polite">
          <div class="butler-bubble assistant">
            Good day. I am IRIS Butler. How may I personalize your experience?
          </div>
        </div>

        <div class="butler-chips-container">
          <div class="butler-chips-label">Quick Voice Commands</div>
          <div class="butler-chips">
            <button type="button" class="butler-chip" data-cmd="Change the page to baby pink">🌸 Baby Pink</button>
            <button type="button" class="butler-chip" data-cmd="Make this easier to read">📖 Focus Reading</button>
            <button type="button" class="butler-chip" data-cmd="Increase the text size">🔍 Larger Text</button>
            <button type="button" class="butler-chip" data-cmd="Turn on text to speech">🔊 Read Aloud</button>
            <button type="button" class="butler-chip" data-cmd="Enable high contrast">🌓 High Contrast</button>
            <button type="button" class="butler-chip" data-cmd="Make the page simpler">🧘 Simplify</button>
            <button type="button" class="butler-chip" data-cmd="Apply my preferences">👤 My Profile</button>
            <button type="button" class="butler-chip" data-cmd="Reset page">↺ Reset</button>
          </div>
        </div>

        <form id="butler-input-form" class="butler-input-row">
          <input
            id="butler-text-input"
            class="butler-text-input"
            type="text"
            placeholder="Say or type command (e.g. 'baby pink')..."
            aria-label="Butler command input"
          />
          <button type="submit" class="butler-send-btn">Send</button>
        </form>
      </div>

      <button id="butler-fab" class="butler-fab" aria-label="Open IRIS Voice Butler" title="Speak to IRIS Butler">
        🎙️
      </button>
    `;

    document.body.appendChild(widget);
  }

  bindEvents() {
    // Top / Inline button on features.html or any page
    const inlineBtn = document.getElementById("butler-btn");
    if (inlineBtn) {
      inlineBtn.addEventListener("click", () => {
        this.openPanel();
        this.startListening();
      });
    }

    // Floating FAB button
    const fab = document.getElementById("butler-fab");
    if (fab) {
      fab.addEventListener("click", () => {
        const panel = document.getElementById("butler-panel");
        const isOpen = panel && panel.classList.contains("open");
        if (isOpen && !this.isListening) {
          this.startListening();
        } else {
          this.openPanel();
          this.startListening();
        }
      });
    }

    // Close button
    const closeBtn = document.getElementById("butler-close-btn");
    if (closeBtn) {
      closeBtn.addEventListener("click", () => this.closePanel());
    }

    // Text form submission
    const form = document.getElementById("butler-input-form");
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

    // Quick chips
    document.querySelectorAll(".butler-chip").forEach((chip) => {
      chip.addEventListener("click", () => {
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
    if (panel) {
      panel.classList.add("open");
    }
  }

  closePanel() {
    const panel = document.getElementById("butler-panel");
    if (panel) {
      panel.classList.remove("open");
    }
    if (this.isListening && this.recognition) {
      this.recognition.stop();
    }
  }

  focusInput() {
    const input = document.getElementById("butler-text-input");
    if (input) input.focus();
  }

  updateStatus(statusKey, text) {
    this.currentStatus = statusKey;
    const badge = document.getElementById("butler-status");
    if (badge) {
      badge.textContent = text;
      badge.className = `butler-status-badge ${statusKey}`;
    }
  }

  setButtonsListening(listening) {
    const inlineBtn = document.getElementById("butler-btn");
    const fab = document.getElementById("butler-fab");

    if (inlineBtn) {
      inlineBtn.classList.toggle("listening", listening);
      inlineBtn.textContent = listening ? "🛑 Listening..." : "🎙️ Ask Butler";
    }
    if (fab) {
      fab.classList.toggle("listening", listening);
    }
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
}

// Global Butler instance
let irisButler = null;
window.addEventListener("DOMContentLoaded", () => {
  irisButler = new ButlerVoice();
});
