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
import java.text.SimpleDateFormat
import java.util.Date
import java.util.Locale
import javax.inject.Inject

@HiltViewModel
class NewsViewModel @Inject constructor(
    private val newsRepository: NewsRepository
) : ViewModel() {

    private val _newsState = MutableStateFlow<NetworkResult<List<com.codegrey.adawarenewsaggregator.data.NewsArticle>>>(NetworkResult.Loading())
    val newsState: StateFlow<NetworkResult<List<com.codegrey.adawarenewsaggregator.data.NewsArticle>>> = _newsState

    private val _adState = MutableStateFlow<NetworkResult<Ad>>(NetworkResult.Loading())
    val adState: StateFlow<NetworkResult<Ad>> = _adState

    private val _startDate = MutableStateFlow<String?>(null)
    val startDate: StateFlow<String?> = _startDate

    private val _endDate = MutableStateFlow<String?>(null)
    val endDate: StateFlow<String?> = _endDate

    private val _category = MutableStateFlow("Cybersecurity")
    val category: StateFlow<String> = _category

    private val _threatScoreState = MutableStateFlow<NetworkResult<com.codegrey.adawarenewsaggregator.data.ThreatScore>>(NetworkResult.Loading())
    val threatScoreState: StateFlow<NetworkResult<com.codegrey.adawarenewsaggregator.data.ThreatScore>> = _threatScoreState

    init {
        fetchNews()
        fetchAd()
        fetchTodayThreat()
    }

    fun fetchNews(category: String = _category.value, startDate: String? = _startDate.value, endDate: String? = _endDate.value) {
        _category.value = category
        newsRepository.getNews(category, startDate, endDate).onEach {
            _newsState.value = it
        }.launchIn(viewModelScope)
    }

    private fun fetchAd() {
        newsRepository.getAd().onEach {
            _adState.value = it
        }.launchIn(viewModelScope)
    }

    fun setDateRange(start: Date, end: Date) {
        val formatter = SimpleDateFormat("yyyy-MM-dd", Locale.getDefault())
        _startDate.value = formatter.format(start)
        _endDate.value = formatter.format(end)
        fetchNews(category = _category.value, startDate = _startDate.value, endDate = _endDate.value)
    }

    private fun fetchTodayThreat() {
        newsRepository.getTodayThreat().onEach {
            _threatScoreState.value = it
        }.launchIn(viewModelScope)
    }
}
