# Cybersecurity News Aggregator

This is a sample Android application that demonstrates a modern Android architecture using Jetpack Compose, Hilt, Retrofit, and other popular libraries. It uses a Go API for the backend which is hosted as a Render Instance.
Till now only the API part is complete. Android app is still under development. There's more cool stuff coming up, things like Cloud Functions and Gen AI, stay tuned :) 
## Getting Started

1.  **Clone the repository:**
    ```
    git clone https://github.com/your-username/ad-aware-news-aggregator.git
    ```
2.  **Open the project in Android Studio.**
3.  **Add your API key:**
    *   Open the `di/AppModule.kt` file.
    *   Replace `"https://your.api.url/"` with your actual API base URL.
4.  **Run the app.**

## TODO (Ad integration for revenue)

*   Implement the UI for a single news article item in `NewsList.kt`.
*   Implement the ad banner UI in `AdBanner.kt`.
*   Implement proper error handling in the UI.
*   Replace the mock ad API with a real implementation.
