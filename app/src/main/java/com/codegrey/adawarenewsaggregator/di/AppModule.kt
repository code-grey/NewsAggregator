package com.codegrey.adawarenewsaggregator.di

import com.codegrey.adawarenewsaggregator.data.AdApiService
import com.codegrey.adawarenewsaggregator.data.NewsApiService
import com.codegrey.adawarenewsaggregator.domain.NewsRepository
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
object AppModule {

    @Provides
    @Singleton
    fun provideRetrofit(): Retrofit = Retrofit.Builder()
        .baseUrl("http://10.0.2.2:8080/")
        .addConverterFactory(GsonConverterFactory.create())
        .build()

    @Provides
    @Singleton
    fun provideNewsApiService(retrofit: Retrofit): NewsApiService = retrofit.create(NewsApiService::class.java)

    @Provides
    @Singleton
    fun provideAdApiService(retrofit: Retrofit): AdApiService = retrofit.create(AdApiService::class.java)

    @Provides
    @Singleton
    fun provideNewsRepository(newsApiService: NewsApiService, adApiService: AdApiService): NewsRepository =
        NewsRepository(newsApiService, adApiService)
}
