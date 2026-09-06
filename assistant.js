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

    //Fires when the browser has a final transcript
    this.recognition.onresult = (event) => {
      const transcript = event.results[0][0].transcript;
      this.onHeard(transcript);
    };

    this.recognition.onerror = (e) => console.error("Recognition error:", e.error);
  }
  startListening() {
    this.recognition?.start();
  }

  }
