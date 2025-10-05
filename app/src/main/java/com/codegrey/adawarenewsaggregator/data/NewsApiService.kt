package com.codegrey.adawarenewsaggregator.data

import retrofit2.http.GET
import retrofit2.http.Query

interface NewsApiService {
    @GET("news")
    suspend fun getNews(
        @Query("category") category: String? = null,
        @Query("start") startDate: String? = null,
        @Query("end") endDate: String? = null
    ): List<NewsArticle>

    @GET("today-threat")
    suspend fun getTodayThreat(): ThreatScore
}
