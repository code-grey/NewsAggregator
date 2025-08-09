package com.codegrey.adawarenewsaggregator.ui.news

import androidx.compose.foundation.border
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.codegrey.adawarenewsaggregator.data.Ad

@Composable
fun AdBanner(ad: Ad) {
    Box(
        modifier = Modifier
            .fillMaxWidth()
            .padding(8.dp)
            .border(1.dp, MaterialTheme.colorScheme.outline),
        contentAlignment = Alignment.Center
    ) {
        Text(
            text = "Your Ad Here", // Placeholder text
            modifier = Modifier.padding(16.dp),
            style = MaterialTheme.typography.bodyLarge,
            color = MaterialTheme.colorScheme.onSurface
        )
    }
}