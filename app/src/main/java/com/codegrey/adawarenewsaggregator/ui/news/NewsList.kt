package com.codegrey.adawarenewsaggregator.ui.news

import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.runtime.Composable
import com.codegrey.adawarenewsaggregator.data.NewsArticle

@Composable
fun NewsList(articles: List<NewsArticle>) {
    LazyColumn {
        items(articles) {
            // TODO: Create a Composable for a single news article item
        }
    }
}
