package com.codegrey.adawarenewsaggregator.data

import retrofit2.http.GET

interface AdApiService {
    @GET("ad")
    suspend fun getAd(): Ad
}
