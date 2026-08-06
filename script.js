const preview = document.getElementById("preview");

if (preview) {
  const highContrastBtn = document.getElementById("high-contrast-btn");
  const readingBtn = document.getElementById("reading-mode-btn");
  const textToSpeechBtn = document.getElementById("text-to-speech-btn");
  const reduceMotionBtn = document.getElementById("reduce-motion-btn");
  const previewTitle = document.getElementById("previewTitle");
  const previewText = document.getElementById("previewText");
  const previewNote = document.getElementById("previewNote");

  function clearFeatureClasses() {
    preview.classList.remove("high-contrast", "reading-mode", "reduce-motion");
  }

  function setActiveButton(button) {
    document.querySelectorAll(".feature-btn").forEach((btn) => {
      btn.classList.remove("active");
    });
    button.classList.add("active");
  }

  highContrastBtn.addEventListener("click", function () {
    clearFeatureClasses();
    setActiveButton(highContrastBtn);
    preview.classList.add("high-contrast");
    previewTitle.textContent = "High Contrast enabled";
    previewText.textContent =
      "IRIS shifts to a darker palette, making text and buttons stand out clearly against the background.";
    previewNote.textContent =
      "A helpful option for brighter rooms or low-vision environments.";
  });

  readingBtn.addEventListener("click", function () {
    clearFeatureClasses();
    setActiveButton(readingBtn);
    preview.classList.add("reading-mode");
    previewTitle.textContent = "Focus Reading";
    previewText.textContent =
      "Text is easier to scan with larger spacing, cleaner typography, and a calmer layout that reduces visual noise.";
    previewNote.textContent =
      "Designed for longer reading sessions and clearer comprehension.";
  });

  textToSpeechBtn.addEventListener("click", function () {
    clearFeatureClasses();
    setActiveButton(textToSpeechBtn);
    previewTitle.textContent = "Text-to-Speech";
    previewText.textContent =
      "This demo shows how Iris can bring on-screen content to life with spoken feedback and gentle pacing.";
    previewNote.textContent =
      "A concept demonstration: click the button to hear the preview text read aloud.";
    if ("speechSynthesis" in window) {
      const utterance = new SpeechSynthesisUtterance(previewText.textContent);
      utterance.rate = 0.95;
      window.speechSynthesis.cancel();
      window.speechSynthesis.speak(utterance);
    }
  });

  reduceMotionBtn.addEventListener("click", function () {
    clearFeatureClasses();
    setActiveButton(reduceMotionBtn);
    preview.classList.add("reduce-motion");
    previewTitle.textContent = "Reduce Motion";
    previewText.textContent =
      "This mode minimizes movement and animations to create a steadier, calmer page experience.";
    previewNote.textContent =
      "Ideal for users who prefer less motion on screen.";
  });

  window.addEventListener("DOMContentLoaded", function () {
    highContrastBtn.click();
  });
}

const loginForm = document.getElementById("loginForm");
const loginMessage = document.getElementById("loginMessage");

if (loginForm) {
  loginForm.addEventListener("submit", function (event) {
    event.preventDefault();

    const email = document.getElementById("loginEmail").value.trim();
    const password = document.getElementById("loginPassword").value.trim();
    const rememberMe = document.getElementById("rememberMe").checked;

    if (!email || !password) {
      showLoginMessage("Please enter both email and password.", "error");
      return;
    }

    if (!email.includes("@")) {
      showLoginMessage("Please enter a valid email address.", "error");
      return;
    }

    showLoginMessage(
      `Welcome back, ${email.split("@")[0]}. Redirecting to your dashboard...`,
      "success"
    );

    setTimeout(() => {
      window.location.href = "features.html";
    }, 1200);

    if (rememberMe) {
      localStorage.setItem("irisRememberMe", email);
    } else {
      localStorage.removeItem("irisRememberMe");
    }
  });

  const savedEmail = localStorage.getItem("irisRememberMe");
  if (savedEmail) {
    document.getElementById("loginEmail").value = savedEmail;
    document.getElementById("rememberMe").checked = true;
  }
}

function showLoginMessage(message, type) {
  if (!loginMessage) return;

  loginMessage.textContent = message;
  loginMessage.className = "login-message " + type;
  loginMessage.style.display = "block";
}