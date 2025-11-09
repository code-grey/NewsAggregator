package com.codegrey.adawarenewsaggregator.ui.news

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.hilt.navigation.compose.hiltViewModel
import com.codegrey.adawarenewsaggregator.data.NetworkResult

@Composable
fun NewsScreen(viewModel: NewsViewModel = hiltViewModel()) {
    val newsState by viewModel.newsState.collectAsState()
    val adState by viewModel.adState.collectAsState()

    Column(modifier = Modifier.fillMaxSize()) {
        when (val state = newsState) {
            is NetworkResult.Loading -> {
                Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
                    CircularProgressIndicator()
                }
            }
            is NetworkResult.Success -> {
                NewsList(articles = state.data!!)
            }
            is NetworkResult.Error -> {
                Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
                    Text(text = state.message ?: "An error occurred")
                }
            }
        }

        when (val state = adState) {
            is NetworkResult.Loading -> { /* TODO: Handle ad loading state */ }
            is NetworkResult.Success -> {
                AdBanner(ad = state.data!!)
            }
            is NetworkResult.Error -> { /* TODO: Handle ad error state */ }
        }
    }
}
