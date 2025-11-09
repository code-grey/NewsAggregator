package com.codegrey.adawarenewsaggregator.domain

import com.codegrey.adawarenewsaggregator.data.Ad
import com.codegrey.adawarenewsaggregator.data.AdApiService
import com.codegrey.adawarenewsaggregator.data.NewsApiService
import com.codegrey.adawarenewsaggregator.data.NetworkResult
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.flow
import javax.inject.Inject

class NewsRepository @Inject constructor(
    private val newsApiService: NewsApiService,
    private val adApiService: AdApiService
) {

    fun getNews(): Flow<NetworkResult<List<com.codegrey.adawarenewsaggregator.data.NewsArticle>>> = flow {
        emit(NetworkResult.Loading())
        try {
            val news = newsApiService.getNews()
            emit(NetworkResult.Success(news))
        } catch (e: Exception) {
            emit(NetworkResult.Error(e.message ?: "An unknown error occurred"))
        }
    }

    fun getAd(): Flow<NetworkResult<Ad>> = flow {
        emit(NetworkResult.Loading())
        try {
            val ad = adApiService.getAd()
            emit(NetworkResult.Success(ad))
        } catch (e: Exception) {
            emit(NetworkResult.Error(e.message ?: "An unknown error occurred"))
        }
    }
}
