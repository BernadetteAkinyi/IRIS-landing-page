IRIS – Adaptive Accessibility Companion

IRIS is an adaptive accessibility companion designed to help people with invisible disabilities navigate digital services and everyday experiences in a way that suits their individual needs.

Instead of making users adapt to a fixed interface, IRIS adapts the interface to the user.

Core Idea

IRIS learns about a user's accessibility needs during sign-up and uses that information to personalize their experience.

Users can also interact with an AI-powered Butler through voice commands.

For example:

«"Hey Butler, change the page to baby pink."»

The Butler understands the request and communicates with the appropriate parts of IRIS to change the interface.

A user with dyslexia could also say:

«"Hey Butler, make this easier for me to read."»

IRIS can then apply the user's reading preferences.

---

Front-End

The front-end is responsible for the user's interface and how IRIS adapts to their needs.

Technologies

- HTML
- CSS
- JavaScript

The front-end should have:

1. Landing Page

- Introduction to IRIS
- Explanation of what IRIS does
- Accessibility-focused design
- Sign-up and login options

2. Sign-Up

During sign-up, users should be able to provide:

- Accessibility needs
- Conditions that affect how they interact with interfaces
- Reading preferences
- Color preferences
- Other accessibility preferences

This information becomes part of the user's IRIS profile.

3. Personalized Interface

The interface should automatically adapt based on the user's saved preferences.

Examples include:

- Color themes
- Text size
- Dyslexia-friendly text presentation
- Simplified layouts
- Text-to-speech
- High contrast
- Reduced motion

4. User Dashboard

The dashboard should allow users to:

- View their accessibility preferences
- Change their preferences
- Access available services
- View events
- Access healthcare-related services
- Interact with the Butler

5. Accessibility Controls

Users should be able to manually change their experience when needed.

Examples:

- Change color theme
- Change text presentation
- Enable or disable text-to-speech
- Adjust contrast
- Simplify the interface

6. Butler Interface

The front-end should provide a way for users to communicate with the Butler through voice.

The Butler should be able to receive commands such as:

- "Change the page to baby pink."
- "Make the text easier to read."
- "Increase the text size."
- "Turn on text-to-speech."
- "Make the page simpler."

The front-end should apply the changes requested by the Butler.

---

Back-End

The back-end is responsible for storing user information, managing preferences, handling requests, and connecting the different parts of IRIS.

Technology

- Go

The back-end should have:

1. User Management

The back-end should handle:

- User registration
- User login
- User information
- User profiles

2. Accessibility Preferences

The back-end should store each user's accessibility preferences.

This information should be available whenever the user accesses IRIS.

For example:

User
 ├── Accessibility needs
 ├── Color preferences
 ├── Reading preferences
 ├── Text-to-speech preference
 ├── Contrast preference
 └── Interface preferences

3. Personalization

When a user logs in, the back-end should provide their saved preferences to the front-end.

The front-end can then adapt the interface automatically.

4. Services

The back-end should manage information related to:

- Health services
- Community services
- Events

5. User Events

The back-end should manage events associated with users.

This can include information about events a user has interacted with or selected.

6. API

The Go back-end should provide APIs that allow the front-end to:

- Create users
- Retrieve users
- Retrieve preferences
- Update preferences
- Retrieve services
- Retrieve events
- Update user information

---

Butler

The Butler is the voice-based assistant within IRIS.

The Butler acts as an interface between the user and the different parts of IRIS.

Instead of requiring the user to manually search through settings, they can tell the Butler what they need.

Example

The user says:

«"Hey Butler, change the page to baby pink."»

The Butler:

1. Receives the user's request.
2. Understands what the user wants.
3. Identifies the relevant interface setting.
4. Communicates the requested change.
5. The front-end changes the page.

Using Sign-Up Information

The Butler should also use information collected during sign-up.

For example, if a user identifies that they have dyslexia and prefers a specific reading experience, IRIS already has that information.

The Butler can use the user's existing preferences when responding to requests.

For example:

«"Hey Butler, make this easier to read."»

The Butler can use the user's saved accessibility preferences to determine how the interface should change.

Butler Capabilities

The Butler should be able to:

- Change interface preferences
- Change colors
- Change text presentation
- Adjust text size
- Enable text-to-speech
- Enable high contrast
- Simplify the interface
- Apply saved accessibility preferences
- Respond to user requests
- Communicate with the different parts of IRIS

---

IRIS Flow

                 USER
                   │
                   ▼
              SIGN UP
                   │
                   ▼
        ACCESSIBILITY NEEDS
          & PREFERENCES
                   │
                   ▼
              BACK-END
                   │
                   ▼
          SAVED USER PROFILE
                   │
                   ▼
             FRONT-END
                   │
                   ▼
       PERSONALIZED INTERFACE
                   │
                   │
              ┌────┴────┐
              │         │
              ▼         ▼
           USER      BUTLER
                       │
                       ▼
                 USER REQUEST
                       │
                       ▼
                IRIS ACTION
                       │
                       ▼
              FRONT-END CHANGE

Goal

IRIS should provide an experience where accessibility is not something the user has to repeatedly configure.

The user's needs are captured during sign-up, remembered by IRIS, and can be changed at any time through the interface or the Butler.

IRIS — Your world, your way.