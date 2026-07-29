package com.talhacan.kuranai.ui

import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.darkColorScheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color

val KuranGold = Color(0xFFE4B62D)
val KuranGoldSoft = Color(0xFFF7D968)
val KuranBackground = Color(0xFF07120F)
val KuranSurface = Color(0xFF101C18)
val KuranSurfaceRaised = Color(0xFF18251F)
val KuranText = Color(0xFFFFF8E7)
val KuranMuted = Color(0xFFA8B3AD)

private val DarkColors = darkColorScheme(
    primary = KuranGold,
    onPrimary = Color(0xFF171105),
    secondary = KuranGoldSoft,
    background = KuranBackground,
    onBackground = KuranText,
    surface = KuranSurface,
    onSurface = KuranText,
    surfaceVariant = KuranSurfaceRaised,
    onSurfaceVariant = KuranMuted,
)

private val LightColors = lightColorScheme(
    primary = Color(0xFF7A5A00),
    onPrimary = Color.White,
    secondary = Color(0xFF006B57),
    background = Color(0xFFFFFBF1),
    onBackground = Color(0xFF17150F),
    surface = Color.White,
    onSurface = Color(0xFF17150F),
    surfaceVariant = Color(0xFFF0EAD9),
    onSurfaceVariant = Color(0xFF5F5A4D),
)

@Composable
fun KuranAITheme(
    darkTheme: Boolean = isSystemInDarkTheme(),
    content: @Composable () -> Unit,
) {
    MaterialTheme(
        colorScheme = if (darkTheme) DarkColors else LightColors,
        content = content,
    )
}
