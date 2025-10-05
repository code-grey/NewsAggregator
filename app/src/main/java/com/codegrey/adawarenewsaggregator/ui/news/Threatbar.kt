package com.codegrey.adawarenewsaggregator.ui.news

import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.LinearProgressIndicator
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import com.codegrey.adawarenewsaggregator.data.NetworkResult
import com.codegrey.adawarenewsaggregator.data.ThreatScore

@Composable
fun Threatbar(threatScoreResult: NetworkResult<ThreatScore>) {
    when (threatScoreResult) {
        is NetworkResult.Loading -> {
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp),
                contentAlignment = Alignment.Center
            ) {
                LinearProgressIndicator(modifier = Modifier.fillMaxWidth())
            }
        }
        is NetworkResult.Success -> {
            threatScoreResult.data?.let {
                ThreatScoreIndicator(threatScore = it)
            }
        }
        is NetworkResult.Error -> {
            // Handle error state, maybe show a message
        }
    }
}

@Composable
fun ThreatScoreIndicator(threatScore: ThreatScore) {
    val progress by animateFloatAsState(targetValue = (threatScore.score / 5.0).toFloat(), label = "")

    val color = when {
        progress <= 0.3 -> Color.Green
        progress <= 0.7 -> Color.Yellow
        else -> Color.Red
    }

    Column(
        modifier = Modifier
            .fillMaxWidth()
            .padding(16.dp)
            .clip(RoundedCornerShape(12.dp))
            .background(MaterialTheme.colorScheme.surfaceVariant)
            .padding(16.dp)
    ) {
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.Space-Between,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Text(
                text = "Today's Threat Level",
                style = MaterialTheme.typography.titleMedium
            )
            Text(
                text = threatScore.phrase,
                style = MaterialTheme.typography.titleMedium,
                color = color
            )
        }
        Spacer(modifier = Modifier.height(8.dp))
        Row(verticalAlignment = Alignment.CenterVertically) {
            LinearProgressIndicator(
                progress = { progress },
                modifier = Modifier.weight(1f).height(12.dp).clip(RoundedCornerShape(6.dp)),
                color = color
            )
            Spacer(modifier = Modifier.width(8.dp))
            Text(
                text = String.format("%.2f", threatScore.score),
                style = MaterialTheme.typography.bodyMedium
            )
        }
    }
}