package com.codegrey.adawarenewsaggregator.ui.news

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.DateRange
import androidx.compose.material3.Button
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.DateRangePicker
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.material3.rememberDateRangePickerState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import com.codegrey.adawarenewsaggregator.data.NetworkResult
import java.util.Date

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun NewsScreen(viewModel: NewsViewModel = hiltViewModel()) {
    val newsState by viewModel.newsState.collectAsState()
    val adState by viewModel.adState.collectAsState()
    val threatScoreState by viewModel.threatScoreState.collectAsState()
    var showDatePicker by remember { mutableStateOf(false) }
    val dateRangePickerState = rememberDateRangePickerState()

    Column(modifier = Modifier.fillMaxSize()) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(horizontal = 16.dp, vertical = 8.dp),
            horizontalArrangement = Arrangement.Space-Between,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Text(text = "ThreatFeed")
            Button(onClick = { showDatePicker = true }) {
                Icon(Icons.Default.DateRange, contentDescription = "Select Date Range")
            }
        }

        Threatbar(threatScoreResult = threatScoreState)

        if (showDatePicker) {
            androidx.compose.material3.DatePickerDialog(
                onDismissRequest = { showDatePicker = false },
                confirmButton = {
                    TextButton(onClick = {
                        showDatePicker = false
                        val start = dateRangePickerState.selectedStartDateMillis?.let { Date(it) }
                        val end = dateRangePickerState.selectedEndDateMillis?.let { Date(it) }
                        if (start != null && end != null) {
                            viewModel.setDateRange(start, end)
                        }
                    }) {
                        Text("OK")
                    }
                },
                dismissButton = {
                    TextButton(onClick = { showDatePicker = false }) {
                        Text("Cancel")
                    }
                }
            ) {
                DateRangePicker(state = dateRangePickerState)
            }
        }

        when (val state = newsState) {
            is NetworkResult.Loading -> {
                Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center
                ) {
                    CircularProgressIndicator()
                }
            }

            is NetworkResult.Success -> {
                NewsList(articles = state.data!!)
            }

            is NetworkResult.Error -> {
                Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center
                ) {
                    Text(text = state.message ?: "An error occurred")
                }
            }
        }

        when (val state = adState) {
            is NetworkResult.Loading -> { /* TODO: Handle ad loading state */
            }

            is NetworkResult.Success -> {
                AdBanner(ad = state.data!!)
            }

            is NetworkResult.Error -> { /* TODO: Handle ad error state */
            }
        }
    }
}
