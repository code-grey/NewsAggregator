package com.codegrey.adawarenewsaggregator.ui.news

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.codegrey.adawarenewsaggregator.data.Ad
import com.codegrey.adawarenewsaggregator.data.NetworkResult
import com.codegrey.adawarenewsaggregator.domain.NewsRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.launchIn
import kotlinx.coroutines.flow.onEach
import javax.inject.Inject

@HiltViewModel
class NewsViewModel @Inject constructor(
    private val newsRepository: NewsRepository
) : ViewModel() {

    private val _newsState = MutableStateFlow<NetworkResult<List<com.codegrey.adawarenewsaggregator.data.NewsArticle>>>(NetworkResult.Loading())
    val newsState: StateFlow<NetworkResult<List<com.codegrey.adawarenewsaggregator.data.NewsArticle>>> = _newsState

    private val _adState = MutableStateFlow<NetworkResult<Ad>>(NetworkResult.Loading())
    val adState: StateFlow<NetworkResult<Ad>> = _adState

    init {
        fetchNews()
        fetchAd()
    }

    private fun fetchNews() {
        newsRepository.getNews().onEach {
            _newsState.value = it
        }.launchIn(viewModelScope)
    }

    private fun fetchAd() {
        newsRepository.getAd().onEach {
            _adState.value = it
        }.launchIn(viewModelScope)
    }
}
