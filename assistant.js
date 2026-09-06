class ButlerVoice {
  constructor() {
    // SpeechRecognition is prefixed in Chrome/Edge
    const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition;
    if (!SpeechRecognition) {
      console.warn("Speech recognition not supported in this browser.");
      return;
    }

    this.recognition = new SpeechRecognition();
    this.recognition.continuous = false;   
    this.recognition.lang = "en-US";
    this.recognition.interimResults = false;

    this.synth = window.speechSynthesis;

  }
}