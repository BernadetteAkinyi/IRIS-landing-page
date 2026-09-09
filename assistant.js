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

   async onHeard(transcript) {
  try {
    const res = await fetch("http://localhost:8080/api/assistant", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ transcript }),
    });
    const data = await res.json();
    this.speak(data.reply);
  } catch (err) {
    console.error("Assistant backend error:", err);
    this.speak("My apologies, I seem to have lost my train of thought.");
  }
}

   speak(text) {
    const utterance = new SpeechSynthesisUtterance(text);
    utterance.lang = "en-GB";   
    utterance.rate = 0.95;      // slightly slower = more composed
    utterance.pitch = 0.9;      // slightly lower = more formal
    this.synth.speak(utterance);
  }
}

const butler = new ButlerVoice();
document.getElementById("butler-btn").addEventListener("click", () => {
  butler.startListening();
});
