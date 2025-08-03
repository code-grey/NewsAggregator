package com.codegrey.adawarenewsaggregator.data

import retrofit2.http.GET

interface NewsApiService {
    @GET("news")
    suspend fun getNews(): List<NewsArticle>
}
